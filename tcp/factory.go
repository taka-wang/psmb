package tcp

import (
	"fmt"
	"reflect"

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

// factories factory container
var factories = make(map[string]interface{})

// createDS real factory method
func createDS(conf map[string]string, key string) (interface{}, error) {
	defaultKey := ""
	if got, ok := conf[key]; ok {
		defaultKey = got
	}

	engineFactory, ok := factories[defaultKey]
	if !ok {
		availableDatastores := make([]string, len(factories))
		for k := range factories {
			availableDatastores = append(availableDatastores, k)
		}
		return nil, ErrInvalidDataStoreName
	}
	return engineFactory, nil
}

// createWriterDS real factory method
func createWriterDS(conf map[string]string) (psmb.IWriterTaskDataStore, error) {
	ef, _ := createDS(conf, writerDS)

	if ef != nil {
		fmt.Println(reflect.TypeOf(ef))
		if f, ok := ef.(func(map[string]string) (interface{}, error)); ok {
			got, ok2 := f(conf)
			// todo: check casting
			return got.(psmb.IWriterTaskDataStore), ok2
		}
		err := ErrCasting
		log.WithFields(log.Fields{"err": err}).Error("Create writer data store")
		return nil, err
	}
	return nil, ErrInvalidDataStoreName
}

// createReaderDS real factory method
func createReaderDS(conf map[string]string) (psmb.IReaderTaskDataStore, error) {
	ef, _ := createDS(conf, readerDS)

	if ef != nil {
		fmt.Println(reflect.TypeOf(ef))
		if f, ok := ef.(func(map[string]string) (interface{}, error)); ok {
			got, ok2 := f(conf)
			// todo: check casting
			return got.(psmb.IReaderTaskDataStore), ok2
		}
		err := ErrCasting
		log.WithFields(log.Fields{"err": err}).Error("Create reader data store")
		return nil, err
	}
	return nil, ErrInvalidDataStoreName
}

// Register register factory methods
func Register(name string, factory interface{}) {
	if factory == nil {
		err := ErrDataStoreNotExist
		log.WithFields(log.Fields{"Name": name}).Error(err.Error())
	}
	_, registered := factories[name]
	if registered {
		err := ErrDataStoreExist
		log.WithFields(log.Fields{"Name": name}).Error(err.Error())
	}
	factories[name] = factory
}

// WriterDataStoreCreator factory method to create writer task data store
func WriterDataStoreCreator(driver string) (psmb.IWriterTaskDataStore, error) {
	return createWriterDS(map[string]string{writerDS: driver})
}

// ReaderDataStoreCreator factory method to create reader task data store
func ReaderDataStoreCreator(driver string) (psmb.IReaderTaskDataStore, error) {
	return createReaderDS(map[string]string{readerDS: driver})
}
