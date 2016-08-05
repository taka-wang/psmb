package psmb

const (
	defaultConfigPath  = "/etc/psmbtcp" // environment variable backup
	defaultBackendName = "consul"       // remote backend name
	keyConfigName      = "config"       // config file name
	keyConfigType      = "toml"         // config file extension
)

// environment variables
const (
	envConfPSMBTCP     = "CONF_PSMBTCP"
	envBackendEndpoint = "EP_BACKEND"
)

// logs
const (
	keyLogEnableDebug      = "log.debug"
	keyLogToJSONFormat     = "log.json"
	keyLogToFile           = "log.to_file"
	keyLogFileName         = "log.filename"
	defaultLogEnableDebug  = true
	defaultLogToJSONFormat = false
	defaultLogToFile       = false
	defaultLogFileName     = "/var/log/psmbtcp.log"
)
