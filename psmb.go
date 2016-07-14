package main

import (
	"github.com/facebookgo/inject"
	"github.com/taka-wang/gocron"
)

var defaultProactiveService = NewPSMBTCP()

// Start start bridge
func Start() {
	sch := gocron.NewScheduler()
	err := inject.Populate(&defaultProactiveService, sch)
	if err != nil {
		//panic(err)
	}
	defaultProactiveService.Start()
}

// Stop stop bridge
func Stop() {
	defaultProactiveService.Stop()
}
