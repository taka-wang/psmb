// Package history an redis-based data store for history.
//
// By taka@cmwang.net
//
package history

import (
	"encoding/json"
	"net"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	log "github.com/takawang/logrus"
)

var (
	// RedisPool redis connection pool
	RedisPool  *redis.Pool
	hashName   string
	zsetPrefix string
)

func loadConf(path, endpoint string) {
	// setup viper
	viper.SetConfigName(KeyConfigName)
	viper.SetConfigType(KeyConfigType)

	// local or remote
	if endpoint == "" {
		log.Debug("redis-history: Try to load local config file")
		if path == "" {
			log.Warn("Config environment variable not found, set to default")
			path = DefaultConfigPath
		}
		viper.AddConfigPath(path)
		err := viper.ReadInConfig()
		if err != nil {
			log.Warn("Local config file not found!")
		} else {
			log.Info("Read local config file successfully")
		}
	} else {
		log.Debug("redis-history: Try to load remote config file")
		//log.WithFields(log.Fields{"backend": backend, "endpoint": endpoint, "path": path}).Debug("remote debug")
		viper.AddRemoteProvider(DefaultBackendName, endpoint, path)
		err := viper.ReadRemoteConfig()
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Warn("Remote config file not found!")
		} else {
			log.Info("Read remote config file successfully")
		}
	}

	// set default log values
	viper.SetDefault(keyLogEnableDebug, defaultLogEnableDebug)
	viper.SetDefault(keyLogToJSONFormat, defaultLogToJSONFormat)
	viper.SetDefault(keyLogToFile, defaultLogToFile)
	viper.SetDefault(keyLogFileName, defaultLogFileName)

	// set default redis values
	viper.SetDefault(keyRedisServer, defaultRedisServer)
	viper.SetDefault(keyRedisPort, defaultRedisPort)
	viper.SetDefault(keyRedisMaxIdel, defaultRedisMaxIdel)
	viper.SetDefault(keyRedisMaxActive, defaultRedisMaxActive)
	viper.SetDefault(keyRedisIdelTimeout, defaultRedisIdelTimeout)

	// set default redis-history values
	viper.SetDefault(keyHashName, defaultHashName)
	viper.SetDefault(keySetPrefix, defaultSetPrefix)

	// Note: for docker environment
	// lookup redis server
	host, err := net.LookupHost(defaultRedisDocker)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("local run")
	} else {
		log.WithFields(log.Fields{"hostname": host[0]}).Info("docker run")
		viper.Set(keyRedisServer, host[0]) // override default
	}
}

func initLogger() {
	// set debug level
	if viper.GetBool(keyLogEnableDebug) {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	// set log formatter
	if viper.GetBool(keyLogToJSONFormat) {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
	}
	// set log output
	if viper.GetBool(keyLogToFile) {
		f, err := os.OpenFile(viper.GetString(keyLogFileName), os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Debug("Fail to create log file")
			f = os.Stdout
		}
		log.SetOutput(f)
	} else {
		log.SetOutput(os.Stdout)
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true}) // before init logger
	log.SetLevel(log.DebugLevel)                            // ...

	loadConf(os.Getenv(EnvConfPSMBTCP), os.Getenv(EnvBackendEndpoint)) // load config
	initLogger()                                                       // init logger from config

	hashName = viper.GetString(keyHashName)
	zsetPrefix = viper.GetString(keySetPrefix)

	RedisPool = &redis.Pool{
		MaxIdle:     viper.GetInt(keyRedisMaxIdel),
		MaxActive:   viper.GetInt(keyRedisMaxActive), // When zero, there is no limit on the number of connections in the pool.
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
