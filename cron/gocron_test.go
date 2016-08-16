package cron

import (
	"fmt"
	"testing"
	"time"

	"github.com/takawang/sugar"
)

func task() {
	fmt.Println("I am runnning task.")
}

func task2() {
	fmt.Println("I am runnning task.")
}

func task3() {
	fmt.Println("I am runnning task.")
}

func taskWithParams(a int, b string) {
	t := time.Now()
	fmt.Println(a, b, t.Format("2006-01-02 15:04:05.000"))
}

func TestScheduler(t *testing.T) {

	s := sugar.New(t)

	s.Title("Job order test")

	s.Assert("`RemoveWithName()` should not raise a deadlock", func(log sugar.Log) bool {
		s := scheduler{
			jobMap:    make(map[string]*Job),
			isStopped: make(chan bool),
			location:  time.Local,
		}
		s.EveryWithName(1, "hello").Seconds().Do(taskWithParams, 1, "1s-hello")
		s.EveryWithName(1, "world").Seconds().Do(taskWithParams, 1, "1s-world")
		s.EveryWithName(2, "hello").Seconds().Do(taskWithParams, 1, "2s-hello")
		fmt.Println("Enable scheduler")

		for i, v := range s.jobMap {
			fmt.Println(i, v.enabled)
		}

		fmt.Println(len(s.jobs))

		s.Start()
		b := s.IsRunning()
		log(b)
		time.Sleep(3 * time.Second)
		fmt.Println("RemoveWithName")
		s.RemoveWithName("world")

		time.Sleep(3 * time.Second)

		for i, v := range s.jobMap {
			fmt.Println(i, v.enabled)
		}

		fmt.Println(len(s.jobs))

		time.Sleep(2 * time.Second)
		s.Stop()
		return true
	})

	s.Assert("`EveryWithName()` should update interval", func(log sugar.Log) bool {
		s := scheduler{
			jobMap:    make(map[string]*Job),
			isStopped: make(chan bool),
			location:  time.Local,
		}
		s.EveryWithName(2, "hello").Seconds().Do(taskWithParams, 2, "2s-hello")
		s.EveryWithName(2, "world").Seconds().Do(taskWithParams, 2, "2s-world")
		fmt.Println("Enable scheduler")
		/*
			for i, v := range s.jobMap {
				fmt.Println(i, v.enabled)
			}
		*/

		for i, v := range s.jobMap {
			fmt.Println(i, v.interval)
		}

		s.Start()
		time.Sleep(4 * time.Second)
		fmt.Println("update job `world` from 2 secs to 3 secs", time.Now().Format("2006-01-02 15:04:05.000"))
		s.EveryWithName(3, "world").Seconds().Do(taskWithParams, 3, "3s-world")
		s.UpdateIntervalWithName("hello1", 3)
		s.UpdateIntervalWithName("hello", 1)
		time.Sleep(3 * time.Second)

		for i, v := range s.jobMap {
			fmt.Println(i, v.interval)
		}

		/*
			for i, v := range s.jobMap {
				fmt.Println(i, v.enabled)
			}
		*/
		fmt.Println("job lengh", len(s.jobs))
		time.Sleep(10 * time.Second)
		s.Stop()
		return true
	})

	s.Assert("`Pause()` and `Resume()` should work", func(log sugar.Log) bool {
		s := scheduler{
			jobMap:    make(map[string]*Job),
			isStopped: make(chan bool),
			location:  time.Local,
		}
		s.EveryWithName(2, "hello").Seconds().Do(taskWithParams, 2, "2s-hello")
		s.EveryWithName(2, "world").Seconds().Do(taskWithParams, 2, "2s-world")
		fmt.Println("Enable scheduler")
		/*
			for i, v := range s.jobMap {
				fmt.Println(i, v.enabled)
			}
		*/

		s.Start()
		time.Sleep(4 * time.Second)
		fmt.Println("pause job `hello`", time.Now().Format("2006-01-02 15:04:05.000"))
		s.PauseWithName("hello")
		/*
			for i, v := range s.jobMap {
				fmt.Println(i, v.enabled)
			}
		*/
		time.Sleep(10 * time.Second)
		fmt.Println("resume job `hello`", time.Now().Format("2006-01-02 15:04:05.000"))
		s.ResumeWithName("hello")
		/*
			for i, v := range s.jobMap {
				fmt.Println(i, v.enabled)
			}
		*/
		time.Sleep(10 * time.Second)
		s.Stop()
		return true
	})

	s.Assert("`PauseAll()` and `ResumeAll()` should work", func(log sugar.Log) bool {
		s := scheduler{
			jobMap:    make(map[string]*Job),
			isStopped: make(chan bool),
			location:  time.Local,
		}
		s.EveryWithName(2, "hello").Seconds().Do(taskWithParams, 2, "2s-hello")
		s.EveryWithName(2, "world").Seconds().Do(taskWithParams, 2, "2s-world")
		fmt.Println("Enable scheduler")
		/*
			for i, v := range s.jobMap {
				fmt.Println(i, v.enabled)
			}
		*/

		s.Start()
		time.Sleep(4 * time.Second)
		fmt.Println("pause all jobs", time.Now().Format("2006-01-02 15:04:05.000"))
		s.PauseAll()
		/*
			for i, v := range s.jobMap {
				fmt.Println(i, v.enabled)
			}
		*/
		time.Sleep(10 * time.Second)
		fmt.Println("resume all job", time.Now().Format("2006-01-02 15:04:05.000"))
		s.ResumeAll()
		/*
			for i, v := range s.jobMap {
				fmt.Println(i, v.enabled)
			}
		*/
		time.Sleep(10 * time.Second)
		s.Stop()
		return true
	})

	s.Assert("`Every()` should append job with order", func(log sugar.Log) bool {

		s := scheduler{
			jobMap:    make(map[string]*Job),
			isStopped: make(chan bool),
			location:  time.Local,
		}

		s.Every(3).Seconds().Do(taskWithParams, 1, "3s")
		s.Every(2).Seconds().Do(taskWithParams, 2, "2s")
		s.Every(5).Seconds().Do(taskWithParams, 3, "5s")
		s.EveryWithName(1, "hello").Seconds().Do(taskWithParams, 4, "1s-4")
		s.EveryWithName(1, "world").Seconds().Do(taskWithParams, 5, "1s-5")
		s.Every(500).Seconds().Do(taskWithParams, 6, "500s")

		s.Every(10).Seconds().Do(taskWithParams, 7, "10s")
		for _, job := range s.jobs {
			log("@interval: %d, param: %s", job.interval, job.tasksParams[0])
		}

		for i, v := range s.jobMap {
			fmt.Printf("map: %s, %p\n", i, v)
		}

		s.Start()

		fmt.Println("add emergency job 8", time.Now().Format("2006-01-02 15:04:05.000"))
		s.Emergency().Do(taskWithParams, 8, "emergency")

		time.Sleep(5 * time.Second)
		for _, job := range s.jobs {
			log("@@interval: %d, param: %s", job.interval, job.tasksParams[0])
		}

		fmt.Println("add emergency job 9", time.Now().Format("2006-01-02 15:04:05.000"))
		s.Emergency().Do(taskWithParams, 9, "emergency")

		//s.Every(1).Seconds().Do(taskWithParams, 10, "1s-10")
		time.Sleep(5 * time.Second)

		// debug
		for _, job := range s.jobs {
			log("interval: %d, param: %s", job.interval, job.tasksParams[0])
		}

		if s.jobs[1].interval == 1 {
			s.Stop()
			return true
		}
		s.Stop()
		return false
	})

	s.Assert("`Remove()` should delete desired job", func(log sugar.Log) bool {
		s := scheduler{
			isStopped: make(chan bool),
			location:  time.Local,
		}

		// add three jobs
		s.Every(3).Seconds().Do(task)
		item := s.Every(2).Seconds().Do(task2)
		s.Every(1).Seconds().Do(task3)

		// debug
		for _, job := range s.jobs {
			log("interval: %d", job.interval)
		}

		// remove one job
		s.Remove(item)

		// remove nil
		bb := s.Remove(nil)
		log(bb)
		bbb := s.RemoveWithName("hello world")
		log(bbb)
		bbbb := s.PauseWithName("hello world")
		log(bbbb)
		bbbbb := s.ResumeWithName("hello world")
		log(bbbbb)

		// debug
		for _, job := range s.jobs {
			log("@interval: %d", job.interval)
		}

		if s.Len() == 2 {
			s.Stop()
			return true
		}
		s.Stop()
		return false
	})

	s.Assert("Only start the sch if not started", func(log sugar.Log) bool {
		s := scheduler{
			isStopped: make(chan bool),
			location:  time.Local,
		}
		s.isRunning = true
		s.Start()
		s.Clear()
		s.Location(time.Local)
		s.RunPending()
		s.NextRun()
		NewScheduler(map[string]string{"hello": "driver"})

		j := newJob(10)
		j.pause()
		j.resume()
		j = j.Second()
		j = j.Seconds()
		j = j.Minute()
		j = j.Minutes()
		j = j.Hour()
		j = j.Hours()
		j = j.Day()
		j = j.Days()
		j.init(time.Now())
		j = j.Weekday(time.Sunday)
		j = j.Sunday()
		j.init(time.Now())
		j = j.Monday()
		j = j.Tuesday()
		j = j.Wednesday()
		j = j.Thursday()
		j = j.Friday()
		j = j.Saturday()
		j = j.Week()
		j.init(time.Now())
		j = j.Weeks()
		j.updateInterval(100)
		j = j.At("10:30")
		b := j.isInit()
		log(b)
		j.init(time.Now())
		j.run()
		j = j.At("24:30")
		i := newJob(0)
		log(i)
		return true
	})

}
