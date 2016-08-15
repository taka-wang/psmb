package history

import (
	"testing"

	psmbtcp "github.com/taka-wang/psmb/tcp"
	"github.com/taka-wang/psmb/viper-conf"
	"github.com/takawang/sugar"
)

var (
	hostName string
)

func init() {
	psmbtcp.Register("History", NewDataStore)
}

func TestHistoryMap(t *testing.T) {
	s := sugar.New(t)

	s.Assert("`add` task to history", func(logf sugar.Log) bool {
		historyMap, err := psmbtcp.HistoryDataStoreCreator("History")

		logf(err)
		if err != nil {
			return false
		}

		if err := historyMap.Add("hello", "[0,1,2]"); err != nil {
			return false
		}

		data1 := []uint16{1, 2, 3, 4, 5}
		if err := historyMap.Add("hello", data1); err != nil {
			return false
		}
		data2 := []uint16{2, 3, 4, 5, 6}
		if err := historyMap.Add("hello", data2); err != nil {
			return false
		}

		data3 := []uint16{3, 4, 5, 6, 7}
		if err := historyMap.Add("hello", data3); err != nil {
			return false
		}

		data4 := []uint16{4, 5, 6, 7, 8}
		if err := historyMap.Add("hello", data4); err != nil {
			return false
		}

		if ret, err := historyMap.GetLatest("hello"); err != nil {
			logf(err)
			return false
		}
		logf(ret)

		if ret, err := historyMap.GetLatest("hello1"); err != nil {
			logf(err)
		} else {
			logf(ret)
		}

		if ret, err := historyMap.GetAll("hello"); err != nil {
			logf(err)
			return false
		} else {
			logf(ret)
			for k, v := range ret {
				logf(k, v)
			}
		}

		if ret, err := historyMap.GetAll("hello1"); err != nil {
			logf(err)
		} else {
			logf(ret)
		}

		if ret, err := historyMap.Get("hello", 2); err != nil {
			logf(err)
			return false
		}
		logf(ret)

		return true
	})

	s.Assert("Test fail cases", func(logf sugar.Log) bool {
		setDefaults()
		conf.Set(keyMongoEnableAuth, true)
		if _, err := psmbtcp.HistoryDataStoreCreator("History"); err != nil {
			logf(err)
		}

		conf.Set(defaultMongoDocker, "hello")
		setDefaults()
		if a, err := psmbtcp.HistoryDataStoreCreator("History"); err != nil {
			logf(err)
		} else {
			a.Get("1", 1)
			a.Add("2", "2")
			a.GetAll("23")
			a.GetLatest("12")
		}

		_, err := marshal(func() {})
		logf(err)
		return true
	})
}
