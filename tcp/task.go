package tcp

import (
	"sync"

	psmb "github.com/taka-wang/psmb"
)

//
// Implementations
//

// @Implement IReaderTaskMap contract implicitly

// ReaderTask read/poll task request
type ReaderTask struct {
	// Name task name
	Name string
	// Cmd zmq frame 1
	Cmd string
	// Req request structure
	Req interface{}
}

// ReaderTaskType read/poll task map type
type ReaderTaskType struct {
	sync.RWMutex
	// idName (tid, name)
	idName map[string]string
	// nameID (name, tid)
	nameID map[string]string
	// idMap (tid, ReaderTask)
	idMap map[string]ReaderTask
	// nameMap (name, ReaderTask)
	nameMap map[string]ReaderTask
}

// NewReaderMap instantiate mbtcp read task map
func NewReaderMap() psmb.IReaderTaskMap {
	return &ReaderTaskType{
		idName:  make(map[string]string),
		nameID:  make(map[string]string),
		idMap:   make(map[string]ReaderTask),
		nameMap: make(map[string]ReaderTask),
	}
}

// Add add request to read/poll task map
func (s *ReaderTaskType) Add(name, tid, cmd string, req interface{}) {
	if name == "" { // read task instead of poll task
		name = tid
	}

	s.Lock()
	s.idName[tid] = name
	s.nameID[name] = tid
	s.nameMap[name] = ReaderTask{name, cmd, req}
	s.idMap[tid] = s.nameMap[name]
	s.Unlock()
}

// GetTaskByID get request via TID from read/poll task map
// interface{}: ReaderTask
func (s *ReaderTaskType) GetTaskByID(tid string) (interface{}, bool) {
	s.RLock()
	task, ok := s.idMap[tid]
	s.RUnlock()
	return task, ok
}

// GetTaskByName get request via poll name from read/poll task map
// 	interface{}: ReaderTask
func (s *ReaderTaskType) GetTaskByName(name string) (interface{}, bool) {
	s.RLock()
	task, ok := s.nameMap[name]
	s.RUnlock()
	return task, ok
}

// GetAll get all requests from read/poll task map
func (s *ReaderTaskType) GetAll() interface{} {
	arr := []psmb.MbtcpPollStatus{}
	s.RLock()
	for _, v := range s.nameMap {
		// type casting check!
		if item, ok := v.Req.(psmb.MbtcpPollStatus); ok {
			arr = append(arr, item)
		}
	}
	s.RUnlock()
	return arr
}

// DeleteAll remove all requests from read/poll task map
func (s *ReaderTaskType) DeleteAll() {
	s.Lock()
	s.idName = make(map[string]string)
	s.nameID = make(map[string]string)
	s.idMap = make(map[string]ReaderTask)
	s.nameMap = make(map[string]ReaderTask)
	s.Unlock()
}

// DeleteTaskByID remove request from via TID from read/poll task map
func (s *ReaderTaskType) DeleteTaskByID(tid string) {
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

// DeleteTaskByName remove request via poll name from read/poll task map
func (s *ReaderTaskType) DeleteTaskByName(name string) {
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

// UpdateIntervalByName update poll request interval
func (s *ReaderTaskType) UpdateIntervalByName(name string, interval uint64) error {
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
	s.nameMap[name] = ReaderTask{name, task.Cmd, req} // update nameMap table
	s.idMap[tid] = s.nameMap[name]                    // update idMap table
	s.Unlock()
	return nil
}

// UpdateToggleByName update poll request enabled flag
func (s *ReaderTaskType) UpdateToggleByName(name string, toggle bool) error {
	s.RLock()
	tid, _ := s.nameID[name]
	task, ok := s.nameMap[name]
	s.RUnlock()

	if !ok {
		return ErrInvalidPollName
	}
	// type casting check!
	req, ok2 := task.Req.(psmb.MbtcpPollStatus)
	if !ok2 {
		return ErrInvalidPollName
	}

	req.Enabled = toggle // update flag
	s.Lock()
	s.nameMap[name] = ReaderTask{name, task.Cmd, req} // update nameMap table
	s.idMap[tid] = s.nameMap[name]                    // update idMap table
	s.Unlock()
	return nil
}

// UpdateAllTogglesByName update all poll request enabled flag
func (s *ReaderTaskType) UpdateAllTogglesByName(toggle bool) {
	s.Lock()
	for name, task := range s.nameMap {
		// type casting check!
		if req, ok := task.Req.(psmb.MbtcpPollStatus); ok {
			req.Enabled = toggle                              // update flag
			s.nameMap[name] = ReaderTask{name, task.Cmd, req} // update nameMap table
			tid, _ := s.nameID[name]                          // get Tid
			s.idMap[tid] = s.nameMap[name]                    // update idMap table
		}
	}
	s.Unlock()
}
