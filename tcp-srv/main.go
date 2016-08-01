package main

import (
	"github.com/taka-wang/gocron"
	mr "github.com/taka-wang/psmb/memreader"
	mw "github.com/taka-wang/psmb/memwriter"
	mbtcp "github.com/taka-wang/psmb/tcp"
)

func init() {
	mbtcp.RegisterWriterTask("memory", mw.NewDataStore)
	mbtcp.RegisterReaderTask("memory", mr.NewDataStore)
}

func main() {
	readerDataStore, _ := mbtcp.CreateReaderTaskDataStore(map[string]string{
		"ReaderDataStore": "memory",
	})
	writerDataStore, _ := mbtcp.CreateWriterTaskDataStore(map[string]string{
		"WriterDataStore": "memory",
	})

	srv := mbtcp.NewService(
		readerDataStore,       // readerMap
		writerDataStore,       // writerMap
		gocron.NewScheduler(), // scheduler
	)
	srv.Start()
}
