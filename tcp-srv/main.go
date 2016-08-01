package main

import (
	"github.com/taka-wang/gocron"
	"github.com/taka-wang/psmb/mrds"
	"github.com/taka-wang/psmb/mwds"
	psmbtcp "github.com/taka-wang/psmb/tcp"
)

func init() {
	psmbtcp.Register("Reader", mrds.NewDataStore)
	psmbtcp.Register("Writer", mwds.NewDataStore)
}

func main() {
	// Factory
	readerDataStore, _ := psmbtcp.ReaderDataStoreCreator("Reader")
	writerDataStore, _ := psmbtcp.WriterDataStoreCreator("Writer")

	// DI
	srv := psmbtcp.NewService(
		readerDataStore,       // readerMap
		writerDataStore,       // writerMap
		gocron.NewScheduler(), // scheduler
	)
	srv.Start()
}
