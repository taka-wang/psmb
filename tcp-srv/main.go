package main

import (
	cron "github.com/taka-wang/psmb/cron"
	mr "github.com/taka-wang/psmb/mrds"
	mw "github.com/taka-wang/psmb/mwds"
	mbtcp "github.com/taka-wang/psmb/tcp"
)

func init() {
	// register data stores from packages
	mbtcp.Register("MemReader", mr.NewDataStore)
	mbtcp.Register("MemWriter", mw.NewDataStore)
	mbtcp.Register("Cron", cron.NewScheduler)
}

func main() {
	// dependency injection & factory pattern
	if srv, _ := mbtcp.NewService(
		"MemReader", // Reader Data Store
		"MemWriter", // Writer Data Store
		"Cron",      // Scheduler
	); srv != nil {
		srv.Start()
	}
}
