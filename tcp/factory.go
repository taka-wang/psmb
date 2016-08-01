package tcp

import (
	"fmt"
	"reflect"

	psmb "github.com/taka-wang/psmb"
	log "github.com/takawang/logrus"
)

//
// Factory Pattern
//

// Factories factory container
var Factories = make(map[string]interface{})

// Factory factory method
type Factory func(conf map[string]string) (interface{}, error)

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

/*
// Create create factory
func Create(conf map[string]string, key string) (interface{}, error) {
	defaultDS := "memory"
	if got, ok := conf["WriterDataStore"]; ok {
		defaultDS = got
	}

	engineFactory, ok := Factories[defaultDS]
	if !ok {
		availableDatastores := make([]string, len(Factories))
		for k := range Factories {
			availableDatastores = append(availableDatastores, k)
		}
		return nil, ErrInvalidDataStoreName
	}
	return engineFactory, nil
}
*/

// CreateWriterTaskDataStore create writer task data store
func CreateWriterTaskDataStore(conf map[string]string) (psmb.IWriterTaskDataStore, error) {

	defaultDS := "memory"
	if got, ok := conf["WriterDataStore"]; ok {
		defaultDS = got
	}

	engineFactory, ok := Factories[defaultDS]
	if !ok {
		availableDatastores := make([]string, len(Factories))
		for k := range Factories {
			availableDatastores = append(availableDatastores, k)
		}
		return nil, ErrInvalidDataStoreName
	}

	a := reflect.ValueOf(engineFactory).Call(conf)
	fmt.Println(a)

	fmt.Println("BEGIN")
	fmt.Println("type1: ", reflect.TypeOf(engineFactory))
	// type1:  func(map[string]string) (psmb.IReaderTaskDataStore, error)

	if f, _ := engineFactory.(psmb.IWriterTaskDataStore); f != nil {
		fmt.Println("type: ", reflect.TypeOf(f))
		//return f(conf)
	}
	fmt.Println("END")

	return nil, ErrInvalidDataStoreName
}

// CreateReaderTaskDataStore create reader task data store
func CreateReaderTaskDataStore(conf map[string]string) (psmb.IReaderTaskDataStore, error) {

	defaultDS := "memory"
	if got, ok := conf["ReaderDataStore"]; ok {
		defaultDS = got
	}

	engineFactory, ok := Factories[defaultDS]
	if !ok {
		availableDatastores := make([]string, len(Factories))
		for k := range Factories {
			availableDatastores = append(availableDatastores, k)
		}
		return nil, ErrInvalidDataStoreName
	}

	if f, _ := engineFactory.(psmb.IReaderTaskDataStore); f != nil {
		fmt.Println("type: ", reflect.TypeOf(f))
		//return f(conf)
	}

	return nil, ErrInvalidDataStoreName
}
