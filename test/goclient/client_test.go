package main

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/marksalpeter/sugar"
	zmq "github.com/taka-wang/zmq3"
)

// generic tcp publisher
func publisher(cmd, json string) {
	t := time.Now()
	fmt.Println("Req:", t.Format("2006-01-02 15:04:05.000"))
	sender, _ := zmq.NewSocket(zmq.PUB)
	defer sender.Close()
	sender.Connect("ipc:///tmp/to.psmb")

	for {
		time.Sleep(time.Duration(1) * time.Second)
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

	s.Title("psmb test")

	s.Assert("`FC1` read bits test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := MbtcpOnceReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    1,
			Slave: 1,
			Addr:  10,
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
		var r2 MbtcpOnceReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check reponse
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC1` read bits test: port 503", func(log sugar.Log) bool {
		// send request
		readReq := MbtcpOnceReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum2,
			FC:    1,
			Slave: 1,
			Addr:  10,
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
		var r2 MbtcpOnceReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check reponse
		if r2.Status != "ok" {
			return false
		}
		return true
	})
}
