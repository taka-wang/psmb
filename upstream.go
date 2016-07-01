package main

import (
	"fmt"
	"strings"
)

// services to psbm structures - upstream

// JSONableByteSlice jsonable uint8 array
type JSONableByteSlice []byte

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

// MbtcpOnceReadReq read coil/register request (1.1)
type MbtcpOnceReadReq struct {
	Tid   int64       `json:"tid"`
	From  string      `json:"from,omitempty"`
	FC    int         `json:"fc"`
	IP    string      `json:"ip"`
	Port  string      `json:"port,omitempty"`
	Slave uint8       `json:"slave"`
	Addr  uint16      `json:"addr"`
	Len   uint16      `json:"len,omitempty"`
	Type  RegDataType `json:"type,omitempty"`
	Order Endian      `json:"order,omitempty"`
	Range struct {
		OriginLow  int `json:"a"`
		OriginHigh int `json:"b"`
		TargetLow  int `json:"c"`
		TargetHigh int `json:"d"`
	} `json:"range,omitempty"`
}

// MbtcpOnceReadRes read coil/register response (1.1)
type MbtcpOnceReadRes struct {
	Tid    int64       `json:"tid"`
	Status string      `json:"status"`
	Type   RegDataType `json:"type,omitempty"`
	// Bytes FC3, FC4 and Type 2~8 only
	Bytes JSONableByteSlice `json:"bytes,omitempty"`
	Data  interface{}       `json:"data,omitempty"`
	// []uint16, []int16, []uint32, []int32, []float32, string
}

// MbtcpTimeoutReq set/get TCP connection timeout request (1.3, 1.4)
type MbtcpTimeoutReq struct {
	Tid  int64  `json:"tid"`
	From string `json:"from,omitempty"`
	Data int64  `json:"timeout,omitempty"`
}

// MbtcpTimeoutRes set/get TCP connection timeout response (1.3, 1.4)
type MbtcpTimeoutRes struct {
	Tid    int64  `json:"tid"`
	Status string `json:"status"`
	Data   int64  `json:"timeout,omitempty"`
}

// MbtcpSimpleReq generic modbus tcp response
type MbtcpSimpleReq struct {
	Tid  int64  `json:"tid"`
	From string `json:"from,omitempty"`
}

// MbtcpSimpleRes generic modbus tcp response
type MbtcpSimpleRes struct {
	Tid    int64  `json:"tid"`
	Status string `json:"status"`
}
