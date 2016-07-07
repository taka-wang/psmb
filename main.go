package main

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/taka-wang/gocron"
	zmq "github.com/taka-wang/zmq3"
)

var taskMap map[string]interface{}
var taskMap2 map[string]string
var sch gocron.Scheduler

// RequestParser handle message from services
func RequestParser(socket *zmq.Socket, msg []string) (interface{}, error) {
	// Check the length of multi-part message
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

		// add to task map
		taskMap[cmd.Tid] = cmd
		taskMap2[cmd.Tid] = msg[0]
		sch.Emergency().Do(Task, socket, cmd)
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
func ResponseParser(socket *zmq.Socket, msg []string) {
	// Check the length of multi-part message
	if len(msg) != 2 {
		log.Println("Request parser failed: invalid message length")
		return
		//return "", errors.New("Invalid message length")
	}

	log.Println("Parsing response:", msg[0])

	// Convert zframe 1: command number
	cmdType, _ := strconv.Atoi(msg[0])
	var cmdStr []byte
	var TidStr string

	switch cmdType {
	case 50, 51:
		var res DMbtcpTimeout
		if err := json.Unmarshal([]byte(msg[1]), &res); err != nil {
			log.Println("json err:", err)
			return
		}
		log.Println(res)
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
			log.Println("json err:", err)
			return
		}
		log.Println(res)
		TidStr = res.Tid
		tid, _ := strconv.ParseInt(res.Tid, 10, 64)
		command := MbtcpOnceReadRes{
			Tid:    tid,
			Status: res.Status,
			Data:   res.Data,
		}
		cmdStr, _ = json.Marshal(command)
	}

	log.Println("Conversion for upstream complete")
	log.Println(string(cmdStr))

	// todo: check msg[0], should be web
	if frame, ok := taskMap2[TidStr]; ok {
		delete(taskMap2, TidStr)
		socket.Send(frame, zmq.SNDMORE) // frame 1
		socket.Send(string(cmdStr), 0)  // convert to string; frame 2
	} else {
		socket.Send("null", zmq.SNDMORE) // frame 1
		socket.Send(string(cmdStr), 0)   // convert to string; frame 2
	}

	t := time.Now()
	log.Println("zrecv:" + t.Format("2006-01-02 15:04:05.000"))
}

// RequestCommandBuilder build command to modbusd
func RequestCommandBuilder() {

}

// ResponseCommandBuilder build command to services
// Todo: filter, handle
func ResponseCommandBuilder() {

}

// Task for gocron
func Task(socket *zmq.Socket, m interface{}) {
	str, err := json.Marshal(m) // marshal to json string
	if err != nil {
		log.Println("Marshal request failed:", err)
		return
	}
	socket.Send("tcp", zmq.SNDMORE) // frame 1
	socket.Send(string(str), 0)     // convert to string; frame 2
}

func main() {

	taskMap = make(map[string]interface{})
	taskMap2 = make(map[string]string)

	sch = gocron.NewScheduler()
	sch.Start()
	// s.Every(1).Seconds().Do(publisher)

	// upstream subscriber
	fromService, _ := zmq.NewSocket(zmq.SUB)
	defer fromService.Close()
	fromService.Bind("ipc:///tmp/to.psmb")
	fromService.SetSubscribe("")

	// upstream publisher
	toService, _ := zmq.NewSocket(zmq.PUB) // to upstream
	defer toService.Close()
	toService.Bind("ipc:///tmp/from.psmb")

	// downstream subscriber
	fromModbusd, _ := zmq.NewSocket(zmq.SUB)
	defer fromModbusd.Close()
	fromModbusd.Connect("ipc:///tmp/from.modbus")
	fromModbusd.SetSubscribe("")

	// downstream publisher
	toModbusd, _ := zmq.NewSocket(zmq.PUB)
	defer toModbusd.Close()
	toModbusd.Connect("ipc:///tmp/to.modbus")

	// initialize poll set
	poller := zmq.NewPoller()
	poller.Add(fromService, zmq.POLLIN)
	poller.Add(fromModbusd, zmq.POLLIN)

	// process messages from both sockets
	for {
		sockets, _ := poller.Poll(-1)
		for _, socket := range sockets {
			switch s := socket.Socket; s {
			case fromService:
				msg, _ := fromService.RecvMessage(0)
				log.Println("receive from upstream", msg[0], msg[1])
				RequestParser(toModbusd, msg)
			case fromModbusd:
				msg, _ := fromModbusd.RecvMessage(0)
				log.Println("receive from modbusd", msg[0], msg[1])
				ResponseParser(toService, msg)
			}
		}
	}
}

//t := time.Now()
//fmt.Println("zrecv" + t.Format("2006-01-02 15:04:05.000"))
