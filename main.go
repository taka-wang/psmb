package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/taka-wang/gocron"
	log "github.com/takawang/logrus"
	zmq "github.com/takawang/zmq3"
)

var taskMap map[string]interface{}
var taskMap2 map[string]string
var sch gocron.Scheduler

// Init init func before main
func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	//log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	//log.SetLevel(log.WarnLevel)

	/*
		if Environment == "production" {
			log.SetFormatter(&log.JSONFormatter{})
		} else {
			// The TextFormatter is default, you don't actually have to do this.
			log.SetFormatter(&log.TextFormatter{})
		}
	*/
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetLevel(log.DebugLevel)
	//log.SetLevel(log.ErrorLevel)
}

// RequestParser handle message from services
func RequestParser(socket *zmq.Socket, msg []string) error {
	// Check the length of multi-part message
	if len(msg) != 2 {
		log.Error("Request parser failed: invalid message length")
		return errors.New("Invalid message length")
	}

	log.WithFields(log.Fields{"msg[0]": msg[0]}).Debug("Parsing request:")

	switch msg[0] {
	case "mbtcp.once.read":
		var req MbtcpReadReq
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal request failed:")
			return err
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
		taskMap[cmd.Tid] = req
		taskMap2[cmd.Tid] = msg[0]
		sch.Emergency().Do(Task, socket, cmd)
		//sch.Every(1).Seconds().Do(modbusTask, socket, cmd)
		return nil
	case "mbtcp.once.write":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.timeout.read":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.timeout.update":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.poll.create":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.poll.update":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.poll.read":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.poll.delete":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.poll.toggle":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.polls.read":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.polls.delete":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.polls.toggle":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.polls.import":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.polls.export":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.poll.history":
		log.Warn("TODO")
		return errors.New("TODO")
	default:
		log.Error("unsupport")
		return errors.New("TODO")
	}
}

// ResponseParser handle message from modbusd
func ResponseParser(socket *zmq.Socket, msg []string) error {
	// Check the length of multi-part message
	if len(msg) != 2 {
		log.Error("Request parser failed: invalid message length")
		return errors.New("Invalid message length")
	}

	log.WithFields(log.Fields{"msg[0]": msg[0]}).Debug("Parsing response:")

	// Convert zframe 1: command number
	cmdType, err := strconv.Atoi(msg[0])
	if err != nil {
		return err
	}
	var cmdStr []byte
	var TidStr string

	switch cmdType {
	case 50, 51:
		var res DMbtcpTimeout
		if err := json.Unmarshal([]byte(msg[1]), &res); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal failed:")
			return err
		}
		//log.Debug(res)
		TidStr = res.Tid
		tid, _ := strconv.ParseInt(res.Tid, 10, 64)
		command := MbtcpTimeoutRes{
			Tid:    tid,
			Status: res.Status,
			Data:   res.Timeout,
		}
		cmdStr, _ = json.Marshal(command)
	case 1, 2, 3, 4:
		var res DMbtcpRes
		if err := json.Unmarshal([]byte(msg[1]), &res); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal failed:")
			return err
		}
		tid, _ := strconv.ParseInt(res.Tid, 10, 64)
		TidStr = res.Tid
		if cmd, ok := taskMap2[TidStr]; ok {
			switch cmd {
			case "mbtcp.once.read":
				if req, ok := taskMap[TidStr]; ok {
					readReq := req.(MbtcpReadReq)
					log.Println(readReq)
					log.WithFields(log.Fields{"Req type": readReq.Type}).Debug("Request type:")
					switch readReq.Type {
					case 2:
						//
					case 3:
						//
					case 4:
						//
					case 5:
						//
					case 6:
						//
					case 7:
						//
					case 8:
						//
					default:
						//
					}
				} else {
					return errors.New("req not in map")
				}

				command := MbtcpReadRes{
					Tid:    tid,
					Status: res.Status,
					Data:   res.Data,
				}
				cmdStr, _ = json.Marshal(command)
			default:
				//
			}

		} else {
			return errors.New("req command not in map")
		}

	case 5, 6, 15, 16:
		// todo
		return errors.New("TODO")
	default:
		// todo
		return errors.New("TODO")
	}

	log.WithFields(log.Fields{"JSON": string(cmdStr)}).Debug("Conversion for service complete:")
	// todo: check msg[0], should be web
	if frame, ok := taskMap2[TidStr]; ok {
		log.WithFields(log.Fields{"JSON": string(cmdStr)}).Debug("Send response to service:")
		delete(taskMap2, TidStr)
		socket.Send(frame, zmq.SNDMORE) // frame 1
		socket.Send(string(cmdStr), 0)  // convert to string; frame 2
	} else {
		log.WithFields(log.Fields{"JSON": string(cmdStr)}).Debug("Send response to service:")
		socket.Send("null", zmq.SNDMORE) // frame 1
		socket.Send(string(cmdStr), 0)   // convert to string; frame 2
	}

	t := time.Now()
	log.WithFields(log.Fields{"timestamp": t.Format("2006-01-02 15:04:05.000")}).Info("End ResponseParser:")
	return nil
}

// ResponseCommandBuilder build command to services
// Todo: filter, handle
func ResponseCommandBuilder() {

}

// Task for gocron
func Task(socket *zmq.Socket, m interface{}) {
	str, err := json.Marshal(m) // marshal to json string
	if err != nil {
		log.Error("Marshal request failed:", err)
		return
	}
	log.WithFields(log.Fields{"JSON": string(str)}).Debug("Send request to modbusd:")

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
				log.WithFields(log.Fields{"msg[0]": msg[0], "msg[1]": msg[1]}).Debug("Receive from service:")
				RequestParser(toModbusd, msg)
			case fromModbusd:
				msg, _ := fromModbusd.RecvMessage(0)
				log.WithFields(log.Fields{"msg[0]": msg[0], "msg[1]": msg[1]}).Debug("Receive from modbusd:")
				ResponseParser(toService, msg)
			}
		}
	}
}

//t := time.Now()
//fmt.Println("zrecv" + t.Format("2006-01-02 15:04:05.000"))
