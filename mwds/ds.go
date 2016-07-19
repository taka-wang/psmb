// Package mwds an in-memory data store for writer.
//
// By taka@cmwang.net
//
package mwds

import "sync"

// @Implement IWriterTaskDataStore contract implicitly

// writerTaskDataStore write task map type
type writerTaskDataStore struct {
	sync.RWMutex
	// m key-value map: (tid, command)
	m map[string]string
}

// NewDataStore instantiate mbtcp write task map
func NewDataStore(conf map[string]string) (interface{}, error) {
	return &writerTaskDataStore{
		m: make(map[string]string),
	}, nil
}

// Add add request to write task map
func (ds *writerTaskDataStore) Add(tid, cmd string) {
	ds.Lock()
	ds.m[tid] = cmd
	ds.Unlock()
}

// Get get request from write task map
func (ds *writerTaskDataStore) Get(tid string) (string, bool) {
	ds.RLock()
	cmd, ok := ds.m[tid]
	ds.RUnlock()
	return cmd, ok
}

// Delete remove request from write task map
func (ds *writerTaskDataStore) Delete(tid string) {
	ds.Lock()
	delete(ds.m, tid)
	ds.Unlock()
}
