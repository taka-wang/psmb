package filter

import "errors"

var (
	// ErrInvalidFilterName is the error when the name is invalid
	ErrInvalidFilterName = errors.New("Invalid Filter name")
	// ErrNoData is the error when the return is empty
	ErrNoData = errors.New("Data does not exist.")
	// ErrOutOfCapacity is the error when the store capacity is full
	ErrOutOfCapacity = errors.New("Filter data store run out of capacity!")
)
