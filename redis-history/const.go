package history

const (
	defaultConfigPath  = "/etc/psmbtcp" // environment variable backup
	defaultBackendName = "consul"       // remote backend name
	// config
	keyConfigName = "config"
	keyConfigType = "toml"
	// environment variable name
	envConfPSMBTCP     = "CONF_PSMBTCP"
	envBackendEndpoint = "EP_BACKEND"
)

// [log]
const (
	keyLogDebug        = "log.debug"
	keyLogJSON         = "log.json"
	keyLogToFile       = "log.to_file"
	keyLogFileName     = "log.filename"
	defaultLogDebug    = true
	defaultLogJSON     = false
	defaultLogToFile   = false
	defaultLogFileName = "/var/log/psmbtcp.log"
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
	defaultRedisPort        = "6379"
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
