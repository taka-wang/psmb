package main

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	psmb "github.com/taka-wang/psmb"
	"github.com/takawang/sugar"
	zmq "github.com/takawang/zmq3"
)

// generic tcp publisher
func publisher(cmd, json string) {
	sender, _ := zmq.NewSocket(zmq.PUB)
	defer sender.Close()
	sender.Connect("ipc:///tmp/to.psmb")

	for {
		time.Sleep(time.Duration(10) * time.Millisecond)
		t := time.Now()
		fmt.Println("Req:", t.Format("2006-01-02 15:04:05.000"))
		sender.Send(cmd, zmq.SNDMORE) // frame 1
		sender.Send(json, 0)          // convert to string; frame 2
		// send the exit loop
		break
	}
}

// generic subscribe
func subscriber() (string, string) {
	receiver, _ := zmq.NewSocket(zmq.SUB)
	defer receiver.Close()
	receiver.Connect("ipc:///tmp/from.psmb")
	filter := ""
	receiver.SetSubscribe(filter) // filter frame 1
	for {
		fmt.Println("listen..")
		msg, _ := receiver.RecvMessage(0)

		t := time.Now()
		fmt.Println("Res:", t.Format("2006-01-02 15:04:05.000"))

		// recv then exit loop
		return msg[0], msg[1]
	}
}

func TestPSMB(t *testing.T) {
	/*
		// start psmb service
		go func() {
			psmb.Start()
		}()
	*/

	time.Sleep(2000 * time.Millisecond)

	s := sugar.New(nil)

	var hostName string
	portNum1 := "502"
	portNum2 := "503"

	// generalize host reslove for docker/local env
	host, err := net.LookupHost("mbd")
	if err != nil {
		fmt.Println("local run")
		hostName = "127.0.0.1"
	} else {
		fmt.Println("docker run")
		hostName = host[0] //docker
	}

	// ---------------------------------------------------------------//
	s.Title("One-Off timeout request tests")

	s.Assert("mbtcp.timeout.update test - invalid value (1)", func(log sugar.Log) bool {
		ReadReq := psmb.MbtcpTimeoutReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
			Data: 1,
		}

		ReadReqStr, _ := json.Marshal(ReadReq)
		cmd := "mbtcp.timeout.update"
		go publisher(cmd, string(ReadReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(ReadReqStr))
		log("set timeout as 200000")
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpTimeoutRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("mbtcp.timeout.read test - invalid value", func(log sugar.Log) bool {
		ReadReq := psmb.MbtcpTimeoutReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
		}

		ReadReqStr, _ := json.Marshal(ReadReq)
		cmd := "mbtcp.timeout.read"
		go publisher(cmd, string(ReadReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(ReadReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpTimeoutRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" || r2.Data != 200000 {
			return false
		}
		return true
	})

	s.Assert("mbtcp.timeout.update test - valid value (212345)", func(log sugar.Log) bool {
		ReadReq := psmb.MbtcpTimeoutReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
			Data: 212345,
		}

		ReadReqStr, _ := json.Marshal(ReadReq)
		cmd := "mbtcp.timeout.update"
		go publisher(cmd, string(ReadReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(ReadReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpTimeoutRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("mbtcp.timeout.read test - valid value", func(log sugar.Log) bool {
		ReadReq := psmb.MbtcpTimeoutReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
		}

		ReadReqStr, _ := json.Marshal(ReadReq)
		cmd := "mbtcp.timeout.read"
		go publisher(cmd, string(ReadReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(ReadReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpTimeoutRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" || r2.Data != 212345 {
			return false
		}
		return true
	})

	// ---------------------------------------------------------------//
	s.Title("`mbtcp.once.write` tests")

	s.Assert("`FC5` write bit test: port 502", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    5,
			Slave: 1,
			Addr:  10,
			Data:  1,
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}
		//return true

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    1,
			Slave: 1,
			Addr:  10,
			Len:   3,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		if r2.Data.(int) != 1 {
			return false
		}

		return true
	})

	s.Assert("`FC6` write `DEC` register test: port 502", func(log sugar.Log) bool {
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    6,
			Slave: 1,
			Addr:  10,
			Hex:   false,
			Data:  "22",
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC6` write `HEX` register test: port 502", func(log sugar.Log) bool {
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    6,
			Slave: 1,
			Addr:  10,
			Hex:   true,
			Data:  "ABCD",
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC15` write multiple bits test: port 502", func(log sugar.Log) bool {
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    15,
			Slave: 1,
			Addr:  10,
			Len:   4,
			Data:  []uint16{1, 0, 1, 0},
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC16` write `DEC` registers test: port 502", func(log sugar.Log) bool {
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    16,
			Slave: 1,
			Addr:  10,
			Len:   4,
			Hex:   false,
			Data:  "11,22,33,44",
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC16` write `HEX` registers test: port 502", func(log sugar.Log) bool {
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    16,
			Slave: 1,
			Addr:  10,
			Len:   4,
			Hex:   true,
			Data:  "ABCD1234EFAB1234",
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	// ---------------------------------------------------------------//
	s.Title("`mbtcp.once.read` tests")

	s.Assert("`FC1` read bits test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    1,
			Slave: 1,
			Addr:  3,
			Len:   7,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC1` read bits test: port 503", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum2,
			FC:    1,
			Slave: 1,
			Addr:  3,
			Len:   7,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 1 test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  psmb.RegisterArray,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 2 test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  psmb.HexString,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 3 length 4 test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   4,
			Type:  psmb.Scale,
			Range: &psmb.ScaleRange{0, 65535, 100, 500},
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 3 length 7 test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  psmb.Scale,
			Range: &psmb.ScaleRange{0, 65535, 100, 500},
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "Conversion failed" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 4 length 4 test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   4,
			Type:  psmb.UInt16,
			Order: psmb.AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 4 length 7 test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  psmb.UInt16,
			Order: psmb.AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 6 length 8 test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.UInt32,
			Order: psmb.BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 6 length 7 test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  psmb.UInt32,
			Order: psmb.BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 8 length 8 test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.Float32,
			Order: psmb.AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 8 length 7 test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  psmb.Float32,
			Order: psmb.AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Title("Poll request tests")

	s.Assert("mbtcp.poll.create `FC1` read bits test: port 503", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpPollStatus{
			From:     "web",
			Tid:      time.Now().UTC().UnixNano(),
			Name:     "LED_11",
			Interval: 1,
			Enabled:  true,
			IP:       hostName,
			Port:     portNum1,
			FC:       1,
			Slave:    1,
			Addr:     3,
			Len:      7,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.poll.create"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("mbtcp.poll.create `FC3` read bytes Type 1 test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpPollStatus{
			From:     "web",
			Tid:      time.Now().UTC().UnixNano(),
			Name:     "LED_1",
			Interval: 1,
			Enabled:  true,
			IP:       hostName,
			Port:     portNum1,
			FC:       3,
			Slave:    1,
			Addr:     3,
			Len:      7,
			Type:     psmb.RegisterArray,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.poll.create"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

}
