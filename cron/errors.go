package cron

import "errors"

var (
	// ErrTaskIsNotAFuncError is the error panicked when a task passed to `Job.Do`
	ErrTaskIsNotAFuncError = errors.New("the `task` your a scheduling must be of type func")

	// ErrMissmatchedTaskParams is the error panicked when someone passes too many or too few params to `Job.Do`
	ErrMissmatchedTaskParams = errors.New("the `task` your a scheduling must be of type func")

	// ErrJobIsNotInitialized is the error panicked when a job is scheduled that was not initialized
	ErrJobIsNotInitialized = errors.New("this job was not intialized")

	// ErrIncorrectTimeFormat is the error panicked when `At` is passed an incorrect time
	ErrIncorrectTimeFormat = errors.New("the time format is incorrect")

	// ErrIntervalNotValid error panicked when the interval is not valid
	ErrIntervalNotValid = errors.New("the interval must be greater than 0")
)
