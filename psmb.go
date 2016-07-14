package main

var defaultBridge = NewMbtcpBridge()

// Start start bridge
func Start() {
	defaultBridge.Start()
}

// Stop stop bridge
func Stop() {
	defaultBridge.Stop()
}
