package filter

import (
	"sync"

	"github.com/taka-wang/psmb"
	log "github.com/takawang/logrus"
)

//@Implement IFilterDataStore implicitly

// NewDataStore instantiate filter map
func NewDataStore(conf map[string]string) (interface{}, error) {
	return &dataStore{
		m: make(map[string]interface{}),
	}, nil
}

// dataStore filter map
type dataStore struct {
	sync.RWMutex
	// m key-value map: (name, psmb.MbtcpFilterStatus)
	m map[string]interface{}
}

// Add add request to filter map
func (ds *dataStore) Add(name string, req interface{}) {
	ds.Lock()
	ds.m[name] = req
	ds.Unlock()
}

// Get get request from filter map
func (ds *dataStore) Get(name string) (interface{}, bool) {
	ds.RLock()
	req, ok := ds.m[name]
	ds.RUnlock()
	return req, ok
}

// GetAll get all requests from filter map
func (ds *dataStore) GetAll(name string) interface{} {
	arr := []psmb.MbtcpFilterStatus{}
	ds.RLock()
	for _, v := range ds.m {
		arr = append(arr, v.(psmb.MbtcpFilterStatus))
	}
	ds.RUnlock()

	if len(arr) == 0 {
		err := ErrNoData
		log.WithFields(log.Fields{"err": err}).Error("GetAll")
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
