package psmb

//
// Interfaces
//

// ProactiveService proactive service contracts,
// all services should implement the following methods.
type ProactiveService interface {
	// Start enable proactive service
	Start()
	// Stop disable proactive service
	Stop()
	// parseRequest parse requests from services
	parseRequest(msg []string) (interface{}, error)
	// handleRequest handle requests from services
	handleRequest(cmd string, r interface{}) error
	// parseResponse parse responses from modbusd
	parseResponse(msg []string) (interface{}, error)
	// handleResponse handle responses from modbusd
	handleResponse(cmd string, r interface{}) error
}

// MbtcpWriteTask mbtcp write task interface
//	(Tid, Command) map
type MbtcpWriteTask interface {
	// Add add request to write task map,
	// params: TID, CMD strings.
	Add(tid, cmd string)

	// Get get request from write task map,
	// params: TID string,
	// return: cmd string, exist flag.
	Get(tid string) (string, bool)

	// Delete remove request from write task map
	// params: TID string.
	Delete(tid string)
}

// MbtcpReadTask mbtcp read task interface
type MbtcpReadTask interface {
	// Add add request to read/poll task map
	Add(name, tid, cmd string, req interface{})

	// GetByTID get request via TID from read/poll task map
	//GetByTID(tid string) (mbtcpReadTask, bool)
	GetByTID(tid string) (interface{}, bool)

	// GetByName get request via poll name from read/poll task map
	//GetByName(name string) (mbtcpReadTask, bool)
	GetByName(name string) (interface{}, bool)

	// GetAll get all requests from read/poll task map
	GetAll() []MbtcpPollStatus

	// DeleteAll remove all requests from read/poll task map
	DeleteAll()

	// DeleteByTID remove request from via TID from read/poll task map
	DeleteByTID(tid string)

	// DeleteByName remove request via poll name from read/poll task map
	DeleteByName(name string)

	// UpdateInterval update poll request interval
	UpdateInterval(name string, interval uint64) error

	// UpdateToggle update poll request enabled flag
	UpdateToggle(name string, toggle bool) error

	// UpdateAllToggles update all poll request enabled flag
	UpdateAllToggles(toggle bool)
}
