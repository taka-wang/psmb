// Package writer an in-memory data store for writer.
//
// By taka@cmwang.net
//
package writer

import "sync"

// @Implement IWriterTaskDataStore contract implicitly

// dataStore write task map type
type dataStore struct {
	sync.RWMutex
	// m key-value map: (tid, command)
	m map[string]string
}

// NewDataStore instantiate mbtcp write task map
func NewDataStore(conf map[string]string) (interface{}, error) {
	return &dataStore{
		m: make(map[string]string),
	}, nil
}

// Add add request to write task map
func (ds *dataStore) Add(tid, cmd string) {
	ds.Lock()
	ds.m[tid] = cmd
	ds.Unlock()
}

// Get get request from write task map
func (ds *dataStore) Get(tid string) (string, bool) {
	ds.RLock()
	cmd, ok := ds.m[tid]
	ds.RUnlock()
	return cmd, ok
}

// Delete remove request from write task map
func (ds *dataStore) Delete(tid string) {
	ds.Lock()
	delete(ds.m, tid)
	ds.Unlock()
}
