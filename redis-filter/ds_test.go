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

	s.Assert("``add` task to map", func(log sugar.Log) bool {
		filterMap, err := psmbtcp.FilterDataStoreCreator("Filter")
		log(err)
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
		log("Add A item")
		filterMap.Add(a.Name, a)
		log("Add B item")
		filterMap.Add(b.Name, b)

		// GET
		log("GET A item")
		if r, b := filterMap.Get(a.Name); b != false {
			log(r)
		} else {
			return false
		}

		// TOGGLE A
		log("Toggle A item")
		if err := filterMap.UpdateToggle(a.Name, false); err != nil {
			return false
		}
		// GET
		log("GET A item")
		if r, b := filterMap.Get(a.Name); b != false {
			log(r)
		} else {
			return false
		}

		// GET ALL
		log("Get all items")
		if r := filterMap.GetAll(); r != nil {
			log(r)
		} else {
			return false
		}

		// Toggle all
		log("Toggle all items")
		filterMap.UpdateAllToggles(false)

		// GET ALL
		log("Get all items")
		if r := filterMap.GetAll(); r != nil {
			log(r)
		} else {
			return false
		}

		// DELETE
		log("Remove A item")
		filterMap.Delete(a.Name)

		// GET ALL
		log("Get all items")
		if r := filterMap.GetAll(); r != nil {
			log(r)
		} else {
			return false
		}

		// out of capacity test
		for i := 0; i < 50; i++ {
			s := strconv.Itoa(i)
			if err := filterMap.Add(s, a); err != nil {
				log(err, i)
			} else {
				log("ok", i)
			}
		}

		// DELETe ALL
		log("Delete all items")
		filterMap.DeleteAll()

		// GET ALL
		log("Get all items")
		if r := filterMap.GetAll(); r == nil {
			log("empty")
			return true
		}

		return false

	})

	s.Assert("Test Redis pool", func(log sugar.Log) bool {
		conf.Set(defaultRedisDocker, "hello")
		setDefaults()

		conf.Set(keyRedisServer, "hello")
		init()
		filterMap, err := psmbtcp.FilterDataStoreCreator("Filter")
		log(err)
		filterMap.Add("123", "123")
		filterMap.Get("123")
		filterMap.Delete("123")
		return true
	})

}
