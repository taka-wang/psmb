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

	s.Assert("``add` task to map", func(logf sugar.Log) bool {
		reader, err := psmbtcp.ReaderDataStoreCreator("Reader")
		logf(err)
		if err != nil {
			return false
		}

		// add null
		e1 := reader.Add("", "1000", "1000", nil)
		logf(e1)

		req := psmb.MbtcpPollStatus{Tid: 12345, From: "web"}

		for i := 0; i < 50; i++ {
			s := strconv.Itoa(i)
			if err := reader.Add(s, s, s, req); err != nil {
				logf(err, i)
			} else {
				logf("ok", i)
			}
		}

		if r := reader.GetAll(); r != nil {
			logf(r)
		}

		if _, b := reader.GetTaskByID("10"); b == false {
			logf(b)
		}
		if _, b := reader.GetTaskByName("10"); b == false {
			logf(b)
		}

		if r := reader.GetAll(); r != nil {
			logf(r)
		}

		if err := reader.Add("10000", "10000", "10000", nil); err != nil {
			logf(err)
		}

		if err := reader.UpdateIntervalByName("10000", 1); err != nil {
			logf(err)
		}
		if err := reader.UpdateIntervalByName("10", 1); err != nil {
			logf(err)
		}
		if err := reader.UpdateToggleByName("11", true); err != nil {
			logf(err)
		}
		reader.DeleteTaskByID("10")
		reader.DeleteTaskByName("10")

		if err := reader.UpdateToggleByName("10", true); err != nil {
			logf(err)
		}

		reader.UpdateAllToggles(true)
		reader.DeleteAll()
		if r2 := reader.GetAll(); r2 != nil {
			logf(r2)
		}

		return true
	})
}
