package tcp

import (
	"github.com/taka-wang/psmb"
	"github.com/taka-wang/psmb/cron"
	//"github.com/taka-wang/psmb/mini-conf"
	"github.com/taka-wang/psmb/viper-conf"
)

//
// Factory Pattern
//

// factories factory container
var factories = make(map[string]interface{})

// createPlugin real factory method
func createPlugin(cnf map[string]string, key string) (interface{}, error) {
	defaultKey := ""
	if got, ok := cnf[key]; ok {
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
func createScheduler(cnf map[string]string) (cron.Scheduler, error) {
	ef, _ := createPlugin(cnf, schedulerPluginName)

	if ef != nil {
		//fmt.Println(reflect.TypeOf(ef)) // debug
		if fn, ok := ef.(func(map[string]string) (cron.Scheduler, error)); ok {
			if ds, _ := fn(cnf); ds != nil { // casting
				return ds.(cron.Scheduler), nil
			}
		}
		err := ErrCasting
		conf.Log.WithError(err).Error("Create scheduler")
		return nil, err
	}
	return nil, ErrInvalidPluginName
}

// createFilterDS real factory method
func createFilterDS(cnf map[string]string) (psmb.IFilterDataStore, error) {
	ef, _ := createPlugin(cnf, filterPluginName)

	if ef != nil {
		//fmt.Println(reflect.TypeOf(ef)) // debug
		if fn, ok := ef.(func(map[string]string) (interface{}, error)); ok {
			if ds, _ := fn(cnf); ds != nil { // casting
				return ds.(psmb.IFilterDataStore), nil
			}
		}
		err := ErrCasting
		conf.Log.WithError(err).Error("Create filter data store")
		return nil, err
	}
	return nil, ErrInvalidPluginName
}

// createWriterDS real factory method
func createHistoryDS(cnf map[string]string) (psmb.IHistoryDataStore, error) {
	ef, _ := createPlugin(cnf, historyPluginName)

	if ef != nil {
		//fmt.Println(reflect.TypeOf(ef)) // debug
		if fn, ok := ef.(func(map[string]string) (interface{}, error)); ok {
			if ds, _ := fn(cnf); ds != nil { // casting
				return ds.(psmb.IHistoryDataStore), nil
			}
		}
		err := ErrCasting
		conf.Log.WithError(err).Error("Create history data store")
		return nil, err
	}
	return nil, ErrInvalidPluginName
}

// createWriterDS real factory method
func createWriterDS(cnf map[string]string) (psmb.IWriterTaskDataStore, error) {
	ef, _ := createPlugin(cnf, writerPluginName)

	if ef != nil {
		//fmt.Println(reflect.TypeOf(ef)) // debug
		if fn, ok := ef.(func(map[string]string) (interface{}, error)); ok {
			if ds, _ := fn(cnf); ds != nil { // casting
				return ds.(psmb.IWriterTaskDataStore), nil
			}
		}
		err := ErrCasting
		conf.Log.WithError(err).Error("Create writer data store")
		return nil, err
	}
	return nil, ErrInvalidPluginName
}

// createReaderDS real factory method
func createReaderDS(cnf map[string]string) (psmb.IReaderTaskDataStore, error) {
	ef, _ := createPlugin(cnf, readerPluginName)

	if ef != nil {
		//fmt.Println(reflect.TypeOf(ef)) // debug
		if fn, ok := ef.(func(map[string]string) (interface{}, error)); ok {
			if ds, _ := fn(cnf); ds != nil { // casting
				return ds.(psmb.IReaderTaskDataStore), nil
			}
		}
		err := ErrCasting
		conf.Log.WithError(err).Error("Create reader data store")
		return nil, err
	}
	return nil, ErrInvalidPluginName
}

// Register register factory methods
func Register(name string, factory interface{}) {
	if factory == nil {
		conf.Log.WithError(ErrPluginNotExist).Error("Register: " + name)
	}
	_, registered := factories[name]
	if registered {
		conf.Log.WithError(ErrPluginExist).Error("Register: " + name)
	}
	factories[name] = factory
}

// FilterDataStoreCreator concrete creator to create filter data store
func FilterDataStoreCreator(driver string) (psmb.IFilterDataStore, error) {
	return createFilterDS(map[string]string{filterPluginName: driver})
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
