// Package filter an in-memory data store for filter.
//
// Guideline: if error is one of the return, don't duplicately log to output.
//
// By taka@cmwang.net
//
package filter

import (
	"sync"

	"github.com/taka-wang/psmb"
	//conf "github.com/taka-wang/psmb/mini-conf"
	conf "github.com/taka-wang/psmb/viper-conf"
)

var maxCapacity int

func init() {
	conf.SetDefault(keyMaxCapacity, defaultMaxCapacity)
	maxCapacity = conf.GetInt(keyMaxCapacity)
}

//@Implement IFilterDataStore implicitly

// dataStore filter map
type dataStore struct {
	// read writer mutex
	sync.RWMutex
	// m key-value map: (name, psmb.MbtcpFilterStatus)
	m map[string]interface{}
}

// NewDataStore instantiate filter map
func NewDataStore(conf map[string]string) (interface{}, error) {
	return &dataStore{
		m: make(map[string]interface{}),
	}, nil
}

// Add add request to filter map
func (ds *dataStore) Add(name string, req interface{}) error {
	if name == "" {
		return ErrInvalidFilterName
	}

	ds.RLock()
	boom := len(ds.m)+1 > maxCapacity
	ds.RUnlock()
	if boom {
		return ErrOutOfCapacity
	}

	ds.Lock()
	ds.m[name] = req
	ds.Unlock()
	return nil
}

// Get get request from filter map
func (ds *dataStore) Get(name string) (interface{}, bool) {
	ds.RLock()
	req, ok := ds.m[name]
	ds.RUnlock()
	return req, ok
}

// GetAll get all requests from filter map
func (ds *dataStore) GetAll() interface{} {
	arr := []psmb.MbtcpFilterStatus{}
	ds.RLock()
	for _, v := range ds.m {
		arr = append(arr, v.(psmb.MbtcpFilterStatus))
	}
	ds.RUnlock()

	if len(arr) == 0 {
		err := ErrNoData
		conf.Log.WithError(err).Warn("Fail to get all items from filter data store")
		return nil
	}
	return arr
}

// Delete remove request from filter map
func (ds *dataStore) Delete(name string) {
	ds.Lock()
	delete(ds.m, name)
	ds.Unlock()
}

// DeleteAll delete all filters from filter map
func (ds *dataStore) DeleteAll() {
	ds.Lock()
	ds.m = make(map[string]interface{})
	ds.Unlock()
}

// Toggle toggle request from filter map
func (ds *dataStore) UpdateToggle(name string, toggle bool) error {
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
func (ds *dataStore) UpdateAllToggles(toggle bool) {
	ds.Lock()
	for name, req := range ds.m {
		r := req.(psmb.MbtcpFilterStatus)
		r.Enabled = toggle
		ds.m[name] = r
	}
	ds.Unlock()
}
