package conf

// environment variable names
const (
	envConfPSMBTCP     = "CONF_PSMBTCP" // config path
	envBackendEndpoint = "EP_BACKEND"   // backend endpoint
)

// config
const (
	defaultConfigPath = "/etc/psmbtcp" // environment variable backup
	defaultTempPath   = "/tmp/conf"
	keyConfigName     = "config" // config file name
	keyConfigType     = "toml"   // config file extension
)

// mongo
const (
	keyMongoServer      = "mongo.server"
	keyMongoPort        = "mongo.port"
	keyMongoIsDrop      = "mongo.is_drop"
	keyMongoConnTimeout = "mongo.connection_timeout"
	keyMongoDbName      = "mongo.db_name"
	keyMongoEnableAuth  = "mongo.authentication"
	keyMongoUserName    = "mongo.username"
	keyMongoPassword    = "mongo.password"
)

// mgo_history
const (
	keyDbName         = "mgo-history.db_name"
	keyCollectionName = "mgo-history.collection_name"
)

// redis
const (
	keyRedisServer      = "redis.server"
	keyRedisPort        = "redis.port"
	keyRedisMaxIdel     = "redis.max_idel"
	keyRedisMaxActive   = "redis.max_active"
	keyRedisIdelTimeout = "redis.idel_timeout"
)

// redis-history
const (
	keyHistoryHashName = "redis_history.hash_name"
	keySetPrefix       = "redis_history.zset_prefix"
)

// redis-writer
const (
	keyWriterHashName = "redis_writer.hash_name"
)

// redis-filter
const (
	keyFilterHashName = "redis_filter.hash_name"
)

// tcp
const (
	keyTCPDefaultPort      = "psmbtcp.default_port"
	keyMinConnectionTimout = "psmbtcp.min_connection_timeout"
	keyPollInterval        = "psmbtcp.min_poll_interval"
	keyZmqPubUpstream      = "zmq.pub.upstream"
	keyZmqPubDownstream    = "zmq.pub.downstream"
	keyZmqSubUpstream      = "zmq.sub.upstream"
	keyZmqSubDownstream    = "zmq.sub.downstream"
)
