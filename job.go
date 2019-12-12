package factory

import (
	"time"
)

// JobExecutable defines our spec for jobs, which are functions that return
// an error.
type JobExecutable func() error

// Job has an ID, start time, end time and reference to a JobExecutable.
type Job struct {
	ID         int
	StartTime  time.Time
	EndTime    time.Time
	Executable JobExecutable
}

func CreateNewJob(id int, je JobExecutable) *Job {
	return &Job{
		ID:         id,
		Executable: je,
	}
}
