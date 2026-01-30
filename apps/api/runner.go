package main

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type SimulationRunRequest struct {
	ClientId     string                 `json:"client_id"`
	SimulationId string                 `json:"simulation_id"`
	Parameters   map[string]interface{} `json:"parameters"`
}

type SimulationRunner struct {
	db Persistence
}

func NewSimulationRunner(db Persistence) *SimulationRunner {
	return &SimulationRunner{db: db}
}

func (r *SimulationRunner) ProcessEvent(event SimulationRunRequest) ([]byte, error) {
	results, err := os.ReadFile("tests/archives/simulation_trajectories.zip")
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"client_id": event.ClientId,
		}).Error("failed to read simulation results")
		return nil, err
	}
	return results, nil
}

func (r *SimulationRunner) Start(events chan SimulationRunRequest) error {
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
