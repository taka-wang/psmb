package conf_test

import (
	"testing"

	l "github.com/taka-wang/psmb/viper-conf/Log"
	"github.com/takawang/sugar"
)

var (
	hostName string
)

func TestWriterMap(t *testing.T) {
	s := sugar.New(t)

	s.Assert("`add` task to map", func(log sugar.Log) bool {
		l.Debug("hello")
		return true
	})
}
