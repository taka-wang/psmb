package filter_test

import (
	"testing"

	"github.com/taka-wang/psmb"
	mf "github.com/taka-wang/psmb/mem-filter"
	psmbtcp "github.com/taka-wang/psmb/tcp"
	"github.com/takawang/sugar"
)

func init() {
	psmbtcp.Register("Filter", mf.NewDataStore)
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
			Poll:    "A",
			Name:    "B",
			Enabled: true,
		}
		b := psmb.MbtcpFilterStatus{
			Tid:     123456,
			From:    "test",
			Poll:    "A1",
			Name:    "B1",
			Enabled: true,
		}

		// ADD
		log("Add A item")
		filterMap.Add(a.Poll, a)
		log("Add B item")
		filterMap.Add(b.Poll, b)

		// GET
		log("GET A item")
		if r, b := filterMap.Get(a.Poll); b != false {
			log(r)
		} else {
			return false
		}

		// GETALL
		log("Get all items")
		if r := filterMap.GetAll(a.Poll); r != nil {
			log(r)
		} else {
			return false
		}
		// DELETE
		log("Remove A item")
		filterMap.Delete(a.Poll)

		// GETALL
		log("Get all items")
		if r := filterMap.GetAll(a.Poll); r != nil {
			log(r)
		} else {
			return false
		}
		// TOGGLE

		// TOGGLE ALL

		return true

	})
}
