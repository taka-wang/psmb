// Package history an mongodb-based data store for history.
//
// By taka@cmwang.net
//
package history

import (
	"encoding/json"
	"net"
	"strconv"
	"time"

	"github.com/spf13/viper"
	psmb "github.com/taka-wang/psmb"
	log "github.com/takawang/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	mongoDBDialInfo *mgo.DialInfo
	databaseName    string // mongo database name
	collectionName  string // mongo collection name for history
)

func setDefaults() {
	// set default mongo values
	viper.SetDefault(keyMongoServer, defaultMongoServer)
	viper.SetDefault(keyMongoPort, defaultMongoPort)
	viper.SetDefault(keyMongoIsDrop, defaultMongoIsDrop)
	viper.SetDefault(keyMongoConnTimeout, defaultMongoConnTimeout)
	viper.SetDefault(keyMongoDbName, defaultMongoDbName)
	viper.SetDefault(keyMongoEnableAuth, defaultMongoEnableAuth)
	viper.SetDefault(keyMongoUserName, defaultMongoUserName)
	viper.SetDefault(keyMongoPassword, defaultMongoPassword)

	// set default mongo-history values
	viper.SetDefault(keyDbName, defaultDbName)
	viper.SetDefault(keyCollectionName, defaultCollectionName)

	// Note: for docker environment,
	// lookup mongo server
	host, err := net.LookupHost(defaultMongoDocker)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("local run")
	} else {
		log.WithFields(log.Fields{"hostname": host[0]}).Debug("docker run")
		viper.Set(keyMongoServer, host[0]) // override defaults
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true}) // before init logger
	log.SetLevel(log.DebugLevel)                            // ...

	psmb.InitConfig(packageName) // init based config
	setDefaults()                // set defaults
	psmb.SetLogger(packageName)  // init logger

	databaseName = viper.GetString(keyDbName)
	collectionName = viper.GetString(keyCollectionName)

	if viper.GetBool(keyMongoEnableAuth) {
		// We need this object to establish a session to our MongoDB.
		mongoDBDialInfo = &mgo.DialInfo{
			// allow multiple connection string
			Addrs:    []string{viper.GetString(keyMongoServer) + ":" + viper.GetString(keyMongoPort)},
			Timeout:  viper.GetDuration(keyMongoConnTimeout) * time.Second,
			Database: viper.GetString(keyMongoDbName),
			Username: viper.GetString(keyMongoUserName),
			Password: viper.GetString(keyMongoPassword),
		}
	} else {
		// We need this object to establish a session to our MongoDB.
		mongoDBDialInfo = &mgo.DialInfo{
			// allow multiple connection string
			Addrs:   []string{viper.GetString(keyMongoServer) + ":" + viper.GetString(keyMongoPort)},
			Timeout: viper.GetDuration(keyMongoConnTimeout) * time.Second,
		}
	}
}

// marshal helper function
func marshal(r interface{}) (string, error) {
	bytes, err := json.Marshal(r)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("marshal")
		return "", ErrMarshal
	}
	return string(bytes), nil
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
	if viper.GetBool(keyMongoIsDrop) {
		sessionCopy := pool.Copy()
		err := sessionCopy.DB(databaseName).DropDatabase()
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
	// limit the response
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

	// Check map length
	if len(m) == 0 {
		err := ErrNoData
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

	// Check map length
	if len(m) == 0 {
		err := ErrNoData
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
	ret, err := marshal(result.Data)
	if err != nil {
		return "", err
	}
	return ret, nil
}
