package filter

import (
	"strconv"
	"testing"

	"github.com/taka-wang/psmb"
	psmbtcp "github.com/taka-wang/psmb/tcp"
	"github.com/takawang/sugar"
)

func init() {
	psmbtcp.Register("Filter", NewDataStore)
}

func TestFilter(t *testing.T) {
	s := sugar.New(t)

	s.Assert("``add` task to map", func(logf sugar.Log) bool {
		filterMap, err := psmbtcp.FilterDataStoreCreator("Filter")
		logf(err)
		if err != nil {
			return false
		}

		a := psmb.MbtcpFilterStatus{
			Tid:     1234,
			From:    "test",
			Name:    "B",
			Enabled: true,
		}
		b := psmb.MbtcpFilterStatus{
			Tid:     123456,
			From:    "test",
			Name:    "B1",
			Enabled: true,
		}

		// ADD
		logf("Add A item")
		filterMap.Add(a.Name, a)
		logf("Add B item")
		filterMap.Add(b.Name, b)
		logf("Add null item")
		filterMap.Add("", b)

		// GET
		logf("GET A item")
		if r, b := filterMap.Get(a.Name); b != false {
			logf(r)
		} else {
			return false
		}

		// TOGGLE A
		logf("Toggle A item")
		if err := filterMap.UpdateToggle(a.Name, false); err != nil {
			return false
		}
		logf("Toggle NULL item")
		if err := filterMap.UpdateToggle("D", false); err != nil {
			logf(err)
		}

		// GET
		logf("GET A item")
		if r, b := filterMap.Get(a.Name); b != false {
			logf(r)
		} else {
			return false
		}

		// GET ALL
		logf("Get all items")
		if r := filterMap.GetAll(); r != nil {
			logf(r)
		} else {
			return false
		}

		// Toggle all
		logf("Toggle all items")
		filterMap.UpdateAllToggles(false)

		// GET ALL
		logf("Get all items")
		if r := filterMap.GetAll(); r != nil {
			logf(r)
		} else {
			return false
		}

		// DELETE
		logf("Remove A item")
		filterMap.Delete(a.Name)

		// GET ALL
		logf("Get all items")
		if r := filterMap.GetAll(); r != nil {
			logf(r)
		} else {
			return false
		}

		// out of capacity test
		for i := 0; i < 50; i++ {
			s := strconv.Itoa(i)
			if err := filterMap.Add(s, a); err != nil {
				logf(err, i)
			} else {
				logf("ok", i)
			}
		}

		// DELETe ALL
		logf("Delete all items")
		filterMap.DeleteAll()

		// GET ALL
		logf("Get all items")
		if r := filterMap.GetAll(); r == nil {
			logf("empty")
			return true
		}

		return false

	})
}
