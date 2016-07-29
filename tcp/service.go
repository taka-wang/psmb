/*
Package psmbtcp provides a proactive service library for modbus daemon.
*/
package psmbtcp

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/taka-wang/gocron"
	. "github.com/taka-wang/psmb"
	log "github.com/takawang/logrus"
	zmq "github.com/takawang/zmq3"
)

// @Implement ProactiveService contract implicitly

const (
	// defaultMbtcpPort default modbus slave port number
	defaultMbtcpPort = "502"
	// minMbtcpTimeout minimal modbus tcp connection timeout
	minMbtcpTimeout = 200000
	// minMbtcpPollInterval minimal modbus tcp poll interval
	minMbtcpPollInterval = 1
)

// MbtcpService modbusd tcp proactive service type
type MbtcpService struct {
	// readerMap read/poll task map
	readerMap ReaderTaskMap
	// writerMap write task map
	writerMap WriterTaskMap
	// scheduler gocron scheduler
	scheduler gocron.Scheduler
	// sub ZMQ subscriber endpoints
	sub struct {
		// upstream subscriber from services
		upstream *zmq.Socket
		// downstream subscriber from modbusd
		downstream *zmq.Socket
	}
	// pub ZMQ publisher endpoints
	pub struct {
		// upstream publisher to services
		upstream *zmq.Socket
		// downstream publisher to modbusd
		downstream *zmq.Socket
	}
	// poller ZMQ poller
	poller *zmq.Poller
	// enable poller flag
	enable bool
}

// NewPSMBTCP modbus tcp proactive serivce constructor
func NewPSMBTCP(r ReaderTaskMap, w WriterTaskMap, s gocron.Scheduler) ProactiveService {
	return &MbtcpService{
		enable:    true,
		readerMap: r,
		writerMap: w,
		scheduler: s,
	}
}

// initLogger init logger
func (b *MbtcpService) initLogger() {
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
		// TODO: remove table
		return "", ErrMarshal
	}
	return string(bytes), nil
}

// Task for gocron scheduler
func (b *MbtcpService) Task(socket *zmq.Socket, req interface{}) {
	str, err := Marshal(req)
	if err != nil {
		// TODO: remove table
		return
	}
	log.WithFields(log.Fields{"JSON": str}).Debug("Send request to modbusd:")
	socket.Send("tcp", zmq.SNDMORE) // frame 1
	socket.Send(str, 0)             // convert to string; frame 2
}

// initZMQPub init ZMQ publishers.
// Example:
// 	initZMQPub("ipc:///tmp/from.psmb", "ipc:///tmp/to.modbus")
func (b *MbtcpService) initZMQPub(toServiceEndpoint, toModbusdEndpoint string) {
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
func (b *MbtcpService) initZMQSub(fromServiceEndpoint, fromModbusdEndpoint string) {
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
func (b *MbtcpService) initZMQPoller() {
	log.Debug("Init ZMQ Poller")
	// initialize poll set
	b.poller = zmq.NewPoller()
	b.poller.Add(b.sub.upstream, zmq.POLLIN)
	b.poller.Add(b.sub.downstream, zmq.POLLIN)
}

// naiveResponder naive responder to send message back to upstream.
func (b *MbtcpService) naiveResponder(cmd string, resp interface{}) error {
	respStr, err := Marshal(resp)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Fail to marshal for naive responder!")
		return err
	}

	log.WithFields(log.Fields{"JSON": respStr}).Debug("Send message to services")
	b.pub.upstream.Send(cmd, zmq.SNDMORE) // task command
	b.pub.upstream.Send(respStr, 0)       // convert to string; frame 2
	return nil
}

// ParseRequest parse requests from services,
// only unmarshal request string to corresponding struct
func (b *MbtcpService) ParseRequest(msg []string) (interface{}, error) {
	// Check the length of multi-part message
	if len(msg) != 2 {
		return nil, ErrInvalidMessageLength
	}

	log.WithFields(log.Fields{"msg[0]": msg[0]}).Debug("Parse request from upstream services")

	switch msg[0] {
	case mbtcpGetTimeout, mbtcpSetTimeout:
		var req MbtcpTimeoutReq
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			return nil, ErrUnmarshal
		}
		return req, nil
	case mbtcpOnceWrite:
		// unmarshal partial request (except data field)
		var data json.RawMessage // raw []byte
		req := MbtcpWriteReq{Data: &data}
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			return nil, ErrUnmarshal
		}

		// unmarshal data field
		switch MbtcpCmdType(strconv.Itoa(req.FC)) {
		case fc5: // write single bit; uint16
			var d uint16
			if err := json.Unmarshal(data, &d); err != nil {
				return nil, ErrUnmarshal
			}
			req.Data = d // unmarshal to uint16
			return req, nil
		case fc15: // write multiple bits; []uint16
			var d []uint16
			if err := json.Unmarshal(data, &d); err != nil {
				return nil, ErrUnmarshal
			}
			req.Data = d // unmarshal to uint16 array
			return req, nil
		case fc6: // write single register in dec|hex
			var d string
			if err := json.Unmarshal(data, &d); err != nil { // unmarshal to string
				return nil, ErrUnmarshal
			}
			var dd []uint16
			var err error
			if req.Hex { // check dec or hex
				dd, err = HexStringToRegisters(d)
			} else {
				dd, err = DecimalStringToRegisters(d)
			}
			if err != nil {
				return nil, err
			}
			req.Data = dd[0] // retrieve only one register
			return req, nil
		case fc16: // write multiple register in dec/hex
			var d string
			if err := json.Unmarshal(data, &d); err != nil { // unmarshal to string
				return nil, ErrUnmarshal
			}
			var dd []uint16
			var err error
			if req.Hex { // check dec or hex
				dd, err = HexStringToRegisters(d)
			} else {
				dd, err = DecimalStringToRegisters(d)
			}
			if err != nil {
				return nil, err
			}
			req.Data = dd
			return req, nil
		default: // should not reach here
			return nil, ErrInvalidFunctionCode
		}
	case mbtcpOnceRead:
		var req MbtcpReadReq
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			return nil, ErrUnmarshal
		}
		return req, nil
	case mbtcpCreatePoll:
		var req MbtcpPollStatus
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			return nil, ErrUnmarshal
		}
		return req, nil
	case mbtcpUpdatePoll, mbtcpGetPoll, mbtcpDeletePoll,
		mbtcpTogglePoll, mbtcpGetPolls, mbtcpDeletePolls,
		mbtcpTogglePolls, mbtcpGetPollHistory, mbtcpExportPolls:
		var req MbtcpPollOpReq
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			return nil, ErrUnmarshal
		}
		return req, nil
	case mbtcpImportPolls:
		var req MbtcpPollsStatus
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			return nil, ErrUnmarshal
		}
		return req, nil
	case mbtcpCreateFilter, mbtcpUpdateFilter:
		var req MbtcpFilterStatus
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			return nil, ErrUnmarshal
		}
		return req, nil
	case mbtcpGetFilter, mbtcpDeleteFilter, mbtcpToggleFilter,
		mbtcpGetFilters, mbtcpDeleteFilters, mbtcpToggleFilters, mbtcpExportFilters:
		var req MbtcpFilterOpReq
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			return nil, ErrUnmarshal
		}
		return req, nil
	case mbtcpImportFilters:
		var req MbtcpFiltersStatus
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			return nil, ErrUnmarshal
		}
		return req, nil
	default: // should not reach here!!
		return nil, ErrRequestNotSupport
	}
}

// HandleRequest handle requests from services
// do error checking
func (b *MbtcpService) HandleRequest(cmd string, r interface{}) error {
	log.WithFields(log.Fields{"cmd": cmd}).Debug("Handle request from upstream services")

	switch cmd {
	case mbtcpGetTimeout: // done
		req := r.(MbtcpTimeoutReq)
		TidStr := strconv.FormatInt(req.Tid, 10)         // convert tid to string
		cmdInt, _ := strconv.Atoi(string(getTCPTimeout)) // convert to modbusd command
		command := DMbtcpTimeout{
			Tid: TidStr,
			Cmd: cmdInt,
		}
		// add request to write task map
		b.writerMap.Add(TidStr, cmd)
		// add command to scheduler as emergency request
		b.scheduler.Emergency().Do(b.Task, b.pub.downstream, command)
		return nil
	case mbtcpSetTimeout: // done
		req := r.(MbtcpTimeoutReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string
		cmdInt, _ := strconv.Atoi(string(setTCPTimeout))
		command := DMbtcpTimeout{
			Tid: TidStr,
			Cmd: cmdInt,
		}

		// protect invalid timeout value
		if req.Data < minMbtcpTimeout {
			command.Timeout = minMbtcpTimeout
		} else {
			command.Timeout = req.Data
		}
		// add request to write task map
		b.writerMap.Add(TidStr, cmd)
		// add command to scheduler as emergency request
		b.scheduler.Emergency().Do(b.Task, b.pub.downstream, command)
		return nil
	case mbtcpOnceWrite: // done
		req := r.(MbtcpWriteReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string

		// protect null port
		if req.Port == "" {
			req.Port = defaultMbtcpPort
		}

		// length checker
		switch MbtcpCmdType(strconv.Itoa(req.FC)) {
		case fc15, fc16:
			l := uint16(len(req.Data.([]uint16)))
			if req.Len < l {
				req.Len = l
			}
			// we don't check max length, let modbusd do it.
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
		// add request to write task map
		b.writerMap.Add(TidStr, cmd)
		// add command to scheduler as emergency request
		b.scheduler.Emergency().Do(b.Task, b.pub.downstream, command)
		return nil
	case mbtcpOnceRead: // done
		req := r.(MbtcpReadReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string

		// function code checker
		if req.FC < 1 || req.FC > 4 {
			err := ErrInvalidFunctionCode // invalid read function code
			log.WithFields(log.Fields{"FC": req.FC}).Error(err.Error())
			// send back
			resp := MbtcpSimpleRes{Tid: req.Tid, Status: err.Error()}
			return b.naiveResponder(cmd, resp)
		}

		// protect null port
		if req.Port == "" {
			req.Port = defaultMbtcpPort
		}

		// length checker
		if req.Len < 1 {
			req.Len = 1 // set minimal length of read
		}

		command := DMbtcpReadReq{
			Tid:   TidStr,
			Cmd:   req.FC,
			IP:    req.IP,
			Port:  req.Port,
			Slave: req.Slave,
			Addr:  req.Addr,
			Len:   req.Len,
		}

		// add request to read/poll task map, read task instead of poll task, thus pass null name
		b.readerMap.Add("", TidStr, cmd, req)
		// add command to scheduler as emergency request
		b.scheduler.Emergency().Do(b.Task, b.pub.downstream, command)
		return nil
	case mbtcpCreatePoll: // done
		req := r.(MbtcpPollStatus)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string

		// function code checker
		if req.FC < 1 || req.FC > 4 {
			err := ErrInvalidFunctionCode // invalid read function code
			log.WithFields(log.Fields{"FC": req.FC}).Error(err.Error())
			// send error back
			resp := MbtcpSimpleRes{Tid: req.Tid, Status: err.Error()}
			return b.naiveResponder(cmd, resp)
		}

		// protect null poll name
		if req.Name == "" {
			err := ErrInvalidPollName
			log.WithFields(log.Fields{"Name": req.Name}).Error(err.Error())
			// send back
			resp := MbtcpSimpleRes{Tid: req.Tid, Status: err.Error()}
			return b.naiveResponder(cmd, resp)
		}

		// protect null port
		if req.Port == "" {
			req.Port = defaultMbtcpPort
		}

		// length checker
		if req.Len < 1 {
			req.Len = 1 // set minimal length of read
		}

		// check interval value
		if req.Interval < minMbtcpPollInterval {
			req.Interval = minMbtcpPollInterval
		}

		command := DMbtcpReadReq{
			Tid:   TidStr,
			Cmd:   req.FC,
			IP:    req.IP,
			Port:  req.Port,
			Slave: req.Slave,
			Addr:  req.Addr,
			Len:   req.Len,
		}

		// add task to read/poll task map
		b.readerMap.Add(req.Name, TidStr, cmd, req)
		// add command to scheduler as regular request
		b.scheduler.EveryWithName(req.Interval, req.Name).Seconds().Do(b.Task, b.pub.downstream, command)

		if !req.Enabled { // if not enabled, pause the task
			b.scheduler.PauseWithName(req.Name)
		}
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: "ok"}
		return b.naiveResponder(cmd, resp)
	case mbtcpUpdatePoll: // done
		req := r.(MbtcpPollOpReq)
		status := "ok"

		// check interval value
		if req.Interval < minMbtcpPollInterval {
			req.Interval = minMbtcpPollInterval
		}

		// update task interval
		if ok := b.scheduler.UpdateIntervalWithName(req.Name, req.Interval); !ok {
			err := ErrInvalidPollName // not in scheduler
			log.WithFields(log.Fields{"Name": req.Name}).Error(err.Error())
			status = err.Error() // set error status
		}

		// update read/poll task map
		if status == "ok" {
			if err := b.readerMap.UpdateInterval(req.Name, req.Interval); err != nil {
				log.WithFields(log.Fields{"Name": req.Name}).Error(err.Error())
				status = err.Error() // set error status
			}
		}
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: status}
		return b.naiveResponder(cmd, resp)
	case mbtcpGetPoll: // done
		req := r.(MbtcpPollOpReq)
		t, ok := b.readerMap.GetByName(req.Name)
		task := t.(mbtcpReadTask) // type casting
		if !ok {
			err := ErrInvalidPollName // not in read/poll task map
			log.WithFields(log.Fields{"Name": req.Name}).Error(err.Error())
			// send error back
			resp := MbtcpSimpleRes{Tid: req.Tid, Status: err.Error()}
			return b.naiveResponder(cmd, resp)
		}

		// send back
		request := task.Req.(MbtcpPollStatus)
		resp := MbtcpPollStatus{
			Tid:      request.Tid,
			Name:     request.Name,
			Interval: request.Interval,
			Enabled:  request.Enabled,
			FC:       request.FC,
			IP:       request.IP,
			Port:     request.Port,
			Slave:    request.Slave,
			Addr:     request.Addr,
			Len:      request.Len,
			Type:     request.Type,
			Order:    request.Order,
			Range:    request.Range,
			Status:   "ok",
		}
		return b.naiveResponder(cmd, resp)
	case mbtcpDeletePoll: // done
		req := r.(MbtcpPollOpReq)
		status := "ok"
		// remove task from scheduler
		if ok := b.scheduler.RemoveWithName(req.Name); !ok {
			err := ErrInvalidPollName // not in scheduler
			log.WithFields(log.Fields{"Name": req.Name}).Error(err.Error())
			status = err.Error() // set error status
		}
		// remove task from read/poll map
		b.readerMap.DeleteByName(req.Name)
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: status}
		return b.naiveResponder(cmd, resp)
	case mbtcpTogglePoll: // done
		req := r.(MbtcpPollOpReq)
		status := "ok"
		// update scheduler
		if req.Enabled {
			if ok := b.scheduler.ResumeWithName(req.Name); !ok {
				err := ErrInvalidPollName // not in scheduler
				log.WithFields(log.Fields{"Name": req.Name}).Error(err.Error())
				status = err.Error() // set error status
			}
		} else {
			if ok := b.scheduler.PauseWithName(req.Name); !ok {
				err := ErrInvalidPollName // not in scheduler
				log.WithFields(log.Fields{"Name": req.Name}).Error(err.Error())
				status = err.Error() // set error status
			}
		}
		// update read/poll task map
		if status == "ok" {
			if err := b.readerMap.UpdateToggle(req.Name, req.Enabled); err != nil {
				log.WithFields(log.Fields{"Name": req.Name}).Error(err.Error())
				status = err.Error() // set error status
			}
		}
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: status}
		return b.naiveResponder(cmd, resp)
	case mbtcpGetPolls, mbtcpExportPolls: // done
		req := r.(MbtcpPollOpReq)
		reqs := b.readerMap.GetAll()
		//log.WithFields(log.Fields{"reqs": reqs}).Debug("taka: after GetAll")
		resp := MbtcpPollsStatus{
			Tid:    req.Tid,
			Status: "ok",
			Polls:  reqs,
		}
		// send back
		return b.naiveResponder(cmd, resp)
	case mbtcpDeletePolls: // done
		req := r.(MbtcpPollOpReq)
		b.scheduler.Clear()     // remove all tasks from scheduler
		b.readerMap.DeleteAll() // remove all tasks from read/poll task map
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: "ok"}
		return b.naiveResponder(cmd, resp)
	case mbtcpTogglePolls: // done
		req := r.(MbtcpPollOpReq)
		// update scheduler
		if req.Enabled {
			b.scheduler.ResumeAll()
		} else {
			b.scheduler.PauseAll()
		}
		// update read/poll task map
		b.readerMap.UpdateAllToggles(req.Enabled)
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: "ok"}
		return b.naiveResponder(cmd, resp)
	case mbtcpImportPolls: // done
		request := r.(MbtcpPollsStatus)
		for _, req := range request.Polls {
			// function code checker
			if req.FC < 1 || req.FC > 4 {
				err := ErrInvalidFunctionCode // invalid read function code
				log.WithFields(log.Fields{"FC": req.FC}).Error(err.Error())
				continue // bypass
			}

			// protect null poll name
			if req.Name == "" {
				err := ErrInvalidPollName
				log.WithFields(log.Fields{"Name": req.Name}).Error(err.Error())
				continue // bypass
			}

			// protect null port
			if req.Port == "" {
				req.Port = defaultMbtcpPort
			}

			// length checker
			if req.Len < 1 {
				req.Len = 1 // set minimal length of read
			}

			// check interval value
			if req.Interval < minMbtcpPollInterval {
				req.Interval = minMbtcpPollInterval
			}

			TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string
			command := DMbtcpReadReq{
				Tid:   TidStr,
				Cmd:   req.FC,
				IP:    req.IP,
				Port:  req.Port,
				Slave: req.Slave,
				Addr:  req.Addr,
				Len:   req.Len,
			}

			b.readerMap.Add(req.Name, TidStr, cmd, req)                                                       // Add task to read/poll task map
			b.scheduler.EveryWithName(req.Interval, req.Name).Seconds().Do(b.Task, b.pub.downstream, command) // add command to scheduler as regular request

			if !req.Enabled { // if not enabled, pause the task
				b.scheduler.PauseWithName(req.Name)
			}
		}
		// send back
		resp := MbtcpSimpleRes{Tid: request.Tid, Status: "ok"}
		return b.naiveResponder(cmd, resp)
	case mbtcpGetPollHistory:
		//req := r.(MbtcpPollOpReq)
		return ErrTodo
	case mbtcpCreateFilter, mbtcpUpdateFilter:
		//req := r.(MbtcpFilterStatus)
		return ErrTodo
	case mbtcpGetFilter:
		//req := r.(MbtcpFilterOpReq)
		return ErrTodo
	case mbtcpDeleteFilter:
		//req := r.(MbtcpFilterOpReq)
		return ErrTodo
	case mbtcpToggleFilter:
		//req := r.(MbtcpFilterOpReq)
		return ErrTodo
	case mbtcpGetFilters:
		//req := r.(MbtcpFilterOpReq)
		return ErrTodo
	case mbtcpDeleteFilters:
		//req := r.(MbtcpFilterOpReq)
		return ErrTodo
	case mbtcpToggleFilters:
		//req := r.(MbtcpFilterOpReq)
		return ErrTodo
	case mbtcpImportFilters:
		//request := r.(MbtcpFiltersStatus)
		return ErrTodo
	case mbtcpExportFilters:
		//req := r.(MbtcpFilterOpReq)
		return ErrTodo
	default: // should not reach here!!
		return ErrRequestNotSupport
	}
}

// ParseResponse parse responses from modbusd
// only unmarshal response string to corresponding struct
func (b *MbtcpService) ParseResponse(msg []string) (interface{}, error) {
	// Check the length of multi-part message
	if len(msg) != 2 {
		return nil, ErrInvalidMessageLength
	}

	log.WithFields(log.Fields{"msg[0]": msg[0]}).Debug("Parse response from modbusd")

	switch MbtcpCmdType(msg[0]) {
	case setTCPTimeout, getTCPTimeout:
		var res DMbtcpTimeout
		if err := json.Unmarshal([]byte(msg[1]), &res); err != nil {
			return nil, ErrUnmarshal
		}
		return res, nil
	case fc1, fc2, fc3, fc4, fc5, fc6, fc15, fc16:
		var res DMbtcpRes
		if err := json.Unmarshal([]byte(msg[1]), &res); err != nil {
			return nil, ErrUnmarshal
		}
		return res, nil
	default: // should not reach here!!
		return nil, ErrResponseNotSupport
	}
}

// HandleResponse handle responses from modbusd,
// Todo: filter, handle
func (b *MbtcpService) HandleResponse(cmd string, r interface{}) error {
	log.WithFields(log.Fields{"cmd": cmd}).Debug("Handle response from modbusd")

	switch MbtcpCmdType(cmd) {
	case fc5, fc6, fc15, fc16, setTCPTimeout, getTCPTimeout: // done: one-off requests
		var TidStr string
		var resp interface{}

		switch MbtcpCmdType(cmd) {
		case setTCPTimeout, getTCPTimeout: // one-off timeout requests
			res := r.(DMbtcpTimeout)
			tid, _ := strconv.ParseInt(res.Tid, 10, 64)
			TidStr = res.Tid

			var data int64
			if MbtcpCmdType(cmd) == getTCPTimeout {
				data = res.Timeout
			}

			resp = MbtcpTimeoutRes{
				Tid:    tid,
				Status: res.Status,
				Data:   data, // getTCPTimeout only
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

		// send back one-off task reponse and remove from write task map
		respStr, err := Marshal(resp)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error(err.Error())
			return err
		}

		// check write task map
		if cmd, ok := b.writerMap.Get(TidStr); ok {
			log.WithFields(log.Fields{"JSON": respStr}).Debug("Send message to services")
			b.pub.upstream.Send(cmd, zmq.SNDMORE) // task command
			b.pub.upstream.Send(respStr, 0)       // convert to string; frame 2
			// remove from write task map!
			b.writerMap.Delete(TidStr)
			return nil
		}
		// not in write task map!? should not reach here
		return ErrRequestNotFound
	case fc1, fc2, fc3, fc4: // one-off and polling requests
		res := r.(DMbtcpRes)
		tid, _ := strconv.ParseInt(res.Tid, 10, 64)

		var response interface{}
		var task mbtcpReadTask
		var ok bool
		var t interface{}

		// check read task table
		if t, ok = b.readerMap.GetByTID(res.Tid); !ok {
			return ErrRequestNotFound
		}
		task = t.(mbtcpReadTask) // type casting

		// default response command string
		respCmd := task.Cmd

		switch MbtcpCmdType(cmd) {
		case fc1, fc2: // done: read bits
			var data interface{} // shared variable
			switch task.Cmd {
			case mbtcpOnceRead: // one-off requests
				if res.Status != "ok" {
					data = nil
				} else {
					data = res.Data
				}
				response = MbtcpReadRes{
					Tid:    tid,
					Status: res.Status,
					Data:   data,
				}
				// remove from read/poll table
				b.readerMap.DeleteByTID(res.Tid)
			case mbtcpCreatePoll, mbtcpImportPolls: // poll data
				respCmd = mbtcpData // set as "mbtcp.data"
				if res.Status != "ok" {
					data = nil
				} else {
					data = res.Data
				}
				response = MbtcpPollData{
					TimeStamp: time.Now().UTC().UnixNano(),
					Name:      task.Name,
					Status:    res.Status,
					Data:      data,
				}
				// TODOL if res.Status == "ok" then "add to history"
			default: // should not reach here
				err := ErrResponseNotSupport
				log.WithFields(log.Fields{"cmd": cmd}).Error(err.Error())
				response = MbtcpSimpleRes{
					Tid:    tid,
					Status: err.Error(),
				}
			}
			return b.naiveResponder(respCmd, response)
		case fc3, fc4: // read registers
			// shared variables
			var status string
			var data interface{}

			switch task.Cmd {
			case mbtcpOnceRead: // one-off requests
				readReq := task.Req.(MbtcpReadReq) // type casting

				// check modbus response status
				if res.Status != "ok" {
					response = MbtcpReadRes{
						Tid:    tid,
						Type:   readReq.Type,
						Status: res.Status,
					}
					// remove from read table
					b.readerMap.DeleteByTID(res.Tid)
					return b.naiveResponder(respCmd, response)
				}

				// convert register to byte array
				bytes, err := RegistersToBytes(res.Data)
				if err != nil {
					log.WithFields(log.Fields{"err": err}).Debug("handleResponse: RegistersToBytes failed")
					response = MbtcpReadRes{
						Tid:    tid,
						Type:   readReq.Type,
						Status: err.Error(),
					}
					// remove from read table
					b.readerMap.DeleteByTID(res.Tid)
					return b.naiveResponder(respCmd, response)
				}

				log.WithFields(log.Fields{"Type": readReq.Type}).Debug("Request type:")

				switch readReq.Type {
				case HexString:
					data = BytesToHexString(bytes) // convert byte to hex
					status = res.Status
				case UInt16:
					ret, err := BytesToUInt16s(bytes, readReq.Order) // order
					if err != nil {
						data = nil
						status = err.Error()
					} else {
						data = ret
						status = res.Status
					}
				case Int16:
					ret, err := BytesToInt16s(bytes, readReq.Order) // order
					if err != nil {
						data = nil
						status = err.Error()
					} else {
						data = ret
						status = res.Status
					}
				case Scale, UInt32, Int32, Float32: // 32-bits
					if readReq.Len%2 != 0 {
						err := ErrInvalidLengthToConvert
						data = nil
						status = err.Error()
					} else {
						switch readReq.Type {
						case Scale:
							ret, err := LinearScalingRegisters(
								res.Data,
								readReq.Range.DomainLow,
								readReq.Range.DomainHigh,
								readReq.Range.RangeLow,
								readReq.Range.RangeHigh)
							if err != nil {
								data = nil
								status = err.Error()
							} else {
								data = ret
								status = res.Status
							}
						case UInt32:
							ret, err := BytesToUInt32s(bytes, readReq.Order)
							if err != nil {
								data = nil
								status = err.Error()
							} else {
								data = ret
								status = res.Status
							}
						case Int32:
							ret, err := BytesToInt32s(bytes, readReq.Order)
							if err != nil {
								data = nil
								status = err.Error()
							} else {
								data = ret
								status = res.Status
							}
						case Float32:
							ret, err := BytesToFloat32s(bytes, readReq.Order)
							if err != nil {
								data = nil
								status = err.Error()
							} else {
								data = ret
								status = res.Status
							}
						}
					}
				default: // case 0, 1(RegisterArray)
					data = res.Data
					status = res.Status
				}

				// shared response
				response = MbtcpReadRes{
					Tid:    tid,
					Type:   readReq.Type,
					Bytes:  bytes,
					Data:   data,
					Status: status,
				}

				// remove from read table
				b.readerMap.DeleteByTID(res.Tid)
				return b.naiveResponder(respCmd, response)

			case mbtcpCreatePoll, mbtcpImportPolls: // poll data
				readReq := task.Req.(MbtcpPollStatus) // type casting
				respCmd = mbtcpData                   // set as "mbtcp.data"

				// check modbus response status
				if res.Status != "ok" {
					response = MbtcpPollData{
						TimeStamp: time.Now().UTC().UnixNano(),
						Name:      task.Name,
						Type:      readReq.Type,
						// No Bytes and Data
						Status: res.Status,
					}
					return b.naiveResponder(respCmd, response)
				}

				// convert register to byte array
				bytes, err := RegistersToBytes(res.Data)
				if err != nil {
					response = MbtcpPollData{
						TimeStamp: time.Now().UTC().UnixNano(),
						Name:      task.Name,
						Type:      readReq.Type,
						// No Bytes and Data
						Status: err.Error(),
					}
					return b.naiveResponder(respCmd, response)
				}

				log.WithFields(log.Fields{"Type": readReq.Type}).Debug("Request type:")

				switch readReq.Type {
				case HexString:
					data = BytesToHexString(bytes) // convert byte to hex
					status = res.Status
					// TODO: add to history
				case UInt16:
					ret, err := BytesToUInt16s(bytes, readReq.Order) // order
					if err != nil {
						data = nil
						status = err.Error()
					} else {
						data = ret
						status = res.Status
						// TODO: add to history
					}
				case Int16:
					ret, err := BytesToInt16s(bytes, readReq.Order) // order
					if err != nil {
						data = nil
						status = err.Error()
					} else {
						data = ret
						status = res.Status
						// TODO: add to history
					}
				case Scale, UInt32, Int32, Float32: // 32-Bits
					if readReq.Len%2 != 0 {
						err := ErrInvalidLengthToConvert
						data = nil
						status = err.Error()
					} else {
						switch readReq.Type {
						case Scale:
							ret, err := LinearScalingRegisters(
								res.Data,
								readReq.Range.DomainLow,
								readReq.Range.DomainHigh,
								readReq.Range.RangeLow,
								readReq.Range.RangeHigh)
							if err != nil {
								data = nil
								status = err.Error()
							} else {
								data = ret
								status = res.Status
								// TODO: add to history
							}
						case UInt32:
							ret, err := BytesToUInt32s(bytes, readReq.Order)
							if err != nil {
								data = nil
								status = err.Error()
							} else {
								data = ret
								status = res.Status
								// TODO: add to history
							}
						case Int32:
							ret, err := BytesToInt32s(bytes, readReq.Order)
							if err != nil {
								data = nil
								status = err.Error()
							} else {
								data = ret
								status = res.Status
								// TODO: add to history
							}
						case Float32:
							ret, err := BytesToFloat32s(bytes, readReq.Order)
							if err != nil {
								data = nil
								status = err.Error()
							} else {
								data = ret
								status = res.Status
								// TODO: add to history
							}
						}
					}
				default: // case 0, 1(RegisterArray)
					data = res.Data
					status = res.Status
					// TODO: add to history

				}

				// shared response
				response = MbtcpPollData{
					TimeStamp: time.Now().UTC().UnixNano(),
					Name:      task.Name,
					Type:      readReq.Type,
					Bytes:     bytes,
					Data:      data,
					Status:    status,
				}

				return b.naiveResponder(respCmd, response)
			default: // should not reach here
				err := ErrResponseNotSupport
				log.WithFields(log.Fields{"cmd": task.Cmd}).Error(err.Error())
				response = MbtcpSimpleRes{
					Tid:    tid,
					Status: err.Error(),
				}
				return b.naiveResponder(respCmd, response)
			}
		}
	default: // should not reach here!!
		return ErrResponseNotSupport
	}
	return nil
}

// Start enable proactive service
func (b *MbtcpService) Start() {
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
					"cmd": msg[0],
					"req": msg[1],
				}).Debug("Receive request from upstream services")

				// parse request
				if req, err := b.ParseRequest(msg); req != nil {
					// handle request
					err = b.HandleRequest(msg[0], req)
					if err != nil {
						log.WithFields(log.Fields{
							"cmd": msg[0],
							"err": err,
						}).Error("Fail to handle request")
						// no need to send back again!
					}
				} else {
					log.WithFields(log.Fields{
						"cmd": msg[0],
						"err": err,
					}).Error("Fail to parse request")
					// send error back
					b.naiveResponder(msg[0], MbtcpSimpleRes{Status: err.Error()})
				}
			case b.sub.downstream:
				// receive from modbusd
				msg, _ := b.sub.downstream.RecvMessage(0)
				log.WithFields(log.Fields{
					"cmd":  msg[0],
					"resp": msg[1],
				}).Debug("Receive response from modbusd")

				// parse response
				if res, err := b.ParseResponse(msg); res != nil {
					// handle response
					err = b.HandleResponse(msg[0], res)
					if err != nil {
						log.WithFields(log.Fields{
							"cmd": msg[0],
							"err": err,
						}).Error("Fail to handle response")
						// no need to send back again!
					}
				} else {
					log.WithFields(log.Fields{
						"cmd": msg[0],
						"err": err,
					}).Error("Fail to parse response")
					// no need to send back, we don't know the sender
				}
			}
		}
	}
}

// Stop disable proactive service
func (b *MbtcpService) Stop() {
	log.Debug("Stop proactive service")
	b.scheduler.Stop()
	b.enable = false
	if b.sub.upstream != nil {
		b.sub.upstream.Close()
		b.pub.upstream.Close()
		b.sub.downstream.Close()
		b.pub.downstream.Close()
	}
}
