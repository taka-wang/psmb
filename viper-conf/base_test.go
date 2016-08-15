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

	s.Assert("Init logger", func(log sugar.Log) bool {
		Log.Debug("hello world")
		Log.WithError(ErrFilterNotFound).Error("World hello")

		Log.WithFields(Fields{"err": ErrFilterNotFound, "file name": "Hello"}).Error("Fail to create log file")

		Log.WithFields(Fields{
			"err":       ErrFilterNotFound,
			"file name": "Hello",
		}).Error("Fail to create log file")

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

		return true
	})

	s.Assert("Test setLogger", func(log sugar.Log) bool {
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

	s.Assert("Test Init logger", func(log sugar.Log) bool {
		os.Setenv(envConfPSMBTCP, "")
		os.Setenv(envBackendEndpoint, "")
		base.initConfig()
		os.Setenv(envConfPSMBTCP, "a")
		base.initConfig()
		os.Setenv(envConfPSMBTCP, "a")
		os.Setenv(envBackendEndpoint, "b")
		base.initConfig()
		os.Setenv(envConfPSMBTCP, "a")
		os.Setenv(envBackendEndpoint, "")
		base.initConfig()
		return true
	})

	s.Assert("Test Fail cases", func(log sugar.Log) bool {
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
