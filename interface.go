// Package psmb a proactive service library for modbus daemon.
//
// Define high-level Contracts(interface) and data types.
//
// By taka@cmwang.net
//

package psmb

import "time"

//
// Interfaces
//

// IProactiveService proactive service contracts,
// all services should implement the following methods.
type IProactiveService interface {
	// Start enable proactive service
	Start()
	// Stop disable proactive service
	Stop()
	// ParseRequest parse requests from IT
	ParseRequest(msg []string) (interface{}, error)
	// HandleRequest handle requests from IT
	HandleRequest(cmd string, r interface{}) error
	// ParseResponse parse responses from OT
	ParseResponse(msg []string) (interface{}, error)
	// HandleResponse handle responses from OT
	HandleResponse(cmd string, r interface{}) error
}

// IWriterTaskDataStore write task interface
//	(Tid, Command) map
type IWriterTaskDataStore interface {
	// Add add request to write task map,
	// params: TID, CMD strings.
	Add(tid, cmd string)

	// Get get request from write task map,
	// params: TID string,
	// return: cmd string, exist flag.
	Get(tid string) (string, bool)

	// Delete remove request from write task map
	// params: TID string.
	Delete(tid string)
}

// ReaderTask read/poll task request
type ReaderTask struct {
	// Name task name
	Name string
	// Cmd zmq frame 1
	Cmd string
	// Req request structure
	Req interface{}
}

// IReaderTaskDataStore read task interface
type IReaderTaskDataStore interface {
	// Add add request to read/poll task map
	Add(name, tid, cmd string, req interface{})

	// GetTaskByID get request via TID from read/poll task map
	GetTaskByID(tid string) (interface{}, bool)

	// GetTaskByName get request via poll name from read/poll task map
	GetTaskByName(name string) (interface{}, bool)

	// GetAll get all requests from read/poll task map
	//	ex: mbtcp: []MbtcpPollStatus
	GetAll() interface{}

	// DeleteAll remove all requests from read/poll task map
	DeleteAll()

	// DeleteTaskByID remove request from via TID from read/poll task map
	DeleteTaskByID(tid string)

	// DeleteTaskByName remove request via poll name from read/poll task map
	DeleteTaskByName(name string)

	// UpdateIntervalByName update poll request interval
	UpdateIntervalByName(name string, interval uint64) error

	// UpdateToggleByName update poll request enabled flag
	UpdateToggleByName(name string, toggle bool) error

	// UpdateAllTogglesByName update all poll request enabled flag
	UpdateAllTogglesByName(toggle bool)
}

// IScheduler keeps a slice of jobs that it executes at a regular interval
type IScheduler interface {

	// Clear removes all of the jobs that have been added to the scheduler
	Clear()

	// Emergency create a emergency job, and adds it to the `Scheduler`
	Emergency() *IJob

	// Every creates a new job, and adds it to the `Scheduler`
	Every(interval uint64) *IJob

	// EveryWithName creates a new job, and adds it to the `Scheduler` and job Map
	EveryWithName(interval uint64, name string) *IJob

	// IsRunning returns true if the job  has started
	IsRunning() bool

	// Location sets the default location of every job created with `Every`.
	// The default location is `time.Local`
	Location(*time.Location)

	// NextRun returns the next next job to be run and the time in which
	// it will be run
	NextRun() (*IJob, time.Time)

	// Remove removes an individual job from the scheduler. It returns true
	// if the job was found and removed from the `Scheduler`
	Remove(*IJob) bool

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

	// Depricated: RunPending runs all of the pending jobs
	RunPending()

	// Start starts the scheduler
	Start()

	// Stop stops the scheduler from executing jobs
	Stop()
}

// IJob keeps a slice of jobs that it executes at a regular interval
type IJob interface {
	Do(task interface{}, params ...interface{}) interface{}
	At(t string) interface{}
	Seconds() interface{}
	Second() interface{}
	Minutes() interface{}
	Minute() interface{}
	Hours() interface{}
	Hour() interface{}
	Days() interface{}
	Day() interface{}
	Weekday(weekday time.Weekday) interface{}
	Monday() interface{}
	Tuesday() interface{}
	Wednesday() interface{}
	Thursday() interface{}
	Friday() interface{}
	Saturday() interface{}
	Sunday() interface{}
	Weeks() interface{}
	Week() interface{}
	Location(loc *time.Location) interface{}
}
