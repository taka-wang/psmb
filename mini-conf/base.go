package conf

import (
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/koding/multiconfig"
	log "github.com/takawang/logrus"
)

var base mConf // config instance

// vConf base structu with viper instance
type mConf struct {
	conf *confType
}

// InitConfig int config function
func (b *mConf) initConfig() {
	// get environment variables
	confPath := os.Getenv(envConfPSMBTCP) // config file location
	if confPath == "" {
		confPath = defaultConfigPath
	}
	filePath := path.Join(confPath, keyConfigName) + "." + keyConfigType
	endpoint := os.Getenv(envBackendEndpoint) // backend endpoint, i.e., consul url

	if endpoint == "" {
		log.Debug("Try to load 'local' config file")
	} else {
		log.Debug("Try to load 'remote' config file")
		client, err := api.NewClient(&api.Config{Address: endpoint})
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Warn(": Fail to load 'remote' config file, not found!")
			return
		}
		pair, _, err := client.KV().Get(filePath, nil)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Warn(": Fail to load 'remote' config file, not found!")
			return
		}
		// dump to file
		if err := ioutil.WriteFile(defaultTempPath, pair.Value, 0644); err != nil {
			log.WithFields(log.Fields{"err": err}).Warn(": Fail to load 'remote' config file, not found!")
			return
		}
		filePath = defaultTempPath
	}
	m := multiconfig.NewWithPath(filePath)
	m.MustLoad(b.conf) // Populated the struct
}

// setLogger init logger function
func (b *mConf) setLogger() {
	// set debug level
	if b.conf.Log.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	// set log formatter, JSON or plain text
	if b.conf.Log.JSON {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
	}
	// set log output
	if b.conf.Log.ToFile {
		f, err := os.OpenFile(b.conf.Log.Filename, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Debug(": Fail to create log file")
			f = os.Stdout
		}
		log.SetOutput(f)
	} else {
		log.SetOutput(os.Stdout)
	}
}

func init() {
	base = mConf{conf: new(confType)}
	base.initConfig()
	base.setLogger()
}

//
// Exported API
//

// SetDefault set the default value for this key.
func SetDefault(key string, value interface{}) {
	// do nothing
}

// Set set the value for the key in the override regiser.
func Set(key string, value interface{}) {
	switch key {
	case keyMongoServer:
		base.conf.Mongo.Server = value.(string)
	case keyRedisServer:
		base.conf.Redis.Server = value.(string)
	}
}

// GetInt returns the value associated with the key as an integer
func GetInt(key string) int {
	switch key {
	case keyRedisMaxIdel:
		return base.conf.Redis.MaxIdel
	case keyRedisMaxActive:
		return base.conf.Redis.MaxActive
	case keyPollInterval:
		return base.conf.Psmbtcp.MinPollInterval
	}
	return 0
}

// GetInt64 returns the value associated with the key as an int64
func GetInt64(key string) int64 {
	return base.conf.Psmbtcp.MinConnectionTimeout
}

// GetString returns the value associated with the key as a string
func GetString(key string) string {
	switch key {
	case keyDbName:
		return base.conf.MgoHistory.DbName
	case keyCollectionName:
		return base.conf.MgoHistory.CollectionName
	case keyMongoServer:
		return base.conf.Mongo.Server
	case keyMongoPort:
		return base.conf.Mongo.Port
	case keyMongoDbName:
		return base.conf.Mongo.DbName
	case keyMongoUserName:
		return base.conf.Mongo.Username
	case keyMongoPassword:
		return base.conf.Mongo.Password
	case keyHistoryHashName:
		return base.conf.RedisHistory.HashName
	case keySetPrefix:
		return base.conf.RedisHistory.ZsetPrefix
	case keyRedisServer:
		return base.conf.Redis.Server
	case keyRedisPort:
		return base.conf.Redis.Port
	case keyWriterHashName:
		return base.conf.RedisWriter.HashName
	case keyTCPDefaultPort:
		return base.conf.Psmbtcp.DefaultPort
	case keyZmqPubUpstream:
		return base.conf.Zmq.Pub.Upstream
	case keyZmqPubDownstream:
		return base.conf.Zmq.Pub.Downstream
	case keyZmqSubUpstream:
		return base.conf.Zmq.Sub.Upstream
	case keyZmqSubDownstream:
		return base.conf.Zmq.Sub.Downstream
	}
	return ""
}

// GetBool returns the value associated with the key as a boolean
func GetBool(key string) bool {
	switch key {
	case keyMongoEnableAuth:
		return base.conf.Mongo.Authentication
	case keyMongoIsDrop:
		return base.conf.Mongo.IsDrop
	}
	return false
}

// GetFloat64 returns the value associated with the key as a float64
func GetFloat64(key string) float64 {
	return 0
}

// GetDuration returns the value associated with the key as a duration
func GetDuration(key string) time.Duration {
	switch key {
	case keyMongoConnTimeout:
		return time.Duration(base.conf.Mongo.ConnectionTimeout)
	case keyRedisIdelTimeout:
		return time.Duration(base.conf.Redis.IdelTimeout)
	}
	return 0
}
