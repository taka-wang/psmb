package main

import (
	"github.com/taka-wang/gocron"
	_ "github.com/taka-wang/psmb/memwriter"
	mbtcp "github.com/taka-wang/psmb/tcp"
)

func main() {

	writerDataStore, _ := mbtcp.CreateWriterTaskDataStore(map[string]string{
		"DATASTORE": "memory",
	})

	srv := mbtcp.NewService(
		mbtcp.NewReaderMap(),  // readerMap
		writerDataStore,       // writerMap
		gocron.NewScheduler(), // scheduler
	)
	srv.Start()
}
