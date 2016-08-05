package history

// [mongo]
const (
	defaultMongoDocker      = "mongodb" // mongo service name for link
	keyMongoServer          = "mongo.server"
	keyMongoPort            = "mongo.port"
	keyMongoIsDrop          = "mongo.is_drop"
	keyMongoConnTimeout     = "mongo.connection_timeout"
	keyMongoDbName          = "mongo.db_name"
	keyMongoEnableAuth      = "mongo.authentication"
	keyMongoUserName        = "mongo.username"
	keyMongoPassword        = "mongo.password"
	defaultMongoServer      = "127.0.0.1"
	defaultMongoPort        = "27017"
	defaultMongoIsDrop      = true
	defaultMongoConnTimeout = 60
	defaultMongoDbName      = "test"
	defaultMongoEnableAuth  = false
	defaultMongoUserName    = "username"
	defaultMongoPassword    = "password"
)

// [mgo_history]
const (
	packageName           = "mgo_history"
	keyDbName             = "mgo-history.db_name"
	keyCollectionName     = "mgo-history.collection_name"
	defaultDbName         = "test"
	defaultCollectionName = "mbtcp:history"
)
