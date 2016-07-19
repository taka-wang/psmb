package psmb

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math"
	"strconv"
	"strings"
)

// BitStringToUInt8s converts bits string to uint8 array.
//	source: function code 15
func BitStringToUInt8s(bitString string) ([]uint8, error) {
	var result = []uint8{}
	s := strings.Trim(bitString, ",") // trim left, right
	for _, v := range strings.Split(s, ",") {
		i, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			return nil, ErrBitStringToUInt8s
		}
		result = append(result, uint8(i))
	}
	return result, nil
}

// DecimalStringToRegisters converts decimal string to uint16/words array in big endian order.
// 	limitation: leading space before comma is not allowed.
// 	source: upstream.
func DecimalStringToRegisters(decString string) ([]uint16, error) {
	var result = []uint16{}
	s := strings.Trim(decString, ",") // trim left, right
	for _, v := range strings.Split(s, ",") {
		i, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			return nil, ErrDecimalStringToRegisters
		}
		result = append(result, uint16(i))
	}
	return result, nil
}

// HexStringToRegisters converts hexadecimal string to uint16/words array in big endian order.
// 	source: upstream.
func HexStringToRegisters(hexString string) ([]uint16, error) {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, ErrHexStringToRegisters
	}
	return BytesToUInt16s(bytes, 0)
}

// BytesToHexString converts bytes array to hexadecimal string.
// 	example: 112C004F12345678
func BytesToHexString(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// RegistersToBytes converts registers/uint16 array to byte array in big endian order.
// 	source: downstream.
func RegistersToBytes(data []uint16) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			return nil, ErrRegistersToBytes
		}
	}
	return buf.Bytes(), nil
}

// LinearScalingRegisters scales the registers linearly.
// Equation:
// 	Let domainLow, domainHigh, rangeLow, rangeHigh as a, b, c, d accordingly.
// 	Output = c + (d - c) * (input - a) / (b - a)
func LinearScalingRegisters(data []uint16, domainLow, domainHigh, rangeLow, rangeHigh float64) ([]float32, error) {
	l := len(data)
	low := math.Min(domainLow, domainHigh)
	high := math.Max(domainLow, domainHigh)
	result := make([]float32, l)

	var tmp float64
	for idx := 0; idx < l; idx++ {
		tmp = rangeLow + (rangeHigh-rangeLow)*(math.Min(math.Max(low, float64(data[idx])), high)-low)/(high-low)
		if math.IsNaN(tmp) {
			return nil, ErrNotANumber
		}
		result[idx] = float32(tmp)
	}
	return result, nil
}

// BytesToFloat32s converts byte array to float32 array in four endian orders. i.e.,
//	BigEndian (0),
//	LittleEndian (1)
//	MidBigEndian (2)
//	MidLittleEndian (3)
func BytesToFloat32s(buf []byte, endian Endian) ([]float32, error) {
	l := len(buf)
	if l == 0 || l%4 != 0 {
		return nil, ErrBytesToFloat32s
	}

	result := make([]float32, l/4)
	for idx := 0; idx < l/4; idx++ {
		switch endian {
		case DCBA: // little endian
			result[idx] = math.Float32frombits(uint32(buf[4*idx]) | uint32(buf[4*idx+1])<<8 | uint32(buf[4*idx+2])<<16 | uint32(buf[4*idx+3])<<24)
		case BADC: // mid-big endian
			result[idx] = math.Float32frombits(uint32(buf[4*idx+2]) | uint32(buf[4*idx+3])<<8 | uint32(buf[4*idx])<<16 | uint32(buf[4*idx+1])<<24)
		case CDAB: // mid-little endian
			result[idx] = math.Float32frombits(uint32(buf[4*idx+1]) | uint32(buf[4*idx])<<8 | uint32(buf[4*idx+3])<<16 | uint32(buf[4*idx+2])<<24)
		default: // big endian
			result[idx] = math.Float32frombits(uint32(buf[4*idx+3]) | uint32(buf[4*idx+2])<<8 | uint32(buf[4*idx+1])<<16 | uint32(buf[4*idx])<<24)
		}
	}
	return result, nil
}

// BytesToInt32s converts byte array to Int32 array in four endian orders. i.e.,
//	BigEndian (0),
//	LittleEndian (1)
//	MidBigEndian (2)
//	MidLittleEndian (3)
func BytesToInt32s(buf []byte, endian Endian) ([]int32, error) {
	l := len(buf)
	if l == 0 || l%4 != 0 {
		return nil, ErrBytesToInt32s
	}
	result := make([]int32, l/4)
	for idx := 0; idx < l/4; idx++ {
		switch endian {
		case DCBA: // little endian
			result[idx] = int32(buf[4*idx]) | int32(buf[4*idx+1])<<8 | int32(buf[4*idx+2])<<16 | int32(buf[4*idx+3])<<24
		case BADC: // mid-big endian
			result[idx] = int32(buf[4*idx+2]) | int32(buf[4*idx+3])<<8 | int32(buf[4*idx])<<16 | int32(buf[4*idx+1])<<24
		case CDAB: // mid-little endian
			result[idx] = int32(buf[4*idx+1]) | int32(buf[4*idx])<<8 | int32(buf[4*idx+3])<<16 | int32(buf[4*idx+2])<<24
		default: // big endian
			result[idx] = int32(buf[4*idx+3]) | int32(buf[4*idx+2])<<8 | int32(buf[4*idx+1])<<16 | int32(buf[4*idx])<<24
		}
	}
	return result, nil
}

// BytesToUInt32s converts byte array to UInt32 array in four endian orders. i.e.,
//	BigEndian (0),
//	LittleEndian (1)
//	MidBigEndian (2)
//	MidLittleEndian (3)
func BytesToUInt32s(buf []byte, endian Endian) ([]uint32, error) {
	l := len(buf)
	if l == 0 || l%4 != 0 {
		return nil, ErrBytesToUInt32s
	}
	result := make([]uint32, l/4)
	for idx := 0; idx < l/4; idx++ {
		switch endian {
		case DCBA: // little endian
			result[idx] = uint32(buf[4*idx]) | uint32(buf[4*idx+1])<<8 | uint32(buf[4*idx+2])<<16 | uint32(buf[4*idx+3])<<24
		case BADC: // mid-big endian
			result[idx] = uint32(buf[4*idx+2]) | uint32(buf[4*idx+3])<<8 | uint32(buf[4*idx])<<16 | uint32(buf[4*idx+1])<<24
		case CDAB: // mid-little endian
			result[idx] = uint32(buf[4*idx+1]) | uint32(buf[4*idx])<<8 | uint32(buf[4*idx+3])<<16 | uint32(buf[4*idx+2])<<24
		default: // big endian
			result[idx] = uint32(buf[4*idx+3]) | uint32(buf[4*idx+2])<<8 | uint32(buf[4*idx+1])<<16 | uint32(buf[4*idx])<<24
		}
	}
	return result, nil
}

// BytesToInt16s converts byte array to Int16 array in two endian orders. i.e.,
//	BigEndian (0) or LittleEndian (1)
func BytesToInt16s(buf []byte, endian Endian) ([]int16, error) {
	l := len(buf)
	if l == 0 || l%2 != 0 {
		return nil, ErrBytesToInt16s
	}
	result := make([]int16, l/2)
	for idx := 0; idx < l/2; idx++ {
		if endian == LittleEndian {
			result[idx] = int16(buf[2*idx]) | int16(buf[2*idx+1])<<8
		} else { // BigEndian
			result[idx] = int16(buf[2*idx+1]) | int16(buf[2*idx])<<8
		}
	}
	return result, nil
}

// BytesToUInt16s converts byte array to UInt16 array in two endian orders. i.e.,
// 	BigEndian (0) or LittleEndian (1)
func BytesToUInt16s(buf []byte, endian Endian) ([]uint16, error) {
	l := len(buf)
	if l == 0 || l%2 != 0 {
		return nil, ErrBytesToUInt16s
	}
	result := make([]uint16, l/2)
	for idx := 0; idx < l/2; idx++ {
		if endian == LittleEndian {
			result[idx] = uint16(buf[2*idx]) | uint16(buf[2*idx+1])<<8
		} else { // BigEndian
			result[idx] = uint16(buf[2*idx+1]) | uint16(buf[2*idx])<<8
		}
	}
	return result, nil
}
