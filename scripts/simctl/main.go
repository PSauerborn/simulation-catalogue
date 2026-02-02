package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/urfave/cli/v3"
)

const (
	defaultAPIURL = "https://simulation-catalogue.s31-software.com/api/v1"
)

func main() {
	apiKey := os.Getenv("SIM_API_KEY")
	if len(apiKey) == 0 {
		panic("SIM_API_KEY environment variable is required")
	}
	// create new API client
	client := NewAPIClient(defaultAPIURL, apiKey)

	cmd := &cli.Command{
		Name:  "simctl",
		Usage: "Manager for simulations via the API",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "url",
				Aliases: []string{"u"},
				Usage:   "Base URL of the API",
				Value:   defaultAPIURL,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "health",
				Usage: "Check the health of the API",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return checkHealth(ctx, client, cmd)
				},
			},
			{
				Name:  "create",
				Usage: "Create a new simulation from a JSON file",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return createSimulation(ctx, client, cmd)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Usage:    "Path to the JSON simulation definition file",
						Required: true,
					},
				},
			},
			{
				Name:  "update-binary",
				Usage: "Update the binary for a simulation",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return updateSimulationBinary(ctx, client, cmd)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "ID of the simulation",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "arch",
						Usage:    "CPU architecture (e.g., arm64, amd64)",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Usage:    "Path to the binary file",
						Required: true,
					},
				},
			},
			{
				Name:  "delete",
				Usage: "Delete a simulation",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return deleteSimulation(ctx, client, cmd)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "ID of the simulation",
						Required: true,
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func createSimulation(ctx context.Context, client *APIClient, cmd *cli.Command) error {
	filePath := cmd.String("file")
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Validate JSON structure
	var simReq NewSimulationRequest
	if err := json.Unmarshal(fileContent, &simReq); err != nil {
		return fmt.Errorf("failed to unmarshal simulation JSON: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(simReq); err != nil {
		return fmt.Errorf("failed to validate simulation request: %w", err)
	}

	url := fmt.Sprintf("%s/admin/simulations", client.BaseUrl)
	resp, err := client.DoRequest(ctx, "POST", url, bytes.NewBuffer(fileContent))
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create simulation (status %d): %s", resp.StatusCode, string(body))
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("created but failed to decode response: %w", err)
	}

	fmt.Printf("Simulation created successfully. ID: %s\n", result["id"])
	return nil
}

func updateSimulationBinary(ctx context.Context, client *APIClient, cmd *cli.Command) error {
	simID := cmd.String("id")
	arch := cmd.String("arch")
	filePath := cmd.String("file")

	allowedArchs := []string{"amd64", "arm64"}
	if !slices.Contains(allowedArchs, arch) {
		return fmt.Errorf("invalid architecture: %s", arch)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open binary file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("binary", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}
	writer.Close()

	url := fmt.Sprintf("%s/admin/simulations/%s/binary/%s", client.BaseUrl, simID, arch)
	req, err := http.NewRequestWithContext(ctx, "PUT", url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", client.APIKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update binary (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	fmt.Println("Simulation binary updated successfully.")
	return nil
}

func deleteSimulation(ctx context.Context, client *APIClient, cmd *cli.Command) error {
	simID := cmd.String("id")
	url := fmt.Sprintf("%s/admin/simulations/%s", client.BaseUrl, simID)

	resp, err := client.DoRequest(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete simulation (status %d): %s", resp.StatusCode, string(body))
	}

	fmt.Printf("Simulation %s deleted successfully.\n", simID)
	return nil
}

func checkHealth(ctx context.Context, client *APIClient, cmd *cli.Command) error {
	url := fmt.Sprintf("%s/public/health", client.BaseUrl)
	resp, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("health check failed (status %d): %s", resp.StatusCode, string(body))
	}

	fmt.Println("API is healthy.")
	return nil
}
