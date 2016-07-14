package psmb

import "fmt"

var defaultProactiveService = NewPSMBTCP()

// Start start bridge
func Start() {
	fmt.Println("start")
	defaultProactiveService.Start()
}

// Stop stop bridge
func Stop() {
	defaultProactiveService.Stop()
}
