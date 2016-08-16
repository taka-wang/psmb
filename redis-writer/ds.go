// Package writer an redis-based data store for writer.
//
// By taka@cmwang.net
//
package writer

import (
	"net"
	"time"

	"github.com/garyburd/redigo/redis"
	//cf "github.com/taka-wang/psmb/mini-conf"
	cf "github.com/taka-wang/psmb/viper-conf"
)

var (
	// RedisPool redis connection pool
	RedisPool *redis.Pool
	hashName  string
)

func setDefaults() {
	// set default redis values
	cf.SetDefault(keyRedisServer, defaultRedisServer)
	cf.SetDefault(keyRedisPort, defaultRedisPort)
	cf.SetDefault(keyRedisMaxIdel, defaultRedisMaxIdel)
	cf.SetDefault(keyRedisMaxActive, defaultRedisMaxActive)
	cf.SetDefault(keyRedisIdelTimeout, defaultRedisIdelTimeout)

	// set default redis-writer values
	cf.SetDefault(keyHashName, defaultHashName)

	// Note: for docker environment
	// lookup redis server
	host, err := net.LookupHost(defaultRedisDocker)
	if err != nil {
		cf.Log.WithError(err).Debug("Local run")
	} else {
		cf.Log.WithField("hostname", host[0]).Info("Docker run")
		cf.Set(keyRedisServer, host[0]) // override default
	}
}

func init() {
	setDefaults() // set defaults

	hashName = cf.GetString(keyHashName)
}

// @Implement IWriterTaskDataStore contract implicitly

// dataStore write task map type
type dataStore struct {
	pool *redis.Pool
}

// NewDataStore instantiate mbtcp write task map
func NewDataStore(conf map[string]string) (interface{}, error) {
	return &dataStore{
		pool: &redis.Pool{
			MaxIdle: cf.GetInt(keyRedisMaxIdel),
			// When zero, there is no limit on the number of connections in the pool.
			MaxActive:   cf.GetInt(keyRedisMaxActive),
			IdleTimeout: cf.GetDuration(keyRedisIdelTimeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", cf.GetString(keyRedisServer)+":"+cf.GetString(keyRedisPort))
				if err != nil {
					cf.Log.WithError(err).Error("Redis pool dial error")
				}
				return conn, err
			},
		},
	}, nil
}

// Add add request to write task map
func (ds *dataStore) Add(tid, cmd string) {

	defer ds.pool.Close()
	if _, err := ds.pool.Get().Do("HSET", hashName, tid, cmd); err != nil {
		cf.Log.WithError(err).Warn("Fail to add item to writer data store")
	}
}

// Get get request from write task map
func (ds *dataStore) Get(tid string) (string, bool) {
	defer ds.pool.Close()

	ret, err := redis.String(ds.pool.Get().Do("HGET", hashName, tid))
	if err != nil {
		cf.Log.WithError(err).Warn("Fail to get item from writer data store")
		return "", false
	}
	return ret, true
}

// Delete remove request from write task map
func (ds *dataStore) Delete(tid string) {
	defer ds.pool.Close()

	if _, err := ds.pool.Get().Do("HDEL", hashName, tid); err != nil {
		cf.Log.WithError(err).Error("Fail to delete item from writer data store")
	}
}
