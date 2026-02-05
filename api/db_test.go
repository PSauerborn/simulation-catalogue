package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type MockDB struct {
	healthy        bool
	clients        []Client
	simulations    []SimulationMeta
	simulationRuns []SimulationRun
	binaries       map[string][]byte
	apiKeys        []APIKey
}

func NewMockDB(healthy bool) *MockDB {
	clientsFixtures, err := os.ReadFile("tests/fixtures/clients.json")
	if err != nil {
		panic(err)
	}

	var clients []Client
	if err := json.Unmarshal(clientsFixtures, &clients); err != nil {
		panic(err)
	}

	simulationFixtures, err := os.ReadFile("tests/fixtures/simulations.json")
	if err != nil {
		panic(err)
	}

	var simulations []SimulationMeta
	if err := json.Unmarshal(simulationFixtures, &simulations); err != nil {
		panic(err)
	}

	binaries, err := os.ReadDir("tests/binaries")
	if err != nil {
		panic(err)
	}

	binaryMap := make(map[string][]byte)
	for _, binary := range binaries {
		binaryData, err := os.ReadFile("tests/binaries/" + binary.Name())
		if err != nil {
			panic(err)
		}
		binaryMap[fmt.Sprintf("%s_%s", binary.Name(), "amd64")] = binaryData
	}

	simulationRunFixtures, err := os.ReadFile("tests/fixtures/simulation_runs.json")
	if err != nil {
		panic(err)
	}

	var simulationRuns []SimulationRun
	if err := json.Unmarshal(simulationRunFixtures, &simulationRuns); err != nil {
		panic(err)
	}

	archive, err := os.ReadFile("tests/archives/simulation_trajectories.zip")
	if err != nil {
		panic(err)
	}

	// update output for all completed runs
	for i, run := range simulationRuns {
		if run.Status == "completed" {
			simulationRuns[i].Output = archive
		}
	}

	apiKeysFixtures, err := os.ReadFile("tests/fixtures/api_keys.json")
	if err != nil {
		panic(err)
	}

	var apiKeys []APIKey
	if err := json.Unmarshal(apiKeysFixtures, &apiKeys); err != nil {
		panic(err)
	}

	return &MockDB{
		healthy:        healthy,
		clients:        clients,
		simulations:    simulations,
		simulationRuns: simulationRuns,
		binaries:       binaryMap,
		apiKeys:        apiKeys,
	}
}

func (db *MockDB) SimulationsById() map[string]SimulationMeta {
	simulations := make(map[string]SimulationMeta)
	for _, simulation := range db.simulations {
		simulations[simulation.Id] = simulation
	}
	return simulations
}

func (db *MockDB) SimulationsByName() map[string]SimulationMeta {
	simulations := make(map[string]SimulationMeta)
	for _, simulation := range db.simulations {
		simulations[simulation.Name] = simulation
	}
	return simulations
}

func (db *MockDB) ClientsById() map[string]Client {
	clients := make(map[string]Client)
	for _, client := range db.clients {
		clients[client.Id] = client
	}
	return clients
}

func (db *MockDB) SimulationRunsByClientId() map[string][]SimulationRun {
	runs := make(map[string][]SimulationRun)
	for _, run := range db.simulationRuns {
		runs[run.ClientId] = append(runs[run.ClientId], run)
	}
	return runs
}

func (db *MockDB) HealthCheck() error {
	if !db.healthy {
		return fmt.Errorf("database is not healthy")
	}
	return nil
}

func (db *MockDB) GetClient(id string) (Client, error) {
	for _, client := range db.clients {
		if client.Id == id {
			println("client found")
			return client, nil
		}
	}
	return Client{}, ErrClientNotFound{ID: id}
}

func (db *MockDB) InitClient() (string, error) {
	id := uuid.New().String()
	id = strings.ReplaceAll(id, "-", "")

	client := Client{
		Id: id,
	}
	db.clients = append(db.clients, client)
	return id, nil
}

func (db *MockDB) UpdateClientLastActive(id string) error {
	for i, client := range db.clients {
		if client.Id == id {
			db.clients[i].LastActiveAt = time.Now().UTC()
			return nil
		}
	}
	return ErrClientNotFound{ID: id}
}

// implement persistence
func (db *MockDB) ListSimulations() ([]SimulationMeta, error) {
	return db.simulations, nil
}

func (db *MockDB) GetSimulationMeta(id string) (SimulationMeta, error) {
	for _, simulation := range db.simulations {
		if simulation.Id == id {
			return simulation, nil
		}
	}
	return SimulationMeta{}, ErrSimulationNotFound{ID: id}
}

func (db *MockDB) GetSimulationBinary(id, cpuArch string) ([]byte, error) {
	if _, exists := db.SimulationsById()[id]; !exists {
		return nil, ErrSimulationNotFound{ID: id}
	}

	binary, ok := db.binaries[fmt.Sprintf("%s_%s", id, cpuArch)]
	if !ok {
		return nil, ErrBinaryNotFound{ID: id}
	}
	return binary, nil
}

func (db *MockDB) GetSimulationRun(clientId string) (SimulationRun, error) {
	for _, run := range db.simulationRuns {
		if run.ClientId == clientId {
			return run, nil
		}
	}
	return SimulationRun{}, fmt.Errorf("simulation run not found")
}

func (db *MockDB) CreateSimulationRun(clientId string, simulationId string, parameters map[string]interface{}) (string, error) {
	id := uuid.New().String()
	id = strings.ReplaceAll(id, "-", "")

	// remove existing run for client
	for i, run := range db.simulationRuns {
		if run.ClientId == clientId {
			db.simulationRuns[i].SimulationId = &simulationId
			db.simulationRuns[i].Parameters = parameters
			db.simulationRuns[i].Status = "queued"
			return id, nil
		}
	}

	return id, nil
}

func (db *MockDB) UpdateSimulationRun(clientId string, status string, output []byte) error {
	for i, run := range db.simulationRuns {
		if run.ClientId == clientId {
			db.simulationRuns[i].Status = status
			db.simulationRuns[i].Output = output
			return nil
		}
	}
	return fmt.Errorf("simulation run not found")
}

func (db *MockDB) CreateSimulation(name, description, model string, parameters []SimulationParameter, outputs []SimulationOutput) (string, error) {
	id := uuid.New().String()
	id = strings.ReplaceAll(id, "-", "")

	simulation := SimulationMeta{
		Id:          id,
		Name:        name,
		Description: description,
		Model:       model,
		Parameters:  parameters,
		Outputs:     outputs,
	}
	db.simulations = append(db.simulations, simulation)
	return id, nil
}

func (db *MockDB) UpdateSimulationMeta(id string, name, description, model string, parameters []SimulationParameter, outputs []SimulationOutput) error {
	for i, simulation := range db.simulations {
		if simulation.Id == id {
			db.simulations[i] = SimulationMeta{
				Id:          id,
				Name:        name,
				Description: description,
				Model:       model,
				Parameters:  parameters,
				Outputs:     outputs,
				CreatedAt:   simulation.CreatedAt,
				UpdatedAt:   simulation.UpdatedAt,
			}
			return nil
		}
	}
	return fmt.Errorf("simulation not found")
}

func (db *MockDB) UpdateSimulationBinary(simId, cpuArch string, binary []byte) error {
	db.binaries[fmt.Sprintf("%s_%s", simId, cpuArch)] = binary
	return nil
}

func (db *MockDB) DeleteSimulation(id string) error {
	filtered := make([]SimulationMeta, 0)
	for _, simulation := range db.simulations {
		if simulation.Id != id {
			filtered = append(filtered, simulation)
		}
	}
	db.simulations = filtered
	return nil
}

func (db *MockDB) GetAPIKey(key string) (APIKey, error) {
	for _, apiKey := range db.apiKeys {
		if apiKey.Key == key {
			return apiKey, nil
		}
	}
	return APIKey{}, ErrAPIKeyNotFound{Key: key}
}
