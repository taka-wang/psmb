package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	zmq "github.com/taka-wang/zmq3"
)

func main() {
	go subscriber()
	sender, _ := zmq.NewSocket(zmq.PUB)

	// s.Every(1).Seconds().Do(publisher)
	//s := gocron.NewScheduler()
	//s.Start()

	// upstream subscriber
	receiver, _ := zmq.NewSocket(zmq.SUB)
	defer receiver.Close()
	receiver.Bind("ipc:///tmp/to.psmb")

	filter := ""
	receiver.SetSubscribe(filter) // filter frame 1
	for {
		msg, _ := receiver.RecvMessage(0)
		fmt.Println(msg[0]) // frame 1: method
		fmt.Println(msg[1]) // frame 2: command

		switch msg[0] {
		case "mbtcp.once.read":
			fmt.Println("got")
			var req = MbtcpOnceReadReq
			if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
				fmt.Println("json err:", err)
			}
			sender.Connect("ipc:///tmp/to.modbus")
			command := DMbtcpReadReq{
				Tid:   string(req.Tid),
				Cmd:   req.FC,
				IP:    req.IP,
				Port:  req.Port,
				Slave: req.Slave,
				Addr:  req.Addr,
				Len:   req.Len,
			}

			cmd, err := json.Marshal(command) // marshal to json string
			if err != nil {
				fmt.Println("json err:", err)
			}

			sender.Send("tcp", zmq.SNDMORE) // frame 1
			sender.Send(string(cmd), 0)     // convert to string; frame 2
			sender.Close()

		default:
			fmt.Println("unsupport")
		}

		t := time.Now()
		fmt.Println("zrecv" + t.Format("2006-01-02 15:04:05.000"))
		time.Sleep(300 * time.Millisecond)
	}
}

func subscriber() {
	sender, _ := zmq.NewSocket(zmq.PUB) // to upstream

	receiver, _ := zmq.NewSocket(zmq.SUB)
	defer receiver.Close()
	receiver.Connect("ipc:///tmp/from.modbus")

	filter := ""
	receiver.SetSubscribe(filter) // filter frame 1
	for {
		msg, _ := receiver.RecvMessage(0)
		fmt.Println(msg[0]) // frame 1: method
		fmt.Println(msg[1]) // frame 2: command
		switch msg[0] {
		case "mbtcp.once.read":

			var res = DMbtcpRes
			if err := json.Unmarshal([]byte(msg[1]), &res); err != nil {
				fmt.Println("json err:", err)
			}
			sender.Bind("ipc:///tmp/from.psmb")

			tid, err := strconv.ParseInt(res.Tid, 10, 64)
			command := MbtcpOnceReadUInt16Res{
				MbtcpOnceReadRes{
					Tid:    tid,
					Status: res.Status,
				},
				Data: res.Data,
			}
			cmdStr, _ := json.Marshal(command)
			sender.Send(msg[0], zmq.SNDMORE) // frame 1
			sender.Send(string(cmdStr), 0)   // convert to string; frame 2
			sender.Close()

		default:
			fmt.Println("unsupport")
		}

		t := time.Now()
		fmt.Println("zrecv" + t.Format("2006-01-02 15:04:05.000"))
	}
}

/*
func publisher() {
	t := time.Now()
	fmt.Println("zsend" + t.Format("2006-01-02 15:04:05.000"))
	command := DMbtcpReadReq{
		Tid:   "12345678910",
		Cmd:   1,
		IP:    "172.17.0.2",
		Port:  "502",
		Slave: 1,
		Addr:  10,
		Len:   8,
	}

	cmd, err := json.Marshal(command) // marshal to json string
	if err != nil {
		fmt.Println("json err:", err)
	}

	sender.Send("tcp", zmq.SNDMORE) // frame 1
	sender.Send(string(cmd), 0)     // convert to string; frame 2
}
*/
