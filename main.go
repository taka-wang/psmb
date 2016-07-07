package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/taka-wang/gocron"
	zmq "github.com/taka-wang/zmq3"
)

var taskMap map[string]interface{}
var sch gocron.Scheduler

// RequestParser handle message from services
func RequestParser(socket *zmq.Socket, msg []string) (interface{}, error) {
	// check the length of multi-parts message
	if len(msg) != 2 {
		log.Println("Request parser failed: invalid message length")
		return "", errors.New("Invalid message length")
	}

	log.Println("Parsing request:", msg[0])

	switch msg[0] {
	case "mbtcp.once.read":

		var req MbtcpOnceReadReq
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			log.Println("Unmarshal request failed:", err)
			return "", err
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

		sch.Emergency().Do(Task, socket, cmd)
		cmd["Request"] = msg[0]
		taskMap[cmd.Tid] = cmd // add to map
		//sch.Every(1).Seconds().Do(modbusTask, socket, cmd)
		return cmd, nil

		/*

			cmdStr, err := json.Marshal(cmd) // marshal to json string
			if err != nil {
				log.Println("Marshal request failed:", err)
				return "", err
			}
			return string(cmdStr), nil
		*/

	case "mbtcp.once.write":
		log.Println("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.timeout.read":
		log.Println("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.timeout.update":
		log.Println("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.poll.create":
		log.Println("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.poll.update":
		log.Println("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.poll.read":
		log.Println("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.poll.delete":
		log.Println("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.poll.toggle":
		log.Println("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.polls.read":
		log.Println("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.polls.delete":
		log.Println("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.polls.toggle":
		log.Println("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.polls.import":
		log.Println("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.polls.export":
		log.Println("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.poll.history":
		log.Println("TODO")
		return nil, errors.New("TODO")
	default:
		log.Println("unsupport")
		return nil, errors.New("unsupport request")
	}
}

// ResponseParser handle message from modbusd
func ResponseParser(msg []string) {

}

// RequestCommandBuilder build command to modbusd
func RequestCommandBuilder() {

}

// ResponseCommandBuilder build command to services
// Todo: filter, handle
func ResponseCommandBuilder() {

}

func Task(socket *zmq.Socket, m interface{}) {
	str, err := json.Marshal(m) // marshal to json string
	if err != nil {
		log.Println("Marshal request failed:", err)
		return
	}
	socket.Send("tcp", zmq.SNDMORE) // frame 1
	socket.Send(string(str), 0)     // convert to string; frame 2
}

func subscriber() {

	toWeb, _ := zmq.NewSocket(zmq.PUB) // to upstream
	toWeb.Bind("ipc:///tmp/from.psmb")

	fromModbusd, _ := zmq.NewSocket(zmq.SUB)
	defer fromModbusd.Close()
	fromModbusd.Connect("ipc:///tmp/from.modbus")
	filter := ""
	fromModbusd.SetSubscribe(filter)

	for {
		fmt.Println("listen from modbusd")
		msg, _ := fromModbusd.RecvMessage(0)
		fmt.Println("recv from modbusd", msg[0], msg[1])

		// convert zframe 1: command number
		cmdType, _ := strconv.Atoi(msg[0])
		var cmdStr []byte
		var TidStr string

		switch cmdType {
		case 50, 51:
			var res DMbtcpTimeout
			if err := json.Unmarshal([]byte(msg[1]), &res); err != nil {
				fmt.Println("json err:", err)
			}
			fmt.Println(res)
			TidStr = res.Tid
			tid, _ := strconv.ParseInt(res.Tid, 10, 64)
			command := MbtcpTimeoutRes{
				Tid:    tid,
				Status: res.Status,
				Data:   res.Timeout,
			}
			cmdStr, _ = json.Marshal(command)
		default:
			var res DMbtcpRes
			if err := json.Unmarshal([]byte(msg[1]), &res); err != nil {
				fmt.Println("json err:", err)
			}
			fmt.Println(res)
			TidStr = res.Tid
			tid, _ := strconv.ParseInt(res.Tid, 10, 64)
			command := MbtcpOnceReadRes{
				Tid:    tid,
				Status: res.Status,
				Data:   res.Data,
			}
			cmdStr, _ = json.Marshal(command)
		}

		fmt.Println("convert to upstream complete")
		fmt.Println(string(cmdStr))

		// todo: check msg[0]
		// should be web
		if task, ok := taskMap[TidStr]; ok {
			//delete(taskMap2, TidStr)
			toWeb.Send(task.Request, zmq.SNDMORE) // frame 1
			toWeb.Send(string(cmdStr), 0)         // convert to string; frame 2
		} else {
			toWeb.Send("null", zmq.SNDMORE) // frame 1
			toWeb.Send(string(cmdStr), 0)   // convert to string; frame 2
		}

		t := time.Now()
		fmt.Println("zrecv:" + t.Format("2006-01-02 15:04:05.000"))
	}
}

func main() {

	taskMap = make(map[string]interface{})

	sch = gocron.NewScheduler()
	sch.Start()
	// s.Every(1).Seconds().Do(publisher)

	go subscriber()

	// downstream publish
	toModbusd, _ := zmq.NewSocket(zmq.PUB)
	toModbusd.Connect("ipc:///tmp/to.modbus")
	defer toModbusd.Close()

	// upstream subscriber
	fromUpstream, _ := zmq.NewSocket(zmq.SUB)
	defer fromUpstream.Close()
	fromUpstream.Bind("ipc:///tmp/to.psmb")
	fromUpstream.SetSubscribe("") // filter frame 1

	for {
		// receive multi-parts message
		msg, _ := fromUpstream.RecvMessage(0)
		log.Println("receive from upstream", msg[0], msg[1])
		RequestParser(toModbusd, msg)
		time.Sleep(100 * time.Millisecond)
	}
}

//t := time.Now()
//fmt.Println("zrecv" + t.Format("2006-01-02 15:04:05.000"))
