package psmb

// psbm to modbusd structures - downstream

// DMbtcpRes Modbus tcp generic response (include write response)
type DMbtcpRes struct {
	Tid    int64  `json:"tid"`
	Status string `json:"status"`
}

// DMbtcpReadReq Modbus tcp read request
type DMbtcpReadReq struct {
	Tid   int64  `json:"tid"`
	Cmd   string `json:"cmd"`
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Slave uint8  `json:"slave"`
	Addr  uint16 `json:"addr"`
	Len   uint16 `json:"len"`
}

// DMbtcpReadRes Modbus tcp read bits/registers response
type DMbtcpReadRes struct {
	Tid    int64    `json:"tid"`
	Status string   `json:"status"`
	Data   []uint16 `json:"data"`
}

// DMbtcpTimeoutReq modbus tcp timeout request
type DMbtcpTimeoutReq struct {
	Tid     int64  `json:"tid"`
	Cmd     string `json:"cmd"`
	Timeout int64  `json:"timeout,omitempty"`
}

// DMbtcpTimeoutRes modbus tcp timeout request
type DMbtcpTimeoutRes struct {
	Tid     int64  `json:"tid"`
	Cmd     string `json:"cmd"`
	Status  string `json:"status"`
	Timeout int64  `json:"timeout,omitempty"`
}

// DMbtcpSingleWriteReq Modbus tcp write single bit/register request
type DMbtcpSingleWriteReq struct {
	Tid   int64  `json:"tid"`
	Cmd   string `json:"cmd"`
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Slave uint8  `json:"slave"`
	Addr  uint16 `json:"addr"`
	Data  uint16 `json:"data"`
}

// DMbtcpMultipleWriteReq Modbus tcp write multiple bits/registers request
type DMbtcpMultipleWriteReq struct {
	Tid   int64    `json:"tid"`
	Cmd   string   `json:"cmd"`
	IP    string   `json:"ip"`
	Port  string   `json:"port"`
	Slave uint8    `json:"slave"`
	Addr  uint16   `json:"addr"`
	Len   uint16   `json:"len"`
	Data  []uint16 `json:"data"`
}
