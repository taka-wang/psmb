package conf_test

import (
	"errors"
	"testing"

	conf "github.com/taka-wang/psmb/viper-conf"
	"github.com/takawang/sugar"
)

var (
	ErrFilterNotFound = errors.New("Filter not found")
)

func TestLogger(t *testing.T) {
	s := sugar.New(t)

	s.Assert("Init logger", func(log sugar.Log) bool {
		conf.Log.Debug("hello world")
		conf.Log.WithError(ErrFilterNotFound).Error("World hello")

		conf.Log.WithFields(conf.Fields{"err": ErrFilterNotFound, "file name": "Hello"}).Error("Fail to create log file")

		conf.Log.WithFields(conf.Fields{
			"err":       ErrFilterNotFound,
			"file name": "Hello",
		}).Error("Fail to create log file")

		return true
	})
}
