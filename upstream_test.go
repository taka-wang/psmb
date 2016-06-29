package psmb

import (
	"encoding/json"
	"testing"

	"github.com/marksalpeter/sugar"
)

func TestOneOffStruct(t *testing.T) {

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

		res1 := MbtcpOnceReadUInt16Res{MbtcpOnceReadRes{Tid: 12345, Status: "ok"}, []uint16{1, 1, 0, 1, 0, 1}}
		writeReqStr1, err := json.Marshal(res1)
		if err != nil {
			return false
		}
		log(string(writeReqStr1))

		res2 := MbtcpOnceReadRes{Tid: 12345, Status: "timeout"}
		writeReqStr2, err := json.Marshal(res2)
		if err != nil {
			return false
		}
		log(string(writeReqStr2))

		res3 := MbtcpOnceReadUInt16Res{MbtcpOnceReadRes{Tid: 12345, Status: "ok", Type: 1}, []uint16{255, 1234, 789}}
		writeReqStr3, err := json.Marshal(res3)
		if err != nil {
			return false
		}
		log(string(writeReqStr3))

		res4 := MbtcpOnceReadFloat32Res{MbtcpOnceReadRes{Tid: 12345, Status: "ok", Type: 2, Bytes: []byte{0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34}}, []float32{22.34, 33.12, 44.56}}
		writeReqStr4, err := json.Marshal(res4)
		if err != nil {
			return false
		}
		log(string(writeReqStr4))

		res5 := MbtcpOnceReadUInt16Res{MbtcpOnceReadRes{Tid: 12345, Status: "ok", Type: 2, Bytes: []byte{0XFF, 0X00, 0XAB}}, []uint16{255, 1234, 789}}
		writeReqStr5, err := json.Marshal(res5)
		if err != nil {
			return false
		}
		log(string(writeReqStr5))

		res6 := MbtcpOnceReadRes{Tid: 12345, Status: "conversion fail", Bytes: []byte{0XFF, 0X00, 0XAB}}
		writeReqStr6, err := json.Marshal(res6)
		if err != nil {
			return false
		}
		log(string(writeReqStr6))

		res7 := MbtcpOnceReadStringRes{MbtcpOnceReadRes{Tid: 12345, Status: "ok", Type: 2, Bytes: []byte{0XFF, 0X00, 0XAB}}, "112C004F12345678"}
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
