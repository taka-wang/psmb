package psmb

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/takawang/sugar"
)

func TestUpstreamStructTest(t *testing.T) {

	s := sugar.New(t)

	s.Title("One-off modbus tcp struct tests")

	s.Assert("`mbtcp.once.read` request test", func(logf sugar.Log) bool {
		input :=
			`{
                "from": "web",
                "tid": 123456,
                "fc" : 1,
                "ip": "192.168.0.1",
                "port": "503",
                "slave": 1,
                "addr": 10,
                "len": 4
            }`
		var r1 MbtcpReadReq
		if err := json.Unmarshal([]byte(input), &r1); err != nil {
			logf("json err:", err)
			return false
		}
		logf(r1)

		input2 :=
			`{
                "from": "web",
                "tid": 123456,
                "fc" : 3,
                "ip": "192.168.0.1",
                "port": "503",
                "slave": 1,
                "addr": 10,
                "len": 4,
                "type": 1
            }`
		var r2 MbtcpReadReq
		if err := json.Unmarshal([]byte(input2), &r2); err != nil {
			logf("json err:", err)
			return false
		}
		logf(r2)

		input3 :=
			`{
                "from": "web",
                "tid": 123456,
                "fc" : 3,
                "ip": "192.168.0.1",
                "port": "503",
                "slave": 1,
                "addr": 10,
                "len": 4,
                "type": 3,
                "range": 
                {
                    "a": 0,
                    "b": 65535,
                    "c": 100,
                    "d": 500
                }
            }`
		var r3 MbtcpReadReq
		if err := json.Unmarshal([]byte(input3), &r3); err != nil {
			logf(err)
			return false
		}
		logf(r3)

		input4 :=
			`{
                "from": "web",
                "tid": 123456,
                "fc" : 3,
                "ip": "192.168.0.1",
                "port": "503",
                "slave": 1,
                "addr": 10,
                "len": 4,
                "type": 4,
                "order": 1
            }`
		var r4 MbtcpReadReq
		if err := json.Unmarshal([]byte(input4), &r4); err != nil {
			logf(err)
			return false
		}
		logf(r4)

		input5 :=
			`{
                "from": "web",
                "tid": 123456,
                "fc" : 3,
                "ip": "192.168.0.1",
                "port": "503",
                "slave": 1,
                "addr": 10,
                "len": 4,
                "type": 6,
                "order": 3
            }`
		var r5 MbtcpReadReq
		if err := json.Unmarshal([]byte(input5), &r5); err != nil {
			logf(err)
			return false
		}
		logf(r5)
		return true
	})

	s.Assert("`mbtcp.once.read` response test", func(logf sugar.Log) bool {

		// res1
		res1 := MbtcpReadRes{
			Tid: 12345, Status: "ok",
			Data: []uint16{1, 1, 0, 1, 0, 1},
		}
		writeReqStr1, err := json.Marshal(res1)
		if err != nil {
			return false
		}
		logf(string(writeReqStr1))

		// res2
		res2 := MbtcpReadRes{
			Tid: 12345, Status: "timeout",
		}
		writeReqStr2, err := json.Marshal(res2)
		if err != nil {
			return false
		}
		logf(string(writeReqStr2))

		// res3
		res3 := MbtcpReadRes{
			Tid: 12345, Status: "ok", Type: 1,
			Data: []uint16{255, 1234, 789},
		}
		writeReqStr3, err := json.Marshal(res3)
		if err != nil {
			return false
		}
		logf(string(writeReqStr3))

		// res4
		res4 := MbtcpReadRes{
			Tid: 12345, Status: "ok", Type: 2,
			Bytes: []byte{0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34},
			Data:  []float32{22.34, 33.12, 44.56},
		}
		writeReqStr4, err := json.Marshal(res4)
		if err != nil {
			return false
		}
		logf(string(writeReqStr4))

		// res5
		res5 := MbtcpReadRes{
			Tid: 12345, Status: "ok", Type: 2,
			Bytes: []byte{0XFF, 0X00, 0XAB},
			Data:  []uint16{255, 1234, 789},
		}
		writeReqStr5, err := json.Marshal(res5)
		if err != nil {
			return false
		}
		logf(string(writeReqStr5))

		// res6
		res6 := MbtcpReadRes{
			Tid: 12345, Status: "conversion fail",
			Bytes: []byte{0XFF, 0X00, 0XAB},
		}
		writeReqStr6, err := json.Marshal(res6)
		if err != nil {
			return false
		}
		logf(string(writeReqStr6))

		// res7
		res7 := MbtcpReadRes{
			Tid: 12345, Status: "ok", Type: 2,
			Bytes: []byte{0XFF, 0X00, 0XAB},
			Data:  "112C004F12345678",
		}
		writeReqStr7, err := json.Marshal(res7)
		if err != nil {
			return false
		}
		logf(string(writeReqStr7))

		return true

	})

	s.Title("get/set modbus tcp timeout struct tests")

	s.Assert("`mbtcp.timeout.read` request test", func(logf sugar.Log) bool {
		input :=
			`{
                "from": "web",
                "tid": 123456
            }`
		var r1 MbtcpTimeoutReq
		if err := json.Unmarshal([]byte(input), &r1); err != nil {
			logf(err)
			return false
		}
		return true
	})

	s.Assert("`mbtcp.timeout.read` response test", func(logf sugar.Log) bool {
		res1 := MbtcpTimeoutRes{Tid: 123456, Status: "ok"}
		writeReqStr1, err := json.Marshal(res1)
		if err != nil {
			return false
		}
		logf(string(writeReqStr1))

		res2 := MbtcpTimeoutRes{Tid: 123456, Status: "ok", Data: 210000}
		writeReqStr2, err := json.Marshal(res2)
		if err != nil {
			return false
		}
		logf(string(writeReqStr2))
		return true
	})

	s.Assert("`mbtcp.timeout.update` request test", func(logf sugar.Log) bool {
		input :=
			`{
                "from": "web",
                "tid": 123456,
                "timeout": 210000
            }`
		var r1 MbtcpTimeoutReq
		if err := json.Unmarshal([]byte(input), &r1); err != nil {
			logf(err)
			return false
		}
		return true
	})

	s.Assert("`mbtcp.timeout.update` response test", func(logf sugar.Log) bool {
		res1 := MbtcpTimeoutRes{Tid: 123456, Status: "ok"}
		writeReqStr1, err := json.Marshal(res1)
		if err != nil {
			return false
		}
		logf(string(writeReqStr1))
		return true
	})

	s.Assert("`mbtcp.once.write` request test", func(logf sugar.Log) bool {
		input :=
			`{
				"from": "web",
				"tid": 123456,
				"fc" : 5,
				"ip": "192.168.0.1",
				
				"slave": 1,
				"addr": 10,
				"data": 1
            }`
		var data json.RawMessage
		r1 := MbtcpWriteReq{Data: &data}

		if err := json.Unmarshal([]byte(input), &r1); err != nil {
			logf(err)
			return false
		}
		switch r1.FC {
		case 5:
			var d uint16
			if err := json.Unmarshal(data, &d); err != nil {
				logf(err)
				return false
			}
			r1.Data = d
		case 6:
			//
		}

		logf(r1)

		input2 :=
			`{
				"from": "web",
				"tid": 123456,
				"fc" : 6,
				"ip": "192.168.0.1",
				"port": "503",
				"slave": 1,
				"addr": 10,
				"hex": false,
				"data": "22"
			}`
		var data2 json.RawMessage
		r2 := MbtcpWriteReq{Data: &data2}
		if err := json.Unmarshal([]byte(input2), &r2); err != nil {
			logf(err)
			return false
		}
		switch r2.FC {
		case 5:
			//
		case 6:
			var d string
			if err := json.Unmarshal(data2, &d); err != nil {
				logf(err)
				return false
			}

			if r2.Hex {
				dd, err := HexStringToRegisters(d)
				if err != nil {
					logf(err)
					return false
				}
				logf(len(dd))
				r2.Data = dd[0]
			} else {
				dd, err := strconv.Atoi(d)
				if err != nil {
					logf(err)
					return false
				}
				r2.Data = dd
			}
		}

		logf(r2)

		input3 :=
			`{
				"from": "web",
				"tid": 123456,
				"fc" : 6,
				"ip": "192.168.0.1",
				"port": "503",
				"slave": 1,
				"addr": 10,
				"hex": true,
				"data": "ABCD"
			}`

		var data3 json.RawMessage
		r3 := MbtcpWriteReq{Data: &data3}
		if err := json.Unmarshal([]byte(input3), &r3); err != nil {
			logf(err)
			return false
		}
		switch r3.FC {
		case 5:
			//
		case 6:
			var d string
			if err := json.Unmarshal(data3, &d); err != nil {
				logf(err)
				return false
			}

			if r3.Hex {
				dd, err := HexStringToRegisters(d)
				if err != nil {
					logf(err)
					return false
				}
				logf(len(dd))
				r3.Data = dd[0]
			} else {
				dd, _ := strconv.Atoi(d)
				if err := json.Unmarshal(data3, &d); err != nil {
					logf(err)
					return false
				}
				r3.Data = dd
			}

		}
		logf(r3)

		input4 :=
			`{
				"from": "web",
				"tid": 123456,
				"fc" : 15,
				"ip": "192.168.0.1",
				"port": "503",
				"slave": 1,
				"addr": 10,
				"len": 4,
				"data": [1,0,1,0]
			}`
		var data4 json.RawMessage
		r4 := MbtcpWriteReq{Data: &data4}
		if err := json.Unmarshal([]byte(input4), &r4); err != nil {
			logf(err)
			return false
		}
		switch r4.FC {
		case 5:
		//
		case 6:
		//
		case 15:
			var d []uint16
			if err := json.Unmarshal(data4, &d); err != nil {
				logf(err)
				return false
			}
			r4.Data = d
		}
		logf(r4)

		input5 :=
			`{
				"from": "web",
				"tid": 123456,
				"fc" : 16,
				"ip": "192.168.0.1",
				"port": "503",
				"slave": 1,
				"addr": 10,
				"len": 4,
				"hex": false,
				"data": "11,22,33,44"
			}`

		var data5 json.RawMessage
		r5 := MbtcpWriteReq{Data: &data5}
		if err := json.Unmarshal([]byte(input5), &r5); err != nil {
			logf(err)
			return false
		}
		switch r5.FC {
		case 5:
			//
		case 6:
			//
		case 15:
			//
		case 16:
			var d string
			if err := json.Unmarshal(data5, &d); err != nil {
				logf(err)
				return false
			}
			if r5.Hex {
				dd, err := HexStringToRegisters(d)
				if err != nil {
					logf(err)
					return false
				}
				r5.Data = dd
			} else {
				dd, err := DecimalStringToRegisters(d)
				if err != nil {
					logf(err)
					return false
				}
				r5.Data = dd
			}
		}
		logf(r5)

		input6 :=
			`{
				"from": "web",
				"tid": 123456,
				"fc" : 16,
				"ip": "192.168.0.1",
				"port": "503",
				"slave": 1,
				"addr": 10,
				"len": 4,
				"hex": true,
				"data": "ABCD1234EFAB1234"
			}`
		var data6 json.RawMessage
		r6 := MbtcpWriteReq{Data: &data6}
		if err := json.Unmarshal([]byte(input6), &r6); err != nil {
			logf(err)
			return false
		}
		switch r6.FC {
		case 5:
			//
		case 6:
			//
		case 15:
			//
		case 16:
			var d string
			if err := json.Unmarshal(data6, &d); err != nil {
				logf(err)
				return false
			}
			if r6.Hex {
				dd, err := HexStringToRegisters(d)
				if err != nil {
					logf(err)
					return false
				}
				r6.Data = dd
			} else {
				dd, err := DecimalStringToRegisters(d)
				if err != nil {
					logf(err)
					return false
				}
				r6.Data = dd
			}
		}
		logf(r6)

		return true
	})

}
