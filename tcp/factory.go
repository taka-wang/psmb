package tcp

import (
	psmb "github.com/taka-wang/psmb"
	log "github.com/takawang/logrus"
)

//
// Factory Pattern
//

// Factories factory container
var Factories = make(map[string]interface{})

// Register register
func Register(name string, factory interface{}) {
	if factory == nil {
		err := ErrDataStoreNotExist
		log.WithFields(log.Fields{"Name": name}).Error(err.Error())
	}
	_, registered := Factories[name]
	if registered {
		err := ErrDataStoreExist
		log.WithFields(log.Fields{"Name": name}).Error(err.Error())
	}
	Factories[name] = factory
}

// Create create method
func create(conf map[string]string, key string) (interface{}, error) {
	defaultKey := ""
	if got, ok := conf[key]; ok {
		defaultKey = got
	}

	engineFactory, ok := Factories[defaultKey]
	if !ok {
		availableDatastores := make([]string, len(Factories))
		for k := range Factories {
			availableDatastores = append(availableDatastores, k)
		}
		return nil, ErrInvalidDataStoreName
	}
	return engineFactory, nil
}

// CreateWriterTaskDataStore create writer task data store
func CreateWriterTaskDataStore(conf map[string]string) (psmb.IWriterTaskDataStore, error) {
	ef, _ := create(conf, "WriterDataStore")
	if ef != nil {
		if f, _ := ef.(func(map[string]string) (psmb.IWriterTaskDataStore, error)); f != nil {
			got, ok := f(conf)
			return got.(psmb.IWriterTaskDataStore), ok
		}
	}
	return nil, ErrInvalidDataStoreName
}

// CreateReaderTaskDataStore create reader task data store
func CreateReaderTaskDataStore(conf map[string]string) (psmb.IReaderTaskDataStore, error) {
	ef, _ := create(conf, "ReaderDataStore")
	if ef != nil {
		if f, _ := ef.(func(map[string]string) (psmb.IReaderTaskDataStore, error)); f != nil {
			got, ok := f(conf)
			return got.(psmb.IReaderTaskDataStore), ok
		}
	}
	return nil, ErrInvalidDataStoreName
}
