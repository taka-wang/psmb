package psmb

import (
	"testing"

	"github.com/takawang/sugar"
)

func TestTypes(t *testing.T) {

	s := sugar.New(t)

	s.Assert("Test MarshalJSON", func(logf sugar.Log) bool {
		var a JSONableByteSlice
		if _, err := a.MarshalJSON(); err != nil {
			logf(err)
		}
		return true
	})

}
