package bigopool

import (
	"fmt"
)

type (
	Worker struct {
		// A channel for receiving a job that was dispatched
		jobCh chan Job

		// reporting channels
		errCh    chan error
		resultCh chan Result
	}
)

// NewWorker creates a new worker that can be registered to a workerPool
// and receive jobs
func NewWorker(jobCh chan Job, errCh chan error, resultCh chan Result) Worker {
	return Worker{
		jobCh:    jobCh,
		errCh:    errCh,
		resultCh: resultCh,
	}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		defer func() {
			cause := recover()
			if cause != nil {
				fmt.Println("panic recovered", cause)
			}
		}()
		for {
			select {
			case job, more := <-w.jobCh:
				if job != nil {
					result, err := job.Execute()
					if err != nil {
						w.errCh <- err
					}
					w.resultCh <- result
				}

				// w.jobCh has been closed by Dispatcher.Wait(), so we're done.
				if !more {
					return
				}
			}
		}
	}()
}
