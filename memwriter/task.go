package memwriter

import (
	"sync"

	psmb "github.com/taka-wang/psmb"
	psmbtcp "github.com/taka-wang/psmb/tcp"
)

// init register writer task
func init() {
	psmbtcp.RegisterWriterTask("memory", NewWriterTaskDataStore)
}

// @Implement IWriterTaskDataStore contract implicitly

// WriterTaskDataStore write task map type
type WriterTaskDataStore struct {
	sync.RWMutex
	// m key-value map: (tid, command)
	m map[string]string
}

// NewWriterTaskDataStore instantiate mbtcp write task map
func NewWriterTaskDataStore(conf map[string]string) (psmb.IWriterTaskDataStore, error) {
	return &WriterTaskDataStore{
		m: make(map[string]string),
	}, nil
}

// Add add request to write task map
func (s *WriterTaskDataStore) Add(tid, cmd string) {
	s.Lock()
	s.m[tid] = cmd
	s.Unlock()
}

// Get get request from write task map
func (s *WriterTaskDataStore) Get(tid string) (string, bool) {
	s.RLock()
	cmd, ok := s.m[tid]
	s.RUnlock()
	return cmd, ok
}

// Delete remove request from write task map
func (s *WriterTaskDataStore) Delete(tid string) {
	s.Lock()
	delete(s.m, tid)
	s.Unlock()
}
