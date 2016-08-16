package history

import (
	"testing"
	"time"

	"gopkg.in/mgo.v2"

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

	s.Assert("`add` task to history", func(log sugar.Log) bool {
		historyMap, err := psmbtcp.HistoryDataStoreCreator("History")

		log(err)
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
			log("err:", err)
			return false
		} else {
			log(ret)
		}

		if ret, err := historyMap.GetLatest("hello1"); err != nil {
			log("err:", err)
		} else {
			log(ret)
		}

		if ret, err := historyMap.GetAll("hello"); err != nil {
			log("err:", err)
			return false
		} else {
			log(ret)
			for k, v := range ret {
				log(k, v)
			}
		}

		if ret, err := historyMap.GetAll("hello1"); err != nil {
			log(err)
		} else {
			log(ret)
		}

		if ret, err := historyMap.Get("hello", 2); err != nil {
			log(err)
			return false
		} else {
			log(ret)
		}

		return true
	})

	s.Assert("Test fail cases", func(log sugar.Log) bool {
		conf.Set(keyMongoEnableAuth, true)
		setDefaults()
		mongoDBDialInfo = &mgo.DialInfo{
			// allow multiple connection string
			Addrs:   []string{conf.GetString(keyMongoServer) + ":" + conf.GetString(keyMongoPort)},
			Timeout: conf.GetDuration(keyMongoConnTimeout) * time.Second,
		}
		a, err := NewDataStore(nil)
		log(err)
		b := (*a).(dataStore)
		b.openSession()
		return true
	})

	//marshal
	s.Assert("Test fail to marshal", func(log sugar.Log) bool {
		_, err := marshal("hello;world")
		log(err)
		return true
	})
}
