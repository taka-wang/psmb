package conf

import (
	"os"
	"path"
	"time"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	log "github.com/takawang/logrus"
)

var baseConf Base

// Base base structu with viper instance
type Base struct {
	v *viper.Viper
}

func init() {
	baseConf = Base{v: viper.New()}
}

//
// Exported API
//

// SetLogger generic logger init function
func SetLogger(packageName string) {
	baseConf.setLogger(packageName)
}

// InitConfig generic config function
func InitConfig(packageName string) {
	baseConf.initConfig(packageName)
}

// SetDefault set the default value for this key.
func SetDefault(key string, value interface{}) {
	baseConf.v.SetDefault(key, value)
}

// Set set the value for the key in the override regiser.
func Set(key string, value interface{}) {
	baseConf.v.Set(key, value)
}

// GetInt returns the value associated with the key as an integer
func GetInt(key string) int {
	return baseConf.v.GetInt(key)
}

// GetInt64 returns the value associated with the key as an int64
func GetInt64(key string) int64 {
	return baseConf.v.GetInt64(key)
}

// GetString returns the value associated with the key as a string
func GetString(key string) string {
	return baseConf.v.GetString(key)
}

// GetBool returns the value associated with the key as a boolean
func GetBool(key string) bool {
	return baseConf.v.GetBool(key)
}

// GetFloat64 returns the value associated with the key as a float64
func GetFloat64(key string) float64 {
	return baseConf.v.GetFloat64(key)
}

// GetDuration returns the value associated with the key as a duration
func GetDuration(key string) time.Duration {
	return baseConf.v.GetDuration(key)
}

//
// Internal
//

// setLogger generic logger init function
func (b *Base) setLogger(packageName string) {

	// set debug level
	if b.v.GetBool(keyLogEnableDebug) {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	// set log formatter, JSON or plain text
	if b.v.GetBool(keyLogToJSONFormat) {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
	}
	// set log output
	if b.v.GetBool(keyLogToFile) {
		f, err := os.OpenFile(b.v.GetString(keyLogFileName), os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Debug(packageName + ": Fail to create log file")
			f = os.Stdout
		}
		log.SetOutput(f)
	} else {
		log.SetOutput(os.Stdout)
	}
}

// InitConfig generic config function
func (b *Base) initConfig(packageName string) {
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
		log.Debug(packageName + ": Try to load 'local' config file")
		b.v.AddConfigPath(confPath)
		err := b.v.ReadInConfig() // read config from file
		if err != nil {
			log.Warn(packageName + ": Fail to load 'local' config file, not found!")
		} else {
			log.Info(packageName + ": Read 'local' config file successfully")
		}
	} else {
		log.Debug(packageName + ": Try to load 'remote' config file")
		//log.WithFields(log.Fields{"endpoint": endpoint, "path": confPath}).Debug("remote debug")
		b.v.AddRemoteProvider(defaultBackendName, endpoint, path.Join(confPath, keyConfigName)+"."+keyConfigType)
		err := b.v.ReadRemoteConfig() // read config from backend
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Warn(packageName + ": Fail to load 'remote' config file, not found!")
		} else {
			log.Info(packageName + ": Read remote config file successfully")
		}
	}

	// set default log values
	b.v.SetDefault(keyLogEnableDebug, defaultLogEnableDebug)
	b.v.SetDefault(keyLogToJSONFormat, defaultLogToJSONFormat)
	b.v.SetDefault(keyLogToFile, defaultLogToFile)
	b.v.SetDefault(keyLogFileName, defaultLogFileName)
}
