package main

import (
	cron "github.com/taka-wang/psmb/cron"
	mreader "github.com/taka-wang/psmb/mem-reader"
	mwriter "github.com/taka-wang/psmb/mem-writer"
	history "github.com/taka-wang/psmb/redis-history"
	rwriter "github.com/taka-wang/psmb/redis-writer"
	mbtcp "github.com/taka-wang/psmb/tcp"
)

func init() {
	// register plugins
	mbtcp.Register("MemReader", mreader.NewDataStore)
	mbtcp.Register("MemWriter", mwriter.NewDataStore)
	mbtcp.Register("RedisWriter", rwriter.NewDataStore)
	mbtcp.Register("History", history.NewDataStore)
	mbtcp.Register("Cron", cron.NewScheduler)
}

func main() {
	// dependency injection & factory pattern
	if srv, _ := mbtcp.NewService(
		"MemReader", // Reader Data Store
		"MemWriter", // Writer Data Store
		"History",   // History Data Store
		"Cron",      // Scheduler
	); srv != nil {
		srv.Start()
	}
}
