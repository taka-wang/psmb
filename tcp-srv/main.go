package main

import (
	"github.com/taka-wang/gocron"
	mr "github.com/taka-wang/psmb/memreader"
	mw "github.com/taka-wang/psmb/memwriter"
	mbtcp "github.com/taka-wang/psmb/tcp"
)

func init() {
	mbtcp.Register("Writer", mw.NewDataStore)
	mbtcp.Register("Reader1", mr.NewDataStore)
	mbtcp.Register("Reader2", mr.NewDataStore)
}

func main() {
	readerDataStore, _ := mbtcp.CreateReaderTaskDataStore("Reader2")
	writerDataStore, _ := mbtcp.CreateWriterTaskDataStore("Writer")

	srv := mbtcp.NewService(
		readerDataStore,       // readerMap
		writerDataStore,       // writerMap
		gocron.NewScheduler(), // scheduler
	)
	srv.Start()
}
