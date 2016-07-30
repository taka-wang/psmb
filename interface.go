package psmb

//
// Interfaces
//

// IProactiveService proactive service contracts,
// all services should implement the following methods.
type IProactiveService interface {
	// Start enable proactive service
	Start()
	// Stop disable proactive service
	Stop()
	// ParseRequest parse requests from IT
	ParseRequest(msg []string) (interface{}, error)
	// HandleRequest handle requests from IT
	HandleRequest(cmd string, r interface{}) error
	// ParseResponse parse responses from OT
	ParseResponse(msg []string) (interface{}, error)
	// HandleResponse handle responses from OT
	HandleResponse(cmd string, r interface{}) error
}

// IWriterTaskMap write task interface
//	(Tid, Command) map
type IWriterTaskMap interface {
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

// IReaderTaskMap read task interface
type IReaderTaskMap interface {
	// Add add request to read/poll task map
	Add(name, tid, cmd string, req interface{})

	// GetByTID get request via TID from read/poll task map
	GetByTID(tid string) (interface{}, bool)

	// GetByName get request via poll name from read/poll task map
	GetByName(name string) (interface{}, bool)

	// GetAll get all requests from read/poll task map
	//	ex: mbtcp: []MbtcpPollStatus
	GetAll() interface{}

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
