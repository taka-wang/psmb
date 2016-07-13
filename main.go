package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/taka-wang/gocron"
	log "github.com/takawang/logrus"
	zmq "github.com/takawang/zmq3"
)

const (
	// DefaultPort default modbus slave port number
	DefaultPort = "502"
	// MinTCPTimeout minimal modbus tcp connection timeout
	MinTCPTimeout = 200000
	// MinTCPPollInterval minimal modbus tcp poll interval
	MinTCPPollInterval = 1
)

var sch gocron.Scheduler

// MbtcpReadTaskMap read/poll task map
var MbtcpReadTaskMap = struct {
	sync.RWMutex
	m map[string]mbtcpReadTask // m mbtcpReadTask map
}{m: make(map[string]mbtcpReadTask)}

// GetMbtcpReadTask get task from map
func GetMbtcpReadTask(tid string) (mbtcpReadTask, bool) {
	MbtcpReadTaskMap.RLock()
	task, ok := MbtcpReadTaskMap.m[tid]
	MbtcpReadTaskMap.RUnlock()
	return task, ok
}

// RemoveMbtcpReadTask remove task from  map
func RemoveMbtcpReadTask(tid string) {
	MbtcpReadTaskMap.Lock()
	delete(MbtcpReadTaskMap.m, tid) // remove from  Map!!
	MbtcpReadTaskMap.Unlock()
}

// AddMbtcpReadTask add one-off task to  map
func AddMbtcpReadTask(name, tid, cmd string, req interface{}) {
	MbtcpReadTaskMap.Lock()
	MbtcpReadTaskMap.m[tid] = mbtcpReadTask{name, cmd, req}
	MbtcpReadTaskMap.Unlock()
}

// SimpleTask simple task (except read request)
var SimpleTask = struct {
	sync.RWMutex
	m map[string]string // tid: command
}{m: make(map[string]string)}

// GetSimpleTask get command from map
func GetSimpleTask(tid string) (string, bool) {
	SimpleTask.RLock()
	cmd, ok := SimpleTask.m[tid]
	SimpleTask.RUnlock()
	return cmd, ok
}

// RemoveSimpleTask remove command from map
func RemoveSimpleTask(tid string) {
	SimpleTask.Lock()
	delete(SimpleTask.m, tid)
	SimpleTask.Unlock()
}

// AddSimpleTask add cmd to map
func AddSimpleTask(tid, cmd string) {
	SimpleTask.Lock()
	SimpleTask.m[tid] = cmd
	SimpleTask.Unlock()
}

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

// Task for gocron
func Task(socket *zmq.Socket, req interface{}) {
	str, err := json.Marshal(req) // marshal to json string
	if err != nil {
		log.Error("Marshal request failed:", err)
		// todo: remove table
		return
	}
	log.WithFields(log.Fields{"JSON": string(str)}).Debug("Send request to modbusd:")

	socket.Send("tcp", zmq.SNDMORE) // frame 1
	socket.Send(string(str), 0)     // convert to string; frame 2
}

// SimpleTaskResponser send simple response to service
func SimpleTaskResponser(tid string, resp interface{}, socket *zmq.Socket) error {
	respStr, err := json.Marshal(resp)
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Error("Marshal failed:")
		return err
	}

	if cmd, ok := GetSimpleTask(tid); ok {
		log.WithFields(log.Fields{"JSON": string(respStr)}).Debug("Send response to service:")
		socket.Send(cmd, zmq.SNDMORE)   // task command
		socket.Send(string(respStr), 0) // convert to string; frame 2
		RemoveSimpleTask(tid)           // remove from Map!!
		return nil
	}
	return errors.New("Request command not in map")
}

// ParseRequest parse message from services
// R&R: only unmarshal request string to corresponding struct
func ParseRequest(msg []string) (interface{}, error) {
	// Check the length of multi-part message
	if len(msg) != 2 {
		// should not reach here!!
		log.Error("Request parser failed: invalid message length")
		return nil, errors.New("Invalid message length")
	}

	log.WithFields(log.Fields{"msg[0]": msg[0]}).Debug("Parsing request:")

	switch msg[0] {
	case "mbtcp.once.read": // done
		var req MbtcpReadReq
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal request failed:")
			return nil, err
		}
		return req, nil
	case "mbtcp.once.write": // done
		var data json.RawMessage // raw []byte
		req := MbtcpWriteReq{Data: &data}
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal request failed:")
			return nil, err
		}

		switch req.FC {
		case 5: // single bit; uint16
			var d uint16
			if err := json.Unmarshal(data, &d); err != nil {
				log.WithFields(log.Fields{"Error": err}).Error("Unmarshal FC5 request failed:")
				return nil, err
			}
			req.Data = d
			return req, nil
		case 6: // single register in dec|hex
			var d string
			if err := json.Unmarshal(data, &d); err != nil {
				log.WithFields(log.Fields{"Error": err}).Error("Unmarshal FC6 request failed:")
				return nil, err
			}
			if req.Hex {
				dd, err := HexStringToRegisters(d)
				if err != nil {
					log.WithFields(log.Fields{"Error": err}).Error("Unmarshal FC6 Hex request failed:")
					return nil, err
				}
				req.Data = dd[0] // one register
			} else {
				dd, err := strconv.Atoi(d)
				if err != nil {
					log.WithFields(log.Fields{"Error": err}).Error("Unmarshal FC6 Dec request failed:")
					return nil, err
				}
				req.Data = dd
			}
			return req, nil
		case 15: // multiple bits; []uint16
			var d []uint16
			if err := json.Unmarshal(data, &d); err != nil {
				log.WithFields(log.Fields{"Error": err}).Error("Unmarshal FC15 request failed:")
				return nil, err
			}
			req.Data = d

			// len checker
			l := uint16(len(d))
			if req.Len < l {
				req.Len = l
			}
			return req, nil
		case 16: // multiple register in dec/hex
			var d string
			if err := json.Unmarshal(data, &d); err != nil {
				log.WithFields(log.Fields{"Error": err}).Error("Unmarshal FC16 request failed:")
				return nil, err
			}
			var l uint16 // length
			if req.Hex {
				dd, err := HexStringToRegisters(d)
				if err != nil {
					log.WithFields(log.Fields{"Error": err}).Error("Unmarshal FC16 Hex request failed:")
					return nil, err
				}
				req.Data = dd
				l = uint16(len(dd))
			} else {
				dd, err := DecimalStringToRegisters(d)
				if err != nil {
					log.WithFields(log.Fields{"Error": err}).Error("Unmarshal FC16 Dec request failed:")
					return nil, err
				}
				req.Data = dd
				l = uint16(len(dd))
			}
			// len checker
			if req.Len < l {
				req.Len = l
			}
			return req, nil
		default:
			return nil, errors.New("Request not support")
		}
	case "mbtcp.timeout.read", "mbtcp.timeout.update": // done
		var req MbtcpTimeoutReq
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal request failed:")
			return nil, err
		}
		return req, nil
	case "mbtcp.poll.create": // done
		var req MbtcpPollStatus
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal request failed:")
			return nil, err
		}
		return req, nil
	case "mbtcp.poll.update", "mbtcp.poll.read", "mbtcp.poll.delete", // done
		"mbtcp.polls.read", "mbtcp.poll.toggle", "mbtcp.polls.delete",
		"mbtcp.polls.toggle", "mbtcp.poll.history", "mbtcp.polls.export":
		var req MbtcpPollOpReq
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal request failed:")
			return nil, err
		}
		return req, nil
	case "mbtcp.polls.import": // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filter.create": // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filter.update": // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filter.read": // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filter.delete": // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filter.toggle": // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filters.read": // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filters.delete": // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filters.toggle": // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filters.import": // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filters.export": // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	default: // done
		// should not reach here!!
		log.WithFields(log.Fields{"request": msg[0]}).Warn("Request not support:")
		return nil, errors.New("Request not support")
	}
}

// RequestHandler build command to services
func RequestHandler(cmd string, r interface{}, downSocket, upSocket *zmq.Socket) error {
	log.WithFields(log.Fields{"cmd": cmd}).Debug("Build request command:")

	switch cmd {
	case "mbtcp.once.read": // done
		req := r.(MbtcpReadReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string
		// default port checker
		if req.Port == "" {
			req.Port = DefaultPort
		}

		AddMbtcpReadTask("", TidStr, cmd, req) // add to task map
		command := DMbtcpReadReq{
			Tid:   TidStr,
			Cmd:   req.FC,
			IP:    req.IP,
			Port:  req.Port,
			Slave: req.Slave,
			Addr:  req.Addr,
			Len:   req.Len,
		}
		// add command to scheduler as emergency request
		sch.Emergency().Do(Task, downSocket, command)
		return nil
	case "mbtcp.once.write": // done
		req := r.(MbtcpWriteReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string
		// default port checker
		if req.Port == "" {
			req.Port = DefaultPort
		}
		AddSimpleTask(TidStr, cmd) // add to task map
		command := DMbtcpWriteReq{
			Tid:   TidStr,
			Cmd:   req.FC,
			IP:    req.IP,
			Port:  req.Port,
			Slave: req.Slave,
			Addr:  req.Addr,
			Len:   req.Len,
			Data:  req.Data,
		}

		// add command to scheduler as emergency request
		sch.Emergency().Do(Task, downSocket, command)
		return nil
	case "mbtcp.timeout.read": // done
		req := r.(MbtcpTimeoutReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string
		AddSimpleTask(TidStr, cmd)               // add to task map
		cmdInt, _ := strconv.Atoi(string(getTimeout))
		command := DMbtcpTimeout{
			Tid: TidStr,
			Cmd: cmdInt,
		}
		// add command to scheduler as emergency request
		sch.Emergency().Do(Task, downSocket, command)
		return nil
	case "mbtcp.timeout.update": // done
		req := r.(MbtcpTimeoutReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string
		// protect dummy input
		if req.Data < MinTCPTimeout {
			req.Data = MinTCPTimeout
		}
		AddSimpleTask(TidStr, cmd) // add to task map
		cmdInt, _ := strconv.Atoi(string(setTimeout))
		command := DMbtcpTimeout{
			Tid:     TidStr,
			Cmd:     cmdInt,
			Timeout: req.Data,
		}
		// add command to scheduler as emergency request
		sch.Emergency().Do(Task, downSocket, command)
		return nil
	case "mbtcp.poll.create":

		// TODO! check name
		// TODO! check fc code!!!!!!!
		req := r.(MbtcpPollStatus)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string

		// check interval value
		if req.Interval < MinTCPPollInterval {
			req.Interval = MinTCPPollInterval
		}
		// check port
		if req.Port == "" {
			req.Port = DefaultPort
		}

		AddSimpleTask(TidStr, cmd)                   // simple request
		AddMbtcpReadTask(req.Name, TidStr, cmd, req) // read request

		command := DMbtcpReadReq{
			Tid:   TidStr,
			Cmd:   req.FC,
			IP:    req.IP,
			Port:  req.Port,
			Slave: req.Slave,
			Addr:  req.Addr,
			Len:   req.Len,
		}
		// check name
		// add to polling table
		sch.EveryWithName(req.Interval, req.Name).Seconds().Do(Task, downSocket, command)
		if !req.Enabled {
			sch.PauseWithName(req.Name)
		}
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: "ok"}
		return SimpleTaskResponser(TidStr, resp, upSocket)
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
		// TODO! check name
		req := r.(MbtcpPollOpReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string
		AddSimpleTask(TidStr, cmd)               // add to simple task map

		status := "ok"
		if req.Enabled {
			if ok := sch.ResumeWithName(req.Name); !ok {
				status = "enable poll failed"
			}
		} else {
			if ok := sch.PauseWithName(req.Name); !ok {
				status = "disable poll failed"
			}
		}
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: status}
		return SimpleTaskResponser(TidStr, resp, upSocket)
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
	case "mbtcp.filter.create":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.filter.update":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.filter.read":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.filter.delete":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.filter.toggle":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.filters.read":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.filters.delete":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.filters.toggle":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.filters.import":
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.filters.export":
		log.Warn("TODO")
		return errors.New("TODO")
	default:
		// should not reach here!!
		log.WithFields(log.Fields{"cmd": cmd}).Warn("Request not support:")
		return errors.New("Request not support")
	}
}

// ParseResponse parse message from modbusd
// Done.
func ParseResponse(msg []string) (interface{}, error) {
	// Check the length of multi-part message
	if len(msg) != 2 {
		log.Error("Request parser failed: invalid message length")
		return nil, errors.New("Invalid message length")
	}

	log.WithFields(log.Fields{"msg[0]": msg[0]}).Debug("Parsing response:")

	switch MbtcpCmdType(msg[0]) {
	case setTimeout, getTimeout: // set|get timeout
		var res DMbtcpTimeout
		if err := json.Unmarshal([]byte(msg[1]), &res); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal failed:")
			return nil, err
		}
		return res, nil

	case fc1, fc2, fc3, fc4, fc5, fc6, fc15, fc16:
		var res DMbtcpRes
		if err := json.Unmarshal([]byte(msg[1]), &res); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal failed:")
			return nil, err
		}
		return res, nil
	default:
		// should not reach here!!
		log.WithFields(log.Fields{"response": msg[0]}).Warn("Response not support:")
		return nil, errors.New("Response not support")
	}
}

// ResponseHandler build command to services
// Todo: filter, handle
func ResponseHandler(cmd MbtcpCmdType, r interface{}, socket *zmq.Socket) error {
	log.WithFields(log.Fields{"cmd": cmd}).Debug("Handle response:")

	switch cmd {
	case fc5, fc6, fc15, fc16, setTimeout, getTimeout: // [done]: one-off requests
		var TidStr string
		var resp interface{}

		switch cmd {
		case setTimeout: // one-off timeout requests
			res := r.(DMbtcpTimeout)
			tid, _ := strconv.ParseInt(res.Tid, 10, 64)
			TidStr = res.Tid
			resp = MbtcpTimeoutRes{
				Tid:    tid,
				Status: res.Status,
			}
		case getTimeout: // one-off timeout requests
			res := r.(DMbtcpTimeout)
			tid, _ := strconv.ParseInt(res.Tid, 10, 64)
			TidStr = res.Tid
			resp = MbtcpTimeoutRes{
				Tid:    tid,
				Status: res.Status,
				Data:   res.Timeout,
			}
		case fc5, fc6, fc15, fc16: // one-off write requests
			res := r.(DMbtcpRes)
			tid, _ := strconv.ParseInt(res.Tid, 10, 64)
			TidStr = res.Tid
			resp = MbtcpSimpleRes{
				Tid:    tid,
				Status: res.Status,
			}
		}
		// send back
		return SimpleTaskResponser(TidStr, resp, socket)

	case fc1, fc2, fc3, fc4: // one-off and polling requests
		var cmdStr []byte
		var task mbtcpReadTask
		var ok bool
		switch cmd {
		case fc1, fc2:
			res := r.(DMbtcpRes)
			tid, _ := strconv.ParseInt(res.Tid, 10, 64)

			// check read task table
			if task, ok = GetMbtcpReadTask(res.Tid); !ok {
				return errors.New("req command not in map")
			}

			var response interface{}

			switch task.Cmd {
			case "mbtcp.once.read":
				response = MbtcpReadRes{
					Tid:    tid,
					Status: res.Status,
					Data:   res.Data,
				}
			case "mbtcp.poll.create", "mbtcp.polls.import":
				response = MbtcpPollData{
					TimeStamp: time.Now().UTC().UnixNano(),
					Name:      task.Name,
					Status:    res.Status,
					Data:      res.Data,
				}
				// test
				log.WithFields(log.Fields{"TS": time.Now().Format("2006-01-02 15:04:05.000")}).Error("Time Stamp:")
			default: // should not reach here
				//
				log.Error("Should not reach here")
				response = MbtcpSimpleRes{
					Tid:    tid,
					Status: "not support command",
				}
			}
			cmdStr, _ = json.Marshal(response)

		case fc3, fc4:
			res := r.(DMbtcpRes)
			tid, _ := strconv.ParseInt(res.Tid, 10, 64)
			if task, ok = GetMbtcpReadTask(res.Tid); !ok {
				log.Error("req command not in map")
				return errors.New("req command not in map")
			}

			switch task.Cmd {
			case "mbtcp.once.read":
				// todo: if res.status != "ok" do something
				var command MbtcpReadRes
				readReq := task.Req.(MbtcpReadReq)
				log.WithFields(log.Fields{"Req type": readReq.Type}).Debug("Request type:")
				switch readReq.Type {
				case 2:
					b, err := RegistersToBytes(res.Data)
					if err != nil {
						log.Error(err)
						command = MbtcpReadRes{
							Tid:    tid,
							Status: err.Error(),
							Data:   res.Data,
						}
					}
					command = MbtcpReadRes{
						Tid:    tid,
						Status: res.Status,
						Bytes:  b,
						Data:   BytesToHexString(b),
					}
				case 3:
					if readReq.Len%2 != 0 {
						command = MbtcpReadRes{
							Tid:    tid,
							Status: "Conversion failed",
							Data:   res.Data,
						}
					} else {
						b, err := RegistersToBytes(res.Data)
						if err != nil {
							log.Error(err)
							command = MbtcpReadRes{
								Tid:    tid,
								Status: err.Error(),
								Data:   res.Data,
							}
						}

						// todo: check range values

						f := LinearScalingRegisters(res.Data,
							readReq.Range.DomainLow,
							readReq.Range.DomainHigh,
							readReq.Range.RangeLow,
							readReq.Range.RangeHigh,
						)

						command = MbtcpReadRes{
							Tid:    tid,
							Status: res.Status,
							Bytes:  b,
							Data:   f,
						}
					}

				case 4:
					// order
					command = MbtcpReadRes{
						Tid:    tid,
						Status: res.Status,
						Data:   res.Data,
					}
				case 5:
					// order
					command = MbtcpReadRes{
						Tid:    tid,
						Status: res.Status,
						Data:   res.Data,
					}
				case 6:
					// length, order
					command = MbtcpReadRes{
						Tid:    tid,
						Status: res.Status,
						Data:   res.Data,
					}
				case 7:
					// length, order
					command = MbtcpReadRes{
						Tid:    tid,
						Status: res.Status,
						Data:   res.Data,
					}
				case 8:
					// length, order
					command = MbtcpReadRes{
						Tid:    tid,
						Status: res.Status,
						Data:   res.Data,
					}
				default: // case 0, 1
					command = MbtcpReadRes{
						Tid:    tid,
						Status: res.Status,
						Data:   res.Data,
					}
				}

				cmdStr, _ = json.Marshal(command)
			default:
				//
				log.Debug("Maybe polling request")
			}
		}
		// different handle
		log.WithFields(log.Fields{"JSON": string(cmdStr)}).Debug("Send response to service:")
		socket.Send(task.Cmd, zmq.SNDMORE) // frame 1
		socket.Send(string(cmdStr), 0)     // convert to string; frame 2

		return nil
	default:
		// should not reach here!!
		log.WithFields(log.Fields{"cmd": cmd}).Warn("Response not support:")
		return errors.New("Response not support")
	}
}

func main() {

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

	// process messages from both subscriber sockets
	for {
		sockets, _ := poller.Poll(-1)
		for _, socket := range sockets {
			switch s := socket.Socket; s {
			case fromService:
				// receive from upstream
				msg, _ := fromService.RecvMessage(0)
				log.WithFields(log.Fields{
					"msg[0]": msg[0],
					"msg[1]": msg[1],
				}).Debug("Receive from service:")

				// parse request
				req, err := ParseRequest(msg)
				if err != nil {
					// todo: send error back
				} else {
					err = RequestHandler(msg[0], req, toModbusd, toService)
				}
			case fromModbusd:
				// receive from modbusd
				msg, _ := fromModbusd.RecvMessage(0)
				log.WithFields(log.Fields{
					"msg[0]": msg[0],
					"msg[1]": msg[1],
				}).Debug("Receive from modbusd:")

				// parse response
				res, err := ParseResponse(msg)
				if err != nil {
					// todo: send error back
				} else {
					err = ResponseHandler(MbtcpCmdType(msg[0]), res, toService)
				}
			}
		}
	}
}

//t := time.Now()
//fmt.Println("zrecv" + t.Format("2006-01-02 15:04:05.000"))
//sch.Every(1).Seconds().Do(modbusTask, socket, cmd)
