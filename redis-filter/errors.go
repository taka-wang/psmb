package filter

import "errors"

var (
	// ErrInvalidFilterName is the error when the name is invalid
	ErrInvalidFilterName = errors.New("Invalid Filter name")

	// ErrConnection is the error when the connection failed
	ErrConnection = errors.New("Fail to connect to redis server")

	// ErrInvalidName is the error when the name is invalid
	ErrInvalidName = errors.New("Invalid name")

	// ErrNoData is the error when the return is empty
	ErrNoData = errors.New("Data does not exist.")

	// ErrUnmarshal is the error when unmarshalling JSON string to structure failed.
	ErrUnmarshal = errors.New("Fail to unmarshal!")

	// ErrOutOfCapacity is the error when the store capacity is full
	ErrOutOfCapacity = errors.New("Filter data store run out of capacity!")
)
