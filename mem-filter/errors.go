package filter

import "errors"

var (
	// ErrInvalidFilterName is the error when the name is invalid
	ErrInvalidFilterName = errors.New("Invalid Filter name")
	// ErrNoData is the error when the return is empty
	ErrNoData = errors.New("Data does not exist.")
)
