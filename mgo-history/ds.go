// Package history an redis-based data store for history.
//
// By taka@cmwang.net
//
package history

import (
	"encoding/json"
	"net"
	"strconv"
	"time"

	log "github.com/takawang/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	mongoDBDialInfo   *mgo.DialInfo
	hostName          string
	port              string // "27017"
	isDrop            bool
	dbName            string
	cDataName         string
	connectionTimeout time.Duration
)

func init() {
	// TODO: load config
	isDrop = true
	dbName = "test"
	cDataName = "mbtcp:history"
	connectionTimeout = 60
	port = "27017"

	// lookup IP
	host, err := net.LookupHost("mongodb")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("local run")
		hostName = "127.0.0.1"
	} else {
		log.WithFields(log.Fields{"hostname": host[0]}).Debug("docker run")
		hostName = host[0] //docker
	}

	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo = &mgo.DialInfo{
		Addrs:   []string{hostName + ":" + port}, // allow multiple connection string
		Timeout: connectionTimeout * time.Second,
		//Database: AuthDatabase,
		//Username: AuthUserName,
		//Password: AuthPassword,
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
	if isDrop {
		sessionCopy := pool.Copy()
		err = sessionCopy.DB(dbName).DropDatabase()
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
	c := session.DB(dbName).C(cDataName)
	if err := c.Insert(&blob{Name: name, Data: data, Timestamp: time.Now().UTC().UnixNano()}); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Fail to add to history collection")
		return err
	}

	return nil
}

func (ds *dataStore) Get(name string, start, stop int) (map[string]string, error) {
	session, err := ds.openSession()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Get")
		return nil, err
	}
	defer ds.closeSession(session)

	//
	//
	return nil, nil
}

func (ds *dataStore) GetAll(name string) (map[string]string, error) {
	session, err := ds.openSession()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("GetAll")
		return nil, err
	}
	defer ds.closeSession(session)

	// Collection history
	c := session.DB(dbName).C(cDataName)
	var results []blob
	if err := c.Find(bson.M{"name": name}).Sort("-timestamp").All(&results); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("GetAll")
		return nil, err
	}

	// Convert to map
	m := make(map[string]string)
	for _, v := range results {
		// marshal data to string
		if str, err := marshal(v.Data); err == nil {
			m[str] = strconv.FormatInt(v.Timestamp, 10)
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

func (ds *dataStore) GetLast(name string) (string, error) {
	session, err := ds.openSession()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("GetLast")
		return "", err
	}
	defer ds.closeSession(session)

	// Collection latest
	c := session.DB(dbName).C(cDataName)
	result := blob{}

	// Query latest
	if err := c.Find(bson.M{"name": name}).Sort("-timestamp").One(&result); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("GetLast")
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
