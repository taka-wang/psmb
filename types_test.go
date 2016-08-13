package psmb

import (
	"testing"

	"github.com/takawang/sugar"
)

func TestTypes(t *testing.T) {

	s := sugar.New(t)

	s.Assert("TODO", func(log sugar.Log) bool {
		return true
	})

}
