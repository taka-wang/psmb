// Package filter an redis-based data store for filter.
//
// Guideline: if error is one of the return, don't duplicately log to output.
//
//
// By taka@cmwang.net
//
package filter

import (
	"encoding/json"
	"net"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/taka-wang/psmb"
	//"github.com/taka-wang/psmb/mini-conf"
	"github.com/taka-wang/psmb/viper-conf"
)

var (
	hashName    string
	maxCapacity int
)

func setDefaults() {
	// set default redis values
	conf.SetDefault(keyRedisServer, defaultRedisServer)
	conf.SetDefault(keyRedisPort, defaultRedisPort)
	conf.SetDefault(keyRedisMaxIdel, defaultRedisMaxIdel)
	conf.SetDefault(keyRedisMaxActive, defaultRedisMaxActive)
	conf.SetDefault(keyRedisIdelTimeout, defaultRedisIdelTimeout)
	conf.SetDefault(keyMaxCapacity, defaultMaxCapacity)

	// set default redis-writer values
	conf.SetDefault(keyHashName, defaultHashName)

	// Note: for docker environment
	// lookup redis server
	host, err := net.LookupHost(defaultRedisDocker)
	if err != nil {
		conf.Log.WithError(err).Debug("Local run")
	} else {
		conf.Log.WithField("hostname", host[0]).Info("Docker run")
		conf.Set(keyRedisServer, host[0]) // override default
	}
}

func init() {
	setDefaults() // set defaults
	hashName = conf.GetString(keyHashName)
	maxCapacity = conf.GetInt(keyMaxCapacity)
}

//@Implement IFilterDataStore implicitly

// dataStore filter map
type dataStore struct {
	mutex sync.Mutex
	count int
	pool  *redis.Pool
}

// NewDataStore instantiate filter map
func NewDataStore(c map[string]string) (interface{}, error) {
	return &dataStore{
		pool: &redis.Pool{
			MaxIdle: conf.GetInt(keyRedisMaxIdel),
			// When zero, there is no limit on the number of connections in the pool.
			MaxActive:   conf.GetInt(keyRedisMaxActive),
			IdleTimeout: conf.GetDuration(keyRedisIdelTimeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", conf.GetString(keyRedisServer)+":"+conf.GetString(keyRedisPort))
				if err != nil {
					conf.Log.WithError(err).Error("Redis pool dial error")
				}
				return conn, err
			},
		},
	}, nil
}

// Add add request to filter map
func (ds *dataStore) Add(name string, req interface{}) error {
	if ds.count+1 > maxCapacity {
		return ErrOutOfCapacity
	}

	ds.mutex.Lock() // lock
	conn := ds.pool.Get()
	defer conn.Close()
	defer ds.mutex.Unlock() // unlock

	// marshal
	bytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	if _, err := conn.Do("HSET", hashName, name, string(bytes)); err != nil {
		return err
	}

	ret, err := redis.Int(conn.Do("HLEN", hashName))
	if err != nil {
		return err
	}

	ds.count = ret // update count
	return nil
}

// Get get request from filter map
func (ds *dataStore) Get(name string) (interface{}, bool) {
	if name == "" {
		return nil, false
	}

	ds.mutex.Lock() // lock
	conn := ds.pool.Get()
	defer conn.Close()

	ret, err := redis.String(conn.Do("HGET", hashName, name))
	ds.mutex.Unlock() // unlock
	if err != nil {
		// we intend to suppress this log
		//conf.Log.WithError(err).Warn("Fail to get item from filter map")
		return nil, false
	}
	// unmarshal
	var d psmb.MbtcpFilterStatus
	if err := json.Unmarshal([]byte(ret), &d); err != nil {
		conf.Log.WithError(ErrUnmarshal).Error("Fail to unmarshal items from filter map")
		return nil, false
	}
	return d, true
}

// GetAll get all requests from filter map
func (ds *dataStore) GetAll() interface{} {
	ds.mutex.Lock() // lock
	conn := ds.pool.Get()
	defer conn.Close()

	ret, err := redis.StringMap(conn.Do("HGETALL", hashName))
	ds.mutex.Unlock() // unlock
	if err != nil {
		conf.Log.WithError(err).Warn("Fail to get all items from filter map")
		return nil
	}

	arr := []psmb.MbtcpFilterStatus{}
	for _, v := range ret {
		var d psmb.MbtcpFilterStatus
		if err := json.Unmarshal([]byte(v), &d); err == nil {
			arr = append(arr, d)
		}
	}

	if len(arr) == 0 {
		conf.Log.WithError(ErrNoData).Warn("Filter map is empty")
		return nil
	}
	return arr
}

// Delete remove request from filter map
func (ds *dataStore) Delete(name string) {
	if name == "" {
		conf.Log.WithError(ErrInvalidName).Warn("Fail to delete item from filter map")
		return
	}

	ds.mutex.Lock() // lock
	conn := ds.pool.Get()
	defer conn.Close()
	defer ds.mutex.Unlock() // unlock

	// delete item
	if _, err := conn.Do("HDEL", hashName, name); err != nil {
		conf.Log.WithError(err).Warn("Fail to delete item from filter map")
		return
	}
	// get length
	ret, err := redis.Int(conn.Do("HLEN", hashName))
	if err != nil {
		conf.Log.WithError(err).Warn("Fail to get length from filter map")
		return
	}
	ds.count = ret // update count
}

// DeleteAll delete all filters from filter map
func (ds *dataStore) DeleteAll() {
	ds.mutex.Lock() // lock
	conn := ds.pool.Get()
	defer conn.Close()
	defer ds.mutex.Unlock() // unlock

	if _, err := conn.Do("DEL", hashName); err != nil {
		conf.Log.WithError(err).Warn("Fail to delete all items from filter map")
		return
	}
	ds.count = 0 // reset
}

// Toggle toggle request from filter map
func (ds *dataStore) UpdateToggle(name string, toggle bool) error {
	if name == "" {
		return ErrInvalidName
	}

	ds.mutex.Lock()
	conn := ds.pool.Get()
	defer conn.Close()
	defer ds.mutex.Unlock()

	ret, err := redis.String(conn.Do("HGET", hashName, name))
	if err != nil {
		return err
	}
	// unmarshal
	var d psmb.MbtcpFilterStatus
	if err := json.Unmarshal([]byte(ret), &d); err != nil {
		return err
	}

	// update toggle
	d.Enabled = toggle

	// marshal
	bytes, err := json.Marshal(d)
	if err != nil {
		return err
	}

	if _, err := conn.Do("HSET", hashName, name, string(bytes)); err != nil {
		return err
	}

	return nil
}

// UpdateAllToggles toggle all request from filter map
func (ds *dataStore) UpdateAllToggles(toggle bool) {
	ds.mutex.Lock()
	conn := ds.pool.Get()
	defer conn.Close()

	ret, err := redis.StringMap(conn.Do("HGETALL", hashName))
	ds.mutex.Unlock()

	if err != nil {
		conf.Log.WithError(err).Warn("Fail to get all items from filter map")
		return
	}
	//conf.Log.WithField("data", ret).Debug("UpdateAllToggles")

	for _, v := range ret {
		var d psmb.MbtcpFilterStatus
		if err := json.Unmarshal([]byte(v), &d); err == nil { // unmarshal
			d.Enabled = toggle
			bytes, err := json.Marshal(d) // marshal
			if err != nil {
				conf.Log.WithError(err).Warn("Fail to marshal items from filter map")
			} else {
				if _, err := conn.Do("HSET", hashName, d.Name, string(bytes)); err != nil {
					conf.Log.WithError(err).Warn("Fail to update toggle to filter map")
				}
			}
		}
	}
}
