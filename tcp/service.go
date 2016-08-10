// Package tcp provide proactive service library for modbus `TCP` .
//
// By taka@cmwang.net
//
package tcp

import (
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	. "github.com/taka-wang/psmb"
	cron "github.com/taka-wang/psmb/cron"
	//conf "github.com/taka-wang/psmb/mini-conf"
	conf "github.com/taka-wang/psmb/viper-conf"
	log "github.com/takawang/logrus"
	zmq "github.com/takawang/zmq3"
)

var (
	// defaultMbPort default modbus slave port number
	defaultMbPort string
	// minConnTimeout minimal modbus tcp connection timeout
	minConnTimeout int64
	// minPollInterval minimal modbus tcp poll interval
	minPollInterval uint64
)

func setDefaults() {
	// set default psmbtcp values
	conf.SetDefault(keyTCPDefaultPort, defaultTCPDefaultPort)
	conf.SetDefault(keyMinConnectionTimout, defaultMinConnectionTimout)
	conf.SetDefault(keyPollInterval, defaultPollInterval)
	conf.SetDefault(keyMaxWorker, defaultMaxWorker)
	conf.SetDefault(keyMaxQueue, defaultMaxQueue)
	// set default zmq values
	conf.SetDefault(keyZmqPubUpstream, defaultZmqPubUpstream)
	conf.SetDefault(keyZmqPubDownstream, defaultZmqPubDownstream)
	conf.SetDefault(keyZmqSubUpstream, defaultZmqSubUpstream)
	conf.SetDefault(keyZmqSubDownstream, defaultZmqSubDownstream)
}

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true}) // before init logger
	log.SetLevel(log.DebugLevel)                            // ...
	setDefaults()                                           // set defaults

	defaultMbPort = conf.GetString(keyTCPDefaultPort)
	minConnTimeout = conf.GetInt64(keyMinConnectionTimout)
	minPollInterval = uint64(conf.GetInt(keyPollInterval))
}

// @Implement IProactiveService contract implicitly

type (

	// zSockets zmq sockets
	zSockets struct {
		// upstream subscriber from services
		upstream *zmq.Socket
		// downstream subscriber from modbusd
		downstream *zmq.Socket
	}

	// Service modbusd tcp proactive service type
	Service struct {
		// readerMap read/poll task map
		readerMap IReaderTaskDataStore
		// writerMap write task map
		writerMap IWriterTaskDataStore
		// historyMap history map
		historyMap IHistoryDataStore
		// filterMap filter map
		filterMap IFilterDataStore
		// scheduler cron scheduler
		scheduler cron.Scheduler
		// sub ZMQ subscriber endpoints
		sub zSockets
		// pub ZMQ publisher endpoints
		pub zSockets
		// poller ZMQ poller
		poller *zmq.Poller
		// enable poller flag
		enable bool
	}
)

// NewService modbus tcp proactive serivce constructor
func NewService(reader, writer, history, filter, sch string) (IProactiveService, error) {
	var readerPlugin IReaderTaskDataStore
	var writerPlugin IWriterTaskDataStore
	var historyPlugin IHistoryDataStore
	var filterPlugin IFilterDataStore
	var schedulerPlugin cron.Scheduler
	var err error
	// factory methods
	if readerPlugin, err = ReaderDataStoreCreator(reader); err != nil { // reader factory
		log.WithError(err).Panic("Fail to create reader data store")
		return nil, err
	}

	if writerPlugin, err = WriterDataStoreCreator(writer); err != nil { // writer factory
		log.WithError(err).Panic("Fail to create writer data store")
		return nil, err
	}

	if historyPlugin, err = HistoryDataStoreCreator(history); err != nil { // historian factory
		log.WithError(err).Panic("Fail to create history data store")
		return nil, err
	}

	if filterPlugin, err = FilterDataStoreCreator(filter); err != nil { // filter factory
		log.WithError(err).Panic("Fail to create filter data store")
		return nil, err
	}

	if schedulerPlugin, err = SchedulerCreator(sch); err != nil { // scheduler factory
		log.WithError(err).Panic("Fail to create scheduler")
		return nil, err
	}

	pubUpstream, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		log.WithError(err).Panic("Fail to create upstream publisher")
		return nil, err
	}
	pubDownstream, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		log.WithError(err).Panic("Fail to create downstream publisher")
		return nil, err
	}
	subUpstream, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		log.WithError(err).Panic("Fail to create upstream subscriber")
		return nil, err
	}
	subDownstream, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		log.WithError(err).Panic("Fail to create downstream subscriber")
		return nil, err
	}

	return &Service{
		enable:     true,
		readerMap:  readerPlugin,
		writerMap:  writerPlugin,
		historyMap: historyPlugin,
		filterMap:  filterPlugin,
		scheduler:  schedulerPlugin,
		pub: zSockets{
			upstream:   pubUpstream,
			downstream: pubDownstream,
		},
		sub: zSockets{
			upstream:   subUpstream,
			downstream: subDownstream,
		},
	}, nil
}

// marshal helper function to marshal structure
func marshal(r interface{}) (string, error) {
	bytes, err := json.Marshal(r) // marshal to json string
	if err != nil {
		return "", ErrMarshal
	}
	return string(bytes), nil
}

// addToHistory helper function to add data to history map
func (b *Service) addToHistory(name string, data interface{}) bool {
	// apply filter before logging
	retBool := b.applyFilter(name, data)
	if err := b.historyMap.Add(name, data); err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"name": name,
			"data": data,
		}).Error("Fail to add data to history data store")
	}
	return retBool

	/* debug
	log.WithFields(log.Fields{
		"name": name,
		"data": data,
	}).Debug("Add data to history data store")
	*/
}

// applyFilter apply filter, if no need to filter, return true.
func (b *Service) applyFilter(name string, data interface{}) bool {
	f, ok := b.filterMap.Get(name) // get filter request from map
	if !ok {
		//log.Debug(ErrFilterNotFound.Error())
		return true // no filter
	}
	filter := f.(MbtcpFilterStatus) // casting

	if len(filter.Arg) == 0 {
		log.WithError(ErrInvalidArgs).Debug("Apply filter")
		return true // no args
	}

	var latestStr string // latest history marshalled string
	var err error
	if filter.Type == Change {
		if latestStr, err = b.historyMap.GetLatest(name); err != nil {
			log.WithError(ErrNoLatestData).Debug("Apply filter")
			return true // no latest
		}
	}

	// reflect data interface type
	rVals := reflect.ValueOf(data)
	switch rVals.Kind() {
	case reflect.Array:
		if rVals.Len() == 0 {
			log.WithError(ErrNoData).Debug("Apply filter")
			return true // no data to filter
		}
		var val float32 // first element container in data interface
		switch rVals.Index(0).Kind() {
		case reflect.Uint16, reflect.Uint32: //uint16, uint32:
			val = float32(rVals.Index(0).Uint())
		default: // reflect.Float32
			val = float32(rVals.Index(0).Float())
		}

		switch filter.Type {
		case GreaterEqual: // val >= desired value
			if val >= filter.Arg[0] {
				return true
			}
			return false
		case Greater: // val > desired value
			if val > filter.Arg[0] {
				return true
			}
			return false
		case Equal: // val == desired value
			if val == filter.Arg[0] {
				return true
			}
			return false
		case Less: //  val < desired value
			if val < filter.Arg[0] {
				return true
			}
			return false
		case LessEqual: // val <= desired value
			if val <= filter.Arg[0] {
				return true
			}
			return false
		case InsideRange: // desired 1 < val < desired 2; desired values are sorted.
			if len(filter.Arg) < 2 {
				log.WithError(ErrInvalidArgs).Debug("Apply filter")
				return true
			}
			if val > filter.Arg[0] && val < filter.Arg[1] {
				return true
			}
			return false
		case InsideIncRange: // desired 1 <= val <= desired 2; desired values are sorted.
			if len(filter.Arg) < 2 {
				log.WithError(ErrInvalidArgs).Debug("Apply filter")
				return true
			}
			if val >= filter.Arg[0] && val <= filter.Arg[1] {
				return true
			}
			return false
		case OutsideRange: // val < desired 1 || val > desired 2; desired values are sorted.
			if len(filter.Arg) < 2 {
				log.WithError(ErrInvalidArgs).Debug("Apply filter")
				return true
			}
			if val < filter.Arg[0] || val > filter.Arg[1] {
				return true
			}
			return false
		case OutsideIncRange:
			if len(filter.Arg) < 2 { // val <= desired 1 || val >= desired 2; desired values are sorted.
				log.WithError(ErrInvalidArgs).Debug("Apply filter")
				return true
			}
			if val <= filter.Arg[0] || val >= filter.Arg[1] {
				return true
			}
			return false
		default: // change; compare with the latest history
			// unmarshal latest data
			var float32ArrData []float32
			if err := json.Unmarshal([]byte(latestStr), &float32ArrData); err != nil {
				log.WithError(ErrUnmarshal).Debug("Apply filter")
				return true // fail to unmarshal latest
			}

			if len(float32ArrData) == 0 {
				log.WithError(ErrNoLatestData).Debug("Apply filter")
				return true // empty history
			}
			if val == float32ArrData[0] { // compare
				return true
			}
			return false
		}
	case reflect.String:
		/* we do not intend to support filter on hex string
		str := rVals.String()
		if str == "" {
			return true // no data to filter
		}
		*/
		return true
	default: // should not reach here
		return true
	}
}

// Task task for scheduler
func (b *Service) Task(socket *zmq.Socket, req interface{}) {
	str, err := marshal(req)
	if err != nil {
		log.WithError(err).Error("Task")
		return
	}
	log.WithField("JSON", str).Debug("Send request to modbusd:")
	socket.Send("tcp", zmq.SNDMORE) // frame 1
	socket.Send(str, 0)             // convert to string; frame 2
}

func (b *Service) startZMQ() {
	log.Debug("Start ZMQ")

	// publishers
	if err := b.pub.upstream.Bind(conf.GetString(keyZmqPubUpstream)); err != nil {
		log.WithError(err).Fatal("Fail to bind upstream publisher")
	}
	if err := b.pub.downstream.Connect(conf.GetString(keyZmqPubDownstream)); err != nil {
		log.WithError(err).Fatal("Fail to connect to downstream publisher")
	}

	// subscribers
	if err := b.sub.upstream.Bind(conf.GetString(keyZmqSubUpstream)); err != nil {
		log.WithError(err).Fatal("Fail to bind upstream subscriber")
	}
	if err := b.sub.upstream.SetSubscribe(""); err != nil {
		log.WithError(err).Fatal("Fail to set upstream subscriber's filter")
	}
	if err := b.sub.downstream.Connect(conf.GetString(keyZmqSubDownstream)); err != nil {
		log.WithError(err).Fatal("Fail to connect to downstream subscriber")
	}
	if err := b.sub.downstream.SetSubscribe(""); err != nil {
		log.WithError(err).Fatal("Fail to set downstream subscriber's filter")
	}

	// poller
	b.poller = zmq.NewPoller() // new poller
	b.poller.Add(b.sub.upstream, zmq.POLLIN)
	b.poller.Add(b.sub.downstream, zmq.POLLIN)
}

func (b *Service) stopZMQ() {
	log.Debug("Stop ZMQ")

	// publishers
	if err := b.pub.upstream.Unbind(conf.GetString(keyZmqPubUpstream)); err != nil {
		log.WithError(err).Debug("Fail to unbind upstream publisher")
	}
	if err := b.pub.downstream.Disconnect(conf.GetString(keyZmqPubDownstream)); err != nil {
		log.WithError(err).Debug("Fail to disconnect from downstream publisher")
	}

	// subscribers
	if err := b.sub.upstream.Unbind(conf.GetString(keyZmqSubUpstream)); err != nil {
		log.WithError(err).Debug("Fail to unbind upstream subscriber")
	}
	if err := b.sub.downstream.Disconnect(conf.GetString(keyZmqSubDownstream)); err != nil {
		log.WithError(err).Debug("Fail to disconnect from downstream subscriber")
	}
}

// naiveResponder naive responder to send message back to upstream.
func (b *Service) naiveResponder(cmd string, resp interface{}) error {
	respStr, err := marshal(resp)
	if err != nil {
		log.WithError(err).Error("Fail to marshal for naive responder!")
		return err
	}

	log.WithField("response", respStr).Debug("Send message to services")
	b.pub.upstream.Send(cmd, zmq.SNDMORE) // task command
	b.pub.upstream.Send(respStr, 0)       // convert to string; frame 2
	return nil
}

// ParseRequest parse requests from services,
// 	only unmarshal request string to corresponding struct
func (b *Service) ParseRequest(msg []string) (interface{}, error) {
	// Check the length of multi-part message
	if len(msg) != 2 {
		return nil, ErrInvalidMessageLength
	}

	log.WithField("cmd", msg[0]).Debug("Parse request from upstream services")

	switch msg[0] {
	case mbGetTimeout, mbSetTimeout:
		var req MbtcpTimeoutReq
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			return nil, ErrUnmarshal
		}
		return req, nil
	case mbOnceWrite:
		// unmarshal partial request (except data field)
		var data json.RawMessage // raw []byte
		req := MbtcpWriteReq{Data: &data}
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			return nil, ErrUnmarshal
		}

		var uint16Data uint16
		var stringData string
		var uint16ArrData []uint16

		// unmarshal remaining data field
		switch MbCmdType(strconv.Itoa(req.FC)) {
		case fc5: // write single bit; uint16
			if err := json.Unmarshal(data, &uint16Data); err != nil {
				return nil, ErrUnmarshal
			}
			req.Data = uint16Data // unmarshal to uint16
			return req, nil
		case fc15: // write multiple bits; []uint16
			if err := json.Unmarshal(data, &uint16ArrData); err != nil {
				return nil, ErrUnmarshal
			}
			req.Data = uint16ArrData // unmarshal to uint16 array
			return req, nil
		case fc6: // write single register in dec|hex
			err := json.Unmarshal(data, &stringData)
			if err != nil {
				return nil, ErrUnmarshal
			}

			// check dec or hex
			if req.Hex {
				uint16ArrData, err = HexStringToRegisters(stringData)
			} else {
				uint16ArrData, err = DecimalStringToRegisters(stringData)
			}
			if err != nil {
				return nil, err
			}
			req.Data = uint16ArrData[0] // retrieve only one register
			return req, nil
		case fc16: // write multiple register in dec/hex
			err := json.Unmarshal(data, &stringData)
			if err != nil {
				return nil, ErrUnmarshal
			}

			// check dec or hex
			if req.Hex {
				uint16ArrData, err = HexStringToRegisters(stringData)
			} else {
				uint16ArrData, err = DecimalStringToRegisters(stringData)
			}
			if err != nil {
				return nil, err
			}
			req.Data = uint16ArrData
			return req, nil
		default: // should not reach here
			return nil, ErrInvalidFunctionCode
		}
	case mbOnceRead:
		var req MbtcpReadReq
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			return nil, ErrUnmarshal
		}
		return req, nil
	case mbCreatePoll:
		var req MbtcpPollStatus
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			return nil, ErrUnmarshal
		}
		return req, nil
	case mbUpdatePoll, mbGetPoll, mbDeletePoll,
		mbTogglePoll, mbGetPolls, mbDeletePolls,
		mbTogglePolls, mbGetPollHistory, mbExportPolls:
		var req MbtcpPollOpReq
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			return nil, ErrUnmarshal
		}
		return req, nil
	case mbImportPolls:
		var req MbtcpPollsStatus
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			return nil, ErrUnmarshal
		}
		return req, nil
	case mbCreateFilter, mbUpdateFilter:
		var req MbtcpFilterStatus
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			return nil, ErrUnmarshal
		}
		return req, nil
	case mbGetFilter, mbDeleteFilter, mbToggleFilter,
		mbGetFilters, mbDeleteFilters, mbToggleFilters, mbExportFilters:
		var req MbtcpFilterOpReq
		if err := json.Unmarshal([]byte(msg[1]), &req); err != nil {
			return nil, ErrUnmarshal
		}
		return req, nil
	case mbImportFilters:
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
// 	do error checking
func (b *Service) HandleRequest(cmd string, r interface{}) error {
	log.WithField("cmd", cmd).Debug("Handle request from upstream services")

	switch cmd {
	case mbGetTimeout:
		req := r.(MbtcpTimeoutReq)
		TidStr := strconv.FormatInt(req.Tid, 10)        // convert tid to string
		cmdInt, _ := strconv.Atoi(string(getMbTimeout)) // convert to modbusd command
		command := DMbtcpTimeout{
			Tid: TidStr,
			Cmd: cmdInt,
		}
		// add request to write task map
		b.writerMap.Add(TidStr, cmd)
		// add command to scheduler as emergency request
		b.scheduler.Emergency().Do(b.Task, b.pub.downstream, command)
		return nil
	case mbSetTimeout:
		req := r.(MbtcpTimeoutReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string
		cmdInt, _ := strconv.Atoi(string(setMbTimeout))
		command := DMbtcpTimeout{
			Tid: TidStr,
			Cmd: cmdInt,
		}

		// protect invalid timeout value
		if req.Data < minConnTimeout {
			command.Timeout = minConnTimeout
		} else {
			command.Timeout = req.Data
		}
		// add request to write task map
		b.writerMap.Add(TidStr, cmd)
		// add command to scheduler as emergency request
		b.scheduler.Emergency().Do(b.Task, b.pub.downstream, command)
		return nil
	case mbOnceWrite:
		req := r.(MbtcpWriteReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string

		// protect null port
		if req.Port == "" {
			req.Port = defaultMbPort
		}

		// length checker
		switch MbCmdType(strconv.Itoa(req.FC)) {
		case fc15, fc16:
			l := uint16(len(req.Data.([]uint16)))
			if req.Len < l {
				req.Len = l
			}
			// we don't check max length, let modbusd do it.
		default:
			// do nothing
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
	case mbOnceRead:
		req := r.(MbtcpReadReq)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string

		// function code checker
		if req.FC < 1 || req.FC > 4 {
			err := ErrInvalidFunctionCode // invalid read function code
			log.WithField("FC", req.FC).Warn(err.Error())
			// send back
			resp := MbtcpSimpleRes{Tid: req.Tid, Status: err.Error()}
			return b.naiveResponder(cmd, resp)
		}

		// protect null port
		if req.Port == "" {
			req.Port = defaultMbPort
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

		// add request to read/poll task map,
		// since this is a read task instead of poll task, thus pass null name
		if err := b.readerMap.Add("", TidStr, cmd, req); err != nil {
			log.WithError(err).Warn(mbOnceRead) // maybe out of capacity
			// send error back
			resp := MbtcpSimpleRes{Tid: req.Tid, Status: err.Error()}
			return b.naiveResponder(cmd, resp)
		}
		// add command to scheduler as emergency request
		b.scheduler.Emergency().Do(b.Task, b.pub.downstream, command)
		return nil
	case mbCreatePoll:
		req := r.(MbtcpPollStatus)
		TidStr := strconv.FormatInt(req.Tid, 10) // convert tid to string

		// function code checker
		if req.FC < 1 || req.FC > 4 {
			err := ErrInvalidFunctionCode // invalid read function code
			log.WithField("FC", req.FC).Warn(err.Error())
			// send error back
			resp := MbtcpSimpleRes{Tid: req.Tid, Status: err.Error()}
			return b.naiveResponder(cmd, resp)
		}

		// protect null poll name
		if req.Name == "" {
			err := ErrInvalidPollName
			log.WithField("Name", req.Name).Warn(err.Error())
			// send back
			resp := MbtcpSimpleRes{Tid: req.Tid, Status: err.Error()}
			return b.naiveResponder(cmd, resp)
		}

		// protect null port
		if req.Port == "" {
			req.Port = defaultMbPort
		}

		// length checker
		if req.Len < 1 {
			req.Len = 1 // set minimal length of read
		}

		// check interval value
		if req.Interval < minPollInterval {
			req.Interval = minPollInterval
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
		if err := b.readerMap.Add(req.Name, TidStr, cmd, req); err != nil {
			log.WithError(err).Warn(mbCreatePoll) // maybe out of capacity
			// send error back
			resp := MbtcpSimpleRes{Tid: req.Tid, Status: err.Error()}
			return b.naiveResponder(cmd, resp)
		}

		// add command to scheduler as regular request
		b.scheduler.EveryWithName(req.Interval, req.Name).Seconds().Do(b.Task, b.pub.downstream, command)

		if !req.Enabled { // if not enabled, pause the task
			b.scheduler.PauseWithName(req.Name)
		}
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: "ok"}
		return b.naiveResponder(cmd, resp)
	case mbUpdatePoll:
		req := r.(MbtcpPollOpReq)
		status := "ok"

		// check interval value
		if req.Interval < minPollInterval {
			req.Interval = minPollInterval
		}

		// update task interval
		if ok := b.scheduler.UpdateIntervalWithName(req.Name, req.Interval); !ok {
			err := ErrInvalidPollName // not in scheduler
			log.WithField("Name", req.Name).Warn(err.Error())
			status = err.Error() // set error status
		}

		// update read/poll task map
		if status == "ok" {
			if err := b.readerMap.UpdateIntervalByName(req.Name, req.Interval); err != nil {
				log.WithField("Name", req.Name).Warn(err.Error())
				status = err.Error() // set error status
			}
		}
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: status}
		return b.naiveResponder(cmd, resp)
	case mbGetPoll:
		req := r.(MbtcpPollOpReq)
		t, ok := b.readerMap.GetTaskByName(req.Name)
		task := t.(ReaderTask) // type casting
		if !ok {
			err := ErrInvalidPollName // not in read/poll task map
			log.WithField("Name", req.Name).Warn(err.Error())
			// send error back
			resp := MbtcpSimpleRes{Tid: req.Tid, Status: err.Error()}
			return b.naiveResponder(cmd, resp)
		}

		// send back
		request := task.Req.(MbtcpPollStatus)
		resp := MbtcpPollStatus{
			Tid:      req.Tid,
			Name:     req.Name,
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
	case mbDeletePoll:
		req := r.(MbtcpPollOpReq)
		status := "ok"
		// remove task from scheduler
		if ok := b.scheduler.RemoveWithName(req.Name); !ok {
			err := ErrInvalidPollName // not in scheduler
			log.WithField("Name", req.Name).Warn(err.Error())
			status = err.Error() // set error status
		}
		// remove task from read/poll map
		b.readerMap.DeleteTaskByName(req.Name)
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: status}
		return b.naiveResponder(cmd, resp)
	case mbTogglePoll:
		req := r.(MbtcpPollOpReq)
		status := "ok"

		var ok bool
		if req.Enabled {
			ok = b.scheduler.ResumeWithName(req.Name) // update scheduler
		} else {
			ok = b.scheduler.PauseWithName(req.Name) // update scheduler
		}

		if !ok {
			// not in scheduler
			err := ErrInvalidPollName
			log.WithField("Name", req.Name).Warn(err.Error())
			status = err.Error() // set error status
		} else {
			// update read/poll task map
			if err := b.readerMap.UpdateToggleByName(req.Name, req.Enabled); err != nil {
				log.WithField("Name", req.Name).Warn(err.Error())
				status = err.Error() // set error status
			}
		}
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: status}
		return b.naiveResponder(cmd, resp)
	case mbGetPolls, mbExportPolls:
		req := r.(MbtcpPollOpReq)
		reqs := b.readerMap.GetAll().([]MbtcpPollStatus) // type casting
		//log.WithField("reqs", reqs).Debug("after GetAll")

		// send back
		resp := MbtcpPollsStatus{
			Tid:    req.Tid,
			Status: "ok",
			Polls:  reqs,
		}
		return b.naiveResponder(cmd, resp)
	case mbDeletePolls:
		req := r.(MbtcpPollOpReq)
		b.scheduler.Clear()     // remove all tasks from scheduler
		b.readerMap.DeleteAll() // remove all tasks from read/poll task map
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: "ok"}
		return b.naiveResponder(cmd, resp)
	case mbTogglePolls:
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
	case mbImportPolls:
		request := r.(MbtcpPollsStatus)
		for _, req := range request.Polls {
			// function code checker
			if req.FC < 1 || req.FC > 4 {
				err := ErrInvalidFunctionCode // invalid read function code
				log.WithField("FC", req.FC).Warn(err.Error())
				continue // bypass
			}

			// protect null poll name
			if req.Name == "" {
				err := ErrInvalidPollName
				log.WithField("Name", req.Name).Warn(err.Error())
				continue // bypass
			}

			// protect null port
			if req.Port == "" {
				req.Port = defaultMbPort
			}

			// length checker
			if req.Len < 1 {
				req.Len = 1 // set minimal length of read
			}

			// check interval value
			if req.Interval < minPollInterval {
				req.Interval = minPollInterval
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

			// Add task to read/poll task map
			if err := b.readerMap.Add(req.Name, TidStr, cmd, req); err != nil {
				log.WithError(err).Warn(mbImportPolls) // maybe out of capacity
				// send error back
				resp := MbtcpSimpleRes{Tid: request.Tid, Status: err.Error()}
				return b.naiveResponder(cmd, resp)
			}

			b.scheduler.EveryWithName(req.Interval, req.Name).Seconds().Do(b.Task, b.pub.downstream, command) // add command to scheduler as regular request

			if !req.Enabled { // if not enabled, pause the task
				b.scheduler.PauseWithName(req.Name)
			}
		}
		// send back
		resp := MbtcpSimpleRes{Tid: request.Tid, Status: "ok"}
		return b.naiveResponder(cmd, resp)
	case mbGetPollHistory:
		req := r.(MbtcpPollOpReq)
		resp := MbtcpHistoryData{Tid: req.Tid, Name: req.Name, Status: "ok"}
		ret, err := b.historyMap.GetAll(req.Name)
		if err != nil {
			log.WithField("Name", req.Name).Warn(err.Error())
			resp.Status = err.Error()
			return b.naiveResponder(cmd, resp)
		}
		resp.Data = ret
		return b.naiveResponder(cmd, resp)
	case mbCreateFilter, mbUpdateFilter:
		req := r.(MbtcpFilterStatus)
		status := "ok"
		if req.Name == "" {
			err := ErrInvalidPollName
			log.WithError(err).Warn(mbCreateFilter)
			status = err.Error() // set error status
		} else {
			// swap filter args
			if len(req.Arg) > 1 && req.Arg[0] > req.Arg[1] {
				tmp := req.Arg[1]
				req.Arg[1] = req.Arg[0]
				req.Arg[0] = tmp
			}

			// add or update to filter map
			if err := b.filterMap.Add(req.Name, req); err != nil {
				log.WithError(err).Error(mbImportFilters)
				status = err.Error() // set error status
			}
		}

		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: status}
		return b.naiveResponder(cmd, resp)
	case mbGetFilter:
		req := r.(MbtcpFilterOpReq)
		status := "ok"
		var filter interface{}
		var ok bool

		if req.Name == "" {
			err := ErrInvalidPollName
			log.WithError(err).Error(mbGetFilter)
			status = err.Error() // set error status
			ok = false
		} else {
			// get filter from map
			if filter, ok = b.filterMap.Get(req.Name); !ok {
				err := ErrInvalidPollName
				log.WithError(err).Error(mbGetFilter)
				status = err.Error() // set error status
			}
		}

		if !ok { // send error back
			resp := MbtcpSimpleRes{Tid: req.Tid, Status: status}
			return b.naiveResponder(cmd, resp)
		}

		request := filter.(MbtcpFilterStatus) // type casting
		// send back
		resp := MbtcpFilterStatus{
			Tid:     req.Tid,
			Name:    req.Name,
			Enabled: request.Enabled,
			Type:    request.Type,
			Arg:     request.Arg,
			Status:  status, // "ok"
		}
		return b.naiveResponder(cmd, resp)
	case mbDeleteFilter:
		req := r.(MbtcpFilterOpReq)
		status := "ok"
		if req.Name == "" {
			err := ErrInvalidPollName
			log.WithError(err).Error(mbDeleteFilter)
			status = err.Error() // set error status
		} else {
			b.filterMap.Delete(req.Name)
		}
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: status}
		return b.naiveResponder(cmd, resp)
	case mbToggleFilter:
		req := r.(MbtcpFilterOpReq)
		status := "ok"
		if req.Name == "" {
			err := ErrInvalidPollName
			log.WithError(err).Error(mbToggleFilter)
			status = err.Error() // set error status
		} else {
			if err := b.filterMap.UpdateToggle(req.Name, req.Enabled); err != nil {
				log.WithError(err).Error(mbToggleFilter)
				status = err.Error() // set error status
			}
		}
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: status}
		return b.naiveResponder(cmd, resp)
	case mbGetFilters, mbExportFilters:
		req := r.(MbtcpFilterOpReq)
		if reqs, ok := b.filterMap.GetAll().([]MbtcpFilterStatus); ok {
			resp := MbtcpFiltersStatus{
				Tid:     req.Tid,
				Status:  "ok",
				Filters: reqs,
			}
			// send back
			return b.naiveResponder(cmd, resp)
		}
		// send error back
		err := ErrFiltersNotFound
		log.WithError(err).Error(mbGetFilters)
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: err.Error()}
		return b.naiveResponder(cmd, resp)
	case mbDeleteFilters:
		req := r.(MbtcpFilterOpReq)
		b.filterMap.DeleteAll()
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: "ok"}
		return b.naiveResponder(cmd, resp)
	case mbToggleFilters:
		req := r.(MbtcpFilterOpReq)
		b.filterMap.UpdateAllToggles(req.Enabled)
		// send back
		resp := MbtcpSimpleRes{Tid: req.Tid, Status: "ok"}
		return b.naiveResponder(cmd, resp)
	case mbImportFilters:
		requests, ok := r.(MbtcpFiltersStatus)
		status := "ok"
		if !ok {
			err := ErrCasting
			log.WithError(err).Error(mbImportFilters)
			status = err.Error() // set error status
		} else {
			for _, v := range requests.Filters {
				if v.Name != "" {
					// swap
					if len(v.Arg) > 1 && v.Arg[0] > v.Arg[1] {
						tmp := v.Arg[1]
						v.Arg[1] = v.Arg[0]
						v.Arg[0] = tmp
					}
					// add or update to filter map
					if err := b.filterMap.Add(v.Name, v); err != nil {
						log.WithError(err).Error(mbImportFilters)
						status = err.Error() // set error status
						break                // break the for loop
					}
				}
			}
		}

		// send back
		resp := MbtcpSimpleRes{Tid: requests.Tid, Status: status}
		return b.naiveResponder(cmd, resp)
	default: // should not reach here!!
		return ErrRequestNotSupport
	}
}

// ParseResponse parse responses from modbusd,
// 	only unmarshal response string to corresponding struct
func (b *Service) ParseResponse(msg []string) (interface{}, error) {
	// Check the length of multi-part message
	if len(msg) != 2 {
		return nil, ErrInvalidMessageLength
	}

	log.WithField("cmd", msg[0]).Debug("Parse response from modbusd")

	switch MbCmdType(msg[0]) {
	case setMbTimeout, getMbTimeout:
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

// HandleResponse handle responses from modbusd
func (b *Service) HandleResponse(cmd string, r interface{}) error {
	log.WithField("cmd", cmd).Debug("Handle response from modbusd")

	switch MbCmdType(cmd) {
	case fc5, fc6, fc15, fc16, setMbTimeout, getMbTimeout: // done: one-off requests
		var TidStr string
		var resp interface{}

		switch MbCmdType(cmd) {
		case setMbTimeout, getMbTimeout: // one-off timeout requests
			res := r.(DMbtcpTimeout)
			tid, _ := strconv.ParseInt(res.Tid, 10, 64)
			TidStr = res.Tid
			var int64Data int64

			if MbCmdType(cmd) == getMbTimeout {
				int64Data = res.Timeout
			}

			resp = MbtcpTimeoutRes{
				Tid:    tid,
				Status: res.Status,
				Data:   int64Data, // getMbTimeout only
			}
		case fc5, fc6, fc15, fc16: // one-off write requests
			res := r.(DMbtcpRes)
			tid, _ := strconv.ParseInt(res.Tid, 10, 64)
			TidStr = res.Tid
			resp = MbtcpSimpleRes{
				Tid:    tid,
				Status: res.Status,
			}
		default: // should not reach here
			return ErrResponseNotSupport
		}

		//
		// send back one-off task reponse and remove from write task map
		//
		respStr, err := marshal(resp)
		if err != nil {
			log.WithError(err).Error("marshal")
			return err
		}

		// check write task map
		if cmd, ok := b.writerMap.Get(TidStr); ok {
			log.WithField("JSON", respStr).Debug("Send message to services")
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

		// check read task table
		t, ok := b.readerMap.GetTaskByID(res.Tid)
		if !ok {
			return ErrRequestNotFound
		}

		task := t.(ReaderTask) // type casting
		respCmd := task.Cmd    // default response command string

		var response interface{}
		var data interface{} // shared variable
		status := "ok"       // shared variables
		noFilter := true     // feedback data or not flag

		switch MbCmdType(cmd) {
		case fc1, fc2: // done: read bits

			switch task.Cmd {
			case mbOnceRead: // one-off requests
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
				b.readerMap.DeleteTaskByID(res.Tid)
			case mbCreatePoll, mbImportPolls: // poll data
				respCmd = mbData // set as "mbtcp.data"
				if res.Status != "ok" {
					data = nil
				} else {
					data = res.Data
					noFilter = b.addToHistory(task.Name, data) // add to history; type: []uint16
				}
				response = MbtcpPollData{
					TimeStamp: time.Now().UTC().UnixNano(),
					Name:      task.Name,
					Status:    res.Status,
					Data:      data,
				}
			default: // should not reach here
				err := ErrResponseNotSupport
				log.WithField("cmd", cmd).Error(err.Error())
				response = MbtcpSimpleRes{
					Tid:    tid,
					Status: err.Error(),
				}
			}
			if noFilter {
				return b.naiveResponder(respCmd, response)
			}
			return nil
		case fc3, fc4: // read registers
			switch task.Cmd {
			case mbOnceRead: // one-off requests
				readReq := task.Req.(MbtcpReadReq) // type casting

				// check modbus response status
				if res.Status != "ok" {
					response = MbtcpReadRes{
						Tid:    tid,
						Type:   readReq.Type,
						Status: res.Status,
					}
					// remove from read table
					b.readerMap.DeleteTaskByID(res.Tid)
					return b.naiveResponder(respCmd, response)
				}

				// convert register to byte array
				bytes, err := RegistersToBytes(res.Data)
				if err != nil {
					log.WithError(err).Error("handleResponse: RegistersToBytes failed")
					response = MbtcpReadRes{
						Tid:    tid,
						Type:   readReq.Type,
						Status: err.Error(),
					}
					// remove from read table
					b.readerMap.DeleteTaskByID(res.Tid)
					return b.naiveResponder(respCmd, response)
				}

				log.WithField("Type", readReq.Type).Debug("Request type:")

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
				b.readerMap.DeleteTaskByID(res.Tid)
				return b.naiveResponder(respCmd, response)

			case mbCreatePoll, mbImportPolls: // poll data
				readReq := task.Req.(MbtcpPollStatus) // type casting
				respCmd = mbData                      // set as "mbtcp.data"

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

				log.WithField("Type", readReq.Type).Debug("Request type:")

				switch readReq.Type {
				case HexString:
					data = BytesToHexString(bytes) // convert byte to hex string
					status = res.Status
					noFilter = b.addToHistory(task.Name, data) // add to history; type: string
				case UInt16:
					if ret, err := BytesToUInt16s(bytes, readReq.Order); err != nil {
						data = nil
						status = err.Error()
					} else {
						data = ret
						status = res.Status
						noFilter = b.addToHistory(task.Name, data) // add to history; type: []uint16
					}
				case Int16:
					if ret, err := BytesToInt16s(bytes, readReq.Order); err != nil {
						data = nil
						status = err.Error()
					} else {
						data = ret
						status = res.Status
						noFilter = b.addToHistory(task.Name, data) // add to history; type: []uint16
					}
				case Scale, UInt32, Int32, Float32: // 32-Bits
					if readReq.Len%2 != 0 {
						err := ErrInvalidLengthToConvert
						data = nil
						status = err.Error()
					} else {
						switch readReq.Type {
						case Scale:
							if ret, err := LinearScalingRegisters(
								res.Data,
								readReq.Range.DomainLow,
								readReq.Range.DomainHigh,
								readReq.Range.RangeLow,
								readReq.Range.RangeHigh); err != nil {

								data = nil
								status = err.Error()
							} else {
								data = ret
								status = res.Status
								noFilter = b.addToHistory(task.Name, data) // add to history; type: []float32
							}
						case UInt32:
							if ret, err := BytesToUInt32s(bytes, readReq.Order); err != nil {
								data = nil
								status = err.Error()
							} else {
								data = ret
								status = res.Status
								noFilter = b.addToHistory(task.Name, data) // add to history; type: []uint32
							}
						case Int32:
							if ret, err := BytesToInt32s(bytes, readReq.Order); err != nil {
								data = nil
								status = err.Error()
							} else {
								data = ret
								status = res.Status
								noFilter = b.addToHistory(task.Name, data) // add to history; type: []uint32
							}
						case Float32:
							if ret, err := BytesToFloat32s(bytes, readReq.Order); err != nil {
								data = nil
								status = err.Error()
							} else {
								data = ret
								status = res.Status
								b.addToHistory(task.Name, data) // add to history; type: []float32
								noFilter = b.applyFilter(task.Name, data)
							}
						}
					}
				default: // case 0, 1(RegisterArray)
					data = res.Data
					status = res.Status
					// b.addToHistory(task.Name, HistoryData{Ts: time.Now().UTC().UnixNano(), Data: data})
					noFilter = b.addToHistory(task.Name, data) // add to history; type: []uint16
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
				if noFilter {
					return b.naiveResponder(respCmd, response)
				}
				return nil
			default: // should not reach here
				err := ErrResponseNotSupport
				log.WithField("cmd", task.Cmd).Error(err.Error())
				response = MbtcpSimpleRes{
					Tid:    tid,
					Status: err.Error(),
				}
				return b.naiveResponder(respCmd, response)
			}
		default: // should not reach here
			return ErrResponseNotSupport
		}
	default: // should not reach here!!
		return ErrResponseNotSupport
	}
}

// Start enable proactive service
func (b *Service) Start() {

	log.Debug("Start proactive service")
	b.scheduler.Start()
	b.enable = true
	b.startZMQ()

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
					err := b.HandleRequest(msg[0], req)
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
					err := b.HandleResponse(msg[0], res)
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
func (b *Service) Stop() {
	log.Debug("Stop proactive service")
	b.scheduler.Stop()
	b.enable = false
	b.stopZMQ()
}
