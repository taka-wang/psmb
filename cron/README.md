# cron: A Golang Job Scheduling Package.

Cron is a Golang job scheduling package which lets you run Go functions periodically at pre-determined interval using a simple, human-friendly syntax.


``` go
package main

import (
	"fmt"
	"github.com/taka-wang/psmb/cron"
	"time"
)

func task() {
	fmt.Println("I am runnning task.")
}

func taskWithParams(a int, b string) {
	fmt.Println(a, b)
}

func main() {
	// Do jobs with params
	cron.Every(1).Second().Do(taskWithParams, 1, "hello")

	// Do jobs without params
	mjob := cron.Every(1).Second().Do(task)
	cron.Every(2).Seconds().Do(task)
	cron.Every(1).Minute().Do(task)
	cron.Every(2).Minutes().Do(task)
	cron.Every(1).Hour().Do(task)
	cron.Every(2).Hours().Do(task)
	cron.Every(1).Day().Do(task)
	cron.Every(2).Days().Do(task)

	// Do jobs on specific weekday
	cron.Every(1).Monday().Do(task)
	cron.Every(1).Thursday().Do(task)

	// function At() take a string like 'hour:min'
	cron.Every(1).Day().At("10:30").Do(task)
	cron.Every(1).Monday().At("18:30").Do(task)

	// remove, clear and next_run
	_, time := cron.NextRun()
	fmt.Println(time)

	cron.Remove(mjob)
	cron.Clear()

	// function Start start all the pending jobs
	cron.Start()
	
	// trigger emergency job
	cron.Emergency().Do(taskWithParams, 9, "emergency")

	// also , you can create a your new scheduler,
	// to run two scheduler concurrently
	s := cron.NewScheduler()
	s.Every(3).Seconds().Do(task)
	s.Start()
	for {
		time.Sleep(300 * time.Millisecond)
	}

}
```

---

## UML

![PlantUML model](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/taka-wang/puml/master/cron.puml)
