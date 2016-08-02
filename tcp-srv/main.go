package main

import (
	cron "github.com/taka-wang/psmb/cron"
	history "github.com/taka-wang/psmb/dbhistory"
	dbwriter "github.com/taka-wang/psmb/dbwds"
	reader "github.com/taka-wang/psmb/mrds"
	writer "github.com/taka-wang/psmb/mwds"
	mbtcp "github.com/taka-wang/psmb/tcp"
)

func init() {
	// register data stores from packages
	mbtcp.Register("MemReader", reader.NewDataStore)
	mbtcp.Register("MemWriter", writer.NewDataStore)
	mbtcp.Register("RedisWriter", dbwriter.NewDataStore)
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
