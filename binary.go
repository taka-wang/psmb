package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"log"
	"math"
	"strconv"
	"strings"
)

// Endian defines byte endianness
type Endian int

const (
	// ABCD 32-bit words may be represented in big-endian format
	ABCD Endian = iota
	// DCBA 32-bit words may be represented in little-endian format
	DCBA
	// BADC 32-bit words may be represented in mid-big-endian format
	BADC
	// CDAB 32-bit words may be represented in mid-little-endian format
	CDAB
)
const (
	// AB 16-bit words may be represented in big-endian format
	AB Endian = iota
	// BA 16-bit words may be represented in little-endian format
	BA
)
const (
	// BigEndian 32-bit words may be represented in ABCD format
	BigEndian Endian = iota
	// LittleEndian 32-bit words may be represented in DCBA format
	LittleEndian
	// MidBigEndian 32-bit words may be represented in BADC format
	MidBigEndian
	// MidLittleEndian 32-bit words may be represented in CDAB format
	MidLittleEndian
)

// RegDataType defines how to inteprete registers
type RegDataType int

const (
	// RegisterArray register array, ex: [12345, 23456, 5678]
	RegisterArray RegDataType = iota
	// HexString hexadecimal string, ex: "112C004F12345678"
	HexString
	// Scale linearly scale
	Scale
	// UInt16 uint16 array
	UInt16
	// Int16 int16 array
	Int16
	// UInt32 uint32 array
	UInt32
	// Int32 int32 array
	Int32
	// Float32 float32 array
	Float32
)

// bitStringToUInt8s converts bits string to uint8 array.
// Source: function code 15
func bitStringToUInt8s(bitString string) ([]uint8, error) {
	var result = []uint8{}
	for _, v := range strings.Split(bitString, ",") {
		i, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			return nil, err
		}
		result = append(result, uint8(i))
	}
	return result, nil
}

// DecimalStringToRegisters converts decimal string to uint16/words array in big endian order.
// Limitation: leading space before comma is not allowed.
// Source: upstream.
func DecimalStringToRegisters(decString string) ([]uint16, error) {
	var result = []uint16{}
	for _, v := range strings.Split(decString, ",") {
		i, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			return nil, err
		}
		result = append(result, uint16(i))
	}
	return result, nil
}

// HexStringToRegisters converts hexadecimal string to uint16/words array in big endian order.
// Source: upstream.
func HexStringToRegisters(hexString string) ([]uint16, error) {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}
	return BytesToUInt16s(bytes, 0)
}

// BytesToHexString converts bytes array to hexadecimal string.
// Example: 112C004F12345678
func BytesToHexString(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// RegistersToBytes converts registers/uint16 array to byte array in big endian order.
// Source: downstream.
func RegistersToBytes(data []uint16) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			log.Println("RegistersToBytes failed:", err)
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// LinearScalingRegisters scales the registers linearly.
// Equation:
// 	Let originLow, originHigh, targetLow, targetHigh as a, b, c, d accordingly.
// 	Output = c + (d - c) * (input - a) / (b - a)
func LinearScalingRegisters(data []uint16, originLow, originHigh, targetLow, targetHigh float64) []float32 {
	l := len(data)
	low := math.Min(originLow, originHigh)
	high := math.Max(originLow, originHigh)
	result := make([]float32, l)

	for idx := 0; idx < l; idx++ {
		result[idx] = float32(targetLow + (targetHigh-targetLow)*(math.Min(math.Max(low, float64(data[idx])), high)-low)/(high-low))
	}
	return result
}

// BytesToFloat32s converts byte array to float32 array in four endian orders. i.e.,
//       BigEndian (0),
//       LittleEndian (1)
//       MidBigEndian (2)
//       MidLittleEndian (3)
func BytesToFloat32s(buf []byte, endian Endian) ([]float32, error) {
	l := len(buf)
	if l == 0 || l%4 != 0 {
		return nil, errors.New("Invalid data length")
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
//       BigEndian (0),
//       LittleEndian (1)
//       MidBigEndian (2)
//       MidLittleEndian (3)
func BytesToInt32s(buf []byte, endian Endian) ([]int32, error) {
	l := len(buf)
	if l == 0 || l%4 != 0 {
		return nil, errors.New("Invalid data length")
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
//       BigEndian (0),
//       LittleEndian (1)
//       MidBigEndian (2)
//       MidLittleEndian (3)
func BytesToUInt32s(buf []byte, endian Endian) ([]uint32, error) {
	l := len(buf)
	if l == 0 || l%4 != 0 {
		return nil, errors.New("Invalid data length")
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

// BytesToInt16s converts byte array to Int16 array in two endian orders.
// i.e., BigEndian (0) or LittleEndian (1)
func BytesToInt16s(buf []byte, endian Endian) ([]int16, error) {
	l := len(buf)
	if l == 0 || l%2 != 0 {
		return nil, errors.New("Invalid data length")
	}
	result := make([]int16, l/2)
	for idx := 0; idx < l/2; idx++ {
		if endian == LittleEndian {
			result[idx] = int16(buf[2*idx]) | int16(buf[2*idx+1])<<8
		} else {
			result[idx] = int16(buf[2*idx+1]) | int16(buf[2*idx])<<8
		}
	}
	return result, nil
}

// BytesToUInt16s converts byte array to UInt16 array in two endian orders.
// i.e., BigEndian (0) or LittleEndian (1)
func BytesToUInt16s(buf []byte, endian Endian) ([]uint16, error) {
	l := len(buf)
	if l == 0 || l%2 != 0 {
		return nil, errors.New("Invalid data length")
	}
	result := make([]uint16, l/2)
	for idx := 0; idx < l/2; idx++ {
		if endian == LittleEndian {
			result[idx] = uint16(buf[2*idx]) | uint16(buf[2*idx+1])<<8
		} else {
			result[idx] = uint16(buf[2*idx+1]) | uint16(buf[2*idx])<<8
		}
	}
	return result, nil
}
