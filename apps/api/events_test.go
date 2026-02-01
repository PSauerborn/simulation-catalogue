package main

type MockEventBroker struct {
	Events [][]byte
}

func (broker *MockEventBroker) HealthCheck() error {
	return nil
}

func (broker *MockEventBroker) Publish(event []byte) error {
	broker.Events = append(broker.Events, event)
	return nil
}
