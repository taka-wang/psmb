// Package history an redis-based data store for history.
//
// By taka@cmwang.net
//
package history

import (
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
	cLatestName       string
	cDataName         string
	connectionTimeout time.Duration
)

func init() {
	// TODO: load config
	isDrop = true
	dbName = "test"
	cLatestName = "mbtcp:last"
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
	Name      string
	Data      interface{}
	Timestamp time.Time
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
		log.WithFields(log.Fields{"err": err}).Debug("Add")
		return err
	}
	defer ds.closeSession(session)

	// Collection history
	c1 := session.DB(dbName).C(cDataName)
	if err := c1.Insert(&blob{Name: name, Data: data, Timestamp: time.Now().UTC()}); err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("Fail to add to history collection")
		return err
	}

	// Collection latest
	c2 := session.DB(dbName).C(cLatestName)
	if err := c2.Insert(&blob{Name: name, Data: data, Timestamp: time.Now().UTC()}); err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("Fail to add to latest collection")
		return err
	}
	return nil
}

func (ds *dataStore) Get(name string, start, stop int) (map[string]string, error) {
	session, err := ds.openSession()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("Add")
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
		log.WithFields(log.Fields{"err": err}).Debug("GetAll")
		return nil, err
	}
	defer ds.closeSession(session)

	// Collection history
	c1 := session.DB(dbName).C(cDataName)
	var results []blob
	if err := c1.Find(bson.M{"name": name}).Sort("-timestamp").All(&results); err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("GetAll")
		return nil, err
	}

	// convert
	m := make(map[string]string)
	for _, v := range results {
		if str, ok := v.Data.(string); ok {
			m[str] = strconv.FormatInt(v.Timestamp.UnixNano(), 10)
		}
	}
	return m, nil
}

func (ds *dataStore) GetLast(name string) (string, error) {
	session, err := ds.openSession()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("Add")
		return "", err
	}
	defer ds.closeSession(session)

	// Collection latest
	c2 := session.DB(dbName).C(cLatestName)
	result := blob{}
	if err := c2.Find(bson.M{"name": name}).One(&result); err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("GetLast")
		return "", err
	}

	if str, ok := result.Data.(string); ok {
		return str, nil
	}
	return "", ErrInvalidName // TODO
}
