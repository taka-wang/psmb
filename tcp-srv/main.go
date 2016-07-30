package main

import (
	"github.com/taka-wang/gocron"
	mbtcp "github.com/taka-wang/psmb/tcp"
)

func main() {
	srv := mbtcp.NewService(
		mbtcp.NewReaderMap(),  // readerMap
		mbtcp.NewWriterMap(),  // writerMap
		gocron.NewScheduler(), // scheduler
	)
	srv.Start()
}
