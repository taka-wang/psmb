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
func RequestParser(msg []string) (interface{}, error) {
	// Check the length of multi-part message
	if len(msg) != 2 {
		log.Error("Request parser failed: invalid message length")
		return nil, errors.New("Invalid message length")
	}

	log.WithFields(log.Fields{"msg[0]": msg[0]}).Debug("Parsing request:")

	switch msg[0] {
	case "mbtcp.once.read":
		var req MbtcpReadReq
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal request failed:")
			return nil, err
		}
		return req, nil
	case "mbtcp.once.write":
		// add to Emergency
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case "mbtcp.timeout.read", "mbtcp.timeout.update":
		var req MbtcpTimeoutReq
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal request failed:")
			return nil, err
		}
		return req, nil
	case "mbtcp.poll.create":
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

// RequestCmdBuilder build command to services
func RequestCmdBuilder(cmd string, r interface{}, socket *zmq.Socket) error {
	log.WithFields(log.Fields{"cmd": cmd}).Debug("Build request command:")

	switch cmd {
	case "mbtcp.once.read":
		req := r.(MbtcpReadReq)

		// convert tid to string
		TidStr := strconv.FormatInt(req.Tid, 10)

		// add to task map
		taskMap[TidStr] = req
		taskMap2[TidStr] = cmd

		// build modbusd command
		cmd := DMbtcpReadReq{
			Tid:   TidStr,
			Cmd:   req.FC,
			IP:    req.IP,
			Port:  req.Port,
			Slave: req.Slave,
			Addr:  req.Addr,
			Len:   req.Len,
		}
		// add command to scheduler as emergency request
		sch.Emergency().Do(Task, socket, cmd)
		return nil
	case "mbtcp.once.write":
		// add to Emergency
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.timeout.read":
		// add to Emergency
		log.Warn("TODO")
		return errors.New("TODO")
	case "mbtcp.timeout.update":
		// add to Emergency
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
func ResponseParser(msg []string) (interface{}, error) {
	// Check the length of multi-part message
	if len(msg) != 2 {
		log.Error("Request parser failed: invalid message length")
		return nil, errors.New("Invalid message length")
	}

	log.WithFields(log.Fields{"msg[0]": msg[0]}).Debug("Parsing response:")

	switch msg[0] {
	case "50", "51": // set|get timeout
		var res DMbtcpTimeout
		if err := json.Unmarshal([]byte(msg[1]), &res); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal failed:")
			return nil, err
		}
		return res, nil

	case "1", "2", "3", "4": // read command
		var res DMbtcpRes
		if err := json.Unmarshal([]byte(msg[1]), &res); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal failed:")
			return nil, err
		}
		return res, nil

	case "5", "6": // write single
		// todo
		return nil, errors.New("TODO")
	case "15", "16": // write multiple
		// todo
		return nil, errors.New("TODO")
	default:
		log.WithFields(log.Fields{"response": msg[0]}).Debug("Response not support:")
		return nil, errors.New("Response not support")
	}
}

// ResponseCmdBuilder build command to services
// Todo: filter, handle
func ResponseCmdBuilder(cmd string, r interface{}, socket *zmq.Socket) error {
	log.WithFields(log.Fields{"cmd": cmd}).Debug("Parsing response:")

	var cmdStr []byte
	var TidStr string

	switch msg[0] {
	case "50", "51": // set|get timeout
		res := r.(DMbtcpTimeout)
		tid, _ := strconv.ParseInt(res.Tid, 10, 64)
		TidStr = res.Tid
		command := MbtcpTimeoutRes{
			Tid:    tid,
			Status: res.Status,
			Data:   res.Timeout,
		}
		cmdStr, _ = json.Marshal(command)
	case "1", "2": // read bits
		res := r.(DMbtcpRes)
		tid, _ := strconv.ParseInt(res.Tid, 10, 64)
		TidStr = res.Tid
		if cmd, ok := taskMap2[TidStr]; ok {
			switch cmd {
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
		if cmd, ok := taskMap2[TidStr]; ok {
			switch cmd {
			case "mbtcp.once.read":

				// todo: if res.status != "ok" do something

				var command MbtcpReadRes
				if req, ok := taskMap[TidStr]; ok {
					readReq := req.(MbtcpReadReq)
					//log.Println(readReq)
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
				} else {
					return errors.New("req not in map")
				}

				cmdStr, _ = json.Marshal(command)
			default:
				//
			}

		} else {
			return errors.New("req command not in map")
		}

	case "5", "6": // write single
		// todo
		return errors.New("TODO")
	case "15", "16": // write multiple
		// todo
		return errors.New("TODO")
	default:
		log.WithFields(log.Fields{"response": msg[0]}).Debug("Response not support:")
		return errors.New("Response not support")
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
					// todo
					// send error back
				} else {
					err = RequestCmdBuilder(msg[0], req, toModbusd)
				}
			case fromModbusd:
				msg, _ := fromModbusd.RecvMessage(0)
				log.WithFields(log.Fields{
					"msg[0]": msg[0],
					"msg[1]": msg[1],
				}).Debug("Receive from modbusd:")
				res, err := ResponseParser(msg)
				if err != nil {
					// todo
					// send error back
				} else {
					err = ResponseCmdBuilder(msg[0], res, toService)
				}
			}
		}
	}
}

//t := time.Now()
//fmt.Println("zrecv" + t.Format("2006-01-02 15:04:05.000"))
//sch.Every(1).Seconds().Do(modbusTask, socket, cmd)
