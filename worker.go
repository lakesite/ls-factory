package factory

import (
	log "github.com/sirupsen/logrus"
	"time"
)

// Worker has an ID, internal channel for jobs, internal channel for reporting
// status to the dispatcher, and a channel for quitting.
type Worker struct {
	ID             int32
	jobs           chan *Job
	dispatchStatus chan *DispatchStatus
	workerCmd      chan *WorkerCommand
	log            *log.Logger
}

type WorkerCommand struct {
	ID      int32
	Command string
}

// CreateNewWorker accepts an ID, channel for worker registration, channel for
// jobs, and a channel for dispatch reports.
func CreateNewWorker(id int32, wCmd chan *WorkerCommand, jobQueue chan *Job, dStatus chan *DispatchStatus, l *log.Logger) *Worker {
	w := &Worker{
		ID:             id,
		jobs:           jobQueue,
		dispatchStatus: dStatus,
		workerCmd:      wCmd,
		log:            l,
	}

	return w
}

// Start enables the worker for processing jobs.
func (w *Worker) Start() {
	go func() {
		for {
			select {
			case job := <-w.jobs:
				w.log.WithFields(log.Fields{
					"ID":     w.ID,
					"Job ID": job.ID,
				}).Info("Worker executing job.")
				job.StartTime = time.Now()
				w.dispatchStatus <- &DispatchStatus{
					Type:   DTYPE_JOB,
					ID:     job.ID,
					Status: DSTATUS_START,
				}

				job.Executable()
				job.EndTime = time.Now()
				w.dispatchStatus <- &DispatchStatus{
					Type:   DTYPE_JOB,
					ID:     job.ID,
					Status: DSTATUS_END,
				}

			case wc := <-w.workerCmd:
				if wc.ID == w.ID || w.ID == 0 {
					if wc.Command == WCMD_QUIT {
						w.log.WithFields(log.Fields{
							"ID": w.ID,
						}).Info("Worker received quit command.")
						w.dispatchStatus <- &DispatchStatus{
							Type:   DTYPE_WORKER,
							ID:     w.ID,
							Status: DSTATUS_QUIT,
						}
						return
					}
				}
			}
		}
	}()
}
