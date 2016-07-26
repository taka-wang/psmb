package psmb

import (
	"testing"

	"github.com/takawang/sugar"
)

func TestMbtcpSimpleTask(t *testing.T) {
	s := sugar.New(t)

	s.Assert("`add` task to map", func(log sugar.Log) bool {
		simpleTaskMap := NewMbtcpSimpleTask()
		simpleTaskMap.Add("123456", "12")
		log("add `123456` to table")
		simpleTaskMap.Add("234561", "34")
		log("add `234561` to table")

		r1, b1 := simpleTaskMap.Get("123456")
		log("get `123456` from table")
		if !b1 {
			return false
		}
		if r1 != "12" {
			return false
		}

		_, b2 := simpleTaskMap.Get("1234567")
		log("get `1234567` from table")

		if b2 {
			return false
		}

		simpleTaskMap.Delete("123456")
		log("delete `123456` from table")
		simpleTaskMap.Delete("1234567")
		log("delete `1234567` from table")

		_, b3 := simpleTaskMap.Get("123456")
		log("get `123456` from table")
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
