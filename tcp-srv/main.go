package main

import (
	"fmt"
	"net"

	cron "github.com/taka-wang/psmb/cron"
	dbwds "github.com/taka-wang/psmb/dbwds"
	mr "github.com/taka-wang/psmb/mrds"
	mw "github.com/taka-wang/psmb/mwds"
	mbtcp "github.com/taka-wang/psmb/tcp"
)

var hostName string

func init() {

	// get hostname
	host, err := net.LookupHost("redis")
	if err != nil {
		fmt.Println("local run")
		hostName = "127.0.0.1"
	} else {
		fmt.Println("docker run")
		hostName = host[0] //docker
	}

	// register data stores from packages
	mbtcp.Register("MemReader", mr.NewDataStore)
	mbtcp.Register("MemWriter", mw.NewDataStore)
	mbtcp.Register("Cron", cron.NewScheduler)
	mbtcp.Register("RedisWriter", dbwds.NewDataStore)
}

func main() {
	// dependency injection & factory pattern
	if srv, _ := mbtcp.NewService(
		"MemReader", // Reader Data Store
		"MemWriter", // Writer Data Store
		"Cron",      // Scheduler
	); srv != nil {
		srv.Start()
	}
}
