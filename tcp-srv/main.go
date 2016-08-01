package main

import (
	"github.com/taka-wang/gocron"
	mr "github.com/taka-wang/psmb/memreader"
	mw "github.com/taka-wang/psmb/memwriter"
	psmbtcp "github.com/taka-wang/psmb/tcp"
)

func init() {
	psmbtcp.Register("Reader", mr.NewDataStore)
	psmbtcp.Register("Writer", mw.NewDataStore)
}

func main() {
	readerDataStore, _ := psmbtcp.CreateReaderDataStore("Reader")
	writerDataStore, _ := psmbtcp.CreateWriterDataStore("Writer")

	srv := psmbtcp.NewService(
		readerDataStore,       // readerMap
		writerDataStore,       // writerMap
		gocron.NewScheduler(), // scheduler
	)
	srv.Start()
}
