package mbtcp

import (
	"sync"

	psmb "github.com/taka-wang/psmb"
)

//
// Implementations
//

// @Implement WriterTaskMap contract implicitly

// WriteTaskType write task map type
type WriteTaskType struct {
	sync.RWMutex
	// m key-value map: (tid, command)
	m map[string]string
}

// NewMbtcpWriterMap instantiate mbtcp write task map
func NewMbtcpWriterMap() psmb.WriterTaskMap {
	return &WriteTaskType{
		m: make(map[string]string),
	}
}

// Add add request to write task map
func (s *WriteTaskType) Add(tid, cmd string) {
	s.Lock()
	s.m[tid] = cmd
	s.Unlock()
}

// Get get request from write task map
func (s *WriteTaskType) Get(tid string) (string, bool) {
	s.RLock()
	cmd, ok := s.m[tid]
	s.RUnlock()
	return cmd, ok
}

// Delete remove request from write task map
func (s *WriteTaskType) Delete(tid string) {
	s.Lock()
	delete(s.m, tid)
	s.Unlock()
}

// @Implement ReaderTaskMap contract implicitly

// mbtcpReadTask read/poll task request
type mbtcpReadTask struct {
	// Name task name
	Name string
	// Cmd zmq frame 1
	Cmd string
	// Req request structure
	Req interface{}
}

// ReadTaskType read/poll task map type
type ReadTaskType struct {
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

// NewMbtcpReaderMap instantiate mbtcp read task map
func NewMbtcpReaderMap() psmb.ReaderTaskMap {
	return &ReadTaskType{
		idName:  make(map[string]string),
		nameID:  make(map[string]string),
		idMap:   make(map[string]mbtcpReadTask),
		nameMap: make(map[string]mbtcpReadTask),
	}
}

// Add add request to read/poll task map
func (s *ReadTaskType) Add(name, tid, cmd string, req interface{}) {
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
//func (s *ReadTaskType) GetByTID(tid string) (mbtcpReadTask, bool) {
func (s *ReadTaskType) GetByTID(tid string) (interface{}, bool) {
	s.RLock()
	task, ok := s.idMap[tid]
	s.RUnlock()
	return task, ok
}

// GetByName get request via poll name from read/poll task map
//func (s *ReadTaskType) GetByName(name string) (mbtcpReadTask, bool) {
func (s *ReadTaskType) GetByName(name string) (interface{}, bool) {
	s.RLock()
	task, ok := s.nameMap[name]
	s.RUnlock()
	return task, ok
}

// GetAll get all requests from read/poll task map
func (s *ReadTaskType) GetAll() interface{} {
	arr := []psmb.MbtcpPollStatus{}
	s.RLock()
	for _, v := range s.nameMap {
		if item, ok := v.Req.(psmb.MbtcpPollStatus); ok { // type casting check!
			arr = append(arr, item)
		}
	}
	s.RUnlock()
	return arr
}

// DeleteAll remove all requests from read/poll task map
func (s *ReadTaskType) DeleteAll() {
	s.Lock()
	s.idName = make(map[string]string)
	s.nameID = make(map[string]string)
	s.idMap = make(map[string]mbtcpReadTask)
	s.nameMap = make(map[string]mbtcpReadTask)
	s.Unlock()
}

// DeleteByTID remove request from via TID from read/poll task map
func (s *ReadTaskType) DeleteByTID(tid string) {
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
func (s *ReadTaskType) DeleteByName(name string) {
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
func (s *ReadTaskType) UpdateInterval(name string, interval uint64) error {
	s.RLock()
	tid, _ := s.nameID[name]
	task, ok := s.nameMap[name]
	s.RUnlock()

	if !ok {
		return ErrInvalidPollName
	}

	req, ok2 := task.Req.(psmb.MbtcpPollStatus)
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
func (s *ReadTaskType) UpdateToggle(name string, toggle bool) error {
	s.RLock()
	tid, _ := s.nameID[name]
	task, ok := s.nameMap[name]
	s.RUnlock()

	if !ok {
		return ErrInvalidPollName
	}

	req, ok2 := task.Req.(psmb.MbtcpPollStatus) // type casting check!
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
func (s *ReadTaskType) UpdateAllToggles(toggle bool) {
	s.Lock()
	for name, task := range s.nameMap {
		if req, ok := task.Req.(psmb.MbtcpPollStatus); ok { // type casting check!
			req.Enabled = toggle                                 // update flag
			s.nameMap[name] = mbtcpReadTask{name, task.Cmd, req} // update nameMap table
			tid, _ := s.nameID[name]                             // get Tid
			s.idMap[tid] = s.nameMap[name]                       // update idMap table
		}
	}
	s.Unlock()
}
