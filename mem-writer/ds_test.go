package writer

import (
	"testing"

	psmbtcp "github.com/taka-wang/psmb/tcp"
	"github.com/takawang/sugar"
)

func init() {
	psmbtcp.Register("Writer", NewDataStore)
}

func TestWriterMap(t *testing.T) {
	s := sugar.New(t)

	s.Assert("`add` task to map", func(logf sugar.Log) bool {
		writerMap, err := psmbtcp.WriterDataStoreCreator("Writer")

		logf(err)
		if err != nil {
			return false
		}

		writerMap.Add("123456", "12")
		logf("add `123456` to table")
		writerMap.Add("234561", "34")
		logf("add `234561` to table")

		r1, b1 := writerMap.Get("123456")
		logf("get `123456` from table")
		if !b1 {
			return false
		}
		if r1 != "12" {
			return false
		}

		_, b2 := writerMap.Get("1234567")
		logf("get `1234567` from table")

		if b2 {
			return false
		}

		writerMap.Delete("123456")
		logf("delete `123456` from table")
		writerMap.Delete("1234567")
		logf("delete `1234567` from table")

		_, b3 := writerMap.Get("123456")
		logf("get `123456` from table")
		if b3 {
			return false
		}
		return true
	})
}
