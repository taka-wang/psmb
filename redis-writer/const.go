package writer

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
	defaultRedisMaxIdel     = 3
	defaultRedisMaxActive   = 0
	defaultRedisIdelTimeout = 30
)

// [redis_writer]
const (
	packageName     = "redis_writer"
	keyHashName     = "redis_writer.hash_name"
	defaultHashName = "mbtcp:writer"
)
