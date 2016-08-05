// Package reader an in-memory data store for reader.
//
// By taka@cmwang.net
//
package reader

import (
	"errors"
	"sync"

	psmb "github.com/taka-wang/psmb"
)

// ErrInvalidPollName is the error when the poll name is empty.
var ErrInvalidPollName = errors.New("Invalid poll name!")

// @Implement IReaderTaskDataStore contract implicitly

// readerTaskDataStore read/poll task map type
type readerTaskDataStore struct {
	sync.RWMutex
	// idName (tid, name)
	idName map[string]string
	// nameID (name, tid)
	nameID map[string]string
	// idMap (tid, ReaderTask)
	idMap map[string]psmb.ReaderTask
	// nameMap (name, ReaderTask)
	nameMap map[string]psmb.ReaderTask
}

// NewDataStore instantiate mbtcp read task map
func NewDataStore(conf map[string]string) (interface{}, error) {
	return &readerTaskDataStore{
		idName:  make(map[string]string),
		nameID:  make(map[string]string),
		idMap:   make(map[string]psmb.ReaderTask),
		nameMap: make(map[string]psmb.ReaderTask),
	}, nil
}

// Add add request to read/poll task map
func (ds *readerTaskDataStore) Add(name, tid, cmd string, req interface{}) {
	if name == "" { // read task instead of poll task
		name = tid
	}

	ds.Lock()
	ds.idName[tid] = name
	ds.nameID[name] = tid
	ds.nameMap[name] = psmb.ReaderTask{Name: name, Cmd: cmd, Req: req}
	ds.idMap[tid] = ds.nameMap[name]
	ds.Unlock()
}

// GetTaskByID get request via TID from read/poll task map
// 	interface{}: ReaderTask
func (ds *readerTaskDataStore) GetTaskByID(tid string) (interface{}, bool) {
	ds.RLock()
	task, ok := ds.idMap[tid]
	ds.RUnlock()
	return task, ok
}

// GetTaskByName get request via poll name from read/poll task map
// 	interface{}: ReaderTask
func (ds *readerTaskDataStore) GetTaskByName(name string) (interface{}, bool) {
	ds.RLock()
	task, ok := ds.nameMap[name]
	ds.RUnlock()
	return task, ok
}

// GetAll get all requests from read/poll task map
//	interface{}: []psmb.MbtcpPollStatus
func (ds *readerTaskDataStore) GetAll() interface{} {
	arr := []psmb.MbtcpPollStatus{}
	ds.RLock()
	for _, v := range ds.nameMap {
		// type casting check!
		if item, ok := v.Req.(psmb.MbtcpPollStatus); ok {
			arr = append(arr, item)
		}
	}
	ds.RUnlock()
	return arr
}

// DeleteAll remove all requests from read/poll task map
func (ds *readerTaskDataStore) DeleteAll() {
	ds.Lock()
	ds.idName = make(map[string]string)
	ds.nameID = make(map[string]string)
	ds.idMap = make(map[string]psmb.ReaderTask)
	ds.nameMap = make(map[string]psmb.ReaderTask)
	ds.Unlock()
}

// DeleteTaskByID remove request via TID from read/poll task map
func (ds *readerTaskDataStore) DeleteTaskByID(tid string) {
	ds.RLock()
	name, ok := ds.idName[tid]
	ds.RUnlock()

	ds.Lock()
	delete(ds.idName, tid)
	delete(ds.idMap, tid)
	if ok {
		delete(ds.nameID, name)
		delete(ds.nameMap, name)
	}
	ds.Unlock()
}

// DeleteTaskByName remove request via poll name from read/poll task map
func (ds *readerTaskDataStore) DeleteTaskByName(name string) {
	ds.RLock()
	tid, ok := ds.nameID[name]
	ds.RUnlock()

	ds.Lock()
	if ok {
		delete(ds.idName, tid)
		delete(ds.idMap, tid)
	}
	delete(ds.nameID, name)
	delete(ds.nameMap, name)
	ds.Unlock()
}

// UpdateIntervalByName update poll request interval
func (ds *readerTaskDataStore) UpdateIntervalByName(name string, interval uint64) error {
	ds.RLock()
	tid, _ := ds.nameID[name]
	task, ok := ds.nameMap[name]
	ds.RUnlock()

	if !ok {
		return ErrInvalidPollName
	}

	req, ok2 := task.Req.(psmb.MbtcpPollStatus)
	if !ok2 {
		return ErrInvalidPollName
	}

	req.Interval = interval // update interval
	ds.Lock()
	ds.nameMap[name] = psmb.ReaderTask{Name: name, Cmd: task.Cmd, Req: req} // update nameMap table
	ds.idMap[tid] = ds.nameMap[name]                                        // update idMap table
	ds.Unlock()
	return nil
}

// UpdateToggleByName update poll request enabled flag
func (ds *readerTaskDataStore) UpdateToggleByName(name string, toggle bool) error {
	ds.RLock()
	tid, _ := ds.nameID[name]
	task, ok := ds.nameMap[name]
	ds.RUnlock()

	if !ok {
		return ErrInvalidPollName
	}
	// type casting check!
	req, ok2 := task.Req.(psmb.MbtcpPollStatus)
	if !ok2 {
		return ErrInvalidPollName
	}

	req.Enabled = toggle // update flag
	ds.Lock()
	ds.nameMap[name] = psmb.ReaderTask{Name: name, Cmd: task.Cmd, Req: req} // update nameMap table
	ds.idMap[tid] = ds.nameMap[name]                                        // update idMap table
	ds.Unlock()
	return nil
}

// UpdateAllTogglesByName update all poll request enabled flag
func (ds *readerTaskDataStore) UpdateAllTogglesByName(toggle bool) {
	ds.Lock()
	for name, task := range ds.nameMap {
		// type casting check!
		if req, ok := task.Req.(psmb.MbtcpPollStatus); ok {
			req.Enabled = toggle                                                    // update flag
			ds.nameMap[name] = psmb.ReaderTask{Name: name, Cmd: task.Cmd, Req: req} // update nameMap table
			tid, _ := ds.nameID[name]                                               // get Tid
			ds.idMap[tid] = ds.nameMap[name]                                        // update idMap table
		}
	}
	ds.Unlock()
}
