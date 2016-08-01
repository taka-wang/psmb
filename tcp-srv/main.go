package main

import (
	"fmt"

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
	readerDataStore, err1 := psmbtcp.ReaderDataStoreCreator("Reader")
	if err1 != nil {
		fmt.Println("Fail to create reader ds", err1)
		return
	}
	writerDataStore, err2 := psmbtcp.WriterDataStoreCreator("Writer")
	if err2 != nil {
		fmt.Println("Fail to create writer ds", err2)
		return
	}

	// DI
	srv := psmbtcp.NewService(
		readerDataStore,       // readerMap
		writerDataStore,       // writerMap
		gocron.NewScheduler(), // scheduler
	)
	srv.Start()
}
