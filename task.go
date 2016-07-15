package psmb

import "sync"

//
// Task Interfaces
//

// MbtcpReadTask mbtcp read task interface
type MbtcpReadTask interface {
	// Get get command from read/poll task map
	Get(tid string) (mbtcpReadTask, bool)
	// Delete remove command from read/poll task map
	Delete(tid string)
	// Add add cmd to read/poll task map
	Add(name, tid, cmd string, req interface{})
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
	// m key-value map: (tid, mbtcpReadTask)
	m map[string]mbtcpReadTask
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
		m: make(map[string]mbtcpReadTask),
	}
}

// NewMbtcpSimpleTask instantiate mbtcp simple task map
func NewMbtcpSimpleTask() MbtcpSimpleTask {
	return &mbtcpSimpleTaskType{
		m: make(map[string]string),
	}
}

// Get get command from read/poll task map
func (s *mbtcpReadTaskType) Get(tid string) (mbtcpReadTask, bool) {
	s.RLock()
	task, ok := s.m[tid]
	s.RUnlock()
	return task, ok
}

// Delete remove request from read/poll task map
func (s *mbtcpReadTaskType) Delete(tid string) {
	s.Lock()
	delete(s.m, tid)
	s.Unlock()
}

// Add add request to read/poll task map
func (s *mbtcpReadTaskType) Add(name, tid, cmd string, req interface{}) {
	s.Lock()
	s.m[tid] = mbtcpReadTask{name, cmd, req}
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

// Add add request to simple task map
func (s *mbtcpSimpleTaskType) Add(tid, cmd string) {
	s.Lock()
	s.m[tid] = cmd
	s.Unlock()
}
