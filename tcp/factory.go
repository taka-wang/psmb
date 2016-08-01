package tcp

import (
	psmb "github.com/taka-wang/psmb"
	log "github.com/takawang/logrus"
)

//
// Factory Pattern for writer task data store
//

// WriterTaskDataStoreFactory factory method for writer task data store
type WriterTaskDataStoreFactory func(conf map[string]string) (psmb.IWriterTaskDataStore, error)

// WriterTaskDataStoreFactories factories container
var WriterTaskDataStoreFactories = make(map[string]WriterTaskDataStoreFactory)

// RegisterWriterTask register writer task data store
func RegisterWriterTask(name string, factory WriterTaskDataStoreFactory) {
	if factory == nil {
		err := ErrDataStoreNotExist
		log.WithFields(log.Fields{"Name": name}).Error(err.Error())
	}
	_, registered := WriterTaskDataStoreFactories[name]
	if registered {
		err := ErrDataStoreExist
		log.WithFields(log.Fields{"Name": name}).Error(err.Error())
	}
	WriterTaskDataStoreFactories[name] = factory
}

// CreateWriterTaskDataStore create writer task data store
func CreateWriterTaskDataStore(conf map[string]string) (psmb.IWriterTaskDataStore, error) {

	defaultDS := "memory"
	if got, ok := conf["WriterDataStore"]; ok {
		defaultDS = got
	}

	engineFactory, ok := WriterTaskDataStoreFactories[defaultDS]
	if !ok {
		// Factory has not been registered.
		// Make a list of all available datastore factories for logging.
		availableDatastores := make([]string, len(WriterTaskDataStoreFactories))
		for k := range WriterTaskDataStoreFactories {
			availableDatastores = append(availableDatastores, k)
		}
		return nil, ErrInvalidDataStoreName
	}

	// Run the factory with the configuration.
	return engineFactory(conf)
}

//
// Factory Pattern for reader task data store
//

// ReaderTaskDataStoreFactory factory method for writer task data store
type ReaderTaskDataStoreFactory func(conf map[string]string) (psmb.IReaderTaskDataStore, error)

// ReaderTaskDataStoreFactories factories container
var ReaderTaskDataStoreFactories = make(map[string]ReaderTaskDataStoreFactory)

// RegisterReaderTask register reader task data store
func RegisterReaderTask(name string, factory ReaderTaskDataStoreFactory) {
	if factory == nil {
		err := ErrDataStoreNotExist
		log.WithFields(log.Fields{"Name": name}).Error(err.Error())
	}
	_, registered := ReaderTaskDataStoreFactories[name]
	if registered {
		err := ErrDataStoreExist
		log.WithFields(log.Fields{"Name": name}).Error(err.Error())
	}
	ReaderTaskDataStoreFactories[name] = factory
}

// CreateReaderTaskDataStore create reader task data store
func CreateReaderTaskDataStore(conf map[string]string) (psmb.IReaderTaskDataStore, error) {

	defaultDS := "memory"
	if got, ok := conf["ReaderDataStore"]; ok {
		defaultDS = got
	}

	engineFactory, ok := ReaderTaskDataStoreFactories[defaultDS]
	if !ok {
		// Factory has not been registered.
		// Make a list of all available datastore factories for logging.
		availableDatastores := make([]string, len(ReaderTaskDataStoreFactories))
		for k := range ReaderTaskDataStoreFactories {
			availableDatastores = append(availableDatastores, k)
		}
		return nil, ErrInvalidDataStoreName
	}

	// Run the factory with the configuration.
	return engineFactory(conf)
}
