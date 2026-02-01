package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml/v2"
	log "github.com/sirupsen/logrus"
)

// ValidateParameter validates a single parameter value against its expected type.
// It handles JSON unmarshaling quirks where all numbers are float64.
// Returns an error if the value does not match the expected type.
func ValidateParameter(definedType ParameterTypeEnum, value interface{}) error {
	switch definedType {
	case SimulationParameterTypeInt:
		// check if int type
		if _, ok := value.(int); ok {
			return nil
		}
		// check if float type. note that json unmarshals integers as float64
		// so we need to check for float first, then check if it's an integer
		f, ok := value.(float64)
		if !ok {
			return fmt.Errorf("parameter must be an integer")
		}

		if f != float64(int(f)) {
			return fmt.Errorf("parameter must be an integer (no decimal part)")
		}

	case SimulationParameterTypeFloat:
		if _, ok := value.(float64); !ok {
			return fmt.Errorf("parameter must be a float")
		}

	case SimulationParameterTypeString:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("parameter must be a string")
		}
	case SimulationParameterTypeVector:
		if _, ok := value.([]float64); ok {
			return nil
		}
		// JSON arrays are unmarshaled as []interface{}
		vector, ok := value.([]interface{})
		if !ok {
			return fmt.Errorf("parameter must be a vector (array)")
		}
		// Validate that each element is a float
		for _, elem := range vector {
			if _, ok := elem.(float64); !ok {
				return fmt.Errorf("vector element must be a float")
			}
		}
	}
	return nil
}

// ValidateParameters validates all parameters in the given map against the
// simulation metadata definitions. Returns an error if any parameter is invalid.
func ValidateParameters(meta SimulationMeta, params map[string]interface{}) []error {
	errors := []error{}
	for _, param := range meta.Parameters {
		if err := ValidateParameter(param.Type, params[param.Name]); err != nil {
			errors = append(errors, fmt.Errorf("parameter %s is invalid: %s", param.Name, err.Error()))
		}
	}
	return errors
}

// CastParameters converts parameter values from their JSON-unmarshaled types
// to the correct Go types based on the simulation metadata definitions.
// Integers are cast from float64 to int, and vectors are converted from
// []interface{} to []float64.
func CastParameters(meta SimulationMeta, params map[string]interface{}) (map[string]interface{}, error) {
	castParams := make(map[string]interface{})
	for _, param := range meta.Parameters {
		val, ok := params[param.Name]
		if !ok {
			continue
		}

		switch param.Type {
		case SimulationParameterTypeInt:
			if f, isFloat := val.(float64); isFloat {
				castParams[param.Name] = int(f)
			} else {
				castParams[param.Name] = val
			}

		case SimulationParameterTypeFloat:
			f, ok := val.(float64)
			if !ok {
				return nil, fmt.Errorf("parameter %s is not a float", param.Name)
			}
			castParams[param.Name] = f

		case SimulationParameterTypeVector:
			vec, ok := val.([]interface{})
			if !ok {
				return nil, fmt.Errorf("parameter %s is not a vector", param.Name)
			}
			castVec := make([]float64, len(vec))
			for i, v := range vec {
				f, ok := v.(float64)
				if !ok {
					return nil, fmt.Errorf("parameter %s has non-float element", param.Name)
				}
				castVec[i] = f
			}
			castParams[param.Name] = castVec
		default:
			castParams[param.Name] = val
		}
	}
	return castParams, nil
}

// TempFileEnv represents the temporary file environment for a simulation run.
// It contains paths to the root directory, output directory, and config file.
type TempFileEnv struct {
	RootDir    string // Root temporary directory
	OutputDir  string // Directory where simulation outputs are written
	ConfigPath string // Path to the TOML configuration file
	BinaryPath string // Path to the simulation binary
}

// Cleanup removes the temporary directory and all its contents.
// Logs an error if the removal fails.
func (t TempFileEnv) Cleanup() {
	if err := os.RemoveAll(t.RootDir); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"root_dir": t.RootDir,
		}).Error("failed to remove temp dir")
	}
}

// PrepareTempDir creates a temporary directory structure for running a simulation.
// It creates the root temp directory, output subdirectory, writes the simulation
// binary, generates the TOML configuration file, and returns a TempFileEnv.
func PrepareTempDir(event SimulationRunEvent, meta SimulationMeta, binary []byte) (TempFileEnv, error) {
	tmpDir, err := os.MkdirTemp("", "simulation-*")
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"client_id": event.ClientId,
		}).Error("failed to create temp dir")
		return TempFileEnv{}, err
	}

	outputDir := filepath.Join(tmpDir, "output")
	// create dir if not exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"client_id": event.ClientId,
		}).Error("failed to create output dir")
		return TempFileEnv{}, err
	}

	// write binary to file
	if err := os.WriteFile(filepath.Join(tmpDir, "simulation"), binary, 0755); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"client_id": event.ClientId,
		}).Error("failed to write simulation binary")
		return TempFileEnv{}, err
	}

	params, err := CastParameters(meta, event.Parameters)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"client_id": event.ClientId,
		}).Error("failed to cast parameters")
		return TempFileEnv{}, err
	}

	tomlConfig := map[string]interface{}{
		"config": map[string]interface{}{
			"output_dir": outputDir,
		},
		"parameters": params,
	}

	tomlBytes, err := toml.Marshal(tomlConfig)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"client_id": event.ClientId,
		}).Error("failed to convert params to TOML")
		return TempFileEnv{}, err
	}

	configPath := filepath.Join(tmpDir, "config.toml")
	// write toml to file
	if err := os.WriteFile(configPath, tomlBytes, 0644); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"client_id": event.ClientId,
		}).Error("failed to write config.toml")
		return TempFileEnv{}, err
	}

	return TempFileEnv{
		RootDir:    tmpDir,
		OutputDir:  outputDir,
		ConfigPath: configPath,
		BinaryPath: filepath.Join(tmpDir, "simulation"),
	}, nil
}

// SimulationRunner handles the execution of simulation binaries.
// It processes simulation run events and manages the simulation lifecycle.
type SimulationRunner struct {
	db     Persistence
	config *Config
}

// NewSimulationRunner creates a new SimulationRunner with the given database.
func NewSimulationRunner(db Persistence, config *Config) *SimulationRunner {
	return &SimulationRunner{db: db, config: config}
}

// ProcessEvent executes a single simulation run. It fetches the simulation
// metadata and binary, prepares the temporary environment, runs the simulation,
// and returns the zipped output directory contents.
func (r *SimulationRunner) ProcessEvent(event SimulationRunEvent) ([]byte, error) {
	// Fetch simulation metadata to know parameter types
	meta, err := r.db.GetSimulationMeta(event.SimulationId)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"client_id": event.ClientId,
		}).Error("failed to get simulation meta")
		return nil, err
	}

	binary, err := r.db.GetSimulationBinary(event.SimulationId, r.config.CPUArch)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"client_id": event.ClientId,
		}).Error("failed to get simulation binary")
		return nil, err
	}

	// create temp dir, move simulation binary and config.toml into it
	// and create output directory
	tmpDir, err := PrepareTempDir(event, meta, binary)
	if err != nil {
		return nil, err
	}
	defer tmpDir.Cleanup()

	// create log file for simulation logs
	logFile, err := os.Create(filepath.Join(tmpDir.OutputDir, "simulation.log"))
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"client_id": event.ClientId,
		}).Error("failed to create log file")
		return nil, err
	}
	defer logFile.Close()

	// execute binary. set config path as argument. output to stderr
	cmd := exec.Command(filepath.Join(tmpDir.RootDir, "simulation"), tmpDir.ConfigPath)
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Run(); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"client_id": event.ClientId,
		}).Error("failed to run simulation")
		return nil, err
	}

	zipData, err := ZipDirectory(tmpDir.OutputDir)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"client_id": event.ClientId,
		}).Error("failed to zip output dir")
		return nil, err
	}
	return zipData, nil
}

// Start begins listening for simulation run events on the given channel.
// It processes each event, updates the run status in the database, and
// logs the duration of each simulation run.
func (r *SimulationRunner) Start(events chan SimulationRunEvent) error {
	for event := range events {
		log.WithFields(log.Fields{
			"client_id": event.ClientId,
		}).Info("starting simulation")

		start := time.Now()

		var status string

		output, err := r.ProcessEvent(event)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"client_id": event.ClientId,
			}).Error("failed to process simulation event")

			status = "failed"
		} else {
			status = "completed"
		}

		if err := r.db.UpdateSimulationRun(event.ClientId, status, output); err != nil {
			log.WithError(err).WithFields(log.Fields{
				"client_id": event.ClientId,
			}).Error("failed to update simulation run")
		}

		duration := time.Since(start)

		log.WithFields(log.Fields{
			"client_id":        event.ClientId,
			"duration_seconds": duration.Seconds(),
		}).Info("finished simulation")
	}
	return nil
}
