package main

import (
	"github.com/taka-wang/gocron"
	"github.com/taka-wang/psmb/mrds"
	"github.com/taka-wang/psmb/mwds"
	psmbtcp "github.com/taka-wang/psmb/tcp"
)

func init() {
	// register data stores from packages
	psmbtcp.Register("Reader", mrds.NewDataStore)
	psmbtcp.Register("Writer", mwds.NewDataStore)
}

func main() {
	// DI & Factory
	if srv, err := psmbtcp.NewService(
		"Reader",
		"Writer",
		gocron.NewScheduler(), // scheduler
	); srv != nil {
		srv.Start()
	}
}
