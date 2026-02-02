package main

import (
	"context"
	"io"
	"net/http"
	"time"
)

// Data structures copied/adapted from api/types.go

type APIClient struct {
	BaseUrl string
	APIKey string
	http.Client
}

func (c *APIClient) DoRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	return c.Client.Do(req)
}

func NewAPIClient(baseUrl string, apiKey string) *APIClient {
	baseClient := http.Client{
		Timeout: 10 * time.Second,
	}
	return &APIClient{
		BaseUrl: baseUrl,
		APIKey: apiKey,
		Client: baseClient,
	}
}

// ParameterTypeEnum defines the valid data types for simulation parameters.
type ParameterTypeEnum string

// Supported parameter types for simulation inputs.
const (
	SimulationParameterTypeFloat  ParameterTypeEnum = "float"
	SimulationParameterTypeInt    ParameterTypeEnum = "int"
	SimulationParameterTypeString ParameterTypeEnum = "string"
	SimulationParameterTypeBool   ParameterTypeEnum = "bool"
	SimulationParameterTypeVector ParameterTypeEnum = "vector"
)

// SimulationParameter defines an input parameter for a simulation
type SimulationParameter struct {
	Name        string            `json:"name" validate:"required"`
	Description string            `json:"description" validate:"required"`
	Type        ParameterTypeEnum `json:"type" validate:"required"`
	Default     interface{}       `json:"default" validate:"required"`
}

// OutputTypeEnum defines the valid output types for simulation results.
type OutputTypeEnum string

// Supported output types for simulation results.
const (
	OutputTypeTrajectory OutputTypeEnum = "trajectory"
)

// SimulationOutput defines an output produced by a simulation
type SimulationOutput struct {
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description" validate:"required"`
	Type        OutputTypeEnum `json:"type" validate:"required"`
}

// SimulationMeta contains the metadata for a simulation
type SimulationMeta struct {
	Id          string                `json:"id" validate:"required"`
	Name        string                `json:"name" validate:"required"`
	Description string                `json:"description" validate:"required"`
	Model       string                `json:"model" validate:"required"`
	Parameters  []SimulationParameter `json:"parameters" validate:"required"`
	Outputs     []SimulationOutput    `json:"outputs" validate:"required"`
	CreatedAt   time.Time             `json:"created_at" validate:"required"`
	UpdatedAt   time.Time             `json:"updated_at" validate:"required"`
}

// NewSimulationRequest is the request body for creating a new simulation.
type NewSimulationRequest struct {
	Name        string                `json:"name" validate:"required"`
	Description string                `json:"description" validate:"required"`
	Model       string                `json:"model" validate:"required"`
	Parameters  []SimulationParameter `json:"parameters" validate:"required"`
	Outputs     []SimulationOutput    `json:"outputs" validate:"required"`
}

// PatchSimulationRequest is the request body for updating simulation metadata.
type PatchSimulationRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}
