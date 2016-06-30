package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/taka-wang/gocron"
	psmb "github.com/taka-wang/psmb/psmb"
	zmq "github.com/taka-wang/zmq3"
)

var sender *zmq.Socket

func main() {
	go subscriber()
	sender, _ = zmq.NewSocket(zmq.PUB)
	//defer sender.Close()
	sender.Connect("ipc:///tmp/to.modbus")
	s := gocron.NewScheduler()
	s.Every(1).Seconds().Do(publisher)
	s.Start()
	for {
		time.Sleep(300 * time.Millisecond)
	}
}

func publisher() {
	t := time.Now()
	fmt.Println("zsend" + t.Format("2006-01-02 15:04:05.000"))
	command := psmb.DMbtcpReadReq{
		"172.17.0.2",
		"502",
		1,
		12,
		"fc1",
		10,
		10,
	}

	cmd, err := json.Marshal(command) // marshal to json string
	if err != nil {
		fmt.Println("json err:", err)
	}

	sender.Send("tcp", zmq.SNDMORE) // frame 1
	sender.Send(string(cmd), 0)     // convert to string; frame 2
}

func subscriber() {
	receiver, _ := zmq.NewSocket(zmq.SUB)
	defer receiver.Close()
	receiver.Connect("ipc:///tmp/from.modbus")

	filter := ""
	receiver.SetSubscribe(filter) // filter frame 1
	for {
		msg, _ := receiver.RecvMessage(0)
		fmt.Println(msg[0]) // frame 1: method
		fmt.Println(msg[1]) // frame 2: command
		t := time.Now()
		fmt.Println("zrecv" + t.Format("2006-01-02 15:04:05.000"))
	}
}
