package history

// [redis]
const (
	defaultRedisDocker      = "redis" // redis service name for link
	keyRedisServer          = "redis.server"
	keyRedisPort            = "redis.port"
	keyRedisMaxIdel         = "redis.max_idel"
	keyRedisMaxActive       = "redis.max_active"
	keyRedisIdelTimeout     = "redis.idel_timeout"
	defaultRedisServer      = "127.0.0.1"
	defaultRedisPort        = "6379"
	defaultRedisMaxIdel     = 5
	defaultRedisMaxActive   = 0
	defaultRedisIdelTimeout = 0
)

// [redis_history]
const (
	keyHashName      = "redis_history.hash_name"
	keySetPrefix     = "redis_history.zset_prefix"
	defaultHashName  = "mbtcp:latest"
	defaultSetPrefix = "mbtcp:data:"
)
