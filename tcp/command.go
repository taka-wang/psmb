package tcp

// Upstream & Downstream command tables

// MbCmdType defines modbus tcp command type
type MbCmdType string

// command table for modbusd
const (
	fc1  MbCmdType = "1"
	fc2  MbCmdType = "2"
	fc3  MbCmdType = "3"
	fc4  MbCmdType = "4"
	fc5  MbCmdType = "5"
	fc6  MbCmdType = "6"
	fc15 MbCmdType = "15"
	fc16 MbCmdType = "16"
	// setMbTimeout set TCP connection timeout
	setMbTimeout MbCmdType = "50"
	// getMbTimeout get TCP connection timeout
	getMbTimeout MbCmdType = "51"
)

// command table for upstream services
const (
	mbOnceRead       = "mbtcp.once.read"
	mbOnceWrite      = "mbtcp.once.write"
	mbGetTimeout     = "mbtcp.timeout.read"
	mbSetTimeout     = "mbtcp.timeout.update"
	mbCreatePoll     = "mbtcp.poll.create"
	mbUpdatePoll     = "mbtcp.poll.update"
	mbGetPoll        = "mbtcp.poll.read"
	mbDeletePoll     = "mbtcp.poll.delete"
	mbTogglePoll     = "mbtcp.poll.toggle"
	mbGetPolls       = "mbtcp.polls.read"
	mbDeletePolls    = "mbtcp.polls.delete"
	mbTogglePolls    = "mbtcp.polls.toggle"
	mbImportPolls    = "mbtcp.polls.import"
	mbExportPolls    = "mbtcp.polls.export"
	mbGetPollHistory = "mbtcp.poll.history"
	mbCreateFilter   = "mbtcp.filter.create"
	mbUpdateFilter   = "mbtcp.filter.update"
	mbGetFilter      = "mbtcp.filter.read"
	mbDeleteFilter   = "mbtcp.filter.delete"
	mbToggleFilter   = "mbtcp.filter.toggle"
	mbGetFilters     = "mbtcp.filters.read"
	mbDeleteFilters  = "mbtcp.filters.delete"
	mbToggleFilters  = "mbtcp.filters.toggle"
	mbImportFilters  = "mbtcp.filters.import"
	mbExportFilters  = "mbtcp.filters.export"
	// Poll data
	mbData = "mbtcp.data"
)
