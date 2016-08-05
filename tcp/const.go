package tcp

// plugin name
const (
	readerPluginName    = "ReaderPlugin"
	writerPluginName    = "WriterPlugin"
	schedulerPluginName = "SchedulerPlugin"
	historyPluginName   = "HistoryPlugin"
)

// [psmbtcp]
const (
	packageName                = "tcp"
	keyTCPDefaultPort          = "psmbtcp.default_port"
	keyMinConnectionTimout     = "psmbtcp.min_connection_timeout"
	keyPollInterval            = "psmbtcp.min_poll_interval"
	defaultTCPDefaultPort      = "502"
	defaultMinConnectionTimout = 200000
	defaultPollInterval        = 1
)

// [zmq]
const (
	keyZmqPubUpstream       = "zmq.pub.upstream"
	keyZmqPubDownstream     = "zmq.pub.downstream"
	keyZmqSubUpstream       = "zmq.sub.upstream"
	keyZmqSubDownstream     = "zmq.sub.downstream"
	defaultZmqPubUpstream   = "ipc:///tmp/from.psmb"
	defaultZmqPubDownstream = "ipc:///tmp/to.modbus"
	defaultZmqSubUpstream   = "ipc:///tmp/to.psmb"
	defaultZmqSubDownstream = "ipc:///tmp/from.modbus"
)
