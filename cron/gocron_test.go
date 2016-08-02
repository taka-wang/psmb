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

/*
func TestJob(t *testing.T) {

	// note: we're defining today as the first of the month so we can test an important edge case in the lastRun
	// calculation when a jobs lastRun time occured the previous month from when the job initialized
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), 1, now.Hour(), now.Minute(), now.Second(), 0, now.Location())
	aMinuteAgo := today.Add(-time.Minute)
	aMinuteFromNow := today.Add(time.Minute)
	aMinuteAgoAtTime := fmt.Sprintf("%02d:%02d", aMinuteAgo.Hour(), aMinuteAgo.Minute())
	aMinuteFromNowAtTime := fmt.Sprintf("%02d:%02d", aMinuteFromNow.Hour(), aMinuteFromNow.Minute())
	s := sugar.New(t)

	s.Title("Day")

	s.Assert("`Job.Every(...).Day().At(...)`", func(log sugar.Log) bool {
		// try this with 20 random day intervals
		for i := 20; i > 0; i-- {

			// get a random interval of days [1, 5]
			rand.Seed(time.Now().UnixNano())
			interval := 1 + uint64(rand.Int())%6

			// create and init the job
			job := newJob(interval).Day().At(aMinuteFromNowAtTime)
			job.init(today)

			// jobs last run should be `interval` days from now
			aMinuteFromNowIntervalDaysFromnow := aMinuteFromNow
			if !job.lastRun.Equal(aMinuteFromNow) {
				log("the lastRun did not occur a minute from now days ago")
				log(job.lastRun, aMinuteFromNowIntervalDaysFromnow)
				return false
			}
			//
			// // jobs last run should be `interval` days from now
			// aMinuteAgoIntervalDaysAgo := aMinuteAgo.Add(-1 * Day * time.Duration(interval))
			// if !job.lastRun.Equal(aMinuteAgoIntervalDaysAgo) {
			// 	log("the lastRun did not occur %d days ago", interval)
			// 	log(job.lastRun, aMinuteFromNowIntervalDaysAgo)
			// 	return false
			// }

			// jobs next run is should be today
			if !job.nextRun.Equal(aMinuteFromNow) {
				log("the nextRun will not happen a minute from now")
				log(job.nextRun, aMinuteFromNow)
				return false
			}

			// after run, the nextRun is interval days from the previous nextRun
			job.run()
			aMinutFromNowIntervalDaysAfterNextRun := aMinuteFromNow.Add(Day * time.Duration(interval))
			if !job.nextRun.Equal(aMinutFromNowIntervalDaysAfterNextRun) {
				log("the next nextRun will not happen in %d days", interval)
				log(job.nextRun, aMinutFromNowIntervalDaysAfterNextRun)
				return false
			}

		}

		return true
	})

	s.Assert("`Job.Every(...).Day.At(...)` set to the past", func(log sugar.Log) bool {
		// try this with 20 random day intervals
		for i := 20; i > 0; i-- {

			// get a random interval of days [1, 5]
			rand.Seed(time.Now().UnixNano())
			interval := 1 + uint64(rand.Int())%5

			// create and init the job
			job := newJob(interval).Day().At(aMinuteAgoAtTime)
			job.init(today)

			// jobs last run interval days from tomorrow
			aMinuteAgoIntervalDaysFromTomorrow := aMinuteAgo.Add(Day).Add(-1 * Day * time.Duration(interval))
			if !job.lastRun.Equal(aMinuteAgoIntervalDaysFromTomorrow) {
				log("the lastRun did not %d days from tomorrow", interval)
				log(job.lastRun, aMinuteAgoIntervalDaysFromTomorrow)
				return false
			}

			// jobs next run is tomorrow
			aMinuteAgoTomorrow := aMinuteAgo.Add(Day)
			if !job.nextRun.Equal(aMinuteAgoTomorrow) {
				log("the nextRun will not occur tomorrow")
				log(job.nextRun, aMinuteAgoTomorrow)
				return false
			}

			// after run, the nextRun is interval days from the previous nextRun
			job.run()
			aMinutAgoIntervalDaysAfterNextRun := aMinuteAgoTomorrow.Add(Day * time.Duration(interval))
			if !job.nextRun.Equal(aMinutAgoIntervalDaysAfterNextRun) {
				log("the next nextRun will not happen in %d days", interval)
				log(job.nextRun, aMinutAgoIntervalDaysAfterNextRun)
				return false
			}
		}

		return true
	})

	s.Title("Week")

	s.Assert("`Job.Every(...).Weekday(...).At(...)` set to the past", func(log sugar.Log) bool {
		// try this with 20 random weekdays and week intervals
		for i := 20; i > 0; i-- {

			// get a random interval of weeks [1, 52]
			rand.Seed(time.Now().UnixNano())
			interval := 1 + uint64(rand.Int())%52

			// get a random day of the week that is today or before today
			rand.Seed(time.Now().UnixNano())
			weekday := time.Weekday(rand.Int() % int(today.Weekday()+1))
			durationAfterWeekday := time.Duration(weekday-today.Weekday()) * 24 * time.Hour

			// create and init the job
			job := newJob(interval).Weekday(weekday).At(aMinuteAgoAtTime)
			job.init(today)

			// jobs lastRun was interval weeks ago from next week
			aMinuteAgoIntervalWeeksFromNextWeek := aMinuteAgo.Add(durationAfterWeekday).Add(Week).Add(-1 * Week * time.Duration(interval))
			if !job.lastRun.Equal(aMinuteAgoIntervalWeeksFromNextWeek) {
				log("the lastRun did not occur %d weeks ago", interval+1)
				log(weekday, aMinuteAgoIntervalWeeksFromNextWeek.Weekday(), job.nextRun.Weekday(), job.lastRun, aMinuteAgoIntervalWeeksFromNextWeek)
				return false
			}

			// jobs next run is next week
			aMinuteAgoNextWeek := aMinuteAgo.Add(durationAfterWeekday).Add(Week)
			if !job.nextRun.Equal(aMinuteAgoNextWeek) {
				log("the nextRun will not occur next week")
				log(weekday, aMinuteAgoNextWeek.Weekday(), job.nextRun.Weekday(), job.nextRun, aMinuteAgoNextWeek)
				return false
			}

			// after run, the nextRun is interval weeks from the previous nextRun
			job.run()
			aMinutAgoIntervalWeeksAfterNextRun := aMinuteAgoNextWeek.Add(Week * time.Duration(interval))
			if !job.nextRun.Equal(aMinutAgoIntervalWeeksAfterNextRun) {
				log("the next nextRun will not happen in %d weeks", interval)
				log(job.nextRun, aMinutAgoIntervalWeeksAfterNextRun)
				return false
			}

		}

		return true
	})

	s.Assert("`Job.Every(...).Weekday(...).At(...)` set to the future", func(log sugar.Log) bool {
		// try this with 20 random weekdays and week intervals
		for i := 20; i > 0; i-- {

			// get a random interval of weeks [1, 52]
			rand.Seed(time.Now().UnixNano())
			interval := 1 + uint64(rand.Int())%52

			// get a random day of the week that is today or after today
			rand.Seed(time.Now().UnixNano())
			weekday := time.Weekday(int(today.Weekday()) + rand.Int()%(7-int(today.Weekday())))
			durationUntilWeekday := time.Duration(weekday-today.Weekday()) * 24 * time.Hour

			// create and init the job
			job := newJob(interval).Weekday(weekday).At(aMinuteFromNowAtTime)
			job.init(today)

			// jobs last run was interval weeks ago
			aMinuteFromNowIntervalWeeksAgo := aMinuteFromNow.Add(durationUntilWeekday).Add(-1 * Week * time.Duration(interval))
			if !job.lastRun.Equal(aMinuteFromNowIntervalWeeksAgo) {
				log("the lastRun did not occur %d weeks ago", interval)
				log(weekday, aMinuteFromNowIntervalWeeksAgo.Weekday(), job.nextRun.Weekday(), job.lastRun, aMinuteFromNowIntervalWeeksAgo)
				return false
			}

			// jobs next run is this week
			thisWeekdayAMinuteFromNow := aMinuteFromNow.Add(durationUntilWeekday)
			if !job.nextRun.Equal(thisWeekdayAMinuteFromNow) {
				log("the nextRun will not occur this week")
				log(weekday, thisWeekdayAMinuteFromNow.Weekday(), job.nextRun.Weekday(), job.nextRun, thisWeekdayAMinuteFromNow)
				return false
			}

			// after run, the nextRun is interval weeks from the previous nextRun
			job.run()
			aMinutAgoIntervalWeeksAfterNextRun := thisWeekdayAMinuteFromNow.Add(Week * time.Duration(interval))
			if !job.nextRun.Equal(aMinutAgoIntervalWeeksAfterNextRun) {
				log("the next nextRun will not happen in %d weeks", interval)
				log(job.nextRun, aMinutAgoIntervalWeeksAfterNextRun)
				return false
			}
		}

		return true
	})

	s.Title("Time")

	s.Assert("`Job.Hour()` causes lastRun to be now and nextRun to be `interval` hour(s) from now", func(log sugar.Log) bool {
		// TODO: implement test
		return false
	})

	s.Assert("`Job.Minute()` causes lastRun to be now and nextRun to be `interval` minute(s) from now", func(log sugar.Log) bool {
		// TODO: implement test
		return false
	})

	s.Assert("`Job.Second()` causes lastRun to be now and nextRun to be `interval` second(s) from now", func(log sugar.Log) bool {
		// TODO: implement test
		return false
	})
}

*/

func TestScheduler(t *testing.T) {

	s := sugar.New(t)

	/*
		s.Assert("`runPending(...)` runs all pending jobs", func(log sugar.Log) bool {
			// TODO: implement test
			return false
		})

		s.Assert("`Start()`, `IsRunning()` and `Stop()` perform correctly in asynchrnous environments", func(log sugar.Log) bool {
			// TODO: implement test
			return false
		})

		s.Assert("`Start()` triggers runPending(...) every second", func(log sugar.Log) bool {
			// TODO: implement test
			return false
		})
	*/

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

}
