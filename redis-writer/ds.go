// Package writer an redis-based data store for writer.
//
// By taka@cmwang.net
//
package writer

import (
	"net"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
	psmb "github.com/taka-wang/psmb"
	log "github.com/takawang/logrus"
)

var (
	// RedisPool redis connection pool
	RedisPool *redis.Pool
	hashName  string
)

func setDefaults() {
	// set default redis values
	viper.SetDefault(keyRedisServer, defaultRedisServer)
	viper.SetDefault(keyRedisPort, defaultRedisPort)
	viper.SetDefault(keyRedisMaxIdel, defaultRedisMaxIdel)
	viper.SetDefault(keyRedisMaxActive, defaultRedisMaxActive)
	viper.SetDefault(keyRedisIdelTimeout, defaultRedisIdelTimeout)

	// set default redis-writer values
	viper.SetDefault(keyHashName, defaultHashName)

	// Note: for docker environment
	// lookup redis server
	host, err := net.LookupHost(defaultRedisDocker)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("local run")
	} else {
		log.WithFields(log.Fields{"hostname": host[0]}).Debug("docker run")
		viper.Set("redis.server", host[0]) // override default
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true}) // before init logger
	log.SetLevel(log.DebugLevel)                            // ...

	psmb.InitConfig(packageName) // init based config
	setDefaults()                // set defaults
	psmb.SetLogger(packageName)  // init logger

	hashName = viper.GetString(keyHashName)

	RedisPool = &redis.Pool{
		MaxIdle: viper.GetInt(keyRedisMaxIdel),
		// When zero, there is no limit on the number of connections in the pool.
		MaxActive:   viper.GetInt(keyRedisMaxActive),
		IdleTimeout: viper.GetDuration(keyRedisIdelTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", viper.GetString(keyRedisServer)+":"+viper.GetString(keyRedisPort))
			if err != nil {
				log.WithFields(log.Fields{"err": err}).Error("Redis pool dial error")
			}
			return conn, err
		},
	}
}

// @Implement IWriterTaskDataStore contract implicitly

// writerTaskDataStore write task map type
type writerTaskDataStore struct {
	redis redis.Conn
}

// NewDataStore instantiate mbtcp write task map
func NewDataStore(conf map[string]string) (interface{}, error) {
	// get connection from pool
	conn := RedisPool.Get()
	if nil == conn {
		return nil, ErrConnection
	}

	return &writerTaskDataStore{
		redis: conn,
	}, nil
}

func (ds *writerTaskDataStore) connectRedis() error {
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

func (ds *writerTaskDataStore) closeRedis() {
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

// Add add request to write task map
func (ds *writerTaskDataStore) Add(tid, cmd string) {
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Add")
	}

	if _, err := ds.redis.Do("HSET", hashName, tid, cmd); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Add")
	}
}

// Get get request from write task map
func (ds *writerTaskDataStore) Get(tid string) (string, bool) {
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Get")
	}

	ret, err := redis.String(ds.redis.Do("HGET", hashName, tid))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Get")
		return "", false
	}
	return ret, true
}

// Delete remove request from write task map
func (ds *writerTaskDataStore) Delete(tid string) {
	defer ds.closeRedis()
	if err := ds.connectRedis(); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Delete")
	}
	if _, err := ds.redis.Do("HDEL", hashName, tid); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Delete")
	}
}
