package psmb

import "sync"

//
// Task Interfaces
//

// MbtcpSimpleTask mbtcp simple task interface
type MbtcpSimpleTask interface {
	// Add add request to simple task map
	Add(tid, cmd string)

	// Get get request from simple task map
	Get(tid string) (string, bool)

	// Delete remove request from simple task map
	Delete(tid string)
}

// MbtcpReadTask mbtcp read task interface
type MbtcpReadTask interface {
	// Add add request to read/poll task map
	Add(name, tid, cmd string, req interface{})

	// GetByTID get request via TID from read/poll task map
	GetByTID(tid string) (mbtcpReadTask, bool)

	// GetByName get request via poll name from read/poll task map
	GetByName(name string) (mbtcpReadTask, bool)

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

//
// Task types
//

// mbtcpSimpleTaskType simple task map type
type mbtcpSimpleTaskType struct {
	sync.RWMutex
	// m key-value map: (tid, command)
	m map[string]string
}

// mbtcpReadTask read/poll task request
type mbtcpReadTask struct {
	// Name task name
	Name string
	// Cmd zmq frame 1
	Cmd string
	// Req request structure
	Req interface{}
}

// mbtcpReadTaskType read/poll task map type
type mbtcpReadTaskType struct {
	sync.RWMutex
	// idName (tid, name)
	idName map[string]string
	// nameID (name, tid)
	nameID map[string]string
	// idMap (tid, mbtcpReadTask)
	idMap map[string]mbtcpReadTask
	// nameMap (name, mbtcpReadTask)
	nameMap map[string]mbtcpReadTask
}

//
// Implementations
//

//
// Simple Task
//

// NewMbtcpSimpleTask instantiate mbtcp simple task map
func NewMbtcpSimpleTask() MbtcpSimpleTask {
	return &mbtcpSimpleTaskType{
		m: make(map[string]string),
	}
}

// Add add request to simple task map
func (s *mbtcpSimpleTaskType) Add(tid, cmd string) {
	s.Lock()
	s.m[tid] = cmd
	s.Unlock()
}

// Get get request from simple task map
func (s *mbtcpSimpleTaskType) Get(tid string) (string, bool) {
	s.RLock()
	cmd, ok := s.m[tid]
	s.RUnlock()
	return cmd, ok
}

// Delete remove request from simple task map
func (s *mbtcpSimpleTaskType) Delete(tid string) {
	s.Lock()
	delete(s.m, tid)
	s.Unlock()
}

//
// Read/Poll Task
//

// NewMbtcpReadTask instantiate mbtcp read task map
func NewMbtcpReadTask() MbtcpReadTask {
	return &mbtcpReadTaskType{
		idName:  make(map[string]string),
		nameID:  make(map[string]string),
		idMap:   make(map[string]mbtcpReadTask),
		nameMap: make(map[string]mbtcpReadTask),
	}
}

// Add add request to read/poll task map
func (s *mbtcpReadTaskType) Add(name, tid, cmd string, req interface{}) {
	if name == "" { // read task instead of poll task
		name = tid
	}

	s.Lock()
	s.idName[tid] = name
	s.nameID[name] = tid
	s.nameMap[name] = mbtcpReadTask{name, cmd, req}
	s.idMap[tid] = s.nameMap[name]
	s.Unlock()
}

// GetByTID get request via TID from read/poll task map
func (s *mbtcpReadTaskType) GetByTID(tid string) (mbtcpReadTask, bool) {
	s.RLock()
	task, ok := s.idMap[tid]
	s.RUnlock()
	return task, ok
}

// GetByName get request via poll name from read/poll task map
func (s *mbtcpReadTaskType) GetByName(name string) (mbtcpReadTask, bool) {
	s.RLock()
	task, ok := s.nameMap[name]
	s.RUnlock()
	return task, ok
}

// GetAll get all requests from read/poll task map
func (s *mbtcpReadTaskType) GetAll() []MbtcpPollStatus {
	arr := []MbtcpPollStatus{}
	s.RLock()
	for _, v := range s.nameMap {
		if item, ok := v.Req.(MbtcpPollStatus); ok { // type casting check!
			arr = append(arr, item)
		}
	}
	s.RUnlock()
	return arr
}

// DeleteAll remove all requests from read/poll task map
func (s *mbtcpReadTaskType) DeleteAll() {
	s.Lock()
	s.idName = make(map[string]string)
	s.nameID = make(map[string]string)
	s.idMap = make(map[string]mbtcpReadTask)
	s.nameMap = make(map[string]mbtcpReadTask)
	s.Unlock()
}

// DeleteByTID remove request from via TID from read/poll task map
func (s *mbtcpReadTaskType) DeleteByTID(tid string) {
	s.RLock()
	name, ok := s.idName[tid]
	s.RUnlock()

	s.Lock()
	delete(s.idName, tid)
	delete(s.idMap, tid)
	if ok {
		delete(s.nameID, name)
		delete(s.nameMap, name)
	}
	s.Unlock()
}

// DeleteByName remove request via poll name from read/poll task map
func (s *mbtcpReadTaskType) DeleteByName(name string) {
	s.RLock()
	tid, ok := s.nameID[name]
	s.RUnlock()

	s.Lock()
	if ok {
		delete(s.idName, tid)
		delete(s.idMap, tid)
	}
	delete(s.nameID, name)
	delete(s.nameMap, name)
	s.Unlock()
}

// UpdateInterval update poll request interval
func (s *mbtcpReadTaskType) UpdateInterval(name string, interval uint64) error {
	s.RLock()
	tid, _ := s.nameID[name]
	task, ok := s.nameMap[name]
	s.RUnlock()

	if !ok {
		return ErrInvalidPollName
	}

	req, ok2 := task.Req.(MbtcpPollStatus)
	if !ok2 {
		return ErrInvalidPollName
	}

	req.Interval = interval // update interval
	s.Lock()
	s.nameMap[name] = mbtcpReadTask{name, task.Cmd, req} // update nameMap table
	s.idMap[tid] = s.nameMap[name]                       // update idMap table
	s.Unlock()
	return nil
}

// UpdateToggle update poll request enabled flag
func (s *mbtcpReadTaskType) UpdateToggle(name string, toggle bool) error {
	s.RLock()
	tid, _ := s.nameID[name]
	task, ok := s.nameMap[name]
	s.RUnlock()

	if !ok {
		return ErrInvalidPollName
	}

	req, ok2 := task.Req.(MbtcpPollStatus) // type casting check!
	if !ok2 {
		return ErrInvalidPollName
	}

	req.Enabled = toggle // update flag
	s.Lock()
	s.nameMap[name] = mbtcpReadTask{name, task.Cmd, req} // update nameMap table
	s.idMap[tid] = s.nameMap[name]                       // update idMap table
	s.Unlock()
	return nil
}

// UpdateAllToggles update all poll request enabled flag
func (s *mbtcpReadTaskType) UpdateAllToggles(toggle bool) {
	s.Lock()
	for name, task := range s.nameMap {
		if req, ok := task.Req.(MbtcpPollStatus); ok { // type casting check!
			req.Enabled = toggle                                 // update flag
			s.nameMap[name] = mbtcpReadTask{name, task.Cmd, req} // update nameMap table
			tid, _ := s.nameID[name]                             // get Tid
			s.idMap[tid] = s.nameMap[name]                       // update idMap table
		}
	}
	s.Unlock()
}
