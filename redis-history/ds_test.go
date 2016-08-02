package history_test

import (
	"testing"

	r "github.com/taka-wang/psmb/redis-history"
	psmbtcp "github.com/taka-wang/psmb/tcp"
	"github.com/takawang/sugar"
)

var (
	hostName string
)

func init() {
	psmbtcp.Register("History", r.NewDataStore)
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

		if ret, err := historyMap.GetLast("hello"); err != nil {
			log(err)
			return false
		} else {
			log(ret)
		}

		if ret, err := historyMap.GetLast("hello1"); err != nil {
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

		return true
	})
}
