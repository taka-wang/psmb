package main

import "sync"

// NewMbtcpReadTask init mbtcp read task map
func NewMbtcpReadTask() MbtcpReadTask {
	return &mbtcpReadTaskType{
		m: make(map[string]mbtcpReadTask),
	}
}

// NewMbtcpSimpleTask init mbtcp simple task
func NewMbtcpSimpleTask() MbtcpSimpleTask {
	return &mbtcpSimpleTaskType{
		m: make(map[string]string),
	}
}

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

// mbtcpReadTaskType read/poll task map
type mbtcpReadTaskType struct {
	sync.RWMutex
	m map[string]mbtcpReadTask // m mbtcpReadTask map
}

// mbtcpSimpleTaskType simple task map type
type mbtcpSimpleTaskType struct {
	sync.RWMutex
	m map[string]string // tid: command
}

// Get get command from read/poll task map
func (s *mbtcpReadTaskType) Get(tid string) (mbtcpReadTask, bool) {
	s.RLock()
	task, ok := s.m[tid]
	s.RUnlock()
	return task, ok
}

// Delete remove command from read/poll task map
func (s *mbtcpReadTaskType) Delete(tid string) {
	s.Lock()
	delete(s.m, tid)
	s.Unlock()
}

// Add add cmd to read/poll task map
func (s *mbtcpReadTaskType) Add(name, tid, cmd string, req interface{}) {
	s.Lock()
	s.m[tid] = mbtcpReadTask{name, cmd, req}
	s.Unlock()
}

// Get get command from simple task map
func (s *mbtcpSimpleTaskType) Get(tid string) (string, bool) {
	s.RLock()
	cmd, ok := s.m[tid]
	s.RUnlock()
	return cmd, ok
}

// Delete remove command from simple task map
func (s *mbtcpSimpleTaskType) Delete(tid string) {
	s.Lock()
	delete(s.m, tid)
	s.Unlock()
}

// Add add cmd to simple task map
func (s *mbtcpSimpleTaskType) Add(tid, cmd string) {
	s.Lock()
	s.m[tid] = cmd
	s.Unlock()
}
