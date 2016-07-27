package psmb

// Upstream & Downstream command tables

// command table for modbusd
const (
	fc1  MbtcpCmdType = "1"
	fc2  MbtcpCmdType = "2"
	fc3  MbtcpCmdType = "3"
	fc4  MbtcpCmdType = "4"
	fc5  MbtcpCmdType = "5"
	fc6  MbtcpCmdType = "6"
	fc15 MbtcpCmdType = "15"
	fc16 MbtcpCmdType = "16"
	// setTCPTimeout set TCP connection timeout
	setTCPTimeout MbtcpCmdType = "50"
	// getTCPTimeout get TCP connection timeout
	getTCPTimeout MbtcpCmdType = "51"
)

// command table for upstream services
const (
	mbtcpOnceRead       = "mbtcp.once.read"
	mbtcpOnceWrite      = "mbtcp.once.write"
	mbtcpGetTimeout     = "mbtcp.timeout.read"
	mbtcpSetTimeout     = "mbtcp.timeout.update"
	mbtcpCreatePoll     = "mbtcp.poll.create"
	mbtcpUpdatePoll     = "mbtcp.poll.update"
	mbtcpGetPoll        = "mbtcp.poll.read"
	mbtcpDeletePoll     = "mbtcp.poll.delete"
	mbtcpTogglePoll     = "mbtcp.poll.toggle"
	mbtcpGetPolls       = "mbtcp.polls.read"
	mbtcpDeletePolls    = "mbtcp.polls.delete"
	mbtcpTogglePolls    = "mbtcp.polls.toggle"
	mbtcpImportPolls    = "mbtcp.polls.import"
	mbtcpExportPolls    = "mbtcp.polls.export"
	mbtcpGetPollHistory = "mbtcp.poll.history"
	mbtcpCreateFilter   = "mbtcp.filter.create"
	mbtcpUpdateFilter   = "mbtcp.filter.update"
	mbtcpGetFilter      = "mbtcp.filter.read"
	mbtcpDeleteFilter   = "mbtcp.filter.delete"
	mbtcpToggleFilter   = "mbtcp.filter.toggle"
	mbtcpGetFilters     = "mbtcp.filters.read"
	mbtcpDeleteFilters  = "mbtcp.filters.delete"
	mbtcpToggleFilters  = "mbtcp.filters.toggle"
	mbtcpImportFilters  = "mbtcp.filters.import"
	mbtcpExportFilters  = "mbtcp.filters.export"
	// Poll data
	mbtcpData = "mbtcp.data"
)
