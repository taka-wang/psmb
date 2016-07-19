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

	// ErrInvalidLengthToConvert is the error of invalid length to convert
	ErrInvalidLengthToConvert = errors.New("Invalid length to convert")
)
