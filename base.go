package psmb

import (
	"os"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	log "github.com/takawang/logrus"
)

const (
	defaultConfigPath  = "/etc/psmbtcp" // environment variable backup
	defaultBackendName = "consul"       // remote backend name
	keyConfigName      = "config"
	keyConfigType      = "toml"
	// envs
	envConfPSMBTCP     = "CONF_PSMBTCP"
	envBackendEndpoint = "EP_BACKEND"
	// logs
	keyLogEnableDebug      = "log.debug"
	keyLogToJSONFormat     = "log.json"
	keyLogToFile           = "log.to_file"
	keyLogFileName         = "log.filename"
	defaultLogEnableDebug  = true
	defaultLogToJSONFormat = false
	defaultLogToFile       = false
	defaultLogFileName     = "/var/log/psmbtcp.log"
)

// InitLogger init logger
func InitLogger(packageName string) {
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
			log.WithFields(log.Fields{"err": err}).Debug(packageName + ": Fail to create log file")
			f = os.Stdout
		}
		log.SetOutput(f)
	} else {
		log.SetOutput(os.Stdout)
	}
}

// InitConfig load config
func InitConfig(packageName string) {
	// get environment variables
	path := os.Getenv(envConfPSMBTCP)
	endpoint := os.Getenv(envBackendEndpoint)

	// setup viper
	viper.SetConfigName(keyConfigName)
	viper.SetConfigType(keyConfigType)

	// local or remote
	if endpoint == "" {
		log.Debug(packageName + ": Try to load local config file")
		if path == "" {
			log.Warn(packageName + ": Config environment variable not found, set to default")
			path = defaultConfigPath
		}
		viper.AddConfigPath(path)
		err := viper.ReadInConfig()
		if err != nil {
			log.Warn(packageName + ": Local config file not found!")
		} else {
			log.Info(packageName + ": Read local config file successfully")
		}
	} else {
		log.Debug(packageName + ": Try to load remote config file")
		//log.WithFields(log.Fields{"backend": backend, "endpoint": endpoint, "path": path}).Debug("remote debug")
		viper.AddRemoteProvider(defaultBackendName, endpoint, path)
		err := viper.ReadRemoteConfig()
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Warn("Remote config file not found!")
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
