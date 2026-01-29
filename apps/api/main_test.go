package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}

		controller := NewController(cfg, nil, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/version", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusOK)

		var body map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &body)

		assert.NoError(t, err)
		assert.Equal(t, body["version"], "v1")
	})
}

func TestHealthCheck(t *testing.T) {
	t.Run("healthy", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/health", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusOK)

		var body map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &body)

		assert.NoError(t, err)
		assert.Equal(t, body["status"], "OK")
	})

	t.Run("unhealthy", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(false)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/health", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusInternalServerError)

		var body map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &body)

		assert.NoError(t, err)
		assert.Equal(t, body["error"], "Internal Server Error")
	})
}

func TestGetClient(t *testing.T) {
	t.Run("existing client", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/client", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-Client-Id", "a1b2c3d4e5f6789012345678abcdef01")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)

		var body map[string]Client
		err := json.Unmarshal(response.Body.Bytes(), &body)

		expected := db.ClientsById()["a1b2c3d4e5f6789012345678abcdef01"]

		assert.NoError(t, err)
		assert.Equal(t, expected, body["client"])
	})

	t.Run("non existing client", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/client", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusOK)

		var body map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &body)

		assert.NoError(t, err)
		assert.Equal(t, nil, body["data"])
	})
}

func TestInitClient(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodPost, "/v1/public/client/init", nil)
		response := httptest.NewRecorder()

		assert.Equal(t, 5, len(db.clients))

		router.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusOK)

		var body map[string]string
		err := json.Unmarshal(response.Body.Bytes(), &body)

		assert.NoError(t, err)

		id := body["id"]

		client, exists := db.ClientsById()[id]

		assert.True(t, exists)
		assert.Equal(t, id, client.Id)

		assert.Equal(t, 6, len(db.clients))
	})
}

func TestListSimulations(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/simulations", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusOK)

		var body map[string][]SimulationMeta
		err := json.Unmarshal(response.Body.Bytes(), &body)

		assert.NoError(t, err)
		assert.Equal(t, body["data"], db.simulations)
	})
}

func TestGetSimulationMeta(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/simulations/sim001/meta", nil)
		response := httptest.NewRecorder()

		sim := db.SimulationsById()["sim001"]

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)

		var body SimulationMeta
		err := json.Unmarshal(response.Body.Bytes(), &body)

		assert.NoError(t, err)
		assert.Equal(t, sim, body)
	})

	t.Run("not found", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/simulations/not-found/meta", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestGetSimulationBinary(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/simulations/sim001/binary/amd64", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("simulation not found", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/simulations/not-found/binary", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("binary not found", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/simulations/sim005/binary", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestRunSimulation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)
		events := &MockEventBroker{}

		controller := NewController(cfg, db, events)
		router := NewRouter(controller)

		encoded, _ := json.Marshal(map[string]interface{}{
			"simulation_id": "sim001",
			"parameters": map[string]interface{}{
				"a": "1",
			},
		})

		request := httptest.NewRequest(http.MethodPost, "/v1/public/simulations/run", bytes.NewReader(encoded))
		response := httptest.NewRecorder()

		request.Header.Set("X-Client-Id", "a1b2c3d4e5f6789012345678abcdef01")

		initial := len(db.simulationRuns)
		assert.Equal(t, 0, len(events.Events))

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusCreated, response.Code)

		var body map[string]string
		err := json.Unmarshal(response.Body.Bytes(), &body)

		assert.NoError(t, err)
		runId := body["id"]
		assert.NotEmpty(t, runId)

		final := len(db.simulationRuns)
		assert.Equal(t, initial, final)

		assert.Equal(t, 1, len(events.Events))
	})

	t.Run("simulation not found", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		encoded, _ := json.Marshal(map[string]interface{}{
			"simulation_id": "not-found",
			"parameters": map[string]interface{}{
				"a": "1",
			},
		})

		initial := len(db.simulationRuns)

		request := httptest.NewRequest(http.MethodPost, "/v1/public/simulations/run", bytes.NewReader(encoded))
		response := httptest.NewRecorder()

		request.Header.Set("X-Client-Id", "a1b2c3d4e5f6789012345678abcdef01")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusNotFound, response.Code)

		final := len(db.simulationRuns)
		assert.Equal(t, initial, final)
	})

	t.Run("client not found", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		initial := len(db.simulationRuns)

		encoded, _ := json.Marshal(map[string]interface{}{
			"simulation_id": "sim001",
			"parameters": map[string]interface{}{
				"a": "1",
			},
		})

		request := httptest.NewRequest(http.MethodPost, "/v1/public/simulations/run", bytes.NewReader(encoded))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusBadRequest, response.Code)

		final := len(db.simulationRuns)
		assert.Equal(t, initial, final)
	})

	t.Run("active simulation", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		encoded, _ := json.Marshal(map[string]interface{}{
			"simulation_id": "sim001",
			"parameters": map[string]interface{}{
				"a": "1",
			},
		})

		initial := len(db.simulationRuns)

		request := httptest.NewRequest(http.MethodPost, "/v1/public/simulations/run", bytes.NewReader(encoded))
		response := httptest.NewRecorder()

		request.Header.Set("X-Client-Id", "b2c3d4e5f6789012345678abcdef01a2")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusConflict, response.Code)

		final := len(db.simulationRuns)
		assert.Equal(t, initial, final)
	})
}

func TestGetSimulationRun(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/simulations/run", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-Client-Id", "a1b2c3d4e5f6789012345678abcdef01")

		runs := db.SimulationRunsByClientId()["a1b2c3d4e5f6789012345678abcdef01"]
		assert.Equal(t, 1, len(runs))

		run := runs[0]

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)

		var body SimulationRun
		err := json.Unmarshal(response.Body.Bytes(), &body)

		assert.NoError(t, err)
		assert.Equal(t, run, body)

		runs = db.SimulationRunsByClientId()["a1b2c3d4e5f6789012345678abcdef01"]
		assert.Equal(t, 1, len(runs))
	})

	t.Run("not found", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/simulations/run", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-Client-Id", "e5f6789012345678abcdef01a2b3c4d5")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("no client id", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/simulations/run", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestGetSimulationOutput(t *testing.T) {
	t.Run("json format", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/simulations/output?format=json", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-Client-Id", "a1b2c3d4e5f6789012345678abcdef01")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)

		var body map[string]map[string][]map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &body)
		assert.NoError(t, err)

		output := body["output"]

		_, exists := output["trajectory_position.csv"]
		assert.True(t, exists)

		_, exists = output["trajectory_velocity.csv"]
		assert.True(t, exists)
	})

	t.Run("zip format", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/simulations/output?format=zip", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-Client-Id", "a1b2c3d4e5f6789012345678abcdef01")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("default format", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/simulations/output", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-Client-Id", "a1b2c3d4e5f6789012345678abcdef01")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)

		var body map[string]map[string][]map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &body)
		assert.NoError(t, err)

		output := body["output"]

		_, exists := output["trajectory_position.csv"]
		assert.True(t, exists)

		_, exists = output["trajectory_velocity.csv"]
		assert.True(t, exists)
	})

	t.Run("no run", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/simulations/output", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-Client-Id", "e5f6789012345678abcdef01a2b3c4d5")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)

		var body map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &body)

		assert.NoError(t, err)
		assert.Equal(t, nil, body["output"])
	})

	t.Run("incomplete run", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/simulations/output", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-Client-Id", "b2c3d4e5f6789012345678abcdef01a2")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)

		var body map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &body)

		assert.NoError(t, err)
		assert.Equal(t, nil, body["output"])
	})

	t.Run("invalid format", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodGet, "/v1/public/simulations/output?format=invalid", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-Client-Id", "a1b2c3d4e5f6789012345678abcdef01")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusBadRequest, response.Code)

		var body map[string]string
		err := json.Unmarshal(response.Body.Bytes(), &body)

		assert.NoError(t, err)
		assert.Equal(t, "Invalid format", body["error"])
	})
}

func TestCreateSimulation(t *testing.T) {
	t.Run("no auth", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodPost, "/v1/admin/simulations", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("unknown key", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodPost, "/v1/admin/simulations", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-API-Key", "sk_live_unknown")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("revoked key", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodPost, "/v1/admin/simulations", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-API-Key", "sk_live_c3d4e5f6789012345678abcdef01a2b3")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("expired key", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodPost, "/v1/admin/simulations", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-API-Key", "a1b2c3d4e5f6789012345678abcdef01")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("success", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		body, _ := json.Marshal(NewSimulationRequest{
			Name:        "Bose Hubbard Model",
			Description: "Simulate the Bose-Hubbard model",
			Parameters: []SimulationParameter{
				{
					Name:        "U",
					Description: "Interaction strength",
					Type:        SimulationParameterTypeFloat,
				},
				{
					Name:        "J",
					Description: "Hopping strength",
					Type:        SimulationParameterTypeFloat,
				},
			},
			Outputs: []SimulationOutput{
				{
					Name:        "density",
					Description: "Density profile",
					Type:        OutputTypeTrajectory,
				},
			},
		})

		_, exists := db.SimulationsByName()["Bose Hubbard Model"]
		assert.False(t, exists)

		initialSimulationCount := len(db.simulations)
		assert.Equal(t, 5, initialSimulationCount)

		request := httptest.NewRequest(http.MethodPost, "/v1/admin/simulations", bytes.NewReader(body))
		response := httptest.NewRecorder()

		request.Header.Set("X-API-Key", "sk_live_a1b2c3d4e5f6789012345678abcdef01")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusCreated, response.Code)

		_, exists = db.SimulationsByName()["Bose Hubbard Model"]
		assert.True(t, exists)

		initialSimulationCount = len(db.simulations)
		assert.Equal(t, 6, initialSimulationCount)
	})

}

func TestUpdateSimulationMeta(t *testing.T) {
	t.Run("no auth", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodPut, "/v1/admin/simulations/sim001/meta", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("unknown key", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodPut, "/v1/admin/simulations/sim001/meta", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-API-Key", "sk_live_unknown")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("revoked key", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodPut, "/v1/admin/simulations/sim001/meta", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-API-Key", "sk_live_c3d4e5f6789012345678abcdef01a2b3")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("expired key", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodPut, "/v1/admin/simulations/sim001/meta", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-API-Key", "a1b2c3d4e5f6789012345678abcdef01")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("simulation not found", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		encoded, _ := json.Marshal(PatchSimulationRequest{
			Name:        "Bose Hubbard Model",
			Description: "Simulate the Bose-Hubbard model",
		})

		request := httptest.NewRequest(http.MethodPut, "/v1/admin/simulations/sim999/meta", bytes.NewReader(encoded))
		response := httptest.NewRecorder()

		request.Header.Set("X-API-Key", "sk_live_a1b2c3d4e5f6789012345678abcdef01")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("success", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		encoded, _ := json.Marshal(PatchSimulationRequest{
			Name:        "Bose Hubbard Model",
			Description: "Simulate the Bose-Hubbard model",
		})

		request := httptest.NewRequest(http.MethodPut, "/v1/admin/simulations/sim001/meta", bytes.NewReader(encoded))
		response := httptest.NewRecorder()

		request.Header.Set("X-API-Key", "sk_live_a1b2c3d4e5f6789012345678abcdef01")

		initial := db.SimulationsById()["sim001"]
		assert.Equal(t, "Helical Motion Simulation", initial.Name)
		assert.Equal(t, "Simulates the helical motion of charged particles in a uniform magnetic field with optional electric field perturbations.", initial.Description)

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusNoContent, response.Code)

		updated := db.SimulationsById()["sim001"]
		assert.Equal(t, "Bose Hubbard Model", updated.Name)
		assert.Equal(t, "Simulate the Bose-Hubbard model", updated.Description)
	})
}

func TestUpdateSimulationBinary(t *testing.T) {
	t.Run("no auth", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodPut, "/v1/admin/simulations/sim001/binary/amd64", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("unknown key", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodPut, "/v1/admin/simulations/sim001/binary/amd64", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-API-Key", "sk_live_unknown")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("revoked key", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodPut, "/v1/admin/simulations/sim001/binary/amd64", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-API-Key", "sk_live_c3d4e5f6789012345678abcdef01a2b3")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("expired key", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodPut, "/v1/admin/simulations/sim001/binary/amd64", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-API-Key", "a1b2c3d4e5f6789012345678abcdef01")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("simulation not found", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		data, _ := os.ReadFile("tests/binaries/sim001")

		// create form data
		writer := httptest.NewRecorder()
		form := multipart.NewWriter(writer)

		part, _ := form.CreateFormFile("binary", "sim001")
		if _, err := io.Copy(part, bytes.NewReader(data)); err != nil {
			t.Fatalf("failed to copy binary: %v", err)
		}

		if err := form.Close(); err != nil {
			t.Fatalf("failed to close form: %v", err)
		}

		request := httptest.NewRequest(http.MethodPut, "/v1/admin/simulations/sim999/binary/amd64", bytes.NewReader(writer.Body.Bytes()))
		response := httptest.NewRecorder()

		request.Header.Set("X-API-Key", "sk_live_a1b2c3d4e5f6789012345678abcdef01")
		request.Header.Set("Content-Type", form.FormDataContentType())

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("success", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		contents, _ := os.ReadFile("tests/binaries/sim001")

		buffer := &bytes.Buffer{}
		multipartWriter := multipart.NewWriter(buffer)

		part, err := multipartWriter.CreateFormFile("binary", "sim001")
		if err != nil {
			t.Fatal(err)
		}

		_, err = io.Copy(part, bytes.NewReader(contents))
		if err != nil {
			t.Fatal(err)
		}
		multipartWriter.Close()

		request := httptest.NewRequest(http.MethodPut, "/v1/admin/simulations/sim001/binary/amd64", buffer)
		response := httptest.NewRecorder()

		request.Header.Set("X-API-Key", "sk_live_a1b2c3d4e5f6789012345678abcdef01")
		request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusNoContent, response.Code)
	})
}

func TestDeleteSimulation(t *testing.T) {
	t.Run("no auth", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodDelete, "/v1/admin/simulations/sim001", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("unknown key", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodDelete, "/v1/admin/simulations/sim001", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-API-Key", "sk_live_unknown")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("revoked key", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodDelete, "/v1/admin/simulations/sim001", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-API-Key", "sk_live_c3d4e5f6789012345678abcdef01a2b3")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("expired key", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodDelete, "/v1/admin/simulations/sim001", nil)
		response := httptest.NewRecorder()

		request.Header.Set("X-API-Key", "a1b2c3d4e5f6789012345678abcdef01")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("success", func(t *testing.T) {
		cfg := &Config{
			Version: "v1",
		}
		db := NewMockDB(true)

		controller := NewController(cfg, db, nil)
		router := NewRouter(controller)

		request := httptest.NewRequest(http.MethodDelete, "/v1/admin/simulations/sim001", nil)
		response := httptest.NewRecorder()

		_, exists := db.SimulationsById()["sim001"]
		assert.True(t, exists)

		request.Header.Set("X-API-Key", "sk_live_a1b2c3d4e5f6789012345678abcdef01")

		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusNoContent, response.Code)

		_, exists = db.SimulationsById()["sim001"]
		assert.False(t, exists)
	})
}
