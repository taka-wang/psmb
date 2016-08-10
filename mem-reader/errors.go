package reader

import "errors"

var (
	// ErrInvalidPollName is the error when the poll name is empty.
	ErrInvalidPollName = errors.New("Invalid poll name!")
	// ErrOutOfCapacity is the error when the store capacity is full
	ErrOutOfCapacity = errors.New("Reader data store run out of capacity!")
)
