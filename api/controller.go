package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Controller struct {
	db     Persistence
	events EventBroker
	config *Config
}

// NewController creates a new Controller instance with the provided configuration,
// database persistence layer, and event broker for publishing simulation events.
func NewController(config *Config, db Persistence, events EventBroker) *Controller {
	return &Controller{
		config: config,
		db:     db,
		events: events,
	}
}

// Version returns the current API version from the configuration.
func (cnt *Controller) Version(c *gin.Context) JSONResponse {
	return JSONResponse{Code: http.StatusOK, Body: gin.H{
		"version": cnt.config.Version,
	}}
}

// HealthCheck performs a health check on the database connection
// and returns the current health status of the API.
func (cnt *Controller) HealthCheck(c *gin.Context) JSONResponse {
	if err := cnt.db.HealthCheck(); err != nil {
		log.WithError(err).Error("database health check failed")
		return JSONResponse{Code: http.StatusInternalServerError, Body: gin.H{
			"error": "Internal Server Error",
		}}
	}

	return JSONResponse{Code: http.StatusOK, Body: gin.H{
		"status": "OK",
	}}
}

// GetClient retrieves client information based on the client ID
// stored in the request context. Returns nil if no client ID is present.
func (cnt *Controller) GetClient(c *gin.Context) JSONResponse {
	clientId := c.Request.Header.Get("X-Client-Id")
	if clientId == "" {
		log.Warn("received request for client without client_id")
		return JSONResponse{Code: http.StatusOK, Body: gin.H{
			"client": nil,
		}}
	}

	client, err := cnt.db.GetClient(clientId)
	var errNotFound ErrClientNotInitialized
	if err != nil && errors.As(err, &errNotFound) {
		log.WithError(err).Error("failed to get client")
		return JSONResponse{Code: http.StatusOK, Body: gin.H{
			"client": nil,
		}}
	} else if err != nil {
		log.WithError(err).Error("failed to get client")
		return JSONResponse{Code: http.StatusInternalServerError, Body: gin.H{
			"error": "Internal Server Error",
		}}
	}
	return JSONResponse{Code: http.StatusOK, Body: gin.H{
		"client": client,
	}}
}

// InitClient creates a new client session and returns the generated client ID.
// This ID is used to track simulation runs for anonymous users.
func (cnt *Controller) InitClient(c *gin.Context) JSONResponse {
	id, err := cnt.db.InitClient()
	if err != nil {
		log.WithError(err).Error("failed to init client")
		return JSONResponse{Code: http.StatusInternalServerError, Body: gin.H{
			"error": "Internal Server Error",
		}}
	}
	return JSONResponse{Code: http.StatusOK, Body: gin.H{
		"id": id,
	}}
}

// ListSimulations returns a list of all available simulations
// with their metadata including name, description, and parameters.
func (cnt *Controller) ListSimulations(c *gin.Context) JSONResponse {
	simulations, err := cnt.db.ListSimulations()
	if err != nil {
		log.WithError(err).Error("failed to list simulations")
		return JSONResponse{Code: http.StatusInternalServerError, Body: gin.H{
			"error": "Internal Server Error",
		}}
	}

	return JSONResponse{Code: http.StatusOK, Body: gin.H{
		"data": simulations,
	}}
}

// GetSimulationMeta retrieves the metadata for a specific simulation by ID.
// Returns 404 if the simulation is not found.
func (cnt *Controller) GetSimulationMeta(c *gin.Context) JSONResponse {
	simulation, err := cnt.db.GetSimulationMeta(c.Param("id"))
	var errNotFound ErrSimulationNotFound
	if err != nil && errors.As(err, &errNotFound) {
		log.WithError(err).Error("failed to get simulation meta")
		return JSONResponse{Code: http.StatusNotFound, Body: gin.H{
			"error": "Not Found",
		}}
	} else if err != nil {
		log.WithError(err).Error("failed to get simulation meta")
		return JSONResponse{Code: http.StatusInternalServerError, Body: gin.H{
			"error": "Internal Server Error",
		}}
	}

	return JSONResponse{Code: http.StatusOK, Body: simulation}
}

// GetSimulationBinary downloads the compiled simulation binary for a specific
// CPU architecture. The binary is returned as an octet-stream.
// Returns 404 if the simulation or binary for the specified architecture is not found.
func (cnt *Controller) GetSimulationBinary(c *gin.Context) {
	simId := c.Param("id")
	cpuArch := c.Param("cpu_architecture")

	_, err := cnt.db.GetSimulationMeta(simId)
	var errNotFound ErrSimulationNotFound
	if err != nil && errors.As(err, &errNotFound) {
		log.WithError(err).Error("failed to get simulation binary")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
		return
	} else if err != nil {
		log.WithError(err).Error("failed to get simulation binary")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	binary, err := cnt.db.GetSimulationBinary(simId, cpuArch)
	var errBinaryNotFound ErrBinaryNotFound
	if err != nil && errors.As(err, &errBinaryNotFound) {
		log.WithError(err).Error("failed to get simulation binary")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
		return
	} else if err != nil {
		log.WithError(err).Error("failed to get simulation binary")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", binary)
}

// GetSimulationRun retrieves the current simulation run status for a client.
// Requires a valid client ID in the request context.
func (cnt *Controller) GetSimulationRun(c *gin.Context) JSONResponse {
	clientId := c.GetString("client_id")
	if clientId == "" {
		log.Warn("received request for simulation run without client_id")
		return JSONResponse{Code: http.StatusBadRequest, Body: gin.H{
			"error": "Bad Request",
		}}
	}

	run, err := cnt.db.GetSimulationRun(clientId)
	if err != nil {
		log.WithError(err).Error("failed to get simulation run")
		return JSONResponse{Code: http.StatusInternalServerError, Body: gin.H{
			"error": "Internal Server Error",
		}}
	}
	run.Output = nil

	return JSONResponse{Code: http.StatusOK, Body: run}
}

// GetSimulationOutput retrieves the output of a completed simulation run.
// The output format can be specified via the "format" query parameter:
// - "json" (default): Returns CSV data converted to JSON
// - "zip": Returns the raw zip archive containing CSV files
func (cnt *Controller) GetSimulationOutput(c *gin.Context) {
	clientId := c.GetString("client_id")

	format := strings.ToLower(c.Query("format"))
	if format == "" {
		format = "json"
	}

	allowedFormats := []string{
		"json",
		"zip",
	}

	if !slices.Contains(allowedFormats, format) {
		log.WithFields(log.Fields{
			"format":    format,
			"client_id": clientId,
		}).Warn("invalid output format requested")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid format",
		})
		return
	}

	log.WithFields(log.Fields{
		"client_id": clientId,
		"format":    format,
	}).Info("requesting simulation output")

	// get simulation run
	run, err := cnt.db.GetSimulationRun(clientId)
	if err != nil {
		log.WithError(err).Error("failed to get simulation run")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	if run.Status != "completed" || run.Output == nil {
		c.JSON(http.StatusOK, gin.H{
			"output": nil,
		})
		return
	}

	switch format {
	case "json":
		log.Info("extracting CSV files from zip")

		csvFiles, err := ExtractCSVFilesFromZip(run.Output)
		if err != nil {
			log.WithError(err).Error("failed to extract CSV files from zip")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		log.WithFields(log.Fields{
			"client_id": clientId,
			"files":     len(csvFiles),
		}).Info("extracted CSV files from zip")

		jsonFiles := make(map[string][]map[string]interface{})
		for name, file := range csvFiles {
			records, err := CSVFileToJson(file)
			if err != nil {
				log.WithError(err).Error("failed to parse CSV file")
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
				return
			}

			jsonFiles[name] = DownSampleDataset(records)
		}

		log.WithFields(log.Fields{
			"client_id": clientId,
			"files":     len(jsonFiles),
		}).Info("parsed CSV files to JSON")

		c.JSON(http.StatusOK, gin.H{
			"output": jsonFiles,
		})

	case "zip":
		// output artifacts are already zipped CSV archives
		// so we can just return them as is
		c.Data(http.StatusOK, "application/zip", run.Output)

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid format",
		})
	}
}

// RunSimulation initiates a new simulation run for the client.
// Publishes a simulation event to the event broker for async processing.
// Returns 409 Conflict if the client already has an active simulation run.
func (cnt *Controller) RunSimulation(c *gin.Context) JSONResponse {
	clientId := c.GetString("client_id")

	var request RunSimulationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.WithError(err).Error("failed to bind run simulation request")
		return JSONResponse{Code: http.StatusBadRequest, Body: gin.H{
			"error": "Invalid request",
		}}
	}

	log.WithFields(log.Fields{
		"client_id":     clientId,
		"simulation_id": request.SimulationId,
		"parameters":    request.Parameters,
	}).Info("received request to run simulation")

	meta, err := cnt.db.GetSimulationMeta(request.SimulationId)

	var errNotFound ErrSimulationNotFound
	if err != nil && errors.As(err, &errNotFound) {
		log.WithError(err).Error("failed to get simulation meta")
		return JSONResponse{Code: http.StatusNotFound, Body: gin.H{
			"error": "Not Found",
		}}

	} else if err != nil {
		log.WithError(err).Error("failed to get simulation meta")
		return JSONResponse{Code: http.StatusInternalServerError, Body: gin.H{
			"error": "Internal Server Error",
		}}
	}

	errors := ValidateParameters(meta, request.Parameters)
	for _, err := range errors {
		log.WithError(err).Error("failed to validate parameter")
	}

	if len(errors) > 0 {
		log.WithFields(log.Fields{
			"client_id":     clientId,
			"simulation_id": request.SimulationId,
			"parameters":    request.Parameters,
			"error_count":   len(errors),
		}).Error("failed to validate parameters")

		return JSONResponse{Code: http.StatusBadRequest, Body: gin.H{
			"error": "Invalid parameters",
		}}
	}

	current, err := cnt.db.GetSimulationRun(clientId)
	if err != nil {
		log.WithError(err).Error("failed to get simulation run")
		return JSONResponse{Code: http.StatusInternalServerError, Body: gin.H{
			"error": "Internal Server Error",
		}}
	}

	allowedStatuses := []string{
		"pending",
		"completed",
		"failed",
	}

	if !slices.Contains(allowedStatuses, current.Status) {
		log.WithFields(log.Fields{
			"client_id": clientId,
			"status":    current.Status,
		}).Warn("client already has an active simulation run")
		return JSONResponse{Code: http.StatusConflict, Body: gin.H{
			"error": "Conflict",
		}}
	}

	id, err := cnt.db.CreateSimulationRun(clientId, request.SimulationId, request.Parameters)
	if err != nil {
		log.WithError(err).Error("failed to create simulation run")
		return JSONResponse{Code: http.StatusInternalServerError, Body: gin.H{
			"error": "Internal Server Error",
		}}
	}

	event, _ := json.Marshal(map[string]interface{}{
		"client_id":     clientId,
		"simulation_id": request.SimulationId,
		"parameters":    request.Parameters,
	})

	if err := cnt.events.Publish(event); err != nil {
		log.WithError(err).Error("failed to publish simulation run event")
		return JSONResponse{Code: http.StatusInternalServerError, Body: gin.H{
			"error": "Internal Server Error",
		}}
	}

	return JSONResponse{Code: http.StatusCreated, Body: gin.H{
		"id": id,
	}}
}

// CreateSimulation creates a new simulation with the provided metadata,
// parameters, and output definitions. Requires admin authentication.
func (cnt *Controller) CreateSimulation(c *gin.Context) JSONResponse {
	var body NewSimulationRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		log.WithError(err).Error("failed to bind create simulation request")
		return JSONResponse{Code: http.StatusBadRequest, Body: gin.H{
			"error": "Invalid request",
		}}
	}

	id, err := cnt.db.CreateSimulation(
		body.Name,
		body.Description,
		body.Model,
		body.Parameters,
		body.Outputs,
	)
	if err != nil {
		log.WithError(err).Error("failed to create simulation")
		return JSONResponse{Code: http.StatusInternalServerError, Body: gin.H{
			"error": "Internal Server Error",
		}}
	}

	return JSONResponse{Code: http.StatusCreated, Body: gin.H{
		"id": id,
	}}
}

// UpdateSimulationMeta updates the name and description of an existing simulation.
// Requires admin authentication. Returns 404 if the simulation is not found.
func (cnt *Controller) UpdateSimulationMeta(c *gin.Context) JSONResponse {
	simId := c.Param("id")

	var body PatchSimulationRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		log.WithError(err).Error("failed to bind update simulation meta request")
		return JSONResponse{Code: http.StatusBadRequest, Body: gin.H{
			"error": "Invalid request",
		}}
	}

	_, err := cnt.db.GetSimulationMeta(simId)
	var errNotFound ErrSimulationNotFound
	if err != nil && errors.As(err, &errNotFound) {
		log.WithError(err).Error("failed to get simulation meta")
		return JSONResponse{Code: http.StatusNotFound, Body: gin.H{
			"error": "Not Found",
		}}
	} else if err != nil {
		log.WithError(err).Error("failed to get simulation meta")
		return JSONResponse{Code: http.StatusInternalServerError, Body: gin.H{
			"error": "Internal Server Error",
		}}
	}

	err = cnt.db.UpdateSimulationMeta(simId, body.Name, body.Description, body.Model, body.Parameters, body.Outputs)
	if err != nil {
		log.WithError(err).Error("failed to update simulation meta")
		return JSONResponse{Code: http.StatusInternalServerError, Body: gin.H{
			"error": "Internal Server Error",
		}}
	}

	return JSONResponse{Code: http.StatusNoContent, Body: nil}
}

// UpdateSimulationBinary uploads or updates the compiled binary for a simulation
// for a specific CPU architecture. Requires admin authentication.
// The binary should be uploaded as multipart form data with the field name "binary".
func (cnt *Controller) UpdateSimulationBinary(c *gin.Context) JSONResponse {
	simId := c.Param("id")
	cpuArchitecture := c.Param("cpu_architecture")

	// get payload from form data
	fileHeader, err := c.FormFile("binary")
	if err != nil {
		log.WithError(err).Error("failed to get binary from form data")
		return JSONResponse{Code: http.StatusBadRequest, Body: gin.H{
			"error": "Invalid request",
		}}
	}

	contents, err := fileHeader.Open()
	if err != nil {
		log.WithError(err).Error("error opening file")
		return JSONResponse{Code: http.StatusBadRequest, Body: gin.H{
			"error": "Invalid request",
		}}
	}
	defer contents.Close()

	blob, err := io.ReadAll(contents)
	if err != nil {
		log.WithError(err).Error("error reading file")
		return JSONResponse{Code: http.StatusBadRequest, Body: gin.H{
			"error": "Invalid request",
		}}
	}

	_, err = cnt.db.GetSimulationMeta(simId)
	var errNotFound ErrSimulationNotFound
	if err != nil && errors.As(err, &errNotFound) {
		log.WithError(err).Error("failed to get simulation meta")
		return JSONResponse{Code: http.StatusNotFound, Body: gin.H{
			"error": "Not Found",
		}}
	} else if err != nil {
		log.WithError(err).Error("failed to get simulation meta")
		return JSONResponse{Code: http.StatusInternalServerError, Body: gin.H{
			"error": "Internal Server Error",
		}}
	}

	if err := cnt.db.UpdateSimulationBinary(simId, cpuArchitecture, blob); err != nil {
		log.WithError(err).Error("failed to update simulation binary")
		return JSONResponse{Code: http.StatusInternalServerError, Body: gin.H{
			"error": "Internal Server Error",
		}}
	}

	return JSONResponse{Code: http.StatusNoContent, Body: nil}
}

// DeleteSimulation permanently removes a simulation and all associated data.
// Requires admin authentication.
func (cnt *Controller) DeleteSimulation(c *gin.Context) JSONResponse {
	id := c.Param("id")

	if err := cnt.db.DeleteSimulation(id); err != nil {
		log.WithError(err).Error("failed to delete simulation")
		return JSONResponse{Code: http.StatusInternalServerError, Body: gin.H{
			"error": "Internal Server Error",
		}}
	}

	return JSONResponse{Code: http.StatusNoContent, Body: nil}
}
