package psmb

import "github.com/taka-wang/gocron"

var readerMap = NewMbtcpReaderMap()
var writerMap = NewMbtcpWriterMap()
var defaultScheduler = gocron.NewScheduler()

// Dependency injection
var defaultProactiveService = NewPSMBTCP(readerMap, writerMap, defaultScheduler)

// Start start bridge
func Start() {
	defaultProactiveService.Start()
}

// Stop stop bridge
func Stop() {
	defaultProactiveService.Stop()
}
