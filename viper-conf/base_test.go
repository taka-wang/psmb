package conf_test

import (
	"errors"
	"testing"

	conf "github.com/taka-wang/psmb/viper-conf"
	"github.com/takawang/sugar"
)

var (
	hostName string
)

func TestLogger(t *testing.T) {
	s := sugar.New(t)

	s.Assert("Init logger", func(log sugar.Log) bool {
		conf.Log.Debug("hello world")
		var ErrFilterNotFound = errors.New("Filter not found")
		log.WithError(ErrFilterNotFound).Debug("Hello world")
		return true
	})
}
