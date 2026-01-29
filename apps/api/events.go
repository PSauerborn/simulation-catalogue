package main

import "encoding/json"

// EventBroker defines the interface for publishing simulation events.
// Implementations are responsible for delivering events to consumers
// that process simulation runs asynchronously.
type EventBroker interface {
	// HealthCheck verifies the event broker connection is alive.
	HealthCheck() error
	// Publish sends an event payload to the message queue for processing.
	Publish(event []byte) error
}

type LocalEventBroker struct {
	events chan SimulationRunRequest
}

func NewLocalEventBroker(events chan SimulationRunRequest) *LocalEventBroker {
	return &LocalEventBroker{
		events: events,
	}
}

func (b *LocalEventBroker) HealthCheck() error {
	return nil
}

func (b *LocalEventBroker) Publish(event []byte) error {
	var body SimulationRunRequest
	if err := json.Unmarshal(event, &body); err != nil {
		return err
	}
	b.events <- body
	return nil
}
