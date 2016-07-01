package main_test

import (
	"encoding/json"
	"testing"

	"github.com/marksalpeter/sugar"
	. "github.com/taka-wang/psmb"
)

func TestUpstreamStruct(t *testing.T) {

	s := sugar.New(nil)

	s.Title("One-off modbus tcp struct tests")

	s.Assert("`mbtcp.once.read` request test", func(log sugar.Log) bool {
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
		var r1 MbtcpOnceReadReq
		if err := json.Unmarshal([]byte(input), &r1); err != nil {
			log("json err:", err)
			return false
		}
		log(r1)

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
		var r2 MbtcpOnceReadReq
		if err := json.Unmarshal([]byte(input2), &r2); err != nil {
			log("json err:", err)
			return false
		}
		log(r2)

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
		var r3 MbtcpOnceReadReq
		if err := json.Unmarshal([]byte(input3), &r3); err != nil {
			log("json err:", err)
			return false
		}
		log(r3)

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
		var r4 MbtcpOnceReadReq
		if err := json.Unmarshal([]byte(input4), &r4); err != nil {
			log("json err:", err)
			return false
		}
		log(r4)

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
		var r5 MbtcpOnceReadReq
		if err := json.Unmarshal([]byte(input5), &r5); err != nil {
			log("json err:", err)
			return false
		}
		log(r5)
		return true
	})

	s.Assert("`mbtcp.once.read` response test", func(log sugar.Log) bool {

		// res1
		res1 := MbtcpOnceReadRes{
			Tid: 12345, Status: "ok",
			Data: []uint16{1, 1, 0, 1, 0, 1},
		}
		writeReqStr1, err := json.Marshal(res1)
		if err != nil {
			return false
		}
		log(string(writeReqStr1))

		// res2
		res2 := MbtcpOnceReadRes{
			Tid: 12345, Status: "timeout",
		}
		writeReqStr2, err := json.Marshal(res2)
		if err != nil {
			return false
		}
		log(string(writeReqStr2))

		// res3
		res3 := MbtcpOnceReadRes{
			Tid: 12345, Status: "ok", Type: 1,
			Data: []uint16{255, 1234, 789},
		}
		writeReqStr3, err := json.Marshal(res3)
		if err != nil {
			return false
		}
		log(string(writeReqStr3))

		// res4
		res4 := MbtcpOnceReadRes{
			Tid: 12345, Status: "ok", Type: 2,
			Bytes: []byte{0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34},
			Data:  []float32{22.34, 33.12, 44.56},
		}
		writeReqStr4, err := json.Marshal(res4)
		if err != nil {
			return false
		}
		log(string(writeReqStr4))

		// res5
		res5 := MbtcpOnceReadRes{
			Tid: 12345, Status: "ok", Type: 2,
			Bytes: []byte{0XFF, 0X00, 0XAB},
			Data:  []uint16{255, 1234, 789},
		}
		writeReqStr5, err := json.Marshal(res5)
		if err != nil {
			return false
		}
		log(string(writeReqStr5))

		// res6
		res6 := MbtcpOnceReadRes{
			Tid: 12345, Status: "conversion fail",
			Bytes: []byte{0XFF, 0X00, 0XAB},
		}
		writeReqStr6, err := json.Marshal(res6)
		if err != nil {
			return false
		}
		log(string(writeReqStr6))

		// res7
		res7 := MbtcpOnceReadRes{
			Tid: 12345, Status: "ok", Type: 2,
			Bytes: []byte{0XFF, 0X00, 0XAB},
			Data:  "112C004F12345678",
		}
		writeReqStr7, err := json.Marshal(res7)
		if err != nil {
			return false
		}
		log(string(writeReqStr7))

		return true

	})

	s.Title("get/set modbus tcp timeout struct tests")

	s.Assert("`mbtcp.timeout.read` request test", func(log sugar.Log) bool {
		input :=
			`{
                "from": "web",
                "tid": 123456
            }`
		var r1 MbtcpTimeoutReq
		if err := json.Unmarshal([]byte(input), &r1); err != nil {
			log("json err:", err)
			return false
		}
		return true
	})

	s.Assert("`mbtcp.timeout.read` response test", func(log sugar.Log) bool {
		res1 := MbtcpTimeoutRes{Tid: 123456, Status: "ok"}
		writeReqStr1, err := json.Marshal(res1)
		if err != nil {
			return false
		}
		log(string(writeReqStr1))

		res2 := MbtcpTimeoutRes{Tid: 123456, Status: "ok", Data: 210000}
		writeReqStr2, err := json.Marshal(res2)
		if err != nil {
			return false
		}
		log(string(writeReqStr2))
		return true
	})

	s.Assert("`mbtcp.timeout.update` request test", func(log sugar.Log) bool {
		input :=
			`{
                "from": "web",
                "tid": 123456,
                "timeout": 210000
            }`
		var r1 MbtcpTimeoutReq
		if err := json.Unmarshal([]byte(input), &r1); err != nil {
			log("json err:", err)
			return false
		}
		return true
	})

	s.Assert("`mbtcp.timeout.update` response test", func(log sugar.Log) bool {
		res1 := MbtcpTimeoutRes{Tid: 123456, Status: "ok"}
		writeReqStr1, err := json.Marshal(res1)
		if err != nil {
			return false
		}
		log(string(writeReqStr1))
		return true
	})
}

func TestDownstreamStruct(t *testing.T) {

	s := sugar.New(nil)

	s.Title("modbus tcp downstreamstruct tests")

	s.Assert("`read` request test", func(log sugar.Log) bool {
		req := DMbtcpReadReq{
			Tid:   "123456",
			Cmd:   1,
			IP:    "192.168.3.2",
			Port:  "502",
			Slave: 22,
			Addr:  250,
			Len:   10,
		}
		reqStr, err := json.Marshal(req)
		if err != nil {
			return false
		}
		log(string(reqStr))

		return true
	})

	s.Assert("`single read` response test", func(log sugar.Log) bool {
		input :=
			`{
                "tid": "1",
                "data": [1],
                "status": "ok"
            }`
		var r1 DMbtcpRes
		if err := json.Unmarshal([]byte(input), &r1); err != nil {
			log("json err:", err)
			return false
		}
		log(r1)
		return true
	})

	s.Assert("`multiple read` response test", func(log sugar.Log) bool {
		input :=
			`{
                "tid": "1",
                "data": [1,2,3,4],
                "status": "ok"
            }`
		var r1 DMbtcpRes
		if err := json.Unmarshal([]byte(input), &r1); err != nil {
			log("json err:", err)
			return false
		}
		log(r1)
		return true
	})

	s.Assert("`single write` request test", func(log sugar.Log) bool {
		req := DMbtcpWriteReq{
			Tid:   "123456",
			Cmd:   6,
			IP:    "192.168.3.2",
			Port:  "502",
			Slave: 22,
			Addr:  250,
			Data:  1234,
		}
		reqStr, err := json.Marshal(req)
		if err != nil {
			return false
		}
		log(string(reqStr))

		return true
	})

	s.Assert("`multiple write` request test", func(log sugar.Log) bool {
		req := DMbtcpWriteReq{
			Tid:   "123456",
			Cmd:   6,
			IP:    "192.168.3.2",
			Port:  "502",
			Slave: 22,
			Addr:  250,
			Len:   4,
			Data:  []uint16{1, 2, 3, 4},
		}
		reqStr, err := json.Marshal(req)
		if err != nil {
			return false
		}
		log(string(reqStr))

		return true
	})

	s.Assert("`set timeout` request test", func(log sugar.Log) bool {
		req := DMbtcpTimeout{
			Tid:     "22222",
			Cmd:     50,
			Timeout: 210000,
		}
		reqStr, err := json.Marshal(req)
		if err != nil {
			return false
		}
		log(string(reqStr))

		return true
	})

	s.Assert("`get timeout` response test", func(log sugar.Log) bool {
		input :=
			`{
                "tid": "22222",
                "status": "ok",
                "timeout": 210000
            }`
		var r1 DMbtcpTimeout
		if err := json.Unmarshal([]byte(input), &r1); err != nil {
			log("json err:", err)
			return false
		}
		log(r1)
		return true
	})

}
