package mbtcp

import "errors"

// Service

var (
	// ErrTodo is the error of todo
	ErrTodo = errors.New("TODO")

	// ErrRequestNotFound is the error when the request is not in the write task map.
	ErrRequestNotFound = errors.New("Request not found!")

	// ErrRequestNotSupport is the error when the request is not supported.
	ErrRequestNotSupport = errors.New("Request not support!")

	// ErrResponseNotSupport is the error when the response is not supported.
	ErrResponseNotSupport = errors.New("Response not support!")

	// ErrInvalidMessageLength is the error when the length of message is invalid.
	ErrInvalidMessageLength = errors.New("Invalid message length!")

	// ErrMarshal is the error when marshalling to JSON string failed.
	ErrMarshal = errors.New("Fail to marshal!")

	// ErrUnmarshal is the error when unmarshalling JSON string to structure failed.
	ErrUnmarshal = errors.New("Fail to unmarshal!")

	// ErrInvalidFunctionCode is the error when the function code is not allowed.
	ErrInvalidFunctionCode = errors.New("Invalid function code!")

	// ErrInvalidPollName is the error when the poll name is empty.
	ErrInvalidPollName = errors.New("Invalid poll name!")
)
