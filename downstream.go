package psmb

// To downstream

// dMbtcpRes Modbus tcp generic response
type dMbtcpRes struct {
	Tid    uint64 `json:"tid"`
	Status string `json:"status"`
}

// dMbtcpReadReq Modbus tcp read request
type dMbtcpReadReq struct {
	Tid   uint64 `json:"tid"`
	Cmd   string `json:"cmd"`
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Slave uint8  `json:"slave"`
	Addr  uint16 `json:"addr"`
	Len   uint16 `json:"len"`
}

// dMbtcpReadRes Modbus tcp read bits/registers response
type dMbtcpReadRes struct {
	Tid    uint64   `json:"tid"`
	Status string   `json:"status"`
	Data   []uint16 `json:"data"`
}

// dMbtcpTimeoutReq modbus tcp timeout request
type dMbtcpTimeoutReq struct {
	Tid     uint64 `json:"tid"`
	Cmd     string `json:"cmd"`
	Timeout int64  `json:"timeout,omitempty"`
}

// dMbtcpTimeoutRes modbus tcp timeout request
type dMbtcpTimeoutRes struct {
	Tid     uint64 `json:"tid"`
	Cmd     string `json:"cmd"`
	Status  string `json:"status"`
	Timeout int64  `json:"timeout,omitempty"`
}

// dMbtcpSingleWriteReq Modbus tcp write request
type dMbtcpSingleWriteReq struct {
	Tid   uint64 `json:"tid"`
	Cmd   string `json:"cmd"`
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Slave uint8  `json:"slave"`
	Addr  uint16 `json:"addr"`
	Data  uint16 `json:"data"`
}

// dMbtcpMultipleWriteReq Modbus tcp write request
type dMbtcpMultipleWriteReq struct {
	Tid   uint64   `json:"tid"`
	Cmd   string   `json:"cmd"`
	IP    string   `json:"ip"`
	Port  string   `json:"port"`
	Slave uint8    `json:"slave"`
	Addr  uint16   `json:"addr"`
	Len   uint16   `json:"len"`
	Data  []uint16 `json:"data"`
}
