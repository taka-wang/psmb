package history

import "errors"

// Service

var (
	// ErrConnection is the error when the connection failed
	ErrConnection = errors.New("Fail to connect to redis server")

	// ErrInvalidName is the error when the name is invalid
	ErrInvalidName = errors.New("Invalid name")

	// ErrNoData is the error when the return is empty
	ErrNoData = errors.New("No data exist.")
)
