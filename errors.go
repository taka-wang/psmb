package psmb

import "errors"

// Binary

var (
	// ErrBitStringToUInt8s is the error of BitStringToUInt8s converter
	ErrBitStringToUInt8s = errors.New("Converts bit string to uint8 array failed")

	// ErrDecimalStringToRegisters is the error of DecimalStringToRegisters converter
	ErrDecimalStringToRegisters = errors.New("Converts decimal string to uint16/words array in big endian order failed")

	// ErrHexStringToRegisters is the error of HexStringToRegisters converter
	ErrHexStringToRegisters = errors.New("Converts hexadecimal string to uint16/words array in big endian order failed")

	// ErrRegistersToBytes is the error of RegistersToBytes converter
	ErrRegistersToBytes = errors.New("Converts registers/uint16 array to byte array in big endian order failed")

	// ErrBytesToFloat32s is the error of BytesToFloat32s converter
	ErrBytesToFloat32s = errors.New("Converts byte array to float32 array in four endian orders failed")

	// ErrBytesToInt32s is the error of BytesToInt32s converter
	ErrBytesToInt32s = errors.New("Converts byte array to UInt32 array in four endian orders failed")

	// ErrBytesToUInt32s is the error of BytesToUInt32s converter
	ErrBytesToUInt32s = errors.New("Converts byte array to UInt32 array in four endian orders failed")

	// ErrBytesToInt16s is the error of BytesToInt16s converter
	ErrBytesToInt16s = errors.New("Converts byte array to Int16 array in two endian orders failed")

	// ErrBytesToUInt16s is the error of BytesToUInt16s converter
	ErrBytesToUInt16s = errors.New("Converts byte array to UInt16 array in two endian orders failed")
)

// Service

var (
	// ErrTodo is the error of todo
	ErrTodo = errors.New("TODO")

	// ErrRequestNotFound is the error when the request not in simple task map
	ErrRequestNotFound = errors.New("Request not found!")

	// ErrRequestNotSupport is  the error when the request is not support
	ErrRequestNotSupport = errors.New("Request not support")

	// ErrResponseNotSupport is  the error when the response is not support
	ErrResponseNotSupport = errors.New("Response not support")

	// ErrInvalidMessageLength is the error when the length of message is invalid
	ErrInvalidMessageLength = errors.New("Invalid message length")

	// ErrMarshal is the error when marshal to json string failed
	ErrMarshal = errors.New("Marshal failed")

	// ErrUnmarshal is the error when unmarshal json string to structure failed
	ErrUnmarshal = errors.New("Unmarshal failed")
)
