// Package conf a multiconfig-based config
//
// By taka@cmwang.net
//
package conf

import (
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
	"github.com/hashicorp/consul/api"
	"github.com/koding/multiconfig"
)

// Log logger
var Log *log.Logger
var base mConf // config instance

// vConf base structu with viper instance
type mConf struct {
	m *confType
}

func init() {
	// before load config
	log.SetHandler(text.New(os.Stdout))
	log.SetLevel(log.DebugLevel)
	// init singleton
	base = mConf{m: new(confType)}
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
		base.m.Mongo.Server = value.(string)
	case keyRedisServer:
		base.m.Redis.Server = value.(string)
	}
}

// GetInt returns the value associated with the key as an integer
func GetInt(key string) int {
	switch key {
	case keyRedisMaxIdel:
		return base.m.Redis.MaxIdel
	case keyRedisMaxActive:
		return base.m.Redis.MaxActive
	case keyPollInterval:
		return base.m.Psmbtcp.MinPollInterval
	}
	return 0
}

// GetInt64 returns the value associated with the key as an int64
func GetInt64(key string) int64 {
	return base.m.Psmbtcp.MinConnectionTimeout
}

// GetString returns the value associated with the key as a string
func GetString(key string) string {
	switch key {
	case keyDbName:
		return base.m.MgoHistory.DbName
	case keyCollectionName:
		return base.m.MgoHistory.CollectionName
	case keyMongoServer:
		return base.m.Mongo.Server
	case keyMongoPort:
		return base.m.Mongo.Port
	case keyMongoDbName:
		return base.m.Mongo.DbName
	case keyMongoUserName:
		return base.m.Mongo.Username
	case keyMongoPassword:
		return base.m.Mongo.Password
	case keyHistoryHashName:
		return base.m.RedisHistory.HashName
	case keySetPrefix:
		return base.m.RedisHistory.ZsetPrefix
	case keyRedisServer:
		return base.m.Redis.Server
	case keyRedisPort:
		return base.m.Redis.Port
	case keyWriterHashName:
		return base.m.RedisWriter.HashName
	case keyFilterHashName:
		return base.m.RedisFilter.HashName
	case keyTCPDefaultPort:
		return base.m.Psmbtcp.DefaultPort
	case keyZmqPubUpstream:
		return base.m.Zmq.Pub.Upstream
	case keyZmqPubDownstream:
		return base.m.Zmq.Pub.Downstream
	case keyZmqSubUpstream:
		return base.m.Zmq.Sub.Upstream
	case keyZmqSubDownstream:
		return base.m.Zmq.Sub.Downstream
	}
	return ""
}

// GetBool returns the value associated with the key as a boolean
func GetBool(key string) bool {
	switch key {
	case keyMongoEnableAuth:
		return base.m.Mongo.Authentication
	case keyMongoIsDrop:
		return base.m.Mongo.IsDrop
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
		return time.Duration(base.m.Mongo.ConnectionTimeout)
	case keyRedisIdelTimeout:
		return time.Duration(base.m.Redis.IdelTimeout)
	}
	return 0
}

//
// Internal
//

// setLogger init logger function
func (b *mConf) setLogger() {
	Log = &log.Logger{}

	writer := os.Stdout

	if b.m.Log.ToFile {
		if f, err := os.OpenFile(b.m.Log.Filename, os.O_WRONLY|os.O_CREATE, 0755); err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"file": b.m.Log.Filename,
			}).Debug("Fail to create log file")
		} else {
			writer = f // to file
		}
	}

	// set log formatter, JSON or plain text

	if b.m.Log.JSON {
		Log.Handler = json.New(writer)
	} else {
		Log.Handler = text.New(writer)
	}

	// set debug level
	if b.m.Log.Debug {
		Log.Level = log.DebugLevel
	} else {
		Log.Level = log.InfoLevel
	}
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
		log.WithField("file", filePath).Debug("Try to load 'local' config file")
	} else {
		log.WithField("file", filePath).Debug("Try to load 'remote' config file")
		client, err := api.NewClient(&api.Config{Address: endpoint})
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"file": filePath,
			}).Warn("Fail to load 'remote' config file, backend not found!")
			return
		}
		pair, _, err := client.KV().Get(filePath, nil)
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"file": filePath,
			}).Warn("Fail to load 'remote' config file from backend, value not found!")
			return
		}
		// dump to file
		if err := ioutil.WriteFile(defaultTempPath, pair.Value, 0644); err != nil {
			log.WithFields(log.Fields{
				"err":       err,
				"temp file": defaultTempPath,
			}).Warn("Fail to load 'remote' config file from backend, temp file not found!")
			return
		}
		filePath = defaultTempPath
	}
	m := multiconfig.NewWithPath(filePath)
	m.MustLoad(b.m) // Populated the struct
	log.WithField("file", filePath).Info("Read config file successfully")
}
