package psmb

// command table for upstream services - TCP
const (
	CmdMbtcpOnceRead       = "mbtcp.once.read"
	CmdMbtcpOnceWrite      = "mbtcp.once.write"
	CmdMbtcpGetTimeout     = "mbtcp.timeout.read"
	CmdMbtcpSetTimeout     = "mbtcp.timeout.update"
	CmdMbtcpCreatePoll     = "mbtcp.poll.create"
	CmdMbtcpUpdatePoll     = "mbtcp.poll.update"
	CmdMbtcpGetPoll        = "mbtcp.poll.read"
	CmdMbtcpDeletePoll     = "mbtcp.poll.delete"
	CmdMbtcpTogglePoll     = "mbtcp.poll.toggle"
	CmdMbtcpGetPolls       = "mbtcp.polls.read"
	CmdMbtcpDeletePolls    = "mbtcp.polls.delete"
	CmdMbtcpTogglePolls    = "mbtcp.polls.toggle"
	CmdMbtcpImportPolls    = "mbtcp.polls.import"
	CmdMbtcpExportPolls    = "mbtcp.polls.export"
	CmdMbtcpGetPollHistory = "mbtcp.poll.history"
	CmdMbtcpCreateFilter   = "mbtcp.filter.create"
	CmdMbtcpUpdateFilter   = "mbtcp.filter.update"
	CmdMbtcpGetFilter      = "mbtcp.filter.read"
	CmdMbtcpDeleteFilter   = "mbtcp.filter.delete"
	CmdMbtcpToggleFilter   = "mbtcp.filter.toggle"
	CmdMbtcpGetFilters     = "mbtcp.filters.read"
	CmdMbtcpDeleteFilters  = "mbtcp.filters.delete"
	CmdMbtcpToggleFilters  = "mbtcp.filters.toggle"
	CmdMbtcpImportFilters  = "mbtcp.filters.import"
	CmdMbtcpExportFilters  = "mbtcp.filters.export"
	CmdMbtcpData           = "mbtcp.data" // Poll data
)
