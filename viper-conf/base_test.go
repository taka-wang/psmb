package conf_test

import (
	"testing"

	conf "github.com/taka-wang/psmb/viper-conf"
	"github.com/takawang/sugar"
)

var (
	hostName string
)

func TestWriterMap(t *testing.T) {
	s := sugar.New(t)

	s.Assert("`add` task to map", func(log sugar.Log) bool {
		conf.Log.Debug("hello")
		return true
	})
}
