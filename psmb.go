package main

import (
	"github.com/codegangsta/inject"
	"github.com/taka-wang/gocron"
)

var defaultProactiveService = NewPSMBTCP()

// Start start bridge
func Start() {
	sch := gocron.NewScheduler()
	inj := inject.New()
	inj.Map(sch)
	inj.Apply(&defaultProactiveService)
	defaultProactiveService.Start()
}

// Stop stop bridge
func Stop() {
	defaultProactiveService.Stop()
}
