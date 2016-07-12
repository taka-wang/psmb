package main

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

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

func TestPsmb(t *testing.T) {
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

	s.Title("Oneoff request tests")

	s.Assert("`FC5` write bit test: port 502", func(log sugar.Log) bool {
		writeReq := MbtcpWriteReq{
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
		var r2 MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC1` read bits test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := MbtcpReadReq{
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
		var r2 MbtcpReadRes
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
		readReq := MbtcpReadReq{
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
		var r2 MbtcpReadRes
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
		readReq := MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  RegisterArray,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 MbtcpReadRes
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
		readReq := MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  HexString,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 MbtcpReadRes
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
		readReq := MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   4,
			Type:  Scale,
			Range: &ScaleRange{0, 65535, 100, 500},
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 MbtcpReadRes
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
		readReq := MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  Scale,
			Range: &ScaleRange{0, 65535, 100, 500},
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 MbtcpReadRes
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
		readReq := MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   4,
			Type:  UInt16,
			Order: AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 MbtcpReadRes
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
		readReq := MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  UInt16,
			Order: AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 MbtcpReadRes
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
		readReq := MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  UInt32,
			Order: BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 MbtcpReadRes
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
		readReq := MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  UInt32,
			Order: BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 MbtcpReadRes
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
		readReq := MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  Float32,
			Order: AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 MbtcpReadRes
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
		readReq := MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  Float32,
			Order: AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 MbtcpReadRes
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
