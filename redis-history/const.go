package history

const (
	DefaultConfigPath  = "/etc/psmbtcp" // environment variable backup
	DefaultBackendName = "consul"       // remote backend name
	// config
	KeyConfigName = "config"
	KeyConfigType = "toml"
	// environment variable name
	EnvConfPSMBTCP     = "CONF_PSMBTCP"
	EnvBackendEndpoint = "EP_BACKEND"
)

// [log]
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

// [redis]
const (
	defaultRedisDocker      = "redis"
	keyRedisServer          = "redis.server"
	keyRedisPort            = "redis.port"
	keyRedisMaxIdel         = "redis.max_idel"
	keyRedisMaxActive       = "redis.max_active"
	keyRedisIdelTimeout     = "redis.idel_timeout"
	defaultRedisServer      = "127.0.0.1"
	defaultRedisPort        = "6378"
	defaultRedisMaxIdel     = 3
	defaultRedisMaxActive   = 0
	defaultRedisIdelTimeout = 30
)

// [redis_history]
const (
	keyHashName      = "redis_history.hash_name"
	keySetPrefix     = "redis_history.zset_prefix"
	defaultHashName  = "mbtcp:latest"
	defaultSetPrefix = "mbtcp:data:"
)
