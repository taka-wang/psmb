package main

import (
	"github.com/taka-wang/gocron"
	memwriter "github.com/taka-wang/psmb/memwriter"
	mbtcp "github.com/taka-wang/psmb/tcp"
)

func init() {
	mbtcp.RegisterWriterTask("memory", memwriter.NewWriterTaskDataStore)
}

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
