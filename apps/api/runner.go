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

func (r *SimulationRunner) Start(events chan SimulationRunRequest) error {
	for event := range events {
		log.WithFields(log.Fields{
			"client_id": event.ClientId,
		}).Info("starting simulation")

		start := time.Now()

		time.Sleep(10 * time.Second)

		duration := time.Since(start)

		results, err := os.ReadFile("tests/archives/simulation_trajectories.zip")
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"client_id": event.ClientId,
			}).Error("failed to read simulation results")
			continue
		}

		if err := r.db.UpdateSimulationRun(event.ClientId, "completed", results); err != nil {
			log.WithError(err).WithFields(log.Fields{
				"client_id": event.ClientId,
			}).Error("failed to update simulation run")
		}

		log.WithFields(log.Fields{
			"client_id":        event.ClientId,
			"duration_seconds": duration.Seconds(),
		}).Info("finished simulation")
	}
	return nil
}
