package psmb

import (
	"testing"

	"github.com/takawang/sugar"
)

func TestTypes(t *testing.T) {

	s := sugar.New(t)

	s.Assert("Test MarshalJSON", func(log sugar.Log) bool {
		var a JSONableByteSlice
		if _, err := a.MarshalJSON(); err != nil {
			log(err)
		}
		return true
	})

}
