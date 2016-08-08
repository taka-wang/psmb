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
	conf "github.com/taka-wang/psmb/viper-conf"
	log "github.com/takawang/logrus"
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
		log.WithFields(log.Fields{"err": err}).Debug("local run")
	} else {
		log.WithFields(log.Fields{"hostname": host[0]}).Info("docker run")
		conf.Set(keyRedisServer, host[0]) // override default
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true}) // before init logger
	log.SetLevel(log.DebugLevel)                            // ...
	setDefaults()                                           // set defaults

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
				log.WithFields(log.Fields{"err": err}).Error("Redis pool dial error")
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
		log.Error(err)
		return err
	}
	//log.Debug("connect to redis")
	ds.redis = conn
	return nil
}

func (ds *dataStore) closeRedis() {
	if ds != nil && ds.redis != nil {
		err := ds.redis.Close()
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("Fail to close redis connection")
		}
		/*else {
			log.Debug("Close redis connection")
		}
		*/
	}
}

func (ds *dataStore) Add(name string, data interface{}) error {
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("Add")
	}

	// marshal
	bytes, err := json.Marshal(data)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("marshal")
		return err
	}

	//
	// Time epoch: https://gobyexample.com/epoch
	// nanos := now.UnixNano()
	// fmt.Println(nanos) => 1351700038292387000
	// fmt.Println(time.Unix(0, nanos)) => 2012-10-31 16:13:58.292387 +0000 UTC
	//

	// redis pipeline
	ds.redis.Send("MULTI")
	ds.redis.Send("HSET", hashName, name, string(bytes))                               // latest
	ds.redis.Send("ZADD", zsetPrefix+name, time.Now().UTC().UnixNano(), string(bytes)) // add to zset
	if ret, err := ds.redis.Do("EXEC"); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Add")
		return err
	}
	// TODO: remove debug
	log.WithFields(log.Fields{"ret": ret}).Info("Add")
	return nil
}

func (ds *dataStore) Get(name string, limit int) (map[string]string, error) {
	if name == "" {
		return nil, ErrInvalidName
	}
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("Get")
	}
	// zset limit is inclusive; zrevrange: from lateste to oldest
	ret, err := redis.StringMap(ds.redis.Do("ZREVRANGE", zsetPrefix+name, 0, limit-1, "WITHSCORES"))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Get")
		return nil, err
	}
	if len(ret) == 0 {
		err := ErrNoData
		log.WithFields(log.Fields{"err": err}).Error("Get")
		return nil, err
	}
	return ret, nil
}

func (ds *dataStore) GetAll(name string) (map[string]string, error) {
	if name == "" {
		return nil, ErrInvalidName
	}
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("GetAll")
	}
	// zrevrange: from lateste to oldest
	ret, err := redis.StringMap(ds.redis.Do("ZREVRANGE", zsetPrefix+name, 0, -1, "WITHSCORES"))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("GetAll")
		return nil, err
	}
	if len(ret) == 0 {
		err := ErrNoData
		log.WithFields(log.Fields{"err": err}).Error("GetAll")
		return nil, err
	}
	// TODO: remove
	log.WithFields(log.Fields{"data": ret}).Info("GetAll")
	return ret, nil
}

func (ds *dataStore) GetLatest(name string) (string, error) {
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("GetLatest")
	}

	ret, err := redis.String(ds.redis.Do("HGET", hashName, name))
	if err != nil {
		log.WithFields(log.Fields{"err": "Not Found"}).Error("GetLatest")
		return "", err
	}
	return ret, nil
}
