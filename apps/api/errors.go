package main

import "fmt"

// ErrSimulationNotFound is returned when a requested simulation does not exist.
type ErrSimulationNotFound struct {
	ID string
}

// Error returns a formatted error message including the simulation ID.
func (e ErrSimulationNotFound) Error() string {
	return fmt.Sprintf("simulation %s not found", e.ID)
}

// ErrBinaryNotFound is returned when a simulation binary is not available
// for the requested CPU architecture.
type ErrBinaryNotFound struct {
	ID      string
	CPUArch string
}

// Error returns a formatted error message including the simulation ID.
func (e ErrBinaryNotFound) Error() string {
	return fmt.Sprintf("binary for simulation %s not found", e.ID)
}

// ErrClientNotFound is returned when a requested client session does not exist.
type ErrClientNotFound struct {
	ID string
}

// Error returns a formatted error message including the client ID.
func (e ErrClientNotFound) Error() string {
	return fmt.Sprintf("client %s not found", e.ID)
}

// ErrAPIKeyNotFound is returned when an API key does not exist in the database.
type ErrAPIKeyNotFound struct {
	Key string
}

// Error returns a formatted error message including the API key.
func (e ErrAPIKeyNotFound) Error() string {
	return fmt.Sprintf("API key %s not found", e.Key)
}

type ErrClientNotInitialized struct {
	ID string
}

// Error returns a formatted error message including the client ID.
func (e ErrClientNotInitialized) Error() string {
	return fmt.Sprintf("client %s not initialized", e.ID)
}
