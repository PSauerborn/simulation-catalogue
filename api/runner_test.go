package main

import (
	"os"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
)

func TestValidateParameter(t *testing.T) {
	t.Run("int valid int", func(t *testing.T) {
		err := ValidateParameter(SimulationParameterTypeInt, 1)
		assert.NoError(t, err)
	})

	t.Run("int valid float", func(t *testing.T) {
		err := ValidateParameter(SimulationParameterTypeInt, 1.0)
		assert.NoError(t, err)
	})

	t.Run("int invalid", func(t *testing.T) {
		err := ValidateParameter(SimulationParameterTypeInt, 1.1)
		assert.Error(t, err)
	})

	t.Run("float valid", func(t *testing.T) {
		err := ValidateParameter(SimulationParameterTypeFloat, 1.1)
		assert.NoError(t, err)
	})

	t.Run("float invalid", func(t *testing.T) {
		err := ValidateParameter(SimulationParameterTypeFloat, 1)
		assert.Error(t, err)
	})

	t.Run("string valid", func(t *testing.T) {
		err := ValidateParameter(SimulationParameterTypeString, "hello")
		assert.NoError(t, err)
	})

	t.Run("string invalid", func(t *testing.T) {
		err := ValidateParameter(SimulationParameterTypeString, 1)
		assert.Error(t, err)
	})

	t.Run("vector valid float64", func(t *testing.T) {
		err := ValidateParameter(SimulationParameterTypeVector, []float64{1.0, 2.0, 3.0})
		assert.NoError(t, err)
	})

	t.Run("vector valid interface", func(t *testing.T) {
		err := ValidateParameter(SimulationParameterTypeVector, []interface{}{1.0, 2.0, 3.0})
		assert.NoError(t, err)
	})

	t.Run("vector invalid", func(t *testing.T) {
		err := ValidateParameter(SimulationParameterTypeVector, []string{"not", "a", "vector"})
		assert.Error(t, err)
	})

	t.Run("vector invalid interface", func(t *testing.T) {
		err := ValidateParameter(SimulationParameterTypeVector, []interface{}{"not", "a", "vector"})
		assert.Error(t, err)
	})
}

func TestValidateParameters(t *testing.T) {
	meta := SimulationMeta{
		Parameters: []SimulationParameter{
			{
				Name: "int",
				Type: SimulationParameterTypeInt,
			},
			{
				Name: "float",
				Type: SimulationParameterTypeFloat,
			},
			{
				Name: "string",
				Type: SimulationParameterTypeString,
			},
			{
				Name: "vector",
				Type: SimulationParameterTypeVector,
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		parameters := map[string]interface{}{
			"int":    1,
			"float":  1.1,
			"string": "hello",
			"vector": []float64{1.0, 2.0, 3.0},
		}
		errors := ValidateParameters(meta, parameters)
		assert.Empty(t, errors)
	})

	t.Run("invalid", func(t *testing.T) {
		parameters := map[string]interface{}{
			"int":    1.21,
			"float":  1.1,
			"string": 123,
			"vector": []float64{1.0, 2.0, 3.0},
		}
		err := ValidateParameters(meta, parameters)
		assert.Equal(t, 2, len(err))
	})
}

func TestCastParameters(t *testing.T) {
	meta := SimulationMeta{
		Parameters: []SimulationParameter{
			{
				Name: "int_param",
				Type: SimulationParameterTypeInt,
			},
			{
				Name: "float_param",
				Type: SimulationParameterTypeFloat,
			},
			{
				Name: "string_param",
				Type: SimulationParameterTypeString,
			},
			{
				Name: "vector_param",
				Type: SimulationParameterTypeVector,
			},
		},
	}

	t.Run("casts int from float64", func(t *testing.T) {
		params := map[string]interface{}{
			"int_param": float64(42),
		}
		result, err := CastParameters(meta, params)
		assert.NoError(t, err)
		assert.Equal(t, 42, result["int_param"])
	})

	t.Run("preserves int if already int", func(t *testing.T) {
		params := map[string]interface{}{
			"int_param": 42,
		}
		result, err := CastParameters(meta, params)
		assert.NoError(t, err)
		assert.Equal(t, 42, result["int_param"])
	})

	t.Run("preserves float", func(t *testing.T) {
		params := map[string]interface{}{
			"float_param": 3.14,
		}
		result, err := CastParameters(meta, params)
		assert.NoError(t, err)
		assert.Equal(t, 3.14, result["float_param"])
	})

	t.Run("preserves string", func(t *testing.T) {
		params := map[string]interface{}{
			"string_param": "hello",
		}
		result, err := CastParameters(meta, params)
		assert.NoError(t, err)
		assert.Equal(t, "hello", result["string_param"])
	})

	t.Run("converts vector from interface slice", func(t *testing.T) {
		params := map[string]interface{}{
			"vector_param": []interface{}{1.0, 2.0, 3.0},
		}
		result, err := CastParameters(meta, params)
		assert.NoError(t, err)
		assert.Equal(t, []float64{1.0, 2.0, 3.0}, result["vector_param"])
	})

	t.Run("skips missing parameters", func(t *testing.T) {
		params := map[string]interface{}{
			"int_param": float64(1),
		}
		result, err := CastParameters(meta, params)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Contains(t, result, "int_param")
	})

	t.Run("error on invalid float", func(t *testing.T) {
		params := map[string]interface{}{
			"float_param": "not a float",
		}
		_, err := CastParameters(meta, params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "float_param")
	})

	t.Run("error on invalid vector type", func(t *testing.T) {
		params := map[string]interface{}{
			"vector_param": "not a vector",
		}
		_, err := CastParameters(meta, params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "vector_param")
	})

	t.Run("error on vector with non-float elements", func(t *testing.T) {
		params := map[string]interface{}{
			"vector_param": []interface{}{"a", "b", "c"},
		}
		_, err := CastParameters(meta, params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "non-float element")
	})

	t.Run("all types combined", func(t *testing.T) {
		params := map[string]interface{}{
			"int_param":    float64(10),
			"float_param":  2.71828,
			"string_param": "test",
			"vector_param": []interface{}{1.0, 2.0, 3.0},
		}
		result, err := CastParameters(meta, params)
		assert.NoError(t, err)
		assert.Equal(t, 10, result["int_param"])
		assert.Equal(t, 2.71828, result["float_param"])
		assert.Equal(t, "test", result["string_param"])
		assert.Equal(t, []float64{1.0, 2.0, 3.0}, result["vector_param"])
	})
}

func TestPrepareTempDir(t *testing.T) {
	meta := SimulationMeta{
		Name: "test",
		Parameters: []SimulationParameter{
			{
				Name: "int_param",
				Type: SimulationParameterTypeInt,
			},
			{
				Name: "float_param",
				Type: SimulationParameterTypeFloat,
			},
			{
				Name: "string_param",
				Type: SimulationParameterTypeString,
			},
			{
				Name: "vector_param",
				Type: SimulationParameterTypeVector,
			},
		},
	}

	binary := []byte("go version")

	event := SimulationRunEvent{
		ClientId:     "a1b2c3d4e5f6789012345678abcdef01",
		SimulationId: "sim009",
		Parameters: map[string]interface{}{
			"int_param":    float64(10),
			"float_param":  2.71828,
			"string_param": "test",
			"vector_param": []interface{}{1.0, 2.0, 3.0},
		},
	}

	env, err := PrepareTempDir(event, meta, binary)
	assert.NoError(t, err)

	assert.DirExists(t, env.RootDir)
	assert.FileExists(t, env.BinaryPath)
	assert.FileExists(t, env.ConfigPath)

	// load TOML config
	configContents, err := os.ReadFile(env.ConfigPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, configContents)

	var config map[string]map[string]interface{}
	err = toml.Unmarshal(configContents, &config)
	assert.NoError(t, err)

	assert.Equal(t, int64(10), config["parameters"]["int_param"])
	assert.Equal(t, 2.71828, config["parameters"]["float_param"])
	assert.Equal(t, "test", config["parameters"]["string_param"])
	assert.Equal(t, []interface{}{1.0, 2.0, 3.0}, config["parameters"]["vector_param"])
}

func TestProcessEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		config := &Config{
			CPUArch: "arm64",
		}

		db := NewMockDB(true)
		db.binaries["sim001_arm64"] = []byte(`#!/bin/bash

		arg=$1
		echo $arg`)

		runner := NewSimulationRunner(db, config)

		event := SimulationRunEvent{
			ClientId:     "a1b2c3d4e5f6789012345678abcdef01",
			SimulationId: "sim001",
			Parameters: map[string]interface{}{
				"magnetic_field":  []interface{}{0.0, 0.0, 1.0},
				"electric_field":  []interface{}{0.0, 0.0, 0.0},
				"particle_charge": 1.602e-19,
				"particle_mass":   9.109e-31,
				"timesteps":       float64(1000),
			},
		}

		zipArchive, err := runner.ProcessEvent(event)
		assert.NoError(t, err)
		assert.NotNil(t, zipArchive)
	})
}
