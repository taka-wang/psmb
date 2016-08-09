package psmb

//
// services to psmb structures - upstream
//

type (

	// MbtcpReadReq read coil/register request (1.1).
	// Scale range field example:
	// 	Range: &ScaleRange{1,2,3,4},
	MbtcpReadReq struct {
		Tid   int64        `json:"tid"`
		From  string       `json:"from,omitempty"`
		FC    int          `json:"fc"`
		IP    string       `json:"ip"`
		Port  string       `json:"port,omitempty"`
		Slave uint8        `json:"slave"`
		Addr  uint16       `json:"addr"`
		Len   uint16       `json:"len,omitempty"`
		Type  RegValueType `json:"type,omitempty"`
		Order Endian       `json:"order,omitempty"`
		Range *ScaleRange  `json:"range,omitempty"` // point to struct can be omitted in json encode
	}

	// MbtcpReadRes read coil/register response (1.1).
	// `Data interface` supports:
	//	[]uint16, []int16, []uint32, []int32, []float32, string
	MbtcpReadRes struct {
		Tid    int64        `json:"tid"`
		Status string       `json:"status"`
		Type   RegValueType `json:"type,omitempty"`
		// Bytes FC3, FC4 and Type 2~8 only
		Bytes JSONableByteSlice `json:"bytes,omitempty"`
		Data  interface{}       `json:"data,omitempty"` // universal data container
	}

	// MbtcpWriteReq write coil/register request
	MbtcpWriteReq struct {
		Tid   int64       `json:"tid"`
		From  string      `json:"from,omitempty"`
		FC    int         `json:"fc"`
		IP    string      `json:"ip"`
		Port  string      `json:"port,omitempty"`
		Slave uint8       `json:"slave"`
		Addr  uint16      `json:"addr"`
		Len   uint16      `json:"len,omitempty"`
		Hex   bool        `json:"hex,omitempty"`
		Data  interface{} `json:"data"`
	}

	// MbtcpWriteRes == MbtcpSimpleRes

	// MbtcpTimeoutReq set/get TCP connection timeout request (1.3, 1.4)
	MbtcpTimeoutReq struct {
		Tid  int64  `json:"tid"`
		From string `json:"from,omitempty"`
		Data int64  `json:"timeout,omitempty"`
	}

	// MbtcpTimeoutRes set/get TCP connection timeout response (1.3, 1.4)
	MbtcpTimeoutRes struct {
		Tid    int64  `json:"tid"`
		Status string `json:"status"`
		Data   int64  `json:"timeout,omitempty"`
	}

	/*
	   // MbtcpSimpleReq generic modbus tcp response
	    MbtcpSimpleReq struct {
	   	Tid  int64  `json:"tid"`
	   	From string `json:"from,omitempty"`
	   }
	*/

	// MbtcpSimpleRes generic modbus tcp response
	MbtcpSimpleRes struct {
		Tid    int64  `json:"tid"`
		Status string `json:"status"`
	}

	// MbtcpPollStatus polling coil/register request;
	MbtcpPollStatus struct {
		Tid      int64        `json:"tid"`
		From     string       `json:"from,omitempty"`
		Name     string       `json:"name"`
		Interval uint64       `json:"interval"`
		Enabled  bool         `json:"enabled"`
		FC       int          `json:"fc"`
		IP       string       `json:"ip"`
		Port     string       `json:"port,omitempty"`
		Slave    uint8        `json:"slave"`
		Addr     uint16       `json:"addr"`
		Status   string       `json:"status,omitempty"` // 2.3.2 response only
		Len      uint16       `json:"len,omitempty"`
		Type     RegValueType `json:"type,omitempty"`
		Order    Endian       `json:"order,omitempty"`
		Range    *ScaleRange  `json:"range,omitempty"` // point to struct can be omitted in json encode
	}

	// MbtcpPollRes == MbtcpSimpleRes

	// MbtcpPollData read coil/register response (1.1).
	// `Data interface` supports:
	// 	[]uint16, []int16, []uint32, []int32, []float32, string
	MbtcpPollData struct {
		TimeStamp int64        `json:"ts"`
		Name      string       `json:"name"`
		Status    string       `json:"status"`
		Type      RegValueType `json:"type,omitempty"`
		// Bytes FC3, FC4 and Type 2~8 only
		Bytes JSONableByteSlice `json:"bytes,omitempty"`
		Data  interface{}       `json:"data,omitempty"` // universal data container
	}

	// MbtcpPollOpReq generic modbus tcp poll operation request
	MbtcpPollOpReq struct {
		Tid      int64  `json:"tid"`
		From     string `json:"from,omitempty"`
		Name     string `json:"name,omitempty"`
		Interval uint64 `json:"interval,omitempty"`
		Enabled  bool   `json:"enabled,omitempty"`
	}

	// MbtcpPollsStatus requests status
	MbtcpPollsStatus struct {
		Tid    int64             `json:"tid"`
		From   string            `json:"from,omitempty"`
		Status string            `json:"status,omitempty"`
		Polls  []MbtcpPollStatus `json:"polls"`
	}

	// MbtcpHistoryData history data
	MbtcpHistoryData struct {
		Tid    int64       `json:"tid"`
		Name   string      `json:"name"`
		Status string      `json:"status"`
		Data   interface{} `json:"history,omitempty"` // universal data container
	}

	// HistoryData history data
	HistoryData struct {
		Ts   int64       `json:"ts,omitempty"`
		Data interface{} `json:"data,omitempty"` // universal data container
	}

	// MbtcpFilterStatus filter status
	MbtcpFilterStatus struct {
		Tid     int64        `json:"tid"`
		From    string       `json:"from,omitempty"`
		Name    string       `json:"name"`
		Enabled bool         `json:"enabled"`
		Type    RegValueType `json:"type,omitempty"`
		Arg     []float32    `json:"arg,omitempty"`
		Status  string       `json:"status,omitempty"`
	}

	// MbtcpFilterOpReq generic modbus tcp filter operation request
	MbtcpFilterOpReq struct {
		Tid     int64  `json:"tid"`
		From    string `json:"from,omitempty"`
		Name    string `json:"name,omitempty"`
		Enabled bool   `json:"enabled,omitempty"`
	}

	// MbtcpFiltersStatus requests status
	MbtcpFiltersStatus struct {
		Tid     int64               `json:"tid"`
		From    string              `json:"from,omitempty"`
		Status  string              `json:"status,omitempty"`
		Filters []MbtcpFilterStatus `json:"filters"`
	}
)
