package psmb

import "sync"

//
// Task Interfaces
//

// MbtcpReadTask mbtcp read task interface
type MbtcpReadTask interface {
	// GetTID get command from read/poll task map
	GetTID(tid string) (mbtcpReadTask, bool)
	// GetName get command from read/poll task map
	GetName(name string) (mbtcpReadTask, bool)
	// Delete remove command from read/poll task map
	Delete(tid string)
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
	// m key-value map: (tid, mbtcpReadTask)
	m map[string]mbtcpReadTask
	// p key-value map: (name, tid)
	p map[string]string
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
		p: make(map[string]string),
	}
}

// NewMbtcpSimpleTask instantiate mbtcp simple task map
func NewMbtcpSimpleTask() MbtcpSimpleTask {
	return &mbtcpSimpleTaskType{
		m: make(map[string]string),
	}
}

// GetTID get command from read/poll task map
func (s *mbtcpReadTaskType) GetTID(tid string) (mbtcpReadTask, bool) {
	s.RLock()
	task, ok := s.m[tid]
	s.RUnlock()
	return task, ok
}

// GetName get command from read/poll task map
func (s *mbtcpReadTaskType) GetName(name string) (mbtcpReadTask, bool) {
	s.RLock()
	if id, ok1 := s.p[name]; ok1 {
		task, ok2 := s.m[id]
		s.RUnlock()
		return task, ok2
	}
	s.RUnlock()
	return mbtcpReadTask{}, false
}

// Delete remove request from read/poll task map
func (s *mbtcpReadTaskType) Delete(tid string) {
	s.RLock()
	task, ok := s.m[tid]
	s.RUnlock()

	// remove from p table
	if ok && task.Name != "" {
		s.Lock()
		delete(s.p, task.Name)
		s.Unlock()
	}

	// remove from m table
	s.Lock()
	delete(s.m, tid)
	s.Unlock()
}

// Add add request to read/poll task map
func (s *mbtcpReadTaskType) Add(name, tid, cmd string, req interface{}) {
	s.Lock()
	// add to m table
	s.m[tid] = mbtcpReadTask{name, cmd, req}
	// add to p table
	if name != "" {
		s.p[name] = tid
	}
	s.Unlock()
}

// UpdateInterval update request interval
func (s *mbtcpReadTaskType) UpdateInterval(name string, interval uint64) error {
	s.RLock()
	id, ok1 := s.p[name]
	if !ok1 {
		s.RUnlock()
		return ErrInvalidPollName
	}
	task, ok2 := s.m[id]
	s.RUnlock()

	if !ok2 {
		return ErrInvalidPollName
	}

	req := task.Req.(MbtcpPollStatus)
	req.Interval = interval

	s.Lock()
	s.m[tid] = mbtcpReadTask{name, task.Cmd, req}
	s.Unlock()
	return nil
}

// UpdateToggle update request enable flag
func (s *mbtcpReadTaskType) UpdateToggle(name string, toggle bool) error {
	s.RLock()
	id, ok1 := s.p[name]
	if !ok1 {
		s.RUnlock()
		return ErrInvalidPollName
	}
	task, ok2 := s.m[id]
	s.RUnlock()

	if !ok2 {
		return ErrInvalidPollName
	}

	req := task.Req.(MbtcpPollStatus)
	req.Enabled = toggle

	s.Lock()
	s.m[tid] = mbtcpReadTask{name, task.Cmd, req}
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
