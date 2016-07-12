package psmb

var defaultProactiveService = NewPSMBTCP()

// Start start bridge
func Start() {
	defaultProactiveService.Start()
}

// Stop stop bridge
func Stop() {
	defaultProactiveService.Stop()
}
