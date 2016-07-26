package psmb

import (
	"testing"

	"github.com/marksalpeter/sugar"
)

func TestMbtcpSimpleTask(t *testing.T) {
	s := sugar.New(t)

	s.Assert("`add` task to map", func(log sugar.Log) bool {
		simpleTaskMap := NewMbtcpSimpleTask()
		simpleTaskMap.Add("123456", "12")
		simpleTaskMap.Add("234561", "34")
		r1, b1 := simpleTaskMap.Get("123456")
		if !b1 {
			return false
		}
		if r1 != "12" {
			return false
		}
		_, b2 := simpleTaskMap.Get("1234567")
		if b2 {
			return false
		}

		simpleTaskMap.Delete("123456")
		simpleTaskMap.Delete("1234567")

		_, b3 := simpleTaskMap.Get("123456")
		if b3 {
			return false
		}
		return true
	})
}

func TestMbtcpReadTask(t *testing.T) {
	s := sugar.New(t)

	s.Assert("``add` task to map", func(log sugar.Log) bool {
		// TODO
		return true
	})
}
