// Package history an redis-based data store for history.
//
// By taka@cmwang.net
//
package history

import (
	"encoding/json"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
	log "github.com/takawang/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	mongoDBDialInfo *mgo.DialInfo
	databaseName    string
	collectionName  string
)

func loadConf(path, backend, endpoint string) {
	// setup viper
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	// set default log values
	viper.SetDefault("log.debug", true)
	viper.SetDefault("log.json", false)
	viper.SetDefault("log.to_file", false)
	viper.SetDefault("log.filename", "/var/log/psmbtcp.log")
	// set default mongo values
	viper.SetDefault("mongo.server", "127.0.0.1")
	viper.SetDefault("mongo.port", "27017")
	viper.SetDefault("mongo.is_drop", true)
	viper.SetDefault("mongo.connection_timeout", 60)
	viper.SetDefault("mongo.db_name", "test")
	viper.SetDefault("mongo.authentication", false)
	viper.SetDefault("mongo.username", "username")
	viper.SetDefault("mongo.password", "password")
	// set default mongo-history values
	viper.SetDefault("mgo-history.db_name", "test")
	viper.SetDefault("mgo-history.collection_name", "mbtcp:history")

	// local or remote
	if backend == "" {
		log.Debug("Try to load local config file")
		if path == "" {
			log.Debug("Config environment variable not found, set to default")
			path = "/etc/psmbtcp"
		}
		// ex: viper.AddConfigPath("/go/src/github.com/taka-wang/psmb")
		viper.AddConfigPath(path)
		err := viper.ReadInConfig()
		if err != nil {
			log.Debug("Local config file not found!")
		}
	} else {
		log.Debug("Try to load remote config file")
		if endpoint == "" {
			log.Debug("Endpoint environment variable not found!")
			return
		}
		// ex: viper.AddRemoteProvider("consul", "192.168.33.10:8500", "/etc/psmbtcp.toml")
		viper.AddRemoteProvider(backend, endpoint, path)
		err := viper.ReadRemoteConfig()
		if err != nil {
			log.Debug("Remote config file not found!")
		}
	}

	// Note: for docker environment
	// lookup mongo server
	host, err := net.LookupHost("mongodb")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("local run")
	} else {
		log.WithFields(log.Fields{"hostname": host[0]}).Debug("docker run")
		viper.Set("mongo.server", host[0]) // override
	}
}

func initLogger() {
	// set debug level
	if viper.GetBool("log.debug") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	// set log formatter
	if viper.GetBool("log.json") {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
	}
	// set log output
	if viper.GetBool("log.to_file") {
		f, err := os.OpenFile(viper.GetString("log.filename"), os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Debug("Fail to write to log file")
			f = os.Stdout
		}
		log.SetOutput(f)
	} else {
		log.SetOutput(os.Stdout)
	}
}

func init() {

	// before init logger
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetLevel(log.DebugLevel)
	// load config
	loadConf(os.Getenv("PSMBTCP_CONFIG"), os.Getenv("SD_BACKEND"), os.Getenv("SD_ENDPOINT"))
	// init logger from config
	initLogger()

	databaseName = viper.GetString("mgo-history.db_name")
	collectionName = viper.GetString("mgo-history.collection_name")

	if viper.GetBool("mongo.authentication") {
		// We need this object to establish a session to our MongoDB.
		mongoDBDialInfo = &mgo.DialInfo{
			Addrs:    []string{viper.GetString("mongo.server") + ":" + viper.GetString("mongo.port")}, // allow multiple connection string
			Timeout:  time.Duration(viper.GetInt("mongo.connection_timeout")) * time.Second,
			Database: viper.GetString("mongo.db_name"),
			Username: viper.GetString("mongo.username"),
			Password: viper.GetString("mongo.password"),
		}
	} else {
		// We need this object to establish a session to our MongoDB.
		mongoDBDialInfo = &mgo.DialInfo{
			Addrs:   []string{viper.GetString("mongo.server") + ":" + viper.GetString("mongo.port")}, // allow multiple connection string
			Timeout: time.Duration(viper.GetInt("mongo.connection_timeout")) * time.Second,
		}
	}

}

// @Implement IHistoryDataStore contract implicitly

// blob data object
type blob struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      string        `bson:"name"`
	Data      interface{}   `bson:"data"`
	Timestamp int64         `bson:"timestamp"`
}

// dataStore data store structure
type dataStore struct {
	mongo *mgo.Session
}

// NewDataStore instantiate data store
func NewDataStore(conf map[string]string) (interface{}, error) {
	// Create a session which maintains a pool of socket connections
	pool, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Fail to instantiate data store")
		return nil, err
	}
	//
	pool.SetMode(mgo.Monotonic, true)

	// Drop Database
	if viper.GetBool("mongo.is_drop") {
		sessionCopy := pool.Copy()
		err = sessionCopy.DB(databaseName).DropDatabase()
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("Fail to drop database")
		}
	}

	// Instantiate
	return &dataStore{
		mongo: pool,
	}, nil
}

// openSession create mongo session
func (ds *dataStore) openSession() (*mgo.Session, error) {
	if ds != nil && ds.mongo != nil {
		sessionCopy := ds.mongo.Copy()
		return sessionCopy, nil
	}
	return nil, ErrConnection
}

// closeSession close mongo session
func (ds *dataStore) closeSession(session *mgo.Session) {
	if session != nil {
		session.Close()
	}
}

func (ds *dataStore) Add(name string, data interface{}) error {
	session, err := ds.openSession()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Add")
		return err
	}
	defer ds.closeSession(session)

	// Collection history
	c := session.DB(databaseName).C(collectionName)
	if err := c.Insert(&blob{Name: name, Data: data, Timestamp: time.Now().UTC().UnixNano()}); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Fail to add to history collection")
		return err
	}

	return nil
}

func (ds *dataStore) Get(name string, limit int) (map[string]string, error) {
	session, err := ds.openSession()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Get")
		return nil, err
	}
	defer ds.closeSession(session)

	// Collection history
	c := session.DB(databaseName).C(collectionName)
	var results []blob
	if err := c.Find(bson.M{"name": name}).Sort("-timestamp").Limit(limit).All(&results); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Get")
		return nil, err
	}

	// Convert to map
	m := make(map[string]string)
	for i := 0; i < len(results); i++ {
		// marshal data to string
		if str, err := marshal(results[i].Data); err == nil {
			m[str] = strconv.FormatInt(results[i].Timestamp, 10)
		}
	}

	// Check length
	if len(m) == 0 {
		err = ErrNoData
		log.WithFields(log.Fields{"err": err}).Error("Get")
		return nil, err
	}
	return m, nil
}

func (ds *dataStore) GetAll(name string) (map[string]string, error) {
	session, err := ds.openSession()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("GetAll")
		return nil, err
	}
	defer ds.closeSession(session)

	// Collection history
	c := session.DB(databaseName).C(collectionName)
	var results []blob
	if err := c.Find(bson.M{"name": name}).Sort("-timestamp").All(&results); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("GetAll")
		return nil, err
	}

	// Convert to map
	m := make(map[string]string)
	for i := 0; i < len(results); i++ {
		// marshal data to string
		if str, err := marshal(results[i].Data); err == nil {
			m[str] = strconv.FormatInt(results[i].Timestamp, 10)
		}
	}

	// Check length
	if len(m) == 0 {
		err = ErrNoData
		log.WithFields(log.Fields{"err": err}).Error("GetAll")
		return nil, err
	}
	return m, nil
}

func (ds *dataStore) GetLatest(name string) (string, error) {
	session, err := ds.openSession()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("GetLatest")
		return "", err
	}
	defer ds.closeSession(session)

	// Collection latest
	c := session.DB(databaseName).C(collectionName)
	result := blob{}

	// Query latest
	if err := c.Find(bson.M{"name": name}).Sort("-timestamp").One(&result); err != nil {
		log.WithFields(log.Fields{"err": "Not Found"}).Error("GetLatest")
		return "", err
	}

	// marshal to string
	ret, err1 := marshal(result.Data)
	if err1 != nil {
		return "", err
	}
	return ret, nil
}

func marshal(r interface{}) (string, error) {
	bytes, err := json.Marshal(r)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("marshal")
		return "", ErrMarshal
	}
	return string(bytes), nil
}
