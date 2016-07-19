package mrds_test

import (
	"testing"

	mr "github.com/taka-wang/psmb/mrds"
	psmbtcp "github.com/taka-wang/psmb/tcp"
	"github.com/takawang/sugar"
)

func init() {
	psmbtcp.Register("Reader", mr.NewDataStore)
}

func TestMbtcpReadTask(t *testing.T) {
	s := sugar.New(t)

	s.Assert("``add` task to map", func(log sugar.Log) bool {
		_, err := psmbtcp.ReaderDataStoreCreator("Reader")
		log(err)
		if err != nil {
			return false
		}
		return true
	})
}
