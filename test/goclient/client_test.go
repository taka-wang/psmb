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

var hostName, portNum string

// generic tcp publisher
func publisher(cmd, json string) {

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
		msg, _ := receiver.RecvMessage(0)
		// recv then exit loop
		return msg[0], msg[1]
	}
}

func TestPsmb(t *testing.T) {
	s := sugar.New(nil)

	portNum = "502"

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

	s.Assert("`FC1` test", func(log sugar.Log) bool {
		readReq := MbtcpOnceReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum,
			FC:    1,
			Slave: 1,
			Addr:  10,
			Len:   8,
		}

		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		go publisher("mbtcp.once.read", string(readReqStr))
		_, s2 := subscriber()
		log("req: %s", string(readReqStr))
		log("res: %s", s2)

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
