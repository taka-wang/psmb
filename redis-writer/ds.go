// Package writer an redis-based data store for writer.
//
// By taka@cmwang.net
//
package writer

import (
	"net"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
	log "github.com/takawang/logrus"
)

var (
	// RedisPool redis connection pool
	RedisPool *redis.Pool
	hashName  string
)

func loadConf(path, backend, endpoint string) {
	// setup viper
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	// set default log values
	viper.SetDefault("log.debug", true)
	viper.SetDefault("log.json", false)
	viper.SetDefault("log.to_file", false)
	viper.SetDefault("log.filename", "/var/log/psmbtcp.log")
	// set default redis values
	viper.SetDefault("redis.server", "127.0.0.1")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.max_idel", 3)
	viper.SetDefault("redis.max_active", 0)
	viper.SetDefault("redis.idel_timeout", 30)
	// set default redis-writer values
	viper.SetDefault("redis_writer.hash_name", "mbtcp:writer")

	// local or remote
	if backend == "" {
		log.Debug("Try to load local config file")
		if path == "" {
			log.Debug("Config environment variable not found, set to default")
			path = "/etc/psmbtcp"
		}
		// ex: viper.AddConfigPath("/go/src/github.com/taka-wang/psmb")
		viper.AddConfigPath(path)
		err := viper.ReadInConfig()
		if err != nil {
			log.Debug("Local config file not found!")
		}
	} else {
		log.Debug("Try to load remote config file")
		if endpoint == "" {
			log.Debug("Endpoint environment variable not found!")
			return
		}
		// ex: viper.AddRemoteProvider("consul", "192.168.33.10:8500", "/etc/psmbtcp.toml")
		viper.AddRemoteProvider(backend, endpoint, path)
		err := viper.ReadRemoteConfig()
		if err != nil {
			log.Debug("Remote config file not found!")
		}
	}

	// Note: for docker environment
	// lookup redis server
	host, err := net.LookupHost("redis")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("local run")
	} else {
		log.WithFields(log.Fields{"hostname": host[0]}).Debug("docker run")
		viper.Set("redis.server", host[0]) // override
	}
}

func initLogger() {
	// set debug level
	if viper.GetBool("log.debug") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	// set log formatter
	if viper.GetBool("log.json") {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
	}
	// set log output
	if viper.GetBool("log.to_file") {
		f, err := os.OpenFile(viper.GetString("log.filename"), os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Debug("Fail to write to log file")
			f = os.Stdout
		}
		log.SetOutput(f)
	} else {
		log.SetOutput(os.Stdout)
	}
}

func init() {
	// before init logger
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetLevel(log.DebugLevel)
	// load config
	loadConf(os.Getenv("PSMBTCP_CONFIG"), os.Getenv("SD_BACKEND"), os.Getenv("SD_ENDPOINT"))
	// init logger from config
	initLogger()

	hashName = viper.GetString("redis_writer.hash_name")

	RedisPool = &redis.Pool{
		MaxIdle:     viper.GetInt("redis.max_idel"),
		MaxActive:   viper.GetInt("redis.max_active"), // When zero, there is no limit on the number of connections in the pool.
		IdleTimeout: time.Duration(viper.GetInt("redis.idel_timeout")) * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", viper.GetString("redis.server")+":"+viper.GetString("redis.port"))
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
	// get connection
	conn := RedisPool.Get()
	if nil == conn {
		return nil, ErrConnection
	}

	return &writerTaskDataStore{
		redis: conn,
	}, nil
}

func (ds *writerTaskDataStore) connectRedis() error {
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
