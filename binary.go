package psmb

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Endian Endianness:
type Endian int

const (
	ABCD Endian = iota
	DCBA
	BADC
	CDAB
)
const (
	AB Endian = iota
	BA
)
const (
	BigEndian Endian = iota
	LittleEndian
	MidBigEndian
	MidLittleEndian
)

// BytesToHexString Convert byte array to hex string
func BytesToHexString(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

//DecimalStringToRegisters Convert decimal string to uint16 array in big endian order
// limitation: leading space before comma is not allow
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

// HexStringToRegisters Convert hex string to uint16 array in big endian order
func HexStringToRegisters(hexString string) ([]uint16, error) {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}
	fmt.Println(bytes)
	return BytesToUInt16s(bytes, 0)
}

// RegistersToBytes Convert register(uint16) array to byte array in big endian order
func RegistersToBytes(data []uint16) []byte {
	buf := new(bytes.Buffer)
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			fmt.Println("RegistersToBytes failed:", err)
		}
	}
	return buf.Bytes()
}

// BytesToFloat32s Byte array to float32 array
func BytesToFloat32s(buf []byte, endian Endian) ([]float32, error) {
	l := len(buf)
	if l == 0 || l%4 != 0 {
		return nil, errors.New("Invalid data length")
	}

	result := make([]float32, l/4)
	for idx := 0; idx < l/4; idx++ {
		switch endian {
		case 1: // DCBA, little endian
			result[idx] = math.Float32frombits(uint32(buf[4*idx]) | uint32(buf[4*idx+1])<<8 | uint32(buf[4*idx+2])<<16 | uint32(buf[4*idx+3])<<24)
		case 2: // BADC, mid-big endian
			result[idx] = math.Float32frombits(uint32(buf[4*idx+2]) | uint32(buf[4*idx+3])<<8 | uint32(buf[4*idx])<<16 | uint32(buf[4*idx+1])<<24)
		case 3: // CDAB, mid-little endian
			result[idx] = math.Float32frombits(uint32(buf[4*idx+1]) | uint32(buf[4*idx])<<8 | uint32(buf[4*idx+3])<<16 | uint32(buf[4*idx+2])<<24)
		default: // ABCD, big endian
			result[idx] = math.Float32frombits(uint32(buf[4*idx+3]) | uint32(buf[4*idx+2])<<8 | uint32(buf[4*idx+1])<<16 | uint32(buf[4*idx])<<24)
		}
	}
	return result, nil
}

// BytesToInt32s Byte array to Int32 array
func BytesToInt32s(buf []byte, endian Endian) ([]int32, error) {
	l := len(buf)
	if l == 0 || l%4 != 0 {
		return nil, errors.New("Invalid data length")
	}
	result := make([]int32, l/4)
	for idx := 0; idx < l/4; idx++ {
		switch endian {
		case 1: // DCBA, little endian
			result[idx] = int32(buf[4*idx]) | int32(buf[4*idx+1])<<8 | int32(buf[4*idx+2])<<16 | int32(buf[4*idx+3])<<24
		case 2: // BADC, mid-big endian
			result[idx] = int32(buf[4*idx+2]) | int32(buf[4*idx+3])<<8 | int32(buf[4*idx])<<16 | int32(buf[4*idx+1])<<24
		case 3: // CDAB, mid-little endian
			result[idx] = int32(buf[4*idx+1]) | int32(buf[4*idx])<<8 | int32(buf[4*idx+3])<<16 | int32(buf[4*idx+2])<<24
		default: // ABCD, big endian
			result[idx] = int32(buf[4*idx+3]) | int32(buf[4*idx+2])<<8 | int32(buf[4*idx+1])<<16 | int32(buf[4*idx])<<24
		}
	}
	return result, nil
}

// BytesToUInt32s Byte array to UInt32 array
func BytesToUInt32s(buf []byte, endian Endian) ([]uint32, error) {
	l := len(buf)
	if l == 0 || l%4 != 0 {
		return nil, errors.New("Invalid data length")
	}
	result := make([]uint32, l/4)
	for idx := 0; idx < l/4; idx++ {
		switch endian {
		case 1: // DCBA, little endian
			result[idx] = uint32(buf[4*idx]) | uint32(buf[4*idx+1])<<8 | uint32(buf[4*idx+2])<<16 | uint32(buf[4*idx+3])<<24
		case 2: // BADC, mid-big endian
			result[idx] = uint32(buf[4*idx+2]) | uint32(buf[4*idx+3])<<8 | uint32(buf[4*idx])<<16 | uint32(buf[4*idx+1])<<24
		case 3: // CDAB, mid-little endian
			result[idx] = uint32(buf[4*idx+1]) | uint32(buf[4*idx])<<8 | uint32(buf[4*idx+3])<<16 | uint32(buf[4*idx+2])<<24
		default: // ABCD, big endian
			result[idx] = uint32(buf[4*idx+3]) | uint32(buf[4*idx+2])<<8 | uint32(buf[4*idx+1])<<16 | uint32(buf[4*idx])<<24
		}
	}
	return result, nil
}

// BytesToInt16s Byte array to Int16 array
func BytesToInt16s(buf []byte, endian Endian) ([]int16, error) {
	l := len(buf)
	if l == 0 || l%2 != 0 {
		return nil, errors.New("Invalid data length")
	}
	result := make([]int16, l/2)
	for idx := 0; idx < l/2; idx++ {
		if endian == 1 {
			result[idx] = int16(buf[2*idx]) | int16(buf[2*idx+1])<<8
		} else {
			result[idx] = int16(buf[2*idx+1]) | int16(buf[2*idx])<<8
		}
	}
	return result, nil
}

// BytesToUInt16s Byte array to UInt16 array
func BytesToUInt16s(buf []byte, endian Endian) ([]uint16, error) {
	l := len(buf)
	if l == 0 || l%2 != 0 {
		return nil, errors.New("Invalid data length")
	}
	result := make([]uint16, l/2)
	for idx := 0; idx < l/2; idx++ {
		if endian == 1 {
			result[idx] = uint16(buf[2*idx]) | uint16(buf[2*idx+1])<<8
		} else {
			result[idx] = uint16(buf[2*idx+1]) | uint16(buf[2*idx])<<8
		}
	}
	return result, nil
}
