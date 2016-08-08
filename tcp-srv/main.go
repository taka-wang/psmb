package main

import (
	cron "github.com/taka-wang/psmb/cron"
	mfilter "github.com/taka-wang/psmb/mem-filter"
	mreader "github.com/taka-wang/psmb/mem-reader"
	mwriter "github.com/taka-wang/psmb/mem-writer"
	mgohistory "github.com/taka-wang/psmb/mgo-history"
	history "github.com/taka-wang/psmb/redis-history"
	rwriter "github.com/taka-wang/psmb/redis-writer"
	mbtcp "github.com/taka-wang/psmb/tcp"
)

func init() {
	// register plugins explicitly
	mbtcp.Register("MemReader", mreader.NewDataStore)
	mbtcp.Register("MemWriter", mwriter.NewDataStore)
	mbtcp.Register("RedisWriter", rwriter.NewDataStore)
	mbtcp.Register("History", history.NewDataStore)
	mbtcp.Register("MgoHistory", mgohistory.NewDataStore)
	mbtcp.Register("MemFilter1", mfilter.NewDataStore)
	mbtcp.Register("Cron", cron.NewScheduler)
}

func main() {
	// dependency injection & factory pattern
	if srv, _ := mbtcp.NewService(
		"MemReader", // Reader Data Store
		"MemWriter", // Writer Data Store
		"History",   // History Data Store
		"MemFilter", // Filter Data Store
		"Cron",      // Scheduler
	); srv != nil {
		srv.Start()
	}
}
