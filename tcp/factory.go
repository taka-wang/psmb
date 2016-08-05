package tcp

import (
	"github.com/taka-wang/psmb"
	"github.com/taka-wang/psmb/cron"
	log "github.com/takawang/logrus"
)

//
// Factory Pattern
//

// factories factory container
var factories = make(map[string]interface{})

// createPlugin real factory method
func createPlugin(conf map[string]string, key string) (interface{}, error) {
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
		return nil, ErrInvalidPluginName
	}
	return engineFactory, nil
}

// createScheduler real factory method
func createScheduler(conf map[string]string) (cron.Scheduler, error) {
	ef, _ := createPlugin(conf, schedulerPluginName)

	if ef != nil {
		//fmt.Println(reflect.TypeOf(ef)) // debug
		if fn, ok := ef.(func(map[string]string) (cron.Scheduler, error)); ok {
			if ds, _ := fn(conf); ds != nil { // casting
				return ds.(cron.Scheduler), nil
			}
		}
		err := ErrCasting
		log.WithFields(log.Fields{"err": err}).Error("Create scheduler")
		return nil, err
	}
	return nil, ErrInvalidPluginName
}

// createWriterDS real factory method
func createHistoryDS(conf map[string]string) (psmb.IHistoryDataStore, error) {
	ef, _ := createPlugin(conf, historyPluginName)

	if ef != nil {
		//fmt.Println(reflect.TypeOf(ef)) // debug
		if fn, ok := ef.(func(map[string]string) (interface{}, error)); ok {
			if ds, _ := fn(conf); ds != nil { // casting
				return ds.(psmb.IHistoryDataStore), nil
			}
		}
		err := ErrCasting
		log.WithFields(log.Fields{"err": err}).Error("Create history data store")
		return nil, err
	}
	return nil, ErrInvalidPluginName
}

// createWriterDS real factory method
func createWriterDS(conf map[string]string) (psmb.IWriterTaskDataStore, error) {
	ef, _ := createPlugin(conf, writerPluginName)

	if ef != nil {
		//fmt.Println(reflect.TypeOf(ef)) // debug
		if fn, ok := ef.(func(map[string]string) (interface{}, error)); ok {
			if ds, _ := fn(conf); ds != nil { // casting
				return ds.(psmb.IWriterTaskDataStore), nil
			}
		}
		err := ErrCasting
		log.WithFields(log.Fields{"err": err}).Error("Create writer data store")
		return nil, err
	}
	return nil, ErrInvalidPluginName
}

// createReaderDS real factory method
func createReaderDS(conf map[string]string) (psmb.IReaderTaskDataStore, error) {
	ef, _ := createPlugin(conf, readerPluginName)

	if ef != nil {
		//fmt.Println(reflect.TypeOf(ef)) // debug
		if fn, ok := ef.(func(map[string]string) (interface{}, error)); ok {
			if ds, _ := fn(conf); ds != nil { // casting
				return ds.(psmb.IReaderTaskDataStore), nil
			}
		}
		err := ErrCasting
		log.WithFields(log.Fields{"err": err}).Error("Create reader data store")
		return nil, err
	}
	return nil, ErrInvalidPluginName
}

// Register register factory methods
func Register(name string, factory interface{}) {
	if factory == nil {
		err := ErrPluginNotExist
		log.WithFields(log.Fields{"Name": name}).Error(err.Error())
	}
	_, registered := factories[name]
	if registered {
		err := ErrPluginExist
		log.WithFields(log.Fields{"Name": name}).Error(err.Error())
	}
	factories[name] = factory
}

// SchedulerCreator concrete creator to create scheduler
func SchedulerCreator(driver string) (cron.Scheduler, error) {
	return createScheduler(map[string]string{schedulerPluginName: driver})
}

// HistoryDataStoreCreator concrete creator to create writer task data store
func HistoryDataStoreCreator(driver string) (psmb.IHistoryDataStore, error) {
	return createHistoryDS(map[string]string{historyPluginName: driver})
}

// WriterDataStoreCreator concrete creator to create writer task data store
func WriterDataStoreCreator(driver string) (psmb.IWriterTaskDataStore, error) {
	return createWriterDS(map[string]string{writerPluginName: driver})
}

// ReaderDataStoreCreator concrete creator to create reader task data store
func ReaderDataStoreCreator(driver string) (psmb.IReaderTaskDataStore, error) {
	return createReaderDS(map[string]string{readerPluginName: driver})
}
