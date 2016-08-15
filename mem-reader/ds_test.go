package reader

import (
	"strconv"
	"testing"

	"github.com/taka-wang/psmb"
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

		req := psmb.MbtcpPollStatus{Tid: 12345, From: "web"}

		for i := 0; i < 50; i++ {
			s := strconv.Itoa(i)
			if err := reader.Add(s, s, s, req); err != nil {
				log(err, i)
			} else {
				log("ok", i)
			}
		}

		if r := reader.GetAll(); r != nil {
			log(r)
		}

		if _, b := reader.GetTaskByID("10"); b == false {
			log(b)
		}
		if _, b := reader.GetTaskByName("10"); b == false {
			log(b)
		}

		if r := reader.GetAll(); r != nil {
			log(r)
		}

		reader.DeleteTaskByID("10")
		reader.DeleteTaskByName("10")
		reader.UpdateIntervalByName("10", 1)

		if err := reader.UpdateToggleByName("10", true); err != nil {
			log(err)
		}

		reader.UpdateAllToggles(true)
		reader.DeleteAll()
		if r2 := reader.GetAll(); r2 != nil {
			log(r2)
		}

		return true
	})
}
