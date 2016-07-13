package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"

	"github.com/taka-wang/gocron"
	log "github.com/takawang/logrus"
	zmq "github.com/takawang/zmq3"
)

const (
	// DefaultPort default modbus slave port number
	DefaultPort = "502"
	// MinTCPTimeout minimal modbus tcp connection timeout
	MinTCPTimeout = 200000
)

var sch gocron.Scheduler

// OneOffTask one-off task map
var OneOffTask = struct {
	sync.RWMutex
	m map[string]MbtcpTaskReq // m MbtcpTaskReq map
}{m: make(map[string]MbtcpTaskReq)}

// GetOneOffTask get task from OneOffTask map
func GetOneOffTask(tid string) (MbtcpTaskReq, bool) {
	//log.Debug("GetOneOffTask IN")
	OneOffTask.RLock()
	task, ok := OneOffTask.m[tid]
	OneOffTask.RUnlock()
	//log.Debug("GetOneOffTask OUT")
	return task, ok
}

// RemoveOneOffTask remove task from OneOffTask map
func RemoveOneOffTask(tid string) {
	//log.Debug("RemoveOneOffTask IN")
	OneOffTask.Lock()
	delete(OneOffTask.m, tid) // remove from OneOffTask Map!!
	OneOffTask.Unlock()
	//log.Debug("RemoveOneOffTask OUT")
}

// AddOneOffTask add one-off task to OneOffTask map
func AddOneOffTask(tid, cmd string, req interface{}) {
	//log.Debug("AddOneOffTask IN")
	OneOffTask.Lock()
	OneOffTask.m[tid] = MbtcpTaskReq{cmd, req}
	OneOffTask.Unlock()
	//log.Debug("AddOneOffTask OUT")
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

// ParseRequest parse message from services
func ParseRequest(msg []string) (interface{}, error) {
	// Check the length of multi-part message
	if len(msg) != 2 {
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
		// default port checker
		if req.Port == "" {
			req.Port = DefaultPort
		}
		return req, nil
	case "mbtcp.once.write": // done
		var data json.RawMessage // raw []byte
		req := MbtcpWriteReq{Data: &data}
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal request failed:")
			return nil, err
		}

		// default port checker
		if req.Port == "" {
			req.Port = DefaultPort
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
	case "mbtcp.poll.create":
		// return immediately
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.poll.update":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.poll.read":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.poll.delete":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.poll.toggle":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.polls.read":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.polls.delete":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.polls.toggle":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.polls.import":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.polls.export":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.poll.history":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filter.create":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filter.update":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filter.read":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filter.delete":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filter.toggle":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filters.read":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filters.delete":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filters.toggle":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filters.import":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.filters.export":
		log.Warn("TODO")
		return nil, errors.New("TODO")
	default:
		// should not reach here!!
		log.WithFields(log.Fields{"request": msg[0]}).Warn("Request not support:")
		return nil, errors.New("Request not support")
	}
}

// RequestHandler build command to services
func RequestHandler(cmd string, r interface{}, socket *zmq.Socket) error {
	log.WithFields(log.Fields{"cmd": cmd}).Debug("Build request command:")

	switch cmd {
	case "mbtcp.once.read": // done
		req := r.(MbtcpReadReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string
		AddOneOffTask(TidStr, cmd, req)          // add to task map

		// build modbusd command
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
		sch.Emergency().Do(Task, socket, command)
		return nil
	case "mbtcp.once.write": // done
		req := r.(MbtcpWriteReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string
		AddOneOffTask(TidStr, cmd, req)          // add to task map

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
		sch.Emergency().Do(Task, socket, command)
		return nil
	case "mbtcp.timeout.read", "mbtcp.timeout.update": // done
		req := r.(MbtcpTimeoutReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string
		AddOneOffTask(TidStr, cmd, req)          // add to task map
		command := DMbtcpTimeout{
			Tid: TidStr,
			Cmd: req.FC,
		}

		if cmd == "mbtcp.timeout.update" {
			// protect dummy input
			if req.Data < MinTCPTimeout {
				command.Timeout = MinTCPTimeout
			} else {
				command.Timeout = req.Data
			}
		}
		// add command to scheduler as emergency request
		sch.Emergency().Do(Task, socket, command)
		return nil
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
		case setTimeout, getTimeout: // one-off timeout requests
			res := r.(DMbtcpTimeout)
			tid, _ := strconv.ParseInt(res.Tid, 10, 64)
			TidStr = res.Tid
			resp = MbtcpTimeoutRes{
				Tid:    tid,
				Status: res.Status,
			}
			if cmd == getTimeout {
				resp.Data = res.Timeout
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

		// ----------------- modulize begin ---------------------------------------------

		// marshal response JSON string
		respStr, err := json.Marshal(resp)
		if err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Marshal failed:")
			return err
		}

		if task, ok := GetOneOffTask(TidStr); !ok {
			return errors.New("Request command not in map")
		}

		log.WithFields(log.Fields{"JSON": string(respStr)}).Debug("Send response to service:")
		RemoveOneOffTask(TidStr)           // remove from OneOffTask Map!!
		socket.Send(task.Cmd, zmq.SNDMORE) // task command
		socket.Send(string(respStr), 0)    // convert to string; frame 2
		return nil

		// ----------------- modulize end ---------------------------------------------

	case fc1, fc2, fc3, fc4: // one-off and polling requests
		var cmdStr []byte
		var TidStr string
		var task MbtcpTaskReq
		var ok bool
		switch cmd {
		case fc1, fc2:
			res := r.(DMbtcpRes)
			tid, _ := strconv.ParseInt(res.Tid, 10, 64)
			TidStr = res.Tid

			if task, ok = GetOneOffTask(TidStr); !ok {
				return errors.New("req command not in map")
			}
			switch task.Cmd {
			case "mbtcp.once.read":
				command := MbtcpReadRes{
					Tid:    tid,
					Status: res.Status,
					Data:   res.Data,
				}
				cmdStr, _ = json.Marshal(command)
			default:
				//
			}

		case fc3, fc4:
			res := r.(DMbtcpRes)
			tid, _ := strconv.ParseInt(res.Tid, 10, 64)
			TidStr = res.Tid
			if task, ok = GetOneOffTask(TidStr); !ok {
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
					err = RequestHandler(msg[0], req, toModbusd)
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
