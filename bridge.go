package main

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/taka-wang/gocron"
	log "github.com/takawang/logrus"
	zmq "github.com/takawang/zmq3"
)

// Bridge proactive service
type Bridge interface {
	initPub(serviceEndpoint, modbusdEndpoint string)
	initSub(serviceEndpoint, modbusdEndpoint string)
	initPoller()
	Start()
	Stop()
}

// BridgeType proactive service type
type mbtcpBridge struct {
	enable        bool
	readTaskMap   MbtcpReadTask
	simpleTaskMap MbtcpSimpleTask
	scheduler     gocron.Scheduler
	fromService   *zmq.Socket
	toService     *zmq.Socket
	fromModbusd   *zmq.Socket
	toModbusd     *zmq.Socket
	poller        *zmq.Poller
}

// NewMbtcpBridge init bridge
func NewMbtcpBridge() Bridge {
	return &mbtcpBridge{
		enable:        true,
		readTaskMap:   NewMbtcpReadTask(),
		simpleTaskMap: NewMbtcpSimpleTask(),
		scheduler:     gocron.NewScheduler(),
	}
}

// initPub init zmq publisher
// ex. initPub("ipc:///tmp/from.psmb", "ipc:///tmp/to.modbus")
func (b *mbtcpBridge) initPub(serviceEndpoint, modbusdEndpoint string) {
	// upstream publisher
	b.toService, _ = zmq.NewSocket(zmq.PUB)
	b.toService.Bind(serviceEndpoint)

	// downstream publisher
	b.toModbusd, _ = zmq.NewSocket(zmq.PUB)
	b.toModbusd.Connect(modbusdEndpoint)
}

// initSub init zmq subscriber
// ex. initSub("ipc:///tmp/to.psmb", "ipc:///tmp/from.modbus")
func (b *mbtcpBridge) initSub(serviceEndpoint, modbusdEndpoint string) {
	// upstream subscriber
	b.fromService, _ = zmq.NewSocket(zmq.SUB)
	b.fromService.Bind(serviceEndpoint)
	b.fromService.SetSubscribe("")

	// downstream subscriber
	b.fromModbusd, _ = zmq.NewSocket(zmq.SUB)
	b.fromModbusd.Connect(modbusdEndpoint)
	b.fromModbusd.SetSubscribe("")
}

// initPoller init zmq poller
func (b *mbtcpBridge) initPoller() {
	// initialize poll set
	b.poller = zmq.NewPoller()
	b.poller.Add(b.fromService, zmq.POLLIN)
	b.poller.Add(b.fromModbusd, zmq.POLLIN)
}

func (b *mbtcpBridge) Start() {
	b.scheduler.Start()
	b.initPub("ipc:///tmp/from.psmb", "ipc:///tmp/to.modbus")
	b.initSub("ipc:///tmp/to.psmb", "ipc:///tmp/from.modbus")
	b.initPoller()

	// process messages from both subscriber sockets
	for b.enable {
		sockets, _ := b.poller.Poll(-1)
		for _, socket := range sockets {
			switch s := socket.Socket; s {
			case b.fromService:
				// receive from upstream
				msg, _ := b.fromService.RecvMessage(0)
				log.WithFields(log.Fields{
					"msg[0]": msg[0],
					"msg[1]": msg[1],
				}).Debug("Receive from service:")

				// parse request
				req, err := b.ParseRequest(msg)
				if err != nil {
					// todo: send error back
				} else {
					err = RequestHandler(msg[0], req, b.toModbusd, b.toService)
				}
			case b.fromModbusd:
				// receive from modbusd
				msg, _ := b.fromModbusd.RecvMessage(0)
				log.WithFields(log.Fields{
					"msg[0]": msg[0],
					"msg[1]": msg[1],
				}).Debug("Receive from modbusd:")

				// parse response
				res, err := ParseResponse(msg)
				if err != nil {
					// todo: send error back
				} else {
					err = ResponseHandler(MbtcpCmdType(msg[0]), res, b.toService)
				}
			}
		}
	}
}

func (b *mbtcpBridge) Stop() {
	b.scheduler.Stop()
	b.enable = false
	b.fromService.Close()
	b.toService.Close()
	b.fromModbusd.Close()
	b.toModbusd.Close()
}

func (b *mbtcpBridge) SimpleTaskResponser(tid string, resp interface{}) error {
	respStr, err := json.Marshal(resp)
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Error("Marshal failed:")
		return err
	}

	if cmd, ok := b.simpleTaskMap.Get(tid); ok {
		log.WithFields(log.Fields{"JSON": string(respStr)}).Debug("Send response to service:")
		b.toService.Send(cmd, zmq.SNDMORE)   // task command
		b.toService.Send(string(respStr), 0) // convert to string; frame 2
		b.simpleTaskMap.Delete(tid)          // remove from Map!!
		return nil
	}
	return errors.New("Request command not in map")
}

// ParseRequest parse message from services
// R&R: only unmarshal request string to corresponding struct
func (b *mbtcpBridge) ParseRequest(msg []string) (interface{}, error) {
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
