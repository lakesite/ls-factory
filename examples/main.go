package main

import (
	"time"

	"github.com/lakesite/ls-factory"
	"github.com/sirupsen/logrus"
)

// GenericJobStruct is a generic structure we'll use for submitting a job,
// this can be any structure.
type GenericJobStruct struct {
	Name string
}

// An example function, ProcessWork.
func (gjs *GenericJobStruct) ProcessWork() error {
	// Generic process, let's print out something, then wait for a time:
	log.WithFields(logrus.Fields{"Name": gjs.Name}).Info("Processing some work.")
	time.Sleep(5 * time.Second)
	log.Info("Work finished.")
	return nil
}

var log = logrus.New()

func main() {
	// Create two separate jobs to execute:
	gjs1 := &GenericJobStruct{Name: "Testing1"}
	gjs2 := &GenericJobStruct{Name: "Testing2"}

	// We'll use the convention of simply wrapping the work in a function literal.
	process1 := func() error {
		log.WithFields(logrus.Fields{"Name": gjs1.Name}).Info("Wrapping work in an anonymous function.")
		return gjs1.ProcessWork()
	}

	process2 := func() error {
		log.WithFields(logrus.Fields{"Name": gjs2.Name}).Info("Wrapping work in an anonymous function.")
		return gjs2.ProcessWork()
	}

	// Create a dispatcher:
	dispatcher := factory.CreateNewDispatcher(log)

	// Enqueue n*2 jobs:
	n := 10
	for i:= 0; i<n; i++ {
		dispatcher.QueueJob(process1)
		dispatcher.QueueJob(process2)
	}

	// Tell the dispatcher to start processing jobs with n workers:
	if err := dispatcher.Start(n); err != nil {
		log.Fatal("Dispatch startup error: %s", err.Error())
	}

	for dispatcher.Running() {
		if dispatcher.Finished() {
			log.Info("All jobs finished.")
			break
		}
	}
}