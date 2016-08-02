package cron

import "time"

// Scheduler keeps a slice of jobs that it executes at a regular interval
type Scheduler interface {

	// Clear removes all of the jobs that have been added to the scheduler
	Clear()

	// Emergency create a emergency job, and adds it to the `Scheduler`
	Emergency() *Job

	// Every creates a new job, and adds it to the `Scheduler`
	Every(interval uint64) *Job

	// EveryWithName creates a new job, and adds it to the `Scheduler` and job Map
	EveryWithName(interval uint64, name string) *Job

	// IsRunning returns true if the job  has started
	IsRunning() bool

	// Location sets the default location of every job created with `Every`.
	// The default location is `time.Local`
	Location(*time.Location)

	// NextRun returns the next next job to be run and the time in which
	// it will be run
	NextRun() (*Job, time.Time)

	// Remove removes an individual job from the scheduler. It returns true
	// if the job was found and removed from the `Scheduler`
	Remove(*Job) bool

	// UpdateIntervalWithName update an individual job's interval from the scheduler by name.
	// It returns true if the job was found and update interval
	UpdateIntervalWithName(name string, interval uint64) bool

	// RemoveWithName removes an individual job from the scheduler by name. It returns true
	// if the job was found and removed from the `Scheduler`
	RemoveWithName(string) bool

	// PauseWithName pause an individual job by name. It returns true if the job was found and set enabled
	PauseWithName(string) bool

	// PauseAll disable all jobs
	PauseAll()

	// ResumeWithName resume an individual job by name. It returns true if the job was found and set enabled
	ResumeWithName(string) bool

	// ResumeAll resume all jobs
	ResumeAll()

	// Depricated: RunAll runs all of the jobs regardless of wether or not
	// they are pending
	RunAll()

	// RunAllWithDelay runs all of the jobs regardless of wether or not
	// they are pending with a delay
	RunAllWithDelay(time.Duration)

	// Depricated: RunPending runs all of the pending jobs
	RunPending()

	// Start starts the scheduler
	Start()

	// Stop stops the scheduler from executing jobs
	Stop()
}
