package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

// Persistence defines the interface for database operations.
// Implementations must provide methods for managing clients, simulations,
// simulation runs, and API keys.
type Persistence interface {
	// HealthCheck verifies the database connection is alive.
	HealthCheck() error

	// GetClient retrieves a client by their unique identifier.
	GetClient(id string) (Client, error)
	// InitClient creates a new client and returns their generated ID.
	InitClient() (string, error)
	// UpdateClientLastActive updates the last active timestamp of a client.
	UpdateClientLastActive(id string) error

	// ListSimulations returns all available simulations.
	ListSimulations() ([]SimulationMeta, error)
	// GetSimulationMeta retrieves metadata for a specific simulation.
	GetSimulationMeta(id string) (SimulationMeta, error)
	// GetSimulationBinary retrieves the compiled binary for a simulation and CPU architecture.
	GetSimulationBinary(id, cpuArch string) ([]byte, error)

	// GetSimulationRun retrieves the current simulation run for a client.
	GetSimulationRun(clientId string) (SimulationRun, error)
	// CreateSimulationRun creates a new simulation run with the specified parameters.
	CreateSimulationRun(clientId string, simulationId string, parameters map[string]interface{}) (string, error)
	// UpdateSimulationRun updates the status and output of a simulation run.
	UpdateSimulationRun(clientId string, status string, output []byte) error

	// CreateSimulation creates a new simulation with the provided metadata.
	CreateSimulation(name, description, model string, parameters []SimulationParameter, outputs []SimulationOutput) (string, error)
	// UpdateSimulationMeta updates the name and description of a simulation.
	UpdateSimulationMeta(id string, name, description string) error
	// UpdateSimulationBinary uploads or updates a simulation binary for a CPU architecture.
	UpdateSimulationBinary(simId, cpuArch string, binary []byte) error
	// DeleteSimulation removes a simulation and all associated data.
	DeleteSimulation(id string) error

	// GetAPIKey retrieves an API key by its value.
	GetAPIKey(key string) (APIKey, error)
}

// PostgresDB implements the Persistence interface using PostgreSQL as the backing store.
type PostgresDB struct {
	pool *pgxpool.Pool
}

// NewPostgresDB creates a new PostgresDB instance with a connection pool
// configured using the provided Config settings.
func NewPostgresDB(config *Config) (*PostgresDB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.PostgresHost,
		config.PostgresPort,
		config.PostgresUser,
		config.PostgresPassword,
		config.PostgresDatabase)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{pool: pool}, nil
}

// HealthCheck pings the database to verify the connection is alive.
func (p *PostgresDB) HealthCheck() error {
	return p.pool.Ping(context.Background())
}

// GetClient retrieves a client by their unique identifier.
func (p *PostgresDB) GetClient(id string) (Client, error) {
	query := "SELECT id, created_at, last_active_at FROM base.client WHERE id = $1"
	row := p.pool.QueryRow(context.Background(), query, id)

	var client Client
	if err := row.Scan(&client.Id, &client.CreatedAt, &client.LastActiveAt); err != nil {
		if err == pgx.ErrNoRows {
			return Client{}, ErrClientNotInitialized{ID: id}
		}
		return Client{}, err
	}

	return client, nil
}

// InitClient creates a new client record and returns the generated ID.
func (p *PostgresDB) InitClient() (string, error) {
	ts := time.Now().UTC()

	tx, err := p.pool.Begin(context.Background())
	if err != nil {
		return "", err
	}
	defer func() {
		err := tx.Rollback(context.Background())
		if err != nil && err != pgx.ErrTxClosed {
			log.WithError(err).Warn("failed to rollback transaction")
		}
	}()

	clientId := GenerateId()
	query := "INSERT INTO base.client (id, created_at, last_active_at) VALUES ($1, $2, $2)"
	if _, err := tx.Exec(context.Background(), query, clientId, ts); err != nil {
		return "", err
	}

	query = "INSERT INTO base.simulation_run (client_id, simulation_id, status, parameters, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $5)"
	if _, err := tx.Exec(context.Background(), query, clientId, nil, "pending", nil, ts); err != nil {
		return "", err
	}

	query = "INSERT INTO base.simulation_output (client_id, blob, created_at, updated_at) VALUES ($1, $2, $3, $3)"
	if _, err := tx.Exec(context.Background(), query, clientId, nil, ts); err != nil {
		return "", err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return "", err
	}
	return clientId, nil
}

// UpdateClientLastActive updates the last active timestamp of a client.
func (p *PostgresDB) UpdateClientLastActive(id string) error {
	ts := time.Now().UTC()
	query := "UPDATE base.client SET last_active_at = $1 WHERE id = $2"
	_, err := p.pool.Exec(context.Background(), query, ts, id)
	return err
}

// ListSimulations returns all available simulations with their metadata.
func (p *PostgresDB) ListSimulations() ([]SimulationMeta, error) {
	simulations := make([]SimulationMeta, 0)

	query := "SELECT id, name, description, model, parameters, outputs, created_at, updated_at FROM base.simulation_meta"
	rows, err := p.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var simulation SimulationMeta

		if err := rows.Scan(
			&simulation.Id,
			&simulation.Name,
			&simulation.Description,
			&simulation.Model,
			&simulation.Parameters,
			&simulation.Outputs,
			&simulation.CreatedAt,
			&simulation.UpdatedAt,
		); err != nil {
			return nil, err
		}

		simulations = append(simulations, simulation)
	}

	return simulations, nil
}

// GetSimulationMeta retrieves the metadata for a specific simulation by ID.
func (p *PostgresDB) GetSimulationMeta(id string) (SimulationMeta, error) {
	query := "SELECT id, name, description, model, parameters, outputs, created_at, updated_at FROM base.simulation_meta WHERE id = $1"
	row := p.pool.QueryRow(context.Background(), query, id)

	var simulation SimulationMeta
	if err := row.Scan(
		&simulation.Id,
		&simulation.Name,
		&simulation.Description,
		&simulation.Model,
		&simulation.Parameters,
		&simulation.Outputs,
		&simulation.CreatedAt,
		&simulation.UpdatedAt,
	); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return SimulationMeta{}, ErrSimulationNotFound{
				ID: id,
			}
		default:
			return SimulationMeta{}, err
		}
	}
	return simulation, nil
}

// GetSimulationBinary retrieves the compiled binary for a simulation and CPU architecture.
func (p *PostgresDB) GetSimulationBinary(id, cpuArch string) ([]byte, error) {
	query := "SELECT blob FROM base.simulation_binary WHERE simulation_id = $1 AND cpu_arch = $2"
	row := p.pool.QueryRow(context.Background(), query, id, cpuArch)

	var binary []byte
	if err := row.Scan(&binary); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrBinaryNotFound{
				ID:      id,
				CPUArch: cpuArch,
			}
		default:
			return nil, err
		}
	}

	return binary, nil
}

// GetSimulationRun retrieves the current simulation run for a client.
func (p *PostgresDB) GetSimulationRun(clientId string) (SimulationRun, error) {
	query := `SELECT
		run.client_id,
		run.simulation_id,
		run.status,
		run.parameters,
		output.blob,
		run.created_at,
		run.completed_at
	FROM
		base.simulation_run AS run
	LEFT JOIN
		base.simulation_output AS output
	ON
		output.client_id = run.client_id
	WHERE
		run.client_id = $1`

	row := p.pool.QueryRow(context.Background(), query, clientId)

	var run SimulationRun
	if err := row.Scan(
		&run.ClientId,
		&run.SimulationId,
		&run.Status,
		&run.Parameters,
		&run.Output,
		&run.CreatedAt,
		&run.CompletedAt,
	); err != nil {
		return run, err
	}

	return run, nil
}

// CreateSimulationRun creates a new simulation run record with the specified parameters.
func (p *PostgresDB) CreateSimulationRun(clientId string, simulationId string, parameters map[string]interface{}) (string, error) {
	query := "UPDATE base.simulation_run SET simulation_id = $1, status = $2, parameters = $3, updated_at = $4 WHERE client_id = $5"

	_, err := p.pool.Exec(
		context.Background(),
		query,
		simulationId,
		"queued",
		parameters,
		time.Now().UTC(),
		clientId,
	)
	if err != nil {
		return "", err
	}
	return "", nil
}

// UpdateSimulationRun updates the status and output of a simulation run.
// If output is non-nil, it also updates the simulation output blob.
func (p *PostgresDB) UpdateSimulationRun(clientId string, status string, output []byte) error {
	tx, err := p.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		err := tx.Rollback(context.Background())
		if err != nil && err != pgx.ErrTxClosed {
			log.WithError(err).Error("failed to rollback transaction")
		}
	}()

	query := "UPDATE base.simulation_run SET status = $1, updated_at = $2 WHERE client_id = $3"
	_, err = tx.Exec(context.Background(), query, status, time.Now().UTC(), clientId)
	if err != nil {
		return err
	}

	if output != nil {
		query := "UPDATE base.simulation_output SET blob = $1, updated_at = $2 WHERE client_id = $3"
		_, err = tx.Exec(context.Background(), query, output, time.Now().UTC(), clientId)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}

// CreateSimulation creates a new simulation with the provided metadata, parameters, and outputs.
func (p *PostgresDB) CreateSimulation(name, description, model string, parameters []SimulationParameter, outputs []SimulationOutput) (string, error) {
	id := GenerateId()
	ts := time.Now().UTC()

	tx, err := p.pool.Begin(context.Background())
	if err != nil {
		return "", err
	}
	defer func() {
		err := tx.Rollback(context.Background())
		if err != nil && err != pgx.ErrTxClosed {
			log.WithError(err).Error("failed to rollback transaction")
		}
	}()

	query := "INSERT INTO base.simulation_meta (id, name, description, model, parameters, outputs, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $7)"

	_, err = tx.Exec(
		context.Background(),
		query,
		id,
		name,
		description,
		model,
		parameters,
		outputs,
		ts,
	)
	if err != nil {
		return "", err
	}

	query = "INSERT INTO base.simulation_binary (simulation_id, cpu_arch, created_at, updated_at) VALUES ($1, $2, $3, $3)"

	archs := []string{"amd64", "arm64"}
	for _, arch := range archs {
		_, err = tx.Exec(
			context.Background(),
			query,
			id,
			arch,
			ts,
		)
		if err != nil {
			return "", err
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		return "", err
	}
	return id, nil
}

// UpdateSimulationMeta updates the name and description of an existing simulation.
func (p *PostgresDB) UpdateSimulationMeta(id string, name, description string) error {
	query := "UPDATE base.simulation_meta SET name = $1, description = $2, updated_at = $3 WHERE id = $4"

	_, err := p.pool.Exec(
		context.Background(),
		query,
		name,
		description,
		time.Now().UTC(),
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

// UpdateSimulationBinary uploads or updates the binary for a simulation and CPU architecture.
func (p *PostgresDB) UpdateSimulationBinary(simId, cpuArch string, binary []byte) error {
	query := "UPDATE base.simulation_binary SET blob = $1, revision = revision + 1, updated_at = $2 WHERE simulation_id = $3 AND cpu_arch = $4"

	_, err := p.pool.Exec(
		context.Background(),
		query,
		binary,
		time.Now().UTC(),
		simId,
		cpuArch,
	)
	if err != nil {
		return err
	}
	return nil
}

// DeleteSimulation permanently removes a simulation and all associated data.
func (p *PostgresDB) DeleteSimulation(id string) error {
	query := "DELETE FROM base.simulation_meta WHERE id = $1 CASCADE"

	_, err := p.pool.Exec(
		context.Background(),
		query,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

// GetAPIKey retrieves an API key record by its value.
func (p *PostgresDB) GetAPIKey(key string) (APIKey, error) {
	query := "SELECT id, owner, key, revoked, expires_at, created_at FROM base.api_key WHERE key = $1"
	row := p.pool.QueryRow(context.Background(), query, key)

	var apiKey APIKey
	if err := row.Scan(
		&apiKey.Id,
		&apiKey.Owner,
		&apiKey.Key,
		&apiKey.Revoked,
		&apiKey.ExpiresAt,
		&apiKey.CreatedAt,
	); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return APIKey{}, ErrAPIKeyNotFound{
				Key: key,
			}
		default:
			return APIKey{}, err
		}
	}

	return apiKey, nil
}
