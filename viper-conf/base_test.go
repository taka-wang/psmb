package conf

import (
	"errors"
	"os"
	"testing"

	"github.com/takawang/sugar"
)

var (
	ErrFilterNotFound = errors.New("Filter not found")
)

func TestLogger(t *testing.T) {
	s := sugar.New(t)

	s.Assert("Init logger", func(logf sugar.Log) bool {
		Log.Debug("hello world")
		Log.WithError(ErrFilterNotFound).Error("World hello")

		Log.WithFields(Fields{"err": ErrFilterNotFound, "file name": "Hello"}).Error("Fail to create log file")

		Log.WithFields(Fields{
			"err":       ErrFilterNotFound,
			"file name": "Hello",
		}).Error("Fail to create log file")

		i := GetInt("psmbtcp.max_worker")
		logf(i)
		j := GetInt64("psmbtcp.min_connection_timeout")
		logf(j)
		s := GetString(keyLogFileName)
		logf(s)
		b := GetBool(keyLogToFile)
		logf(b)
		d := GetDuration("redis.idel_timeout")
		logf(d)

		return true
	})

	s.Assert("Test setLogger", func(_ sugar.Log) bool {
		os.Setenv(envBackendEndpoint, "123")
		base.initConfig()
		SetDefault(keyLogEnableDebug, defaultLogEnableDebug)
		Set(keyLogToJSONFormat, true)
		Set(keyLogEnableDebug, false)
		base.setLogger()
		Set(keyLogFileName, "/tmp/abc")
		Set(keyLogToFile, true)
		base.setLogger()
		return true
	})

	s.Assert("Test Init logger", func(_ sugar.Log) bool {
		os.Setenv("CONF_PSMBTCP", "")
		os.Setenv(envBackendEndpoint, "")
		base.initConfig()
		os.Setenv("CONF_PSMBTCP", "a")
		base.initConfig()
		os.Setenv("CONF_PSMBTCP", "a")
		os.Setenv(envBackendEndpoint, "b")
		base.initConfig()
		os.Setenv("CONF_PSMBTCP", "a")
		os.Setenv(envBackendEndpoint, "")
		base.initConfig()
		return true
	})

	s.Assert("Test Fail cases", func(_ sugar.Log) bool {
		os.Setenv(envBackendEndpoint, "123")
		base.initConfig()
		SetDefault(keyLogEnableDebug, defaultLogEnableDebug)
		Set(keyLogToJSONFormat, true)
		Set(keyLogEnableDebug, false)
		base.setLogger()
		Set(keyLogFileName, "/proc/111")
		Set(keyLogToFile, true)
		base.setLogger()
		return true
	})
}
