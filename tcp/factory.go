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
		err := ErrWriterDataStoreNotExist
		log.WithFields(log.Fields{"Name": name}).Error(err.Error())
	}
	_, registered := WriterTaskDataStoreFactories[name]
	if registered {
		err := ErrWriterDataStoreExist
		log.WithFields(log.Fields{"Name": name}).Error(err.Error())
	}
	WriterTaskDataStoreFactories[name] = factory
}

// CreateWriterTaskDataStore create writer task data store
func CreateWriterTaskDataStore(conf map[string]string) (psmb.IWriterTaskDataStore, error) {

	defaultDS := "memory"
	if got, ok := conf["DATASTORE"]; ok {
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
		return nil, ErrInvalidWriterDataStoreName
	}

	// Run the factory with the configuration.
	return engineFactory(conf)
}
