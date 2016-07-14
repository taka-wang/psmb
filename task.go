package main

import "sync"

// MbtcpReadTaskType read/poll task map
type MbtcpReadTaskType struct {
	sync.RWMutex
	m map[string]mbtcpReadTask // m mbtcpReadTask map
}

// Get get command from read/poll task map
func (s *MbtcpReadTaskType) Get(tid string) (mbtcpReadTask, bool) {
	s.RLock()
	task, ok := s.m[tid]
	s.RUnlock()
	return task, ok
}

// Delete remove command from read/poll task map
func (s *MbtcpReadTaskType) Delete(tid string) {
	s.Lock()
	delete(s.m, tid)
	s.Unlock()
}

// Add add cmd to read/poll task map
func (s *MbtcpReadTaskType) Add(name, tid, cmd string, req interface{}) {
	s.Lock()
	s.m[tid] = mbtcpReadTask{name, cmd, req}
	s.Unlock()
}

// SimpleTaskType simple task map type
type SimpleTaskType struct {
	sync.RWMutex
	m map[string]string // tid: command
}

// Get get command from simple task map
func (s *SimpleTaskType) Get(tid string) (string, bool) {
	s.RLock()
	cmd, ok := s.m[tid]
	s.RUnlock()
	return cmd, ok
}

// Delete remove command from simple task map
func (s *SimpleTaskType) Delete(tid string) {
	s.Lock()
	delete(s.m, tid)
	s.Unlock()
}

// Add add cmd to simple task map
func (s *SimpleTaskType) Add(tid, cmd string) {
	s.Lock()
	s.m[tid] = cmd
	s.Unlock()
}
