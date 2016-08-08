package filter

import (
	"errors"
	"sync"

	"github.com/taka-wang/psmb"
)

// ErrInvalidFilterName is the error when the name is invalid
var ErrInvalidFilterName = errors.New("Invalid Filter name")

//@Implement IFilterDataStore implicitly

// NewDataStore instantiate filter map
func NewDataStore(conf map[string]string) (interface{}, error) {
	return &filterDataStore{
		m: make(map[string]interface{}),
	}, nil
}

// filterDataStore filter map
type filterDataStore struct {
	sync.RWMutex
	// m key-value map: (name, psmb.MbtcpFilterStatus)
	m map[string]interface{}
}

// Add add request to filter map
func (ds *filterDataStore) Add(name string, req interface{}) {
	ds.Lock()
	ds.m[name] = req
	ds.Unlock()
}

// Get get request from filter map
func (ds *filterDataStore) Get(name string) (interface{}, bool) {
	ds.RLock()
	req, ok := ds.m[name]
	ds.RUnlock()
	return req, ok
}

// GetAll get request from filter map
func (ds *filterDataStore) GetAll(name string) interface{} {
	arr := []psmb.MbtcpFilterStatus{}
	ds.RLock()
	for _, v := range ds.m {
		arr = append(arr, v.(psmb.MbtcpFilterStatus))
	}
	ds.RUnlock()
	return arr
}

// Delete remove request from filter map
func (ds *filterDataStore) Delete(name string) {
	ds.Lock()
	delete(ds.m, name)
	ds.Unlock()
}

// DeleteAll delete all filters
func (ds *filterDataStore) DeleteAll() {
	ds.Lock()
	ds.m = make(map[string]interface{})
	ds.Unlock()
}

// Toggle toggle request from filter map
func (ds *filterDataStore) UpdateToggle(name string, toggle bool) error {
	ds.RLock()
	req, ok := ds.m[name]
	ds.RUnlock()
	if !ok {
		return ErrInvalidFilterName
	}
	r := req.(psmb.MbtcpFilterStatus)
	r.Enabled = toggle
	ds.Lock()
	ds.m[name] = r
	ds.Unlock()

	return nil
}

// UpdateAllToggles toggle all request from filter map
func (ds *filterDataStore) UpdateAllToggles(toggle bool) {
	ds.Lock()
	for name, req := range ds.m {
		r := req.(psmb.MbtcpFilterStatus)
		r.Enabled = toggle
		ds.m[name] = r
	}
	ds.Unlock()
}
