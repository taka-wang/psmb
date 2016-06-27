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

// BytesToHexStr Convert byte array to hex string
func BytesToHexStr(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

//DecStrToRegs Convert decimal string to uint16 array in big endian order
func DecStrToRegs(decString string) ([]uint16, error) {
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

// HexStrToRegs Convert hex string to uint16 array in big endian order
func HexStrToRegs(hexString string) ([]uint16, error) {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}
	fmt.Println(bytes)
	return BytesToUInt16s(bytes, 0)
}

// RegsToBytes Convert register(uint16) array to byte array in big endian order
func RegsToBytes(data []uint16) []byte {
	buf := new(bytes.Buffer)
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			fmt.Println("RegsToBytes failed:", err)
		}
	}
	return buf.Bytes()
}

// BytesToFloat32s Byte array to float32 array
func BytesToFloat32s(buf []byte, endian int) ([]float32, error) {
	if len(buf)%4 != 0 {
		return nil, errors.New("Invalid data")
	}
	result := []float32{}
	for idx := 0; idx < len(buf); idx = idx + 4 {
		switch endian {
		case 1: // DCBA, little endian
			result = append(result, math.Float32frombits(uint32(buf[idx])|uint32(buf[idx+1])<<8|uint32(buf[idx+2])<<16|uint32(buf[idx+3])<<24))
		case 2: // BADC, mid-big endian
			result = append(result, math.Float32frombits(uint32(buf[idx+2])|uint32(buf[idx+3])<<8|uint32(buf[idx])<<16|uint32(buf[idx+1])<<24))
		case 3: // CDAB, mid-little endian
			result = append(result, math.Float32frombits(uint32(buf[idx+1])|uint32(buf[idx])<<8|uint32(buf[idx+3])<<16|uint32(buf[idx+2])<<24))
		default: // ABCD, big endian
			result = append(result, math.Float32frombits(uint32(buf[idx+3])|uint32(buf[idx+2])<<8|uint32(buf[idx+1])<<16|uint32(buf[idx])<<24))
		}
	}
	return result, nil
}

// BytesToInt32s Byte array to Int32 array
func BytesToInt32s(buf []byte, endian int) ([]int32, error) {
	if len(buf)%4 != 0 {
		return nil, errors.New("Invalid data")
	}
	result := []int32{}
	for idx := 0; idx < len(buf); idx = idx + 4 {
		switch endian {
		case 1: // DCBA, little endian
			result = append(result, int32(buf[idx])|int32(buf[idx+1])<<8|int32(buf[idx+2])<<16|int32(buf[idx+3])<<24)
		case 2: // BADC, mid-big endian
			result = append(result, int32(buf[idx+2])|int32(buf[idx+3])<<8|int32(buf[idx])<<16|int32(buf[idx+1])<<24)
		case 3: // CDAB, mid-little endian
			result = append(result, int32(buf[idx+1])|int32(buf[idx])<<8|int32(buf[idx+3])<<16|int32(buf[idx+2])<<24)
		default: // ABCD, big endian
			result = append(result, int32(buf[idx+3])|int32(buf[idx+2])<<8|int32(buf[idx+1])<<16|int32(buf[idx])<<24)
		}
	}
	return result, nil
}

// BytesToUInt32s Byte array to UInt32 array
func BytesToUInt32s(buf []byte, endian int) ([]uint32, error) {
	if len(buf)%4 != 0 {
		return nil, errors.New("Invalid data")
	}
	result := []uint32{}
	for idx := 0; idx < len(buf); idx = idx + 4 {
		switch endian {
		case 1: // DCBA, little endian
			result = append(result, uint32(buf[idx])|uint32(buf[idx+1])<<8|uint32(buf[idx+2])<<16|uint32(buf[idx+3])<<24)
		case 2: // BADC, mid-big endian
			result = append(result, uint32(buf[idx+2])|uint32(buf[idx+3])<<8|uint32(buf[idx])<<16|uint32(buf[idx+1])<<24)
		case 3: // CDAB, mid-little endian
			result = append(result, uint32(buf[idx+1])|uint32(buf[idx])<<8|uint32(buf[idx+3])<<16|uint32(buf[idx+2])<<24)
		default: // ABCD, big endian
			result = append(result, uint32(buf[idx+3])|uint32(buf[idx+2])<<8|uint32(buf[idx+1])<<16|uint32(buf[idx])<<24)
		}
	}
	return result, nil
}

// BytesToInt16s Byte array to Int16 array
func BytesToInt16s(buf []byte, endian int) ([]int16, error) {
	if len(buf)%2 != 0 {
		return nil, errors.New("Invalid data")
	}
	result := []int16{}

	for idx := 0; idx < len(buf); idx = idx + 2 {
		if endian == 1 {
			result = append(result, int16(buf[idx])|int16(buf[idx+1])<<8)
		} else {
			result = append(result, int16(buf[idx+1])|int16(buf[idx])<<8)
		}
	}
	return result, nil
}

// BytesToUInt16s Byte array to UInt16 array
func BytesToUInt16s(buf []byte, endian int) ([]uint16, error) {
	if len(buf)%2 != 0 {
		return nil, errors.New("Invalid data")
	}
	result := []uint16{}

	for idx := 0; idx < len(buf); idx = idx + 2 {
		if endian == 1 {
			result = append(result, uint16(buf[idx])|uint16(buf[idx+1])<<8)
		} else {
			result = append(result, uint16(buf[idx+1])|uint16(buf[idx])<<8)
		}
	}
	return result, nil
}
