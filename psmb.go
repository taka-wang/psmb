package psmb

import "github.com/taka-wang/gocron"

// Dependency injection
var defaultProactiveService = NewPSMBTCP(
	NewMbtcpReaderMap(),   // readerMap
	NewMbtcpWriterMap(),   // writerMap
	gocron.NewScheduler(), // scheduler
)

// Start start default proactive service
func Start() {
	defaultProactiveService.Start()
}

// Stop stop default proactive service
func Stop() {
	defaultProactiveService.Stop()
}
