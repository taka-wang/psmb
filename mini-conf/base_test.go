package conf

import (
	"errors"
	"testing"

	"github.com/takawang/sugar"
)

var (
	ErrFilterNotFound = errors.New("Filter not found")
)

func TestLogger(t *testing.T) {
	s := sugar.New(t)

	s.Assert("Init logger", func(log sugar.Log) bool {
		Log.Debug("hello world")
		Log.WithError(ErrFilterNotFound).Error("World hello")

		Log.WithFields(Fields{
			"err":       ErrFilterNotFound,
			"file name": "Hello",
		}).Error("Fail to create log file")

		return true
	})
}
