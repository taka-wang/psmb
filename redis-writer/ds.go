// Package writer an redis-based data store for writer.
//
// By taka@cmwang.net
//
package writer

import (
	"net"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	//"github.com/taka-wang/psmb/mini-conf"
	"github.com/taka-wang/psmb/viper-conf"
)

var hashName string

func setDefaults() {
	// set default redis values
	conf.SetDefault(keyRedisServer, defaultRedisServer)
	conf.SetDefault(keyRedisPort, defaultRedisPort)
	conf.SetDefault(keyRedisMaxIdel, defaultRedisMaxIdel)
	conf.SetDefault(keyRedisMaxActive, defaultRedisMaxActive)
	conf.SetDefault(keyRedisIdelTimeout, defaultRedisIdelTimeout)

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
}

// @Implement IWriterTaskDataStore contract implicitly

// dataStore write task map type
type dataStore struct {
	mutex sync.Mutex
	pool  *redis.Pool
}

// NewDataStore instantiate mbtcp write task map
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

// Add add request to write task map
func (ds *dataStore) Add(tid, cmd string) {
	ds.mutex.Lock() // lock
	conn := ds.pool.Get()
	defer conn.Close()
	defer ds.mutex.Unlock() // unlock

	if _, err := conn.Do("HSET", hashName, tid, cmd); err != nil {
		conf.Log.WithError(err).Warn("Fail to add item to writer data store")
	}
}

// Get get request from write task map
func (ds *dataStore) Get(tid string) (string, bool) {
	ds.mutex.Lock() // lock
	conn := ds.pool.Get()
	defer conn.Close()

	ret, err := redis.String(conn.Do("HGET", hashName, tid))
	ds.mutex.Unlock() // unlock
	if err != nil {
		conf.Log.WithError(err).Warn("Fail to get item from writer data store")
		return "", false
	}
	return ret, true
}

// Delete remove request from write task map
func (ds *dataStore) Delete(tid string) {
	ds.mutex.Lock() // lock
	conn := ds.pool.Get()
	defer conn.Close()
	defer ds.mutex.Unlock() // unlock

	if _, err := conn.Do("HDEL", hashName, tid); err != nil {
		conf.Log.WithError(err).Error("Fail to delete item from writer data store")
	}
}
