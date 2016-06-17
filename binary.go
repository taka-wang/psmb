package psmb

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// byte to hex string
//fmt.Println(hex.EncodeToString(buf2))

//DecStrToRegs Convert decimal string to uint16 array in big endian order
func DecStrToRegs(decstr string) ([]uint16, error) {
	var result = []uint16{}
	for _, v := range strings.Split(decstr, ",") {
		i, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			return nil, err
		}
		result = append(result, uint16(i))
	}
	return result, nil
}

// HexStrToRegs Convert hex string to uint16 array in big endian order
func HexStrToRegs(hexstr string) ([]uint16, error) {
	bytes, err := hex.DecodeString(hexstr)
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

// BytesToInt32s Byte array to Int32 array
func BytesToInt32s(buf []byte, order int) ([]int32, error) {
	if len(buf)%4 != 0 {
		return nil, errors.New("Invalid data")
	}
	return []int32{1, 2}, nil
}

// BytesToUInt32s Byte array to UInt32 array
func BytesToUInt32s(buf []byte, order int) ([]uint32, error) {
	if len(buf)%4 != 0 {
		return nil, errors.New("Invalid data")
	}
	return []uint32{1, 2}, nil
}

// BytesToFloat32s Byte array to float32 array
func BytesToFloat32s(buf []byte, order int) ([]float32, error) {
	if len(buf)%4 != 0 {
		return nil, errors.New("Invalid data")
	}
	return []float32{1, 2}, nil
}

// BytesToInt16s Byte array to Int16 array
func BytesToInt16s(buf []byte, order int) ([]int16, error) {
	if len(buf)%2 != 0 {
		return nil, errors.New("Invalid data")
	}
	return []int16{1, 2}, nil
}

// BytesToUInt16s Byte array to UInt16 array
func BytesToUInt16s(buf []byte, order int) ([]uint16, error) {
	if len(buf)%2 != 0 {
		return nil, errors.New("Invalid data")
	}
	return []uint16{1, 2}, nil
}
