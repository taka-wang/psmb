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
	"github.com/spf13/viper"
	//_ "github.com/spf13/viper/remote"
	log "github.com/takawang/logrus"
)

var (
	// RedisPool redis connection pool
	RedisPool  *redis.Pool
	hashName   string // "mbtcp:last"
	zsetPrefix string // "mbtcp:data:"
)

func loadConf(path, remote string) {
	// setup viper
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	// set default values
	viper.SetDefault("redis_history.verbose", false)
	viper.SetDefault("redis_history.port", "6379")
	viper.SetDefault("redis_history.hash_name", "mbtcp:latest")
	viper.SetDefault("redis_history.zset_prefix", "mbtcp:data:")
	viper.SetDefault("redis_history.max_idel", 3)
	viper.SetDefault("redis_history.max_active", 0)
	viper.SetDefault("redis_history.idel_timeout", 30)

	// lookup redis server
	host, err := net.LookupHost("redis")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("local run")
		viper.SetDefault("redis_history.server", "127.0.0.1")
	} else {
		log.WithFields(log.Fields{"hostname": host[0]}).Debug("docker run")
		viper.SetDefault("redis_history.server", host[0])
	}

	// local or remote
	if remote == "" {
		// ex: viper.AddConfigPath("/go/src/github.com/taka-wang/psmb")
		viper.AddConfigPath(path)
		err := viper.ReadInConfig()
		if err != nil {
			log.Debug("Local config file not found!")
		}
	} else {
		// ex: viper.AddRemoteProvider("consul", "192.168.33.10:8500", "/etc/psmbtcp.toml")
		viper.AddRemoteProvider("consul", remote, path)
		err := viper.ReadRemoteConfig()
		if err != nil {
			log.Debug("Remote config file not found!")
		}
	}
}

func init() {

	loadConf("/etc/psmbtcp", "") // load config

	hashName = viper.GetString("redis_history.hash_name")
	zsetPrefix = viper.GetString("redis_history.zset_prefix")

	RedisPool = &redis.Pool{
		MaxIdle:     viper.GetInt("redis_history.max_idel"),
		MaxActive:   viper.GetInt("redis_history.max_active"), // When zero, there is no limit on the number of connections in the pool.
		IdleTimeout: time.Duration(viper.GetInt("redis_history.idel_timeout")) * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", viper.GetString("redis_history.server")+":"+viper.GetString("redis_history.port"))
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

	// MULTI
	ds.redis.Send("MULTI")
	ds.redis.Send("HSET", hashName, name, string(bytes))                               // latest
	ds.redis.Send("ZADD", zsetPrefix+name, time.Now().UTC().UnixNano(), string(bytes)) // add to zset
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
	// zset limit is inclusive
	ret, err := redis.StringMap(ds.redis.Do("ZREVRANGE", zsetPrefix+name, 0, limit-1, "WITHSCORES"))
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
