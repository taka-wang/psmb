package dbhistory

import "errors"

// Service

var (
	// ErrConnection is the error when the connection failed
	ErrConnection = errors.New("Fail to connect to redis server")

	// ErrInvalidName invalid name
	ErrInvalidName = errors.New("Invalid name")
)
