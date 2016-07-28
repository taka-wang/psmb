package psmb

import "errors"

// Binary

var (
	// ErrBitStringToUInt8s is the error of BitStringToUInt8s conversion.
	ErrBitStringToUInt8s = errors.New("Fail to convert bit string to uint8 array")

	// ErrDecimalStringToRegisters is the error of DecimalStringToRegisters conversion.
	ErrDecimalStringToRegisters = errors.New("Fail to convert decimal string to uint16/words array in big endian order")

	// ErrHexStringToRegisters is the error of HexStringToRegisters conversion.
	ErrHexStringToRegisters = errors.New("Fail to convert hexadecimal string to uint16/words array in big endian order")

	// ErrRegistersToBytes is the error of RegistersToBytes conversion.
	ErrRegistersToBytes = errors.New("Fail to convert registers/uint16 array to byte array in big endian order")

	// ErrBytesToFloat32s is the error of BytesToFloat32s conversion.
	ErrBytesToFloat32s = errors.New("Fail to convert byte array to float32 array in four endian orders")

	// ErrNotANumber is the error of LinearScalingRegisters conversion.
	ErrNotANumber = errors.New("Fail to scale the registers linearly")

	// ErrBytesToInt32s is the error of BytesToInt32s conversion.
	ErrBytesToInt32s = errors.New("Fail to convert byte array to UInt32 array in four endian orders")

	// ErrBytesToUInt32s is the error of BytesToUInt32s conversion.
	ErrBytesToUInt32s = errors.New("Fail to convert byte array to UInt32 array in four endian orders")

	// ErrBytesToInt16s is the error of BytesToInt16s conversion.
	ErrBytesToInt16s = errors.New("Fail to convert byte array to Int16 array in two endian orders")

	// ErrBytesToUInt16s is the error of BytesToUInt16s conversion.
	ErrBytesToUInt16s = errors.New("Fail to convert byte array to UInt16 array in two endian orders")
)

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
