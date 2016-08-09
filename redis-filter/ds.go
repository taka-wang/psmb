package filter

import (
	"encoding/json"
	"net"
	"time"

	"github.com/garyburd/redigo/redis"
	//conf "github.com/taka-wang/psmb/mini-conf"
	"github.com/taka-wang/psmb"
	conf "github.com/taka-wang/psmb/viper-conf"
	log "github.com/takawang/logrus"
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
		log.WithFields(log.Fields{"err": err}).Debug("local run")
	} else {
		log.WithFields(log.Fields{"hostname": host[0]}).Debug("docker run")
		conf.Set("redis.server", host[0]) // override default
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true}) // before init logger
	log.SetLevel(log.DebugLevel)                            // ...
	setDefaults()                                           // set defaults

	hashName = conf.GetString(keyHashName)

	RedisPool = &redis.Pool{
		MaxIdle: conf.GetInt(keyRedisMaxIdel),
		// When zero, there is no limit on the number of connections in the pool.
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

//@Implement IFilterDataStore implicitly

// dataStore filter map
type dataStore struct {
	redis redis.Conn
}

// NewDataStore instantiate filter map
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

// Add add request to filter map
func (ds *dataStore) Add(name string, req interface{}) {
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Add")
		return
	}

	// marshal
	bytes, err := json.Marshal(req)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("marshal")
		return
	}

	if _, err := ds.redis.Do("HSET", hashName, name, string(bytes)); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Add")
	}
}

// Get get request from filter map
func (ds *dataStore) Get(name string) (interface{}, bool) {
	if name == "" {
		return nil, false
	}
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("Get")
		return nil, false
	}

	ret, err := redis.String(ds.redis.Do("HGET", hashName, name))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Get")
		return nil, false
	}
	// unmarshal
	var d psmb.MbtcpFilterStatus
	if err := json.Unmarshal([]byte(ret), &d); err != nil {
		log.WithFields(log.Fields{"err": ErrUnmarshal}).Debug("Get")
		return nil, false
	}
	return d, true
}

// GetAll get all requests from filter map
func (ds *dataStore) GetAll(name string) interface{} {
	if name == "" {
		return nil
	}
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("GetAll")
		return nil
	}

	ret, err := redis.StringMap(ds.redis.Do("HGETALL", hashName))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("GetAll")
		return nil
	}
	log.WithFields(log.Fields{"data": ret}).Debug("GetAll")
	// TODO!!
	return nil
}

// Delete remove request from filter map
func (ds *dataStore) Delete(name string) {
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Delete")
		return
	}
	if _, err := ds.redis.Do("HDEL", hashName, name); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Delete")
	}
}

// DeleteAll delete all filters from filter map
func (ds *dataStore) DeleteAll() {
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Delete")
		return
	}
	if _, err := ds.redis.Do("DEL", hashName); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Delete")
	}
}

// Toggle toggle request from filter map
func (ds *dataStore) UpdateToggle(name string, toggle bool) error {
	if name == "" {
		return ErrInvalidName
	}
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("Get")
		return err
	}

	ret, err := redis.String(ds.redis.Do("HGET", hashName, name))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Get")
		return err
	}
	// unmarshal
	var d psmb.MbtcpFilterStatus
	if err := json.Unmarshal([]byte(ret), &d); err != nil {
		log.WithFields(log.Fields{"err": ErrUnmarshal}).Debug("Get")
		return err
	}

	// toggle
	d.Enabled = toggle

	// marshal
	bytes, err := json.Marshal(d)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("marshal")
		return err
	}

	if _, err := ds.redis.Do("HSET", hashName, name, string(bytes)); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Add")
		return err
	}

	return nil
}

// UpdateAllToggles toggle all request from filter map
func (ds *dataStore) UpdateAllToggles(toggle bool) {
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("UpdateAllToggles")
		return
	}

	ret, err := redis.StringMap(ds.redis.Do("HGETALL", hashName))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("UpdateAllToggles")
		return
	}
	log.WithFields(log.Fields{"data": ret}).Debug("UpdateAllToggles")

	// unmarshal
	// TODO!!
	return
}
