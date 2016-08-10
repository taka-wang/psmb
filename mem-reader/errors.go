package reader

import "errors"

var (
	// ErrInvalidPollName is the error when the poll name is empty.
	ErrInvalidPollName = errors.New("Invalid poll name!")

	// ErrNoData is the error when the return is empty
	ErrNoData = errors.New("Data does not exist.")

	// ErrOutOfCapacity is the error when the store capacity is full
	ErrOutOfCapacity = errors.New("Reader data store run out of capacity!")
)
