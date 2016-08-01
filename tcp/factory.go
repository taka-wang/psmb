package tcp

import (
	psmb "github.com/taka-wang/psmb"
	log "github.com/takawang/logrus"
)

const (
	readerDS = "ReaderDataStore"
	writerDS = "WriterDataStore"
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

func createDS(conf map[string]string, key string) (interface{}, error) {
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

func createWriterDS(conf map[string]string) (psmb.IWriterTaskDataStore, error) {
	ef, _ := createDS(conf, writerDS)
	if ef != nil {
		if f, _ := ef.(func(map[string]string) (psmb.IWriterTaskDataStore, error)); f != nil {
			got, ok := f(conf)
			return got.(psmb.IWriterTaskDataStore), ok
		}
	}
	return nil, ErrInvalidDataStoreName
}

func createReaderDS(conf map[string]string) (psmb.IReaderTaskDataStore, error) {
	ef, _ := createDS(conf, readerDS)
	if ef != nil {
		if f, _ := ef.(func(map[string]string) (psmb.IReaderTaskDataStore, error)); f != nil {
			got, ok := f(conf)
			return got.(psmb.IReaderTaskDataStore), ok
		}
	}
	return nil, ErrInvalidDataStoreName
}

// CreateWriterDataStore create writer task data store
func CreateWriterDataStore(driver string) (psmb.IWriterTaskDataStore, error) {
	return createWriterDS(map[string]string{writerDS: driver})
}

// CreateReaderDataStore create reader task data store
func CreateReaderDataStore(driver string) (psmb.IReaderTaskDataStore, error) {
	return createReaderDS(map[string]string{readerDS: driver})
}
