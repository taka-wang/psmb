package main

import "github.com/taka-wang/gocron"

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

// NewBridge init bridge
func NewBridge() Bridge {
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
	b.initPub()
	b.initSub()
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
				req, err := ParseRequest(msg)
				if err != nil {
					// todo: send error back
				} else {
					err = RequestHandler(msg[0], req, toModbusd, toService)
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
					err = ResponseHandler(MbtcpCmdType(msg[0]), res, toService)
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
