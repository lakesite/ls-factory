package factory

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

/*
 * DispatchStatus is a struct for passing job and worker status reports.
 * Type can be: worker, job, ID represents the ID of the worker or job.
 * Status: see WSTATUS constants.
 */
type DispatchStatus struct {
	Type   string
	ID     int
	Status string
}

// Dispatch keeps track of an internal job request queue, a work queue of jobs
// that will be processed, a worker queue of workers, and a channel for status
// reports for jobs and workers.
type Dispatcher struct {
	jobCounter     int                  // internal counter for number of jobs
	jobQueue       chan *Job            // channel of submitted jobs
	dispatchStatus chan *DispatchStatus // channel for job/worker status reports
	workQueue      chan *Job            // channel of work dispatched
	workerQueue    chan *Worker         // channel of workers
	workerCommand  chan *WorkerCommand  // channel for worker commands
	log            *log.Logger          // log API (logrus)
	running        bool                 // Is the dispatcher running
}

// CreateNewDispatcher creates a new dispatcher by making the necessary channels
// for goroutines to communicate and initializes an internal job counter.
func CreateNewDispatcher(l *log.Logger) *Dispatcher {
	return &Dispatcher{
		jobCounter:     0,
		jobQueue:       make(chan *Job),
		dispatchStatus: make(chan *DispatchStatus),
		workQueue:      make(chan *Job),
		workerCommand:  make(chan *WorkerCommand),
		log:            l,
		running:        false,
	}
}

// QueueJob accepts a process (func() error) of JobExecutable, creates a job
// which will be tracked, and adds the job into the internal work queue for
// execution.
// If there's a constraint on the number of jobs, return an error.
func (d *Dispatcher) QueueJob(je JobExecutable) error {
	// Create a new job:
	j := CreateNewJob(d.jobCounter, je)

	// Add the job to the internal queue:
	go func() { d.jobQueue <- j }()

	// Increment the internal job counter:
	d.jobCounter++
	d.log.WithFields(log.Fields{
		"jobCounter": d.jobCounter,
	}).Info("Job Queued.")

	return nil
}

// Finished returns true if we have no more jobs to process, otherwise false.
func (d *Dispatcher) Finished() bool {
	if d.jobCounter < 1 {
		return true
	} else {
		return false
	}
}

// Running returns true if the dispatcher has been issued Start() and is running
func (d *Dispatcher) Running() bool {
	return d.running
}

// Start has the dispatcher create workers to handle jobs, then creates a
// goroutine to handle passing jobs in the queue off to workers and processing
// dispatch status reports.
func (d *Dispatcher) Start(numWorkers int) error {
	if numWorkers < 1 {
		return errors.New("Start requires >= 1 workers.")
	}

	// Create numWorkers:
	for i := 0; i < numWorkers; i++ {
		worker := CreateNewWorker(i, d.workerCommand, d.workQueue, d.dispatchStatus, d.log)
		worker.Start()
	}

	d.running = true

	// wait for work to be added then pass it off.
	go func() {
		for {
			select {
			case job := <-d.jobQueue:
				d.log.WithFields(log.Fields{
					"ID": job.ID,
				}).Info("Adding a new job to the queue for dispatching.")
				// Add the job to the work queue, and don't block the dispatcher.
				go func() { d.workQueue <- job }()

			case ds := <-d.dispatchStatus:
				d.log.WithFields(log.Fields{
					"Type":   ds.Type,
					"ID":     ds.ID,
					"Status": ds.Status,
				}).Info("Received a dispatch status report.")

				if ds.Type == DTYPE_WORKER {
					if ds.Status == DSTATUS_QUIT {
						d.log.WithFields(log.Fields{
							"ID": ds.ID,
						}).Info("Worker quit.")
					}
				}

				if ds.Type == DTYPE_JOB {
					if ds.Status == DSTATUS_START {
						d.log.WithFields(log.Fields{
							"ID": ds.ID,
						}).Info("Job started.")
					}

					if ds.Status == DSTATUS_END {
						d.log.WithFields(log.Fields{
							"ID": ds.ID,
						}).Info("Job finished.")
						d.jobCounter--
					}
				}
			}
		}
	}()
	return nil
}
