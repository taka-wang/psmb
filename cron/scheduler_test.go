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

func TestMbtcpReadTask(t *testing.T) {
	s := sugar.New(t)

	s.Assert("``add` task to map", func(log sugar.Log) bool {
		_, err := psmbtcp.createScheduler("Cron")
		log(err)
		if err != nil {
			return false
		}
		return true
	})
}
