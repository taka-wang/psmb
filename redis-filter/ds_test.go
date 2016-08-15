package filter

import (
	"strconv"
	"testing"

	"github.com/taka-wang/psmb"
	psmbtcp "github.com/taka-wang/psmb/tcp"
	"github.com/taka-wang/psmb/viper-conf"
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

		// GET
		logf("GET A item")
		if r, b := filterMap.Get(a.Name); b != false {
			log(r)
		} else {
			return false
		}

		// TOGGLE A
		logf("Toggle A item")
		if err := filterMap.UpdateToggle(a.Name, false); err != nil {
			return false
		}
		// GET
		logf("GET A item")
		if r, b := filterMap.Get(a.Name); b != false {
			log(r)
		} else {
			return false
		}

		// GET ALL
		logf("Get all items")
		if r := filterMap.GetAll(); r != nil {
			log(r)
		} else {
			return false
		}

		// Toggle all
		logf("Toggle all items")
		filterMap.UpdateAllToggles(false)

		// GET ALL
		logf("Get all items")
		if r := filterMap.GetAll(); r != nil {
			log(r)
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

	s.Assert("Test Fail cases", func(logf sugar.Log) bool {
		filterMap, _ := psmbtcp.FilterDataStoreCreator("Filter")

		a := psmb.MbtcpFilterStatus{
			Tid:     1234,
			From:    "test",
			Name:    "B",
			Enabled: true,
		}

		logf("Add A item")
		filterMap.Add(a.Name, a)

		b := psmb.MbtcpFiltersStatus{
			Tid: 123456,
		}

		logf("Add B item")
		filterMap.Add("h", b)

		logf("Del null item")
		filterMap.Delete("")

		logf("update toggle")
		filterMap.UpdateToggle("", true)
		filterMap.UpdateAllToggles(true)
		filterMap.GetAll()

		return true
	})

	s.Assert("Test Redis pool", func(logf sugar.Log) bool {
		conf.Set(defaultRedisDocker, "hello")
		setDefaults()

		conf.Set(keyRedisServer, "hello")
		setDefaults()
		filterMap, err := psmbtcp.FilterDataStoreCreator("Filter")
		logf(err)
		filterMap.Add("123", "123")
		filterMap.Get("123")
		filterMap.Get("")
		filterMap.Delete("123")
		return true
	})

}
