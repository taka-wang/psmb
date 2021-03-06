package history

import "errors"

// Service

var (
	// ErrConnection is the error when the connection failed
	ErrConnection = errors.New("Fail to connect to mongo server")

	// ErrInvalidName is the error when the name is invalid
	ErrInvalidName = errors.New("Invalid name")

	// ErrNoData is the error when the return is empty
	ErrNoData = errors.New("Data does not exist.")

	// ErrMarshal is the error when marshalling to JSON string failed.
	ErrMarshal = errors.New("Fail to marshal!")
)
