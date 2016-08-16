// Package history an redis-based data store for history.
//
// By taka@cmwang.net
//
package history

import (
	"encoding/json"
	"net"
	"time"

	"github.com/garyburd/redigo/redis"
	//conf "github.com/taka-wang/psmb/mini-conf"
	conf "github.com/taka-wang/psmb/viper-conf"
)

var (
	hashName   string
	zsetPrefix string
)

func setDefaults() {
	// set default redis values
	conf.SetDefault(keyRedisServer, defaultRedisServer)
	conf.SetDefault(keyRedisPort, defaultRedisPort)
	conf.SetDefault(keyRedisMaxIdel, defaultRedisMaxIdel)
	conf.SetDefault(keyRedisMaxActive, defaultRedisMaxActive)
	conf.SetDefault(keyRedisIdelTimeout, defaultRedisIdelTimeout)

	// set default redis-history values
	conf.SetDefault(keyHashName, defaultHashName)
	conf.SetDefault(keySetPrefix, defaultSetPrefix)

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
	zsetPrefix = conf.GetString(keySetPrefix)
}

// @Implement IHistoryDataStore contract implicitly

// dataStore data store
type dataStore struct {
	// pool redis connection pool
	pool *redis.Pool
}

// NewDataStore instantiate data store
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

func (ds *dataStore) Add(name string, data interface{}) error {
	conn := ds.pool.Get()
	defer conn.Close()

	// marshal
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	//
	// Time epoch: https://gobyexample.com/epoch
	// nanos := now.UnixNano()
	// fmt.Println(nanos) => 1351700038292387000
	// fmt.Println(time.Unix(0, nanos)) => 2012-10-31 16:13:58.292387 +0000 UTC
	//

	// redis pipeline
	ts := time.Now().UTC().UnixNano()
	conn.Send("MULTI")
	conn.Send("HSET", hashName, name, string(bytes))      // latest
	conn.Send("ZADD", zsetPrefix+name, ts, string(bytes)) // add to zset
	if _, err := conn.Do("EXEC"); err != nil {
		return err
	}
	// debug
	//conf.Log.WithFields(conf.Fields{"Name": name, "Data": data, "TS": ts}).Debug("Add to redis")
	return nil
}

func (ds *dataStore) Get(name string, limit int) (map[string]string, error) {
	if name == "" {
		return nil, ErrInvalidName
	}

	conn := ds.pool.Get()
	defer conn.Close()

	// zset limit is inclusive; zrevrange: from lateste to oldest
	ret, err := redis.StringMap(conn.Do("ZREVRANGE", zsetPrefix+name, 0, limit-1, "WITHSCORES"))
	if err != nil {
		return nil, err
	}

	// Check map length
	if len(ret) == 0 {
		return nil, ErrNoData
	}
	return ret, nil
}

func (ds *dataStore) GetAll(name string) (map[string]string, error) {
	if name == "" {
		return nil, ErrInvalidName
	}

	conn := ds.pool.Get()
	defer conn.Close()

	// zrevrange: from lateste to oldest
	ret, err := redis.StringMap(conn.Do("ZREVRANGE", zsetPrefix+name, 0, -1, "WITHSCORES"))
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, ErrNoData
	}
	// debug
	//conf.Log.WithField("data", ret).Debug("GetAll")
	return ret, nil
}

func (ds *dataStore) GetLatest(name string) (string, error) {
	conn := ds.pool.Get()
	defer conn.Close()

	ret, err := redis.String(conn.Do("HGET", hashName, name))
	if err != nil {
		return "", err
	}
	return ret, nil
}
