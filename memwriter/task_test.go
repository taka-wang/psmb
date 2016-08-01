package memwriter_test

import (
	"testing"

	memwriter "github.com/taka-wang/psmb/memwriter"
	psmbtcp "github.com/taka-wang/psmb/tcp"
	"github.com/takawang/sugar"
)

func init() {
	psmbtcp.Register("Writer", memwriter.NewDataStore)
}

func TestWriterMap(t *testing.T) {
	s := sugar.New(t)

	s.Assert("`add` task to map", func(log sugar.Log) bool {
		writerMap, err := psmbtcp.CreateWriterDataStore("Writer")

		log(err)
		if err != nil {
			return false
		}

		writerMap.Add("123456", "12")
		log("add `123456` to table")
		writerMap.Add("234561", "34")
		log("add `234561` to table")

		r1, b1 := writerMap.Get("123456")
		log("get `123456` from table")
		if !b1 {
			return false
		}
		if r1 != "12" {
			return false
		}

		_, b2 := writerMap.Get("1234567")
		log("get `1234567` from table")

		if b2 {
			return false
		}

		writerMap.Delete("123456")
		log("delete `123456` from table")
		writerMap.Delete("1234567")
		log("delete `1234567` from table")

		_, b3 := writerMap.Get("123456")
		log("get `123456` from table")
		if b3 {
			return false
		}
		return true
	})
}
