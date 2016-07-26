package psmb

import "sync"

//
// Task Interfaces
//

// MbtcpReadTask mbtcp read task interface
type MbtcpReadTask interface {
	// GetTID get command via TID from read/poll task map
	GetTID(tid string) (mbtcpReadTask, bool)
	// GetName get command via Name from read/poll task map
	GetName(name string) (mbtcpReadTask, bool)
	// DeleteTID remove request via TID from read/poll task map
	DeleteTID(tid string)
	// DeleteName remove request via Name from read/poll task map
	DeleteName(name string)
	// Add add cmd to read/poll task map
	Add(name, tid, cmd string, req interface{})
	// UpdateInterval update request interval
	UpdateInterval(name string, interval uint64) error
	// UpdateToggle update request enable flag
	UpdateToggle(name string, toggle bool) error
}

// MbtcpSimpleTask mbtcp simple task interface
type MbtcpSimpleTask interface {
	// Get get command from simple task map
	Get(tid string) (string, bool)
	// Delete remove command from simple task map
	Delete(tid string)
	// Add add cmd to simple task map
	Add(tid, cmd string)
}

//
// Task types
//

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

// mbtcpSimpleTaskType simple task map type
type mbtcpSimpleTaskType struct {
	sync.RWMutex
	// m key-value map: (tid, command)
	m map[string]string
}

//
// Implementations
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

// NewMbtcpSimpleTask instantiate mbtcp simple task map
func NewMbtcpSimpleTask() MbtcpSimpleTask {
	return &mbtcpSimpleTaskType{
		m: make(map[string]string),
	}
}

// Add add request to read/poll task map
func (s *mbtcpReadTaskType) Add(name, tid, cmd string, req interface{}) {
	if name == "" {
		name = tid
	}
	s.Lock()
	s.idName[tid] = name
	s.nameID[name] = tid
	s.nameMap[name] = mbtcpReadTask{name, cmd, req}
	s.idMap[tid] = s.nameMap[name]
	s.Unlock()
}

// GetTID get command from read/poll task map
func (s *mbtcpReadTaskType) GetTID(tid string) (mbtcpReadTask, bool) {
	s.RLock()
	task, ok := s.idMap[tid]
	s.RUnlock()
	return task, ok
}

// GetName get command from read/poll task map
func (s *mbtcpReadTaskType) GetName(name string) (mbtcpReadTask, bool) {
	s.RLock()
	task, ok := s.nameMap[name]
	s.RUnlock()
	return task, ok
}

// DeleteTID remove request from read/poll task map
func (s *mbtcpReadTaskType) DeleteTID(tid string) {
	s.RLock()
	name, _ := s.idName[tid]
	s.RUnlock()
	s.Lock()
	delete(s.idName, tid)
	delete(s.idMap, tid)
	delete(s.nameID, name)
	delete(s.nameMap, name)
	s.Unlock()
}

// DeleteName remove request from read/poll task map
func (s *mbtcpReadTaskType) DeleteName(name string) {
	s.RLock()
	tid, _ := s.nameID[name]
	s.RUnlock()
	s.Lock()
	delete(s.idName, tid)
	delete(s.idMap, tid)
	delete(s.nameID, name)
	delete(s.nameMap, name)
	s.Unlock()
}

// UpdateInterval update request interval
func (s *mbtcpReadTaskType) UpdateInterval(name string, interval uint64) error {
	s.RLock()
	tid, _ := s.nameID[name]
	task, ok := s.nameMap[name]
	s.RUnlock()
	if !ok {
		return ErrInvalidPollName
	}
	req := task.Req.(MbtcpPollStatus)
	req.Interval = interval
	s.Lock()
	s.nameMap[name] = mbtcpReadTask{name, task.Cmd, req}
	s.idMap[tid] = s.nameMap[name]
	s.Unlock()
	return nil
}

// UpdateToggle update request enable flag
func (s *mbtcpReadTaskType) UpdateToggle(name string, toggle bool) error {
	s.RLock()
	tid, _ := s.nameID[name]
	task, ok := s.nameMap[name]
	s.RUnlock()
	if !ok {
		return ErrInvalidPollName
	}
	req := task.Req.(MbtcpPollStatus)
	req.Enabled = toggle
	s.Lock()
	s.nameMap[name] = mbtcpReadTask{name, task.Cmd, req}
	s.idMap[tid] = s.nameMap[name]
	s.Unlock()
	return nil
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

// Add add request to simple task map
func (s *mbtcpSimpleTaskType) Add(tid, cmd string) {
	s.Lock()
	s.m[tid] = cmd
	s.Unlock()
}
