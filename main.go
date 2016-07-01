package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	zmq "github.com/taka-wang/zmq3"
)

func subscriber() {

	toWeb, _ := zmq.NewSocket(zmq.PUB) // to upstream
	toWeb.Bind("ipc:///tmp/from.psmb")

	fromModbusd, _ := zmq.NewSocket(zmq.SUB)
	defer fromModbusd.Close()
	fromModbusd.Connect("ipc:///tmp/from.modbus")
	filter := ""
	fromModbusd.SetSubscribe(filter)

	for {
		msg, _ := fromModbusd.RecvMessage()
		fmt.Println("recv from modbusd", msg[0], msg[1])

		var res DMbtcpRes
		if err := json.Unmarshal([]byte(msg[1]), &res); err != nil {
			fmt.Println("json err:", err)
		}
		fmt.Println(res)

		tid, _ := strconv.ParseInt(res.Tid, 10, 64)
		command := MbtcpOnceReadUInt16Res{
			ReadRes: MbtcpOnceReadRes{
				tid,
				res.Status,
			},
			Data: res.Data,
		}
		cmdStr, _ := json.Marshal(command)

		fmt.Println("convert to upstream complete")
		fmt.Println(string(cmdStr))

		// todo: check msg[0]
		// should be web
		toWeb.Send(msg[0], zmq.SNDMORE) // frame 1
		toWeb.Send(string(cmdStr), 0)   // convert to string; frame 2
		//toWeb.Close()

		/*
			switch msg[0] {
			case "mbtcp.once.read":

			default:
				fmt.Println("unsupport")
			}
		*/
		t := time.Now()
		fmt.Println("zrecv:" + t.Format("2006-01-02 15:04:05.000"))
	}
}

func main() {

	// s.Every(1).Seconds().Do(publisher)
	//s := gocron.NewScheduler()
	//s.Start()

	go subscriber()

	// downstream pub
	toModbusd, _ := zmq.NewSocket(zmq.PUB)
	toModbusd.Connect("ipc:///tmp/to.modbus")

	// upstream subscriber
	fromWeb, _ := zmq.NewSocket(zmq.SUB)
	defer fromWeb.Close()
	fromWeb.Bind("ipc:///tmp/to.psmb")
	filter := ""
	fromWeb.SetSubscribe(filter) // filter frame 1

	for {
		msg, _ := fromWeb.RecvMessage(0)
		fmt.Println("recv from web", msg[0], msg[1])

		switch msg[0] {
		case "mbtcp.once.read":
			fmt.Println("in: ", msg[0])
			var req MbtcpOnceReadReq
			if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
				fmt.Println("json err:", err)
				break
			}

			cmd := DMbtcpReadReq{
				Tid:   strconv.FormatInt(req.Tid, 10),
				Cmd:   req.FC,
				IP:    req.IP,
				Port:  req.Port,
				Slave: req.Slave,
				Addr:  req.Addr,
				Len:   req.Len,
			}

			cmdStr, err := json.Marshal(cmd) // marshal to json string
			if err != nil {
				fmt.Println("json err:", err)
			}

			fmt.Println("convert to downstream complete")
			fmt.Println(string(cmdStr))

			// send to modbusd
			toModbusd.Send("tcp", zmq.SNDMORE) // frame 1
			toModbusd.Send(string(cmdStr), 0)  // convert to string; frame 2
			toModbusd.Close()

		default:
			fmt.Println("unsupport")
		}

		//t := time.Now()
		//fmt.Println("zrecv" + t.Format("2006-01-02 15:04:05.000"))
		time.Sleep(300 * time.Millisecond)
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
