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

		Log.WithFields(Fields{"err": ErrFilterNotFound, "file name": "Hello"}).Error("Fail to create log file")

		Log.WithFields(Fields{
			"err":       ErrFilterNotFound,
			"file name": "Hello",
		}).Error("Fail to create log file")

		SetDefault(keyLogEnableDebug, defaultLogEnableDebug)
		Set(keyLogEnableDebug, defaultLogEnableDebug)
		i := GetInt("psmbtcp.max_worker")
		log(i)
		j := GetInt64("psmbtcp.min_connection_timeout")
		log(j)
		s := GetString(keyLogFileName)
		log(s)
		b := GetBool(keyLogToFile)
		log(b)

		d := GetDuration("redis.idel_timeout")
		log(d)

		Set(keyLogToFile, true)
		Set(keyLogToJSONFormat, true)
		Set(keyLogEnableDebug, false)
		base.setLogger()

		return true
	})
}
