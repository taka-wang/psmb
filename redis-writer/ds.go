// Package writer an redis-based data store for writer.
//
// By taka@cmwang.net
//
package writer

import (
	"net"
	"time"

	"github.com/garyburd/redigo/redis"
	//conf "github.com/taka-wang/psmb/mini-conf"
	conf "github.com/taka-wang/psmb/viper-conf"
)

var (
	// RedisPool redis connection pool
	RedisPool *redis.Pool
	hashName  string
)

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
		conf.Set("redis.server", host[0]) // override default
	}
}

func init() {
	setDefaults() // set defaults

	hashName = conf.GetString(keyHashName)

	RedisPool = &redis.Pool{
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
	}
}

// @Implement IWriterTaskDataStore contract implicitly

// dataStore write task map type
type dataStore struct {
	redis redis.Conn
}

// NewDataStore instantiate mbtcp write task map
func NewDataStore(conf map[string]string) (interface{}, error) {
	// get connection from pool
	conn := RedisPool.Get()
	if nil == conn {
		return nil, ErrConnection
	}

	return &dataStore{
		redis: conn,
	}, nil
}

func (ds *dataStore) connectRedis() error {
	// get connection from pool
	conn := RedisPool.Get()
	if nil == conn {
		err := ErrConnection
		return err
	}
	ds.redis = conn
	return nil
}

func (ds *dataStore) closeRedis() {
	if ds != nil && ds.redis != nil {
		err := ds.redis.Close()
		if err != nil {
			conf.Log.WithError(err).Error("Fail to close redis connection")
		}
		/*else {
			conf.Log.Debug("Close redis connection")
		}
		*/
	}
}

// Add add request to write task map
func (ds *dataStore) Add(tid, cmd string) {
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		conf.Log.WithError(err).Error("Fail to connect to redis server")
		return
	}

	if _, err := ds.redis.Do("HSET", hashName, tid, cmd); err != nil {
		conf.Log.WithError(err).Error("Fail to add item to writer data store")
	}
}

// Get get request from write task map
func (ds *dataStore) Get(tid string) (string, bool) {
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		conf.Log.WithError(err).Error("Fail to connect to redis server")
		return "", false
	}

	ret, err := redis.String(ds.redis.Do("HGET", hashName, tid))
	if err != nil {
		conf.Log.WithError(err).Error("Fail to get item from writer data store")
		return "", false
	}
	return ret, true
}

// Delete remove request from write task map
func (ds *dataStore) Delete(tid string) {
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		conf.Log.WithError(err).Error("Fail to connect to redis server")
	}
	if _, err := ds.redis.Do("HDEL", hashName, tid); err != nil {
		conf.Log.WithError(err).Error("Fail to delete item from writer data store")
	}
}
