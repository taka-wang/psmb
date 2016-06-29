package psmb

// psbm to modbusd structures - downstream

// DMbtcpRes modbus tcp function code generic response
type DMbtcpRes struct {
	Tid    int64    `json:"tid"`
	Status string   `json:"status"`
	Data   []uint16 `json:"data,omitempty"`
}

// DMbtcpReadReq modbus tcp read request
type DMbtcpReadReq struct {
	Tid   int64  `json:"tid"`
	Cmd   string `json:"cmd"`
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Slave uint8  `json:"slave"`
	Addr  uint16 `json:"addr"`
	Len   uint16 `json:"len"`
}

// DMbtcpSingleWriteReq modbus tcp write single bit/register request
type DMbtcpSingleWriteReq struct {
	Tid   int64  `json:"tid"`
	Cmd   string `json:"cmd"`
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Slave uint8  `json:"slave"`
	Addr  uint16 `json:"addr"`
	Data  uint16 `json:"data"`
}

// DMbtcpMultipleWriteReq modbus tcp write multiple bits/registers request
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

// DMbtcpTimeout modbus tcp set/get timeout request/response
type DMbtcpTimeout struct {
	Tid     int64  `json:"tid"`
	Cmd     string `json:"cmd"`
	Status  string `json:"status,omitempty"`
	Timeout int64  `json:"timeout,omitempty"`
}
