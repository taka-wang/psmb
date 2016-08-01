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

// WriterTaskDataStoreFactory writer factory
type WriterTaskDataStoreFactory func(conf map[string]string) (psmb.IWriterTaskDataStore, error)

// ReaderTaskDataStoreFactory reader factory
type ReaderTaskDataStoreFactory func(conf map[string]string) (psmb.IReaderTaskDataStore, error)

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

// CreateWriterTaskDataStore create writer task data store
func CreateWriterTaskDataStore(conf map[string]string) (psmb.IWriterTaskDataStore, error) {
	/*
		for k, v := range conf {
			fmt.Println(k, v)
		}

		for k, v := range Factories {
			fmt.Println(k, v)
		}
	*/
	defaultDS := "Writer"
	if got, ok := conf["WriterDataStore"]; ok {
		fmt.Println("WriterDataStore exist.")
		log.Debug("WriterDataStore exist.")
		defaultDS = got
	} else {
		fmt.Println("WriterDataStore not exist.")
		log.Debug("WriterDataStore not exist.")
	}

	engineFactory, ok := Factories[defaultDS]
	if !ok {
		availableDatastores := make([]string, len(Factories))
		for k := range Factories {
			availableDatastores = append(availableDatastores, k)
		}
		fmt.Println("Invaliad store name1")
		return nil, ErrInvalidDataStoreName
	}

	fmt.Println(reflect.typeOf(engineFactory))
	fmt.Println(reflect.TypeOf(engineFactory).NumMethod())
	if f, _ := engineFactory.(WriterTaskDataStoreFactory); f != nil {
		return f(conf)
	}
	fmt.Println("Invaliad store name2")
	return nil, ErrInvalidDataStoreName
}

// CreateReaderTaskDataStore create reader task data store
func CreateReaderTaskDataStore(conf map[string]string) (psmb.IReaderTaskDataStore, error) {

	defaultDS := "Reader1"
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

	if f, _ := engineFactory.(ReaderTaskDataStoreFactory); f != nil {
		return f(conf)
	}

	return nil, ErrInvalidDataStoreName
}
