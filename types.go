package psmb

import (
	"fmt"
	"strings"
)

type (

	// ScaleRange defines scale range
	ScaleRange struct {
		DomainLow  float64 `json:"a"`
		DomainHigh float64 `json:"b"`
		RangeLow   float64 `json:"c"`
		RangeHigh  float64 `json:"d"`
	}

	// JSONableByteSlice jsonable uint8 array
	JSONableByteSlice []byte

	// Endian defines byte endianness
	Endian int

	// RegValueType return value type defines how to inteprete registers, i.e.,
	//  for modbus read function codes only
	RegValueType int

	// FilterType filter type
	FilterType int
)

// MarshalJSON implements the Marshaler interface on JSONableByteSlice (i.e., uint8/byte array).
// Ref: http://stackoverflow.com/questions/14177862/how-to-jsonize-a-uint8-slice-in-go
func (u JSONableByteSlice) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}
	return []byte(result), nil
}

// 16-bits Endian
const (
	_ Endian = iota // ignore first value by assigning to blank identifier
	// AB 16-bit words may be represented in big-endian format
	AB
	// BA 16-bit words may be represented in little-endian format
	BA
)

// 32-bits Endian
const (
	_ Endian = iota // ignore first value by assigning to blank identifier
	// ABCD 32-bit words may be represented in big-endian format
	ABCD
	// DCBA 32-bit words may be represented in little-endian format
	DCBA
	// BADC 32-bit words may be represented in mid-big-endian format
	BADC
	// CDAB 32-bit words may be represented in mid-little-endian format
	CDAB
)

// 32-bits Endian
const (
	_ Endian = iota // ignore first value by assigning to blank identifier
	// BigEndian 32-bit words may be represented in ABCD format
	BigEndian
	// LittleEndian 32-bit words may be represented in DCBA format
	LittleEndian
	// MidBigEndian 32-bit words may be represented in BADC format
	MidBigEndian
	// MidLittleEndian 32-bit words may be represented in CDAB format
	MidLittleEndian
)

// Register value type for read function
const (
	_ RegValueType = iota // ignore first value by assigning to blank identifier
	// RegisterArray register array, ex: [12345, 23456, 5678]
	RegisterArray
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

// Filter value type
const (
	// Change change or not
	Change FilterType = iota
	// GreaterEqual greater than or equal
	GreaterEqual
	// Greater greater than
	Greater
	// Equal equal
	Equal
	// Less less than
	Less
	// LessEqual less than or equal
	LessEqual
	// InsideRange inside range
	InsideRange
	// InsideIncRange inside range (include)
	InsideIncRange
	// OutsideRange outside range
	OutsideRange
	// OutsideIncRange outside range (include)
	OutsideIncRange
)
