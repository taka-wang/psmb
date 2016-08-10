package reader_test

import (
	"strconv"
	"testing"

	mreader "github.com/taka-wang/psmb/mem-reader"
	psmbtcp "github.com/taka-wang/psmb/tcp"
	"github.com/takawang/sugar"
)

func init() {
	psmbtcp.Register("Reader", mreader.NewDataStore)
}

func TestMbtcpReadTask(t *testing.T) {
	s := sugar.New(t)

	s.Assert("``add` task to map", func(log sugar.Log) bool {
		reader, err := psmbtcp.ReaderDataStoreCreator("Reader")
		log(err)
		if err != nil {
			return false
		}

		for i := 0; i < 50; i++ {
			s := strconv.Itoa(i)
			if err := reader.Add(s, s, s, nil); err != nil {
				log(err, i)
			} else {
				log("ok", i)
			}
		}

		return true
	})
}
