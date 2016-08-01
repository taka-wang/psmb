package memreader_test

import (
	"testing"

	mr "github.com/taka-wang/psmb/memreader"
	psmbtcp "github.com/taka-wang/psmb/tcp"
	"github.com/takawang/sugar"
)

func init() {
	psmbtcp.Register("Reader", mr.NewDataStore)
}

func TestMbtcpReadTask(t *testing.T) {
	s := sugar.New(t)

	s.Assert("``add` task to map", func(log sugar.Log) bool {
		psmbtcp.CreateWriterTaskDataStore(map[string]string{
			"ReaderDataStore": "Reader",
		})
		return true
	})
}
