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
	// RedisPool redis connection pool
	RedisPool  *redis.Pool
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

	RedisPool = &redis.Pool{
		MaxIdle: conf.GetInt(keyRedisMaxIdel),
		// MaxActive: When zero, there is no limit on the number of connections in the pool.
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

// @Implement IHistoryDataStore contract implicitly

// dataStore data store
type dataStore struct {
	redis redis.Conn
}

// NewDataStore instantiate data store
func NewDataStore(conf map[string]string) (interface{}, error) {
	if conn := RedisPool.Get(); conn != nil {
		return &dataStore{
			redis: conn,
		}, nil
	}
	return nil, ErrConnection
}

func (ds *dataStore) connectRedis() error {
	if conn := RedisPool.Get(); conn != nil {
		ds.redis = conn
		return nil
	}
	return ErrConnection
}

func (ds *dataStore) closeRedis() {
	if ds != nil && ds.redis != nil {
		if err := ds.redis.Close(); err != nil {
			conf.Log.WithError(err).Warn("Fail to close redis connection")
		}
		/*else {
			conf.Log.Debug("Close redis connection")
		}
		*/
	}
}

func (ds *dataStore) Add(name string, data interface{}) error {
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		return err
	}

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
	ds.redis.Send("MULTI")
	ds.redis.Send("HSET", hashName, name, string(bytes))      // latest
	ds.redis.Send("ZADD", zsetPrefix+name, ts, string(bytes)) // add to zset
	if _, err := ds.redis.Do("EXEC"); err != nil {
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

	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		return nil, err
	}
	// zset limit is inclusive; zrevrange: from lateste to oldest
	ret, err := redis.StringMap(ds.redis.Do("ZREVRANGE", zsetPrefix+name, 0, limit-1, "WITHSCORES"))
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

	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		return nil, err
	}

	// zrevrange: from lateste to oldest
	ret, err := redis.StringMap(ds.redis.Do("ZREVRANGE", zsetPrefix+name, 0, -1, "WITHSCORES"))
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
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		return "", err
	}

	ret, err := redis.String(ds.redis.Do("HGET", hashName, name))
	if err != nil {
		return "", err
	}
	return ret, nil
}
