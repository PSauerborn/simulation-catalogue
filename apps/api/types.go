package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

// JSONResponse represents a standardized HTTP JSON response with a status code and body.
type JSONResponse struct {
	Code int
	Body interface{}
}

// Send writes the JSON response to the Gin context with the configured status code.
func (r JSONResponse) Send(c *gin.Context) {
	c.JSON(r.Code, r.Body)
}

// Client represents an anonymous user session for tracking simulation runs.
type Client struct {
	Id           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	LastActiveAt time.Time `json:"last_active_at"`
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

// SimulationParameter defines an input parameter for a simulation,
// including its name, description, and data type.
type SimulationParameter struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        ParameterTypeEnum `json:"type"`
}

// APIKey represents an authentication key for admin API access.
type APIKey struct {
	Id        string    `json:"id"`
	Key       string    `json:"key"`
	Owner     string    `json:"owner"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
}

// IsValid returns true if the API key is not revoked and has not expired.
func (a *APIKey) IsValid() bool {
	return !a.Revoked && a.ExpiresAt.After(time.Now())
}

// OutputTypeEnum defines the valid output types for simulation results.
type OutputTypeEnum string

// Supported output types for simulation results.
const (
	OutputTypeTrajectory OutputTypeEnum = "trajectory"
)

// SimulationOutput defines an output produced by a simulation,
// including its name, description, and data type.
type SimulationOutput struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        OutputTypeEnum `json:"type"`
}

// SimulationMeta contains the metadata for a simulation including its
// name, description, input parameters, and expected outputs.
type SimulationMeta struct {
	Id          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Parameters  []SimulationParameter `json:"parameters"`
	Outputs     []SimulationOutput    `json:"outputs"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

// SimulationRun represents an execution of a simulation with specific parameters.
// It tracks the run status, input parameters, and output data.
type SimulationRun struct {
	ClientId     string                 `json:"client_id"`
	SimulationId *string                `json:"simulation_id"`
	Parameters   map[string]interface{} `json:"parameters"`
	Output       []byte                 `json:"outputs"`
	Status       string                 `json:"status"`
	CreatedAt    time.Time              `json:"created_at"`
	CompletedAt  *time.Time             `json:"completed_at"`
}

// NewSimulationRequest is the request body for creating a new simulation.
type NewSimulationRequest struct {
	Name        string                `json:"name" binding:"required"`
	Description string                `json:"description" binding:"required"`
	Parameters  []SimulationParameter `json:"parameters" binding:"required"`
	Outputs     []SimulationOutput    `json:"outputs" binding:"required"`
}

// PatchSimulationRequest is the request body for updating simulation metadata.
type PatchSimulationRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// RunSimulationRequest is the request body for initiating a simulation run.
type RunSimulationRequest struct {
	SimulationId string                 `json:"simulation_id" binding:"required"`
	Parameters   map[string]interface{} `json:"parameters" binding:"required"`
}
