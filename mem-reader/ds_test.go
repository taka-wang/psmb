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
		err := reader.Add("", "10", "10", nil)
		log(err)

		for i := 0; i < 50; i++ {
			s := strconv.Itoa(i)
			if err := reader.Add(s, s, s, nil); err != nil {
				log(err, i)
			} else {
				log("ok", i)
			}
		}

		_, b := reader.GetTaskByID("10")
		log(b)
		_, b := reader.GetTaskByName("10")
		log(b)
		r := reader.GetAll()
		log(r)

		reader.DeleteAll()

		return true
	})
}
