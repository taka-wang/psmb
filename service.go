package psmb

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/taka-wang/gocron"
	log "github.com/takawang/logrus"
	zmq "github.com/takawang/zmq3"
)

const (
	// defaultMbtcpPort default modbus slave port number
	defaultMbtcpPort = "502"
	// minMbtcpTimeout minimal modbus tcp connection timeout
	minMbtcpTimeout = 200000
	// minMbtcpPollInterval minimal modbus tcp poll interval
	minMbtcpPollInterval = 1
)

// ProactiveService proactive service contracts,
// all services should implement the following methods.
type ProactiveService interface {
	// Start enable proactive service
	Start()
	// Stop disable proactive service
	Stop()

	// parseRequest parse upstream requests
	parseRequest(msg []string) (interface{}, error)
	// handleRequest handle upstream requests
	handleRequest(cmd string, r interface{}) error
	// parseResponse parse downstream responses
	parseResponse(msg []string) (interface{}, error)
	// handleResponse handle downstream responses
	handleResponse(cmd string, r interface{}) error
}

// mbtcpService modbusd tcp proactive service type
type mbtcpService struct {
	// readTaskMap read/poll task map
	readTaskMap MbtcpReadTask
	// simpleTaskMap simple task map
	simpleTaskMap MbtcpSimpleTask
	// scheduler gocron scheduler
	scheduler gocron.Scheduler
	// sub zmq subscriber endpoints
	sub struct {
		// upstream from services
		upstream *zmq.Socket
		// downstream from modbusd
		downstream *zmq.Socket
	}
	// pub zmq publisher endpoints
	pub struct {
		// upstream publish to services
		upstream *zmq.Socket
		// downstream publish to modbusd
		downstream *zmq.Socket
	}
	// poller zmq poller
	poller *zmq.Poller
	// enable poller flag
	enable bool
}

// NewPSMBTCP instantiate modbus tcp proactive serivce
func NewPSMBTCP() ProactiveService {
	return &mbtcpService{
		enable:        true,
		readTaskMap:   NewMbtcpReadTask(),
		simpleTaskMap: NewMbtcpSimpleTask(),
		scheduler:     gocron.NewScheduler(),
	}
}

// initLogger init logger
func (b *mbtcpService) initLogger() {
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

// Marshal helper function to marshal structure
func Marshal(r interface{}) (string, error) {
	bytes, err := json.Marshal(r) // marshal to json string
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Error("Marshal failed:")
		// todo: remove table
		return "", err
	}
	return string(bytes), nil
}

// Task for gocron scheduler
func (b *mbtcpService) Task(socket *zmq.Socket, req interface{}) {
	str, err := Marshal(req)
	if err != nil {
		// todo: remove table
		return
	}
	log.WithFields(log.Fields{"JSON": str}).Debug("Send request to modbusd:")
	socket.Send("tcp", zmq.SNDMORE) // frame 1
	socket.Send(str, 0)             // convert to string; frame 2
}

// initZMQPub init ZMQ publishers.
// Example:
// 	initZMQPub("ipc:///tmp/from.psmb", "ipc:///tmp/to.modbus")
func (b *mbtcpService) initZMQPub(toServiceEndpoint, toModbusdEndpoint string) {
	log.Debug("Init ZMQ Publishers")
	// upstream publisher
	b.pub.upstream, _ = zmq.NewSocket(zmq.PUB)
	b.pub.upstream.Bind(toServiceEndpoint)

	// downstream publisher
	b.pub.downstream, _ = zmq.NewSocket(zmq.PUB)
	b.pub.downstream.Connect(toModbusdEndpoint)
}

// initZMQSub init ZMQ subscribers.
// Example:
// 	initZMQSub("ipc:///tmp/to.psmb", "ipc:///tmp/from.modbus")
func (b *mbtcpService) initZMQSub(fromServiceEndpoint, fromModbusdEndpoint string) {
	log.Debug("Init ZMQ Subscribers")
	// upstream subscriber
	b.sub.upstream, _ = zmq.NewSocket(zmq.SUB)
	b.sub.upstream.Bind(fromServiceEndpoint)
	b.sub.upstream.SetSubscribe("")

	// downstream subscriber
	b.sub.downstream, _ = zmq.NewSocket(zmq.SUB)
	b.sub.downstream.Connect(fromModbusdEndpoint)
	b.sub.downstream.SetSubscribe("")
}

// initZMQPoller init ZMQ poller.
// polling from upstream services and downstream modbusd
func (b *mbtcpService) initZMQPoller() {
	log.Debug("Init ZMQ Poller")
	// initialize poll set
	b.poller = zmq.NewPoller()
	b.poller.Add(b.sub.upstream, zmq.POLLIN)
	b.poller.Add(b.sub.downstream, zmq.POLLIN)
}

// simpleTaskResponser simeple response to upstream
func (b *mbtcpService) simpleTaskResponser(tid string, resp interface{}) error {
	respStr, err := Marshal(resp)
	if err != nil {
		return err
	}

	// check simple task map
	if cmd, ok := b.simpleTaskMap.Get(tid); ok {
		log.WithFields(log.Fields{"JSON": respStr}).Debug("Send response to service:")
		b.pub.upstream.Send(cmd, zmq.SNDMORE) // task command
		b.pub.upstream.Send(respStr, 0)       // convert to string; frame 2
		// remove from map
		b.simpleTaskMap.Delete(tid)
		return nil
	}
	return errors.New("Request not found!")
}

// simpleResponser simple reponser to upstream without checking simple task map
func (b *mbtcpService) simpleResponser(cmd string, resp interface{}) error {
	respStr, err := Marshal(resp)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{"JSON": respStr}).Debug("Send response to service:")
	b.pub.upstream.Send(cmd, zmq.SNDMORE) // task command
	b.pub.upstream.Send(respStr, 0)       // convert to string; frame 2
	return nil
}

// parseRequest parse requests from services
// R&R: only unmarshal request string to corresponding struct
func (b *mbtcpService) parseRequest(msg []string) (interface{}, error) {
	// Check the length of multi-part message
	if len(msg) != 2 {
		// should not reach here!!
		log.Error("Request parser failed: invalid message length")
		return nil, errors.New("Invalid message length")
	}

	log.WithFields(log.Fields{"msg[0]": msg[0]}).Debug("Parsing upstream request:")

	switch msg[0] {
	case mbtcpOnceRead: // done
		var req MbtcpReadReq
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal request failed:")
			return nil, err
		}
		return req, nil
	case mbtcpOnceWrite: // done
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

			var l uint16 // length of registers

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
	case mbtcpGetTimeout, mbtcpSetTimeout: // done
		var req MbtcpTimeoutReq
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal request failed:")
			return nil, err
		}
		return req, nil
	case mbtcpCreatePoll: // done
		var req MbtcpPollStatus
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal request failed:")
			return nil, err
		}
		return req, nil
	case mbtcpUpdatePoll, mbtcpGetPoll, mbtcpDeletePoll, // done
		mbtcpTogglePoll, mbtcpGetPolls, mbtcpDeletePolls,
		mbtcpTogglePolls, mbtcpGetPollHistory, mbtcpExportPolls:
		var req MbtcpPollOpReq
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Unmarshal request failed:")
			return nil, err
		}
		return req, nil
	case mbtcpImportPolls: // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case mbtcpCreateFilter: // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case mbtcpUpdateFilter: // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case mbtcpGetFilter: // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case mbtcpDeleteFilter: // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case mbtcpToggleFilter: // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case mbtcpGetFilters: // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case mbtcpDeleteFilters: // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case mbtcpToggleFilters: // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case mbtcpImportFilters: // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	case mbtcpExportFilters: // todo
		log.Warn("TODO")
		return nil, errors.New("TODO")
	default: // done
		// should not reach here!!
		log.WithFields(log.Fields{"request": msg[0]}).Warn("Request not support:")
		return nil, errors.New("Request not support")
	}
}

// handleRequest handle requests from services
func (b *mbtcpService) handleRequest(cmd string, r interface{}) error {
	log.WithFields(log.Fields{"cmd": cmd}).Debug("Handle upstream request:")

	switch cmd {
	case mbtcpOnceRead: // done
		req := r.(MbtcpReadReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string

		// default port checker
		if req.Port == "" {
			req.Port = defaultMbtcpPort
		}

		// add to read/poll task map
		b.readTaskMap.Add("", TidStr, cmd, req)

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
		b.scheduler.Emergency().Do(b.Task, b.pub.downstream, command)
		return nil
	case mbtcpOnceWrite: // done
		req := r.(MbtcpWriteReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string
		// default port checker
		if req.Port == "" {
			req.Port = defaultMbtcpPort
		}

		// add to simple task map
		b.simpleTaskMap.Add(TidStr, cmd)

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
		b.scheduler.Emergency().Do(b.Task, b.pub.downstream, command)
		return nil
	case mbtcpGetTimeout: // done
		req := r.(MbtcpTimeoutReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string
		b.simpleTaskMap.Add(TidStr, cmd)         // add to simple task map

		cmdInt, _ := strconv.Atoi(string(getTCPTimeout))
		command := DMbtcpTimeout{
			Tid: TidStr,
			Cmd: cmdInt,
		}
		// add command to scheduler as emergency request
		b.scheduler.Emergency().Do(b.Task, b.pub.downstream, command)
		return nil
	case mbtcpSetTimeout: // done
		req := r.(MbtcpTimeoutReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string
		// protect dummy input
		if req.Data < minMbtcpTimeout {
			req.Data = minMbtcpTimeout
		}
		b.simpleTaskMap.Add(TidStr, cmd) // add to simple task map

		cmdInt, _ := strconv.Atoi(string(setTCPTimeout))
		command := DMbtcpTimeout{
			Tid:     TidStr,
			Cmd:     cmdInt,
			Timeout: req.Data,
		}
		// add command to scheduler as emergency request
		b.scheduler.Emergency().Do(b.Task, b.pub.downstream, command)
		return nil
	case mbtcpCreatePoll:

		req := r.(MbtcpPollStatus)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string
		b.simpleTaskMap.Add(TidStr, cmd)         // simple request

		// check fc
		if req.FC < 1 || req.FC > 4 {
			log.WithFields(log.Fields{"FC": req.FC}).Error("Invalid function code")
			// send back
			resp := MbtcpSimpleRes{Tid: req.Tid, Status: "Invalid function code"}
			return b.simpleTaskResponser(TidStr, resp)
		}

		// check name
		if req.Name == "" {
			log.WithFields(log.Fields{"Name": req.Name}).Error("Invalid poll name")
			// send back
			resp := MbtcpSimpleRes{Tid: req.Tid, Status: "Invalid poll name"}
			return b.simpleTaskResponser(TidStr, resp)
		}

		// todo:check name in map or not?

		// check interval value
		if req.Interval < minMbtcpPollInterval {
			req.Interval = minMbtcpPollInterval
		}
		// check port
		if req.Port == "" {
			req.Port = defaultMbtcpPort
		}

		b.readTaskMap.Add(req.Name, TidStr, cmd, req) // read request

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
		b.scheduler.EveryWithName(req.Interval, req.Name).Seconds().Do(b.Task, b.pub.downstream, command)
		if !req.Enabled {
			b.scheduler.PauseWithName(req.Name)
		}
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: "ok"}
		return b.simpleTaskResponser(TidStr, resp)
	case mbtcpUpdatePoll:
		log.Warn("TODO")
		return errors.New("TODO")
	case mbtcpGetPoll:
		log.Warn("TODO")
		return errors.New("TODO")
	case mbtcpDeletePoll:
		log.Warn("TODO")
		return errors.New("TODO")
	case mbtcpTogglePoll:
		// TODO! check name
		req := r.(MbtcpPollOpReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string
		b.simpleTaskMap.Add(TidStr, cmd)         // add to simple task map

		status := "ok"
		if req.Enabled {
			if ok := b.scheduler.ResumeWithName(req.Name); !ok {
				status = "enable poll failed"
			}
		} else {
			if ok := b.scheduler.PauseWithName(req.Name); !ok {
				status = "disable poll failed"
			}
		}
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: status}
		return b.simpleTaskResponser(TidStr, resp)
	case mbtcpGetPolls:
		log.Warn("TODO")
		return errors.New("TODO")
	case mbtcpDeletePolls:
		log.Warn("TODO")
		return errors.New("TODO")
	case mbtcpTogglePolls:
		log.Warn("TODO")
		return errors.New("TODO")
	case mbtcpImportPolls:
		log.Warn("TODO")
		return errors.New("TODO")
	case mbtcpExportPolls:
		log.Warn("TODO")
		return errors.New("TODO")
	case mbtcpGetPollHistory:
		log.Warn("TODO")
		return errors.New("TODO")
	case mbtcpCreateFilter:
		log.Warn("TODO")
		return errors.New("TODO")
	case mbtcpUpdateFilter:
		log.Warn("TODO")
		return errors.New("TODO")
	case mbtcpGetFilter:
		log.Warn("TODO")
		return errors.New("TODO")
	case mbtcpDeleteFilter:
		log.Warn("TODO")
		return errors.New("TODO")
	case mbtcpToggleFilter:
		log.Warn("TODO")
		return errors.New("TODO")
	case mbtcpGetFilters:
		log.Warn("TODO")
		return errors.New("TODO")
	case mbtcpDeleteFilters:
		log.Warn("TODO")
		return errors.New("TODO")
	case mbtcpToggleFilters:
		log.Warn("TODO")
		return errors.New("TODO")
	case mbtcpImportFilters:
		log.Warn("TODO")
		return errors.New("TODO")
	case mbtcpExportFilters:
		log.Warn("TODO")
		return errors.New("TODO")
	default:
		// should not reach here!!
		log.WithFields(log.Fields{"cmd": cmd}).Warn("Request not support:")
		return errors.New("Request not support")
	}
}

// parseResponse parse responses from modbusd
func (b *mbtcpService) parseResponse(msg []string) (interface{}, error) { // done.
	// Check the length of multi-part message
	if len(msg) != 2 {
		log.Error("Request parser failed: invalid message length")
		return nil, errors.New("Invalid message length")
	}

	log.WithFields(log.Fields{"msg[0]": msg[0]}).Debug("Parsing downstream response:")

	switch MbtcpCmdType(msg[0]) {
	case setTCPTimeout, getTCPTimeout: // set|get timeout
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

// handleResponse handle response from modbusd, Todo: filter, handle
func (b *mbtcpService) handleResponse(cmd string, r interface{}) error {
	log.WithFields(log.Fields{"cmd": cmd}).Debug("Handle downstream response:")

	switch MbtcpCmdType(cmd) {
	case fc5, fc6, fc15, fc16, setTCPTimeout, getTCPTimeout: // done: one-off requests
		var TidStr string
		var resp interface{}

		switch MbtcpCmdType(cmd) {
		case setTCPTimeout: // one-off timeout requests
			res := r.(DMbtcpTimeout)
			tid, _ := strconv.ParseInt(res.Tid, 10, 64)
			TidStr = res.Tid
			resp = MbtcpTimeoutRes{
				Tid:    tid,
				Status: res.Status,
			}
		case getTCPTimeout: // one-off timeout requests
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
		return b.simpleTaskResponser(TidStr, resp)

	case fc1, fc2, fc3, fc4: // one-off and polling requests
		res := r.(DMbtcpRes)
		tid, _ := strconv.ParseInt(res.Tid, 10, 64)

		var response interface{}
		var task mbtcpReadTask
		var ok bool
		// check read task table
		if task, ok = b.readTaskMap.Get(res.Tid); !ok {
			log.Error("request not found")
			return errors.New("request not found")
		}

		respCmd := task.Cmd

		switch MbtcpCmdType(cmd) {
		case fc1, fc2: // done: todo: enhancement
			switch task.Cmd {
			case mbtcpOnceRead:
				if res.Status != "ok" {
					response = MbtcpSimpleRes{
						Tid:    tid,
						Status: res.Status,
					}
				} else {
					response = MbtcpReadRes{
						Tid:    tid,
						Status: res.Status,
						Data:   res.Data,
					}
				}
				// remove from read/poll table
				b.readTaskMap.Delete(res.Tid)
			case mbtcpCreatePoll, mbtcpImportPolls: // data
				respCmd = mbtcpData // set as "mbtcp.data"
				if res.Status != "ok" {
					response = MbtcpPollData{
						TimeStamp: time.Now().UTC().UnixNano(),
						Name:      task.Name,
						Status:    res.Status,
					}
				} else {
					// TODO: add to history
					response = MbtcpPollData{
						TimeStamp: time.Now().UTC().UnixNano(),
						Name:      task.Name,
						Status:    res.Status,
						Data:      res.Data,
					}
				}

			default: // should not reach here
				log.Error("Should not reach here")
				response = MbtcpSimpleRes{
					Tid:    tid,
					Status: "Command not support",
				}
			}
			return b.simpleResponser(respCmd, response)
		case fc3, fc4:

			readReq := task.Req.(MbtcpReadReq)

			switch task.Cmd {
			case mbtcpOnceRead:
				// check modbus response status
				if res.Status != "ok" {
					response = MbtcpReadRes{
						Tid:    tid,
						Type:   readReq.Type,
						Status: res.Status,
					}
					// remove from read table
					b.readTaskMap.Delete(res.Tid)
					return b.simpleResponser(respCmd, response)
				}

				// convert register to byte array
				bytes, err := RegistersToBytes(res.Data)
				if err != nil {
					log.Error(err)
					response = MbtcpReadRes{
						Tid:    tid,
						Type:   readReq.Type,
						Status: err.Error(),
					}
					// remove from read table
					b.readTaskMap.Delete(res.Tid)
					return b.simpleResponser(respCmd, response)
				}

				log.WithFields(log.Fields{"Type": readReq.Type}).Debug("Request type:")

				switch readReq.Type {
				case HexString:
					response = MbtcpReadRes{
						Tid:    tid,
						Status: res.Status,
						Type:   readReq.Type,
						Bytes:  bytes,
						Data:   BytesToHexString(bytes), // convert byte to hex
					}

				case UInt16:
					// order
					ret, err := BytesToUInt16s(bytes, readReq.Order)
					if err != nil {
						response = MbtcpReadRes{
							Tid:    tid,
							Type:   readReq.Type,
							Bytes:  bytes,
							Status: err.Error(),
						}
					} else {
						response = MbtcpReadRes{
							Tid:    tid,
							Status: res.Status,
							Type:   readReq.Type,
							Bytes:  bytes,
							Data:   ret,
						}
					}
				case Int16:
					// order
					ret, err := BytesToInt16s(bytes, readReq.Order)
					if err != nil {
						response = MbtcpReadRes{
							Tid:    tid,
							Type:   readReq.Type,
							Bytes:  bytes,
							Status: err.Error(),
						}
					} else {
						response = MbtcpReadRes{
							Tid:    tid,
							Type:   readReq.Type,
							Bytes:  bytes,
							Data:   ret,
							Status: res.Status,
						}
					}

				case Scale, UInt32, Int32, Float32: // 32-bits
					if readReq.Len%2 != 0 {
						response = MbtcpReadRes{
							Tid:    tid,
							Type:   readReq.Type,
							Bytes:  bytes,
							Status: "Conversion failed",
						}
					} else {
						switch readReq.Type {
						case Scale:
							// todo: check range values
							f := LinearScalingRegisters(
								res.Data,
								readReq.Range.DomainLow,
								readReq.Range.DomainHigh,
								readReq.Range.RangeLow,
								readReq.Range.RangeHigh)

							response = MbtcpReadRes{
								Tid:    tid,
								Type:   readReq.Type,
								Bytes:  bytes,
								Data:   f,
								Status: res.Status,
							}
						case UInt32:
							ret, err := BytesToUInt32s(bytes, readReq.Order)
							if err != nil {
								response = MbtcpReadRes{
									Tid:    tid,
									Type:   readReq.Type,
									Bytes:  bytes,
									Status: err.Error(),
								}
							} else {
								response = MbtcpReadRes{
									Tid:    tid,
									Type:   readReq.Type,
									Bytes:  bytes,
									Data:   ret,
									Status: res.Status,
								}
							}
						case Int32:
							ret, err := BytesToInt32s(bytes, readReq.Order)
							if err != nil {
								response = MbtcpReadRes{
									Tid:    tid,
									Type:   readReq.Type,
									Bytes:  bytes,
									Status: err.Error(),
								}
							} else {
								response = MbtcpReadRes{
									Tid:    tid,
									Status: res.Status,
									Type:   readReq.Type,
									Bytes:  bytes,
									Data:   ret,
								}
							}
						case Float32:
							ret, err := BytesToFloat32s(bytes, readReq.Order)
							if err != nil {
								response = MbtcpReadRes{
									Tid:    tid,
									Type:   readReq.Type,
									Bytes:  bytes,
									Status: err.Error(),
								}
							} else {
								response = MbtcpReadRes{
									Tid:    tid,
									Status: res.Status,
									Type:   readReq.Type,
									Bytes:  bytes,
									Data:   ret,
								}
							}

						}
					}
				default: // case 0, 1(RegisterArray)
					response = MbtcpReadRes{
						Tid:    tid,
						Status: res.Status,
						Type:   readReq.Type,
						Bytes:  bytes,
						Data:   res.Data,
					}
				}
				// remove from read table
				b.readTaskMap.Delete(res.Tid)
				return b.simpleResponser(respCmd, response)
			case mbtcpCreatePoll, mbtcpImportPolls: // data
				respCmd = mbtcpData // set as "mbtcp.data"

				if res.Status != "ok" {
					response = MbtcpPollData{
						TimeStamp: time.Now().UTC().UnixNano(),
						Name:      task.Name,
						Type:      readReq.Type,
						Status:    res.Status,
					}
					return b.simpleResponser(respCmd, response)
				}

				// convert register to byte array
				bytes, err := RegistersToBytes(res.Data)
				if err != nil {
					response = MbtcpPollData{
						TimeStamp: time.Now().UTC().UnixNano(),
						Name:      task.Name,
						Type:      readReq.Type,
						Status:    err.Error(),
					}
					return b.simpleResponser(respCmd, response)
				}

				log.WithFields(log.Fields{"Type": readReq.Type}).Debug("Request type:")

				/*
					response = MbtcpPollData{
						TimeStamp: time.Now().UTC().UnixNano(),
						Name:      task.Name,
						Type:      readReq.Type,
						Bytes:     bytes,
					}
				*/
				switch readReq.Type {
				case HexString:
					response = MbtcpPollData{
						TimeStamp: time.Now().UTC().UnixNano(),
						Name:      task.Name,
						Type:      readReq.Type,
						Bytes:     bytes,
						Data:      BytesToHexString(bytes), // convert byte to hex
						Status:    res.Status,
					}
				case UInt16:
					// order
					ret, err := BytesToUInt16s(bytes, readReq.Order)
					if err != nil {
						response = MbtcpPollData{
							TimeStamp: time.Now().UTC().UnixNano(),
							Name:      task.Name,
							Type:      readReq.Type,
							Bytes:     bytes,
							// Data
							Status: err.Error(),
						}
					} else {
						response = MbtcpPollData{
							TimeStamp: time.Now().UTC().UnixNano(),
							Name:      task.Name,
							Type:      readReq.Type,
							Bytes:     bytes,
							Data:      ret,
							Status:    res.Status,
						}
					}
				case Int16:
					// order
					ret, err := BytesToInt16s(bytes, readReq.Order)
					if err != nil {
						response = MbtcpPollData{
							TimeStamp: time.Now().UTC().UnixNano(),
							Name:      task.Name,
							Type:      readReq.Type,
							Bytes:     bytes,
							// Data
							Status: err.Error(),
						}
					} else {
						response = MbtcpPollData{
							TimeStamp: time.Now().UTC().UnixNano(),
							Name:      task.Name,
							Type:      readReq.Type,
							Bytes:     bytes,
							Data:      ret,
							Status:    res.Status,
						}
					}

				case Scale, UInt32, Int32, Float32: // 32-Bits
					if readReq.Len%2 != 0 {
						response = MbtcpPollData{
							TimeStamp: time.Now().UTC().UnixNano(),
							Name:      task.Name,
							Type:      readReq.Type,
							Bytes:     bytes,
							// Data
							Status: "Conversion failed",
						}
					} else {
						switch readReq.Type {
						case Scale:
							// todo: check range values
							f := LinearScalingRegisters(
								res.Data,
								readReq.Range.DomainLow,
								readReq.Range.DomainHigh,
								readReq.Range.RangeLow,
								readReq.Range.RangeHigh)

							response = MbtcpPollData{
								TimeStamp: time.Now().UTC().UnixNano(),
								Name:      task.Name,
								Type:      readReq.Type,
								Bytes:     bytes,
								Data:      f,
								Status:    res.Status,
							}
						case UInt32:
							ret, err := BytesToUInt32s(bytes, readReq.Order)
							if err != nil {
								response = MbtcpPollData{
									TimeStamp: time.Now().UTC().UnixNano(),
									Name:      task.Name,
									Type:      readReq.Type,
									Bytes:     bytes,
									// Data
									Status: err.Error(),
								}
							} else {
								response = MbtcpPollData{
									TimeStamp: time.Now().UTC().UnixNano(),
									Name:      task.Name,
									Type:      readReq.Type,
									Bytes:     bytes,
									Data:      ret,
									Status:    res.Status,
								}
							}
						case Int32:
							ret, err := BytesToInt32s(bytes, readReq.Order)
							if err != nil {
								response = MbtcpPollData{
									TimeStamp: time.Now().UTC().UnixNano(),
									Name:      task.Name,
									Type:      readReq.Type,
									Bytes:     bytes,
									// Data
									Status: err.Error(),
								}
							} else {
								response = MbtcpPollData{
									TimeStamp: time.Now().UTC().UnixNano(),
									Name:      task.Name,
									Type:      readReq.Type,
									Bytes:     bytes,
									Data:      ret,
									Status:    res.Status,
								}
							}
						case Float32:
							ret, err := BytesToFloat32s(bytes, readReq.Order)
							if err != nil {
								response = MbtcpPollData{
									TimeStamp: time.Now().UTC().UnixNano(),
									Name:      task.Name,
									Type:      readReq.Type,
									Bytes:     bytes,
									// Data
									Status: err.Error(),
								}
							} else {
								response = MbtcpPollData{
									TimeStamp: time.Now().UTC().UnixNano(),
									Name:      task.Name,
									Type:      readReq.Type,
									Bytes:     bytes,
									Data:      ret,
									Status:    res.Status,
								}
							}
						}
					}
				default: // case 0, 1(RegisterArray)
					response = MbtcpPollData{
						TimeStamp: time.Now().UTC().UnixNano(),
						Name:      task.Name,
						Type:      readReq.Type,
						Bytes:     bytes,
						Data:      res.Data,
						Status:    res.Status,
					}
				}

				// TODO: add to history

				return b.simpleResponser(respCmd, response)
			default:
				log.Error("Should not reach here")
				response = MbtcpSimpleRes{
					Tid:    tid,
					Status: "not support command",
				}
				return b.simpleResponser(respCmd, response)
			}
		}
	default: // should not reach here!!
		log.WithFields(log.Fields{"cmd": cmd}).Warn("Response not support:")
		return errors.New("Response not support")
	}
	return nil
}

// Start enable proactive service
func (b *mbtcpService) Start() {
	b.initLogger()

	log.Debug("Start proactive service")
	b.scheduler.Start()
	b.initZMQPub("ipc:///tmp/from.psmb", "ipc:///tmp/to.modbus")
	b.initZMQSub("ipc:///tmp/to.psmb", "ipc:///tmp/from.modbus")
	b.initZMQPoller()

	// process messages from both subscriber sockets
	for b.enable {
		sockets, _ := b.poller.Poll(-1)
		for _, socket := range sockets {
			switch s := socket.Socket; s {
			case b.sub.upstream:
				// receive from upstream
				msg, _ := b.sub.upstream.RecvMessage(0)
				log.WithFields(log.Fields{
					"msg[0]": msg[0],
					"msg[1]": msg[1],
				}).Debug("Receive from service:")

				// parse request
				req, err := b.parseRequest(msg)
				if err != nil {
					// todo: send error back
				} else {
					err = b.handleRequest(msg[0], req)
				}
			case b.sub.downstream:
				// receive from modbusd
				msg, _ := b.sub.downstream.RecvMessage(0)
				log.WithFields(log.Fields{
					"msg[0]": msg[0],
					"msg[1]": msg[1],
				}).Debug("Receive from modbusd:")

				// parse response
				res, err := b.parseResponse(msg)
				if err != nil {
					// todo: send error back
				} else {
					err = b.handleResponse(msg[0], res)
				}
			}
		}
	}
}

// Stop disable proactive service
func (b *mbtcpService) Stop() {
	log.Debug("Stop proactive service")
	b.scheduler.Stop()
	b.enable = false
	if b.sub.upstream {
		b.sub.upstream.Close()
		b.pub.upstream.Close()
		b.sub.downstream.Close()
		b.pub.downstream.Close()
	}
}
