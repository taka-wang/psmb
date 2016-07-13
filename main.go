package main

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/taka-wang/gocron"
	log "github.com/takawang/logrus"
	zmq "github.com/takawang/zmq3"
)

const (
	// DefaultPort default modbus slave port number
	DefaultPort = "502"
)

// MbTaskReq task request
type MbTaskReq struct {
	Cmd string
	Req interface{}
}

var taskMap map[string]MbTaskReq
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

// RequestParser handle message from services
func RequestParser(msg []string) (interface{}, error) {
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
		log.WithFields(log.Fields{"request": msg[0]}).Debug("Request not support:")
		return nil, errors.New("Request not support")
	}
}

// RequestHandler build command to services
func RequestHandler(cmd string, r interface{}, socket *zmq.Socket) error {
	log.WithFields(log.Fields{"cmd": cmd}).Debug("Build request command:")

	switch cmd {
	case "mbtcp.once.read": // done
		req := r.(MbtcpReadReq)

		// convert tid to string
		TidStr := strconv.FormatInt(req.Tid, 10)

		// add to task map
		taskMap[TidStr] = MbTaskReq{
			Cmd: cmd,
			Req: req,
		}
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
		// convert tid to string
		TidStr := strconv.FormatInt(req.Tid, 10)

		// add to task map
		taskMap[TidStr] = MbTaskReq{
			Cmd: cmd,
			Req: req,
		}

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
	case "mbtcp.timeout.read": // todo
		// add to Emergency
		return errors.New("TODO")
	case "mbtcp.timeout.update": // todo
		// add to Emergency
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
		log.WithFields(log.Fields{"cmd": cmd}).Debug("Request not support:")
		return errors.New("Request not support")
	}
}

// ResponseParser handle message from modbusd
// Done.
func ResponseParser(msg []string) (interface{}, error) {
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
		log.WithFields(log.Fields{"response": msg[0]}).Debug("Response not support:")
		return nil, errors.New("Response not support")
	}
}

// ResponseHandler build command to services
// Todo: filter, handle
func ResponseHandler(cmd MbtcpCmdType, r interface{}, socket *zmq.Socket) error {
	log.WithFields(log.Fields{"cmd": cmd}).Debug("Parsing response:")

	var cmdStr []byte
	var TidStr string
	var task MbTaskReq
	var ok bool

	switch cmd {
	case fc5, fc6, fc15, fc16, setTimeout, getTimeout:
		//
	case fc1, fc2, fc3, fc4:
		//
	default:
		//

	}

	switch cmd {
	case "50", "51": // set|get timeout
		res := r.(DMbtcpTimeout)
		tid, _ := strconv.ParseInt(res.Tid, 10, 64)
		TidStr = res.Tid
		if task, ok = taskMap[TidStr]; ok {
			command := MbtcpTimeoutRes{
				Tid:    tid,
				Status: res.Status,
				Data:   res.Timeout,
			}
			cmdStr, _ = json.Marshal(command)
		} else {
			return errors.New("req command not in map")
		}

	case "1", "2": // read bits
		res := r.(DMbtcpRes)
		tid, _ := strconv.ParseInt(res.Tid, 10, 64)
		TidStr = res.Tid
		if task, ok = taskMap[TidStr]; ok {
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

		} else {
			return errors.New("req command not in map")
		}

	case "3", "4": // read registers
		res := r.(DMbtcpRes)
		tid, _ := strconv.ParseInt(res.Tid, 10, 64)
		TidStr = res.Tid
		if task, ok = taskMap[TidStr]; ok {
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

		} else {
			return errors.New("req command not in map")
		}

	case "5", "6": // write single
		res := r.(DMbtcpRes)
		tid, _ := strconv.ParseInt(res.Tid, 10, 64)
		TidStr = res.Tid
		if task, ok = taskMap[TidStr]; ok {
			switch task.Cmd {
			case "mbtcp.once.write":
				command := MbtcpSimpleRes{
					Tid:    tid,
					Status: res.Status,
				}
				cmdStr, _ = json.Marshal(command)
			default:
				//
			}

		} else {
			return errors.New("req command not in map")
		}
	case "15", "16": // write multiple
		res := r.(DMbtcpRes)
		tid, _ := strconv.ParseInt(res.Tid, 10, 64)
		TidStr = res.Tid
		if task, ok = taskMap[TidStr]; ok {
			switch task.Cmd {
			case "mbtcp.once.write":
				command := MbtcpSimpleRes{
					Tid:    tid,
					Status: res.Status,
				}
				cmdStr, _ = json.Marshal(command)
			default:
				//
			}

		} else {
			return errors.New("req command not in map")
		}
	default:
		log.WithFields(log.Fields{"cmd": cmd}).Debug("Response not support:")
		return errors.New("Response not support")
	}

	log.WithFields(log.Fields{"JSON": string(cmdStr)}).Debug("Send response to service:")
	//delete(taskMap2, TidStr)
	socket.Send(task.Cmd, zmq.SNDMORE) // frame 1
	socket.Send(string(cmdStr), 0)     // convert to string; frame 2

	return nil
}

func main() {

	taskMap = make(map[string]MbTaskReq)

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
				msg, _ := fromService.RecvMessage(0)
				log.WithFields(log.Fields{
					"msg[0]": msg[0],
					"msg[1]": msg[1],
				}).Debug("Receive from service:")
				req, err := RequestParser(msg)
				if err != nil {
					// todo: send error back
				} else {
					err = RequestHandler(msg[0], req, toModbusd)
				}
			case fromModbusd:
				msg, _ := fromModbusd.RecvMessage(0)
				log.WithFields(log.Fields{
					"msg[0]": msg[0],
					"msg[1]": msg[1],
				}).Debug("Receive from modbusd:")
				res, err := ResponseParser(msg)
				if err != nil {
					// todo: send error back
				} else {
					err = ResponseHandler(msg[0], res, toService)
				}
			}
		}
	}
}

//t := time.Now()
//fmt.Println("zrecv" + t.Format("2006-01-02 15:04:05.000"))
//sch.Every(1).Seconds().Do(modbusTask, socket, cmd)
