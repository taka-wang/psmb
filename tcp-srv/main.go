package main

import (
	"github.com/taka-wang/gocron"
	mwriter "github.com/taka-wang/psmb/memoryWriter"
	mbtcp "github.com/taka-wang/psmb/tcp"
)

func main() {

	writerDataStore, _ := mwriter.CreateWriterTaskDataStore(map[string]string{
		"DATASTORE": "memory",
	})

	srv := mbtcp.NewService(
		mbtcp.NewReaderMap(),  // readerMap
		writerDataStore,       // writerMap
		gocron.NewScheduler(), // scheduler
	)
	srv.Start()
}
