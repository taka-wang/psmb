// Package conf a viper-based config
//
// By taka@cmwang.net
//
package conf

import (
	"os"
	"path"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type Fields log.Fields

/*
// Fielder is an interface for providing fields to custom types.
type Fielder interface {
	Fields() Fields
}

// Fields represents a map of entry level data used for structured logging.
type Fields map[string]interface{}

// Fields implements Fielder.
func (f Fields) Fields() Fields {
	return f
}
*/

// Fields represents a map of entry level data used for structured logging.
/*
type Fields map[string]interface{}

// Fields implements Fielder.
func (f Fields) Fields() log.Fields {
	return f.(log.Fields)

	var a log.Fields
	a = make(map[string]interface{})

	return f
}
*/

type mLog struct {
	*log.Logger
}

// Log logger
var Log mLog
var base vConf // config instance

// vConf base structu with viper instance
type vConf struct {
	v *viper.Viper
}

func init() {
	// before load config
	log.SetHandler(text.New(os.Stdout))
	log.SetLevel(log.DebugLevel)
	// init singleton
	base = vConf{v: viper.New()}
	base.initConfig()
	base.setLogger()
}

//
// Exported API
//

// SetDefault set the default value for this key.
func SetDefault(key string, value interface{}) {
	base.v.SetDefault(key, value)
}

// Set set the value for the key in the override regiser.
func Set(key string, value interface{}) {
	base.v.Set(key, value)
}

// GetInt returns the value associated with the key as an integer
func GetInt(key string) int {
	return base.v.GetInt(key)
}

// GetInt64 returns the value associated with the key as an int64
func GetInt64(key string) int64 {
	return base.v.GetInt64(key)
}

// GetString returns the value associated with the key as a string
func GetString(key string) string {
	return base.v.GetString(key)
}

// GetBool returns the value associated with the key as a boolean
func GetBool(key string) bool {
	return base.v.GetBool(key)
}

// GetFloat64 returns the value associated with the key as a float64
func GetFloat64(key string) float64 {
	return base.v.GetFloat64(key)
}

// GetDuration returns the value associated with the key as a duration
func GetDuration(key string) time.Duration {
	return base.v.GetDuration(key)
}

//
// Internal
//

// setLogger init logger function
func (b *vConf) setLogger() {
	//Log = &log.Logger{}

	Log = mLog{Logger: &log.Logger{}}

	writer := os.Stdout
	if b.v.GetBool(keyLogToFile) {
		if f, err := os.OpenFile(b.v.GetString(keyLogFileName), os.O_WRONLY|os.O_CREATE, 0755); err != nil {
			log.WithFields(log.Fields{
				"err":       err,
				"file name": b.v.GetString(keyLogFileName),
			}).Error("Fail to create log file")

		} else {
			writer = f // to file
		}
	}

	// set log formatter, JSON or plain text
	if b.v.GetBool(keyLogToJSONFormat) {
		Log.Handler = json.New(writer)
	} else {
		Log.Handler = text.New(writer)
	}

	// set debug level
	if b.v.GetBool(keyLogEnableDebug) {
		Log.Level = log.DebugLevel
	} else {
		Log.Level = log.InfoLevel
	}
}

// initConfig int config function
func (b *vConf) initConfig() {
	// get environment variables
	confPath := os.Getenv(envConfPSMBTCP) // config file location
	if confPath == "" {
		confPath = defaultConfigPath
	}
	endpoint := os.Getenv(envBackendEndpoint) // backend endpoint, i.e., consul url

	// setup config filename and extension
	b.v.SetConfigName(keyConfigName)
	b.v.SetConfigType(keyConfigType)

	// local or remote config
	if endpoint == "" {
		log.WithField("file path", confPath).Debug("Try to load 'local' config file")
		b.v.AddConfigPath(confPath)
		err := b.v.ReadInConfig() // read config from file
		if err != nil {
			log.WithField("file path", confPath).Warn("Fail to load 'local' config file, not found!")
		} else {
			log.WithField("file path", confPath).Info("Read 'local' config file successfully")
		}
	} else {
		log.WithFields(log.Fields{
			"endpoint":  endpoint,
			"file path": confPath,
			"file name": keyConfigName,
			"file type": keyConfigType,
		}).Debug("Try to load 'remote' config file")

		b.v.AddRemoteProvider(defaultBackendName, endpoint, path.Join(confPath, keyConfigName)+"."+keyConfigType)
		err := b.v.ReadRemoteConfig() // read config from backend
		if err != nil {
			log.WithFields(log.Fields{
				"err":       err,
				"endpoint":  endpoint,
				"file path": confPath,
				"file name": keyConfigName,
				"file type": keyConfigType,
			}).Error("Fail to load 'remote' config file, not found!")
		} else {
			log.WithFields(log.Fields{
				"endpoint":  endpoint,
				"file path": confPath,
				"file name": keyConfigName,
				"file type": keyConfigType,
			}).Info("Read 'remote' config file successfully")
		}
	}

	// set default log values
	b.v.SetDefault(keyLogEnableDebug, defaultLogEnableDebug)
	b.v.SetDefault(keyLogToJSONFormat, defaultLogToJSONFormat)
	b.v.SetDefault(keyLogToFile, defaultLogToFile)
	b.v.SetDefault(keyLogFileName, defaultLogFileName)
}
