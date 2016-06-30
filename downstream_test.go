package psmb

import (
	"encoding/json"
	"testing"

	"github.com/marksalpeter/sugar"
)

func TestDownstream(t *testing.T) {

	s := sugar.New(nil)

	s.Title("modbus tcp downstreamstruct tests")

	s.Assert("`read` request test", func(log sugar.Log) bool {
		req := DMbtcpReadReq{
			Tid:   123456,
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
                "tid": 1,
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
                "tid": 1,
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
		req := DMbtcpSingleWriteReq{
			Tid:   123456,
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
		req := DMbtcpMultipleWriteReq{
			Tid:   123456,
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
			Tid:     22222,
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
                "tid": 22222,
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
