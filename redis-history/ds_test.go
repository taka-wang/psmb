package history_test

import (
	"testing"

	rhistory "github.com/taka-wang/psmb/redis-history"
	psmbtcp "github.com/taka-wang/psmb/tcp"
	"github.com/takawang/sugar"
)

var (
	hostName string
)

func init() {
	psmbtcp.Register("History", rhistory.NewDataStore)
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
			log(err)
			return false
		} else {
			log(ret)
		}

		if ret, err := historyMap.GetLatest("hello1"); err != nil {
			log(err)
		} else {
			log(ret)
		}

		if ret, err := historyMap.GetAll("hello"); err != nil {
			log(err)
			return false
		} else {
			log(ret)
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
}
