// Package history an redis-based data store for history.
//
// By taka@cmwang.net
//
package history

import (
	"net"
	"time"

	"github.com/garyburd/redigo/redis"
	log "github.com/takawang/logrus"
)

var (
	// RedisPool redis connection pool
	RedisPool  *redis.Pool
	hostName   string
	port       string // "6379"
	hashName   string // "mbtcp:last"
	zsetPrefix string // "mbtcp:data:"
)

func init() {
	// TODO: load config
	port = "6379"
	hashName = "mbtcp:last"
	zsetPrefix = "mbtcp:data:"

	// lookup IP
	host, err := net.LookupHost("redis")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("local run")
		hostName = "127.0.0.1"
	} else {
		log.WithFields(log.Fields{"hostname": host[0]}).Debug("docker run")
		hostName = host[0] //docker
	}

	RedisPool = &redis.Pool{
		MaxIdle:     3,
		MaxActive:   0, // When zero, there is no limit on the number of connections in the pool.
		IdleTimeout: 30 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", hostName+":"+port)
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
	// get connection
	conn := RedisPool.Get()
	if nil == conn {
		return nil, ErrConnection
	}

	return &dataStore{
		redis: conn,
	}, nil
}

func (ds *dataStore) connectRedis() error {
	// get connection
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

	//
	// Time epoch: https://gobyexample.com/epoch
	// nanos := now.UnixNano()
	// fmt.Println(nanos) => 1351700038292387000
	// fmt.Println(time.Unix(0, nanos)) => 2012-10-31 16:13:58.292387 +0000 UTC
	//

	// MULTI
	ds.redis.Send("MULTI")
	ds.redis.Send("HSET", hashName, name, data)                               // latest
	ds.redis.Send("ZADD", zsetPrefix+name, time.Now().UTC().UnixNano(), data) // add to zset
	if _, err := ds.redis.Do("EXEC"); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Add")
		return err
	}
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
	ret, err := redis.StringMap(ds.redis.Do("ZREVRANGE", zsetPrefix+name, 0, limit, "WITHSCORES"))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Get")
		return nil, err
	}
	if len(ret) == 0 {
		err = ErrNoData
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
	ret, err := redis.StringMap(ds.redis.Do("ZREVRANGE", zsetPrefix+name, 0, -1, "WITHSCORES"))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("GetAll")
		return nil, err
	}
	if len(ret) == 0 {
		err = ErrNoData
		log.WithFields(log.Fields{"err": err}).Error("GetAll")
		return nil, err
	}
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
