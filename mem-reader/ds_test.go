package reader

import (
	"strconv"
	"testing"

	psmbtcp "github.com/taka-wang/psmb/tcp"
	"github.com/takawang/sugar"
)

func init() {
	psmbtcp.Register("Reader", NewDataStore)
}

func TestMbtcpReadTask(t *testing.T) {
	s := sugar.New(t)

	s.Assert("``add` task to map", func(log sugar.Log) bool {
		reader, err := psmbtcp.ReaderDataStoreCreator("Reader")
		log(err)
		if err != nil {
			return false
		}

		// add null
		e1 := reader.Add("", "10", "10", nil)
		log(e1)

		for i := 0; i < 50; i++ {
			s := strconv.Itoa(i)
			if err := reader.Add(s, s, s, nil); err != nil {
				log(err, i)
			} else {
				log("ok", i)
			}
		}

		r := reader.GetAll()
		log(r)

		_, b1 := reader.GetTaskByID("10")
		log(b1)
		_, b2 := reader.GetTaskByName("10")
		log(b2)
		r := reader.GetAll()
		log(r)

		reader.DeleteTaskByID("10")
		reader.DeleteTaskByName("10")
		reader.UpdateIntervalByName("10")

		err := reader.UpdateToggleByName("10", true)
		log(err)

		reader.UpdateAllToggles(true)
		reader.DeleteAll()
		r2 := reader.GetAll()
		log(r2)

		return true
	})
}
