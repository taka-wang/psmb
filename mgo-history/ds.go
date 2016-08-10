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

	//conf "github.com/taka-wang/psmb/mini-conf"
	conf "github.com/taka-wang/psmb/viper-conf"
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
	conf.SetDefault(keyMongoServer, defaultMongoServer)
	conf.SetDefault(keyMongoPort, defaultMongoPort)
	conf.SetDefault(keyMongoIsDrop, defaultMongoIsDrop)
	conf.SetDefault(keyMongoConnTimeout, defaultMongoConnTimeout)
	conf.SetDefault(keyMongoDbName, defaultMongoDbName)
	conf.SetDefault(keyMongoEnableAuth, defaultMongoEnableAuth)
	conf.SetDefault(keyMongoUserName, defaultMongoUserName)
	conf.SetDefault(keyMongoPassword, defaultMongoPassword)

	// set default mongo-history values
	conf.SetDefault(keyDbName, defaultDbName)
	conf.SetDefault(keyCollectionName, defaultCollectionName)

	// Note: for docker environment,
	// lookup mongo server
	host, err := net.LookupHost(defaultMongoDocker)
	if err != nil {
		conf.Log.WithError(err).Debug("local run")
	} else {
		conf.Log.WithField("hostname", host[0]).Debug("docker run")
		conf.Set(keyMongoServer, host[0]) // override defaults
	}
}

func init() {
	setDefaults() // set defaults

	databaseName = conf.GetString(keyDbName)
	collectionName = conf.GetString(keyCollectionName)

	if conf.GetBool(keyMongoEnableAuth) {
		// We need this object to establish a session to our MongoDB.
		mongoDBDialInfo = &mgo.DialInfo{
			// allow multiple connection string
			Addrs:    []string{conf.GetString(keyMongoServer) + ":" + conf.GetString(keyMongoPort)},
			Timeout:  conf.GetDuration(keyMongoConnTimeout) * time.Second,
			Database: conf.GetString(keyMongoDbName),
			Username: conf.GetString(keyMongoUserName),
			Password: conf.GetString(keyMongoPassword),
		}
	} else {
		// We need this object to establish a session to our MongoDB.
		mongoDBDialInfo = &mgo.DialInfo{
			// allow multiple connection string
			Addrs:   []string{conf.GetString(keyMongoServer) + ":" + conf.GetString(keyMongoPort)},
			Timeout: conf.GetDuration(keyMongoConnTimeout) * time.Second,
		}
	}
}

// marshal helper function
func marshal(r interface{}) (string, error) {
	bytes, err := json.Marshal(r)
	if err != nil {
		//conf.Log.WithError(err).Error("Fail to marshal")
		return "", ErrMarshal
	}
	return string(bytes), nil
}

// @Implement IHistoryDataStore contract implicitly

type (

	// blob data object
	blob struct {
		ID        bson.ObjectId `bson:"_id,omitempty"`
		Name      string        `bson:"name"`
		Data      interface{}   `bson:"data"`
		Timestamp int64         `bson:"timestamp"`
	}

	// dataStore data store structure
	dataStore struct {
		mongo *mgo.Session
	}
)

// NewDataStore instantiate data store
func NewDataStore(c map[string]string) (interface{}, error) {
	// Create a session which maintains a pool of socket connections
	pool, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		//conf.Log.WithError(err).Error("Fail to instantiate data store")
		return nil, err
	}
	//
	pool.SetMode(mgo.Monotonic, true)

	// Drop Database
	if conf.GetBool(keyMongoIsDrop) {
		sessionCopy := pool.Copy()
		err := sessionCopy.DB(databaseName).DropDatabase()
		if err != nil {
			// we intend to log here
			conf.Log.WithError(err).Debug("Fail to drop database")
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
		//conf.Log.WithError(err).Warn("Add")
		return err
	}
	defer ds.closeSession(session)

	ts := time.Now().UTC().UnixNano()
	// Collection history
	c := session.DB(databaseName).C(collectionName)
	// update or insert by name and data
	if _, err := c.Upsert(bson.M{"name": name, "data": data}, &blob{Name: name, Data: data, Timestamp: ts}); err != nil {
		//conf.Log.WithError(err).Warn("Fail to add to history collection")
		return err
	}
	// debug
	//conf.Log.WithFields(log.Fields{"Name": name, "Data": data, "TS": ts}).Debug("Add to mongo")
	return nil
}

func (ds *dataStore) Get(name string, limit int) (map[string]string, error) {
	session, err := ds.openSession()
	if err != nil {
		//conf.Log.WithError(err).Warn("Get")
		return nil, err
	}
	defer ds.closeSession(session)

	// Collection history
	c := session.DB(databaseName).C(collectionName)
	var results []blob
	// limit the response
	if err := c.Find(bson.M{"name": name}).Sort("-timestamp").Limit(limit).All(&results); err != nil {
		//conf.Log.WithError(err).Debug("Get")
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
		conf.Log.WithError(err).Debug("Get")
		return nil, err
	}
	return m, nil
}

func (ds *dataStore) GetAll(name string) (map[string]string, error) {
	session, err := ds.openSession()
	if err != nil {
		conf.Log.WithError(err).Warn("GetAll")
		return nil, err
	}
	defer ds.closeSession(session)

	// Collection history
	c := session.DB(databaseName).C(collectionName)
	var results []blob
	if err := c.Find(bson.M{"name": name}).Sort("-timestamp").All(&results); err != nil {
		conf.Log.WithError(err).Debug("GetAll")
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
		conf.Log.WithError(err).Debug("GetAll")
		return nil, err
	}
	return m, nil
}

func (ds *dataStore) GetLatest(name string) (string, error) {
	session, err := ds.openSession()
	if err != nil {
		conf.Log.WithError(err).Error("GetLatest")
		return "", err
	}
	defer ds.closeSession(session)

	// Collection latest
	c := session.DB(databaseName).C(collectionName)
	result := blob{}

	// Query latest
	if err := c.Find(bson.M{"name": name}).Sort("-timestamp").One(&result); err != nil {
		conf.Log.WithError(err).Error("GetLatest not found")
		return "", err
	}

	// marshal to string
	ret, err := marshal(result.Data)
	if err != nil {
		return "", err
	}
	return ret, nil
}
