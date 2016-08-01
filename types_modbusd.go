package psmb

//
// psmb to modbusd structures - downstream
//

// DMbtcpRes downstream modbus tcp read/write response
type DMbtcpRes struct {
	// Tid unique transaction id in string format
	Tid string `json:"tid"`
	// Status response status string
	Status string `json:"status"`
	// Data for 'read function code' only.
	Data []uint16 `json:"data,omitempty"`
}

// DMbtcpReadReq downstream modbus tcp read request
type DMbtcpReadReq struct {
	// Tid unique transaction id in `string` format
	Tid string `json:"tid"`
	// Cmd modbusd command type: https://github.com/taka-wang/modbusd#command-mapping-table
	Cmd int `json:"cmd"`
	// IP ip address or hostname of the modbus tcp slave
	IP string `json:"ip"`
	// Port port number of the modbus tcp slave
	Port string `json:"port"`
	// Slave device id of the modbus tcp slave
	Slave uint8 `json:"slave"`
	// Addr start address for read
	Addr uint16 `json:"addr"`
	// Len the length of registers or bits
	Len uint16 `json:"len"`
}

// DMbtcpWriteReq downstream modbus tcp write single bit/register request
type DMbtcpWriteReq struct {
	// Tid unique transaction id in `string` format
	Tid string `json:"tid"`
	// Cmd modbusd command type: https://github.com/taka-wang/modbusd#command-mapping-table
	Cmd int `json:"cmd"`
	// IP ip address or hostname of the modbus tcp slave
	IP string `json:"ip"`
	// Port port number of the modbus tcp slave
	Port string `json:"port"`
	// Slave device id of the modbus tcp slave
	Slave uint8 `json:"slave"`
	// Addr start address for write
	Addr uint16 `json:"addr"`
	// Len omit for fc5, fc6
	Len uint16 `json:"len,omitempty"`
	// Data should be []uint16, uint16 (FC5, FC6)
	Data interface{} `json:"data"`
}

// DMbtcpTimeout downstream modbus tcp set/get timeout request/response
type DMbtcpTimeout struct {
	// Tid unique transaction id in `string` format
	Tid string `json:"tid"`
	// Cmd modbusd command type: https://github.com/taka-wang/modbusd#command-mapping-table
	Cmd int `json:"cmd"`
	// Status for response only.
	Status string `json:"status,omitempty"`
	// Timeout set timeout request and get timeout response only.
	Timeout int64 `json:"timeout,omitempty"`
}
