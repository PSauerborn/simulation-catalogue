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

// LocalEventBroker implements EventBroker using an in-memory channel.
// It is used for local development and testing where a message queue is not needed.
type LocalEventBroker struct {
	events chan SimulationRunEvent
}

// NewLocalEventBroker creates a new LocalEventBroker that sends events
// to the provided channel for processing by simulation runners.
func NewLocalEventBroker(events chan SimulationRunEvent) *LocalEventBroker {
	return &LocalEventBroker{
		events: events,
	}
}

// HealthCheck always returns nil for the local broker since there is
// no external connection to verify.
func (b *LocalEventBroker) HealthCheck() error {
	return nil
}

// Publish unmarshals the event payload and sends it to the events channel
// for asynchronous processing by simulation runners.
func (b *LocalEventBroker) Publish(event []byte) error {
	var body SimulationRunEvent
	if err := json.Unmarshal(event, &body); err != nil {
		return err
	}
	b.events <- body
	return nil
}
