package psmb

import (
	"os"
	"path"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	log "github.com/takawang/logrus"
)

// SetLogger generic logger init function
func SetLogger(packageName string) {
	// set debug level
	if viper.GetBool(keyLogEnableDebug) {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	// set log formatter, JSON or plain text
	if viper.GetBool(keyLogToJSONFormat) {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
	}
	// set log output
	if viper.GetBool(keyLogToFile) {
		f, err := os.OpenFile(viper.GetString(keyLogFileName), os.O_WRONLY|os.O_CREATE, 0755)
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
func InitConfig(packageName string) {
	// get environment variables
	confPath := os.Getenv(envConfPSMBTCP) // config file location
	if confPath == "" {
		confPath = defaultConfigPath
	}
	endpoint := os.Getenv(envBackendEndpoint) // backend endpoint, i.e., consul url

	// setup config filename and extension
	viper.SetConfigName(keyConfigName)
	viper.SetConfigType(keyConfigType)

	// local or remote config
	if endpoint == "" {
		log.Debug(packageName + ": Try to load 'local' config file")
		viper.AddConfigPath(confPath)
		err := viper.ReadInConfig() // read config from file
		if err != nil {
			log.Warn(packageName + ": Fail to load 'local' config file, not found!")
		} else {
			log.Info(packageName + ": Read 'local' config file successfully")
		}
	} else {
		log.Debug(packageName + ": Try to load 'remote' config file")
		//log.WithFields(log.Fields{"endpoint": endpoint, "path": confPath}).Debug("remote debug")
		viper.AddRemoteProvider(defaultBackendName, endpoint, path.Join(confPath, keyConfigName)+"."+keyConfigType)
		err := viper.ReadRemoteConfig() // read config from backend
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Warn(packageName + ": Fail to load 'remote' config file, not found!")
		} else {
			log.Info(packageName + ": Read remote config file successfully")
		}
	}

	// set default log values
	viper.SetDefault(keyLogEnableDebug, defaultLogEnableDebug)
	viper.SetDefault(keyLogToJSONFormat, defaultLogToJSONFormat)
	viper.SetDefault(keyLogToFile, defaultLogToFile)
	viper.SetDefault(keyLogFileName, defaultLogFileName)
}
