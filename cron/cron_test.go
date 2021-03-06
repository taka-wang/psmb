package cron_test

import (
	"testing"

	"github.com/taka-wang/psmb/cron"
	psmbtcp "github.com/taka-wang/psmb/tcp"
	"github.com/takawang/sugar"
)

func init() {
	psmbtcp.Register("Cron", cron.NewScheduler)
}

func TestCron(t *testing.T) {
	s := sugar.New(t)

	s.Assert("``add` task to map", func(logf sugar.Log) bool {
		_, err := psmbtcp.SchedulerCreator("Cron")
		logf(err)
		if err != nil {
			return false
		}

		return true
	})
}
