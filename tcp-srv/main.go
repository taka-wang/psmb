package main

import (
	"github.com/taka-wang/gocron"
	mbtcp "github.com/taka-wang/psmb/tcp"
)

func main() {
	srv := mbtcp.NewPSMBTCP(
		mbtcp.NewMbtcpReaderMap(), // readerMap
		mbtcp.NewMbtcpWriterMap(), // writerMap
		gocron.NewScheduler(),     // scheduler
	)
	srv.Start()
}
