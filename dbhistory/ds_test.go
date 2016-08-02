package dbhistory_test

import (
	"testing"

	dbhistory "github.com/taka-wang/psmb/dbhistory"
	psmbtcp "github.com/taka-wang/psmb/tcp"
	"github.com/takawang/sugar"
)

var (
	hostName string
)

func init() {
	psmbtcp.Register("History", dbhistory.NewDataStore)
}

func TestHistoryMap(t *testing.T) {
	s := sugar.New(t)

	s.Assert("`add` task to history", func(log sugar.Log) bool {
		historyMap, err := psmbtcp.HistoryDataStoreCreator("History")

		log(err)
		if err != nil {
			return false
		}

		if err := historyMap.Add("hello", "[1,2,3,4,5]"); err != nil {
			return false
		}
		if err := historyMap.Add("hello", "[2,3,4,5,6]"); err != nil {
			return false
		}

		if ret, err := historyMap.GetLast("hello"); ret != "" {
			log(ret)
		} else {
			log(err)
			return false
		}

		if ret, err := historyMap.GetLast("hello1"); ret != "" {
			log(ret)
			return false
		} else {
			log(err)
		}

		if ret, err := historyMap.GetAll("hello"); ret != nil {
			log(ret)
		} else {
			log(err)
			return false
		}

		if ret, err := historyMap.GetAll("hello1"); ret != nil {
			log(ret)
			return false
		} else {
			log(err)

		}
		return true
	})
}
