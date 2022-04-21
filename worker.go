package bigopool

import (
	"fmt"
	"runtime/debug"
	"sync"
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
func (w Worker) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer func() {
			// we probably don't want this because this will essentially hide panics that happen in
			// the worker job.Execute() but with the WaitGroup negative error there will probably be a panic that happens on
			// writing to a closed channel
			if cause := recover(); cause != nil {
				fmt.Println(string(debug.Stack()), cause)
			}
		}()
		defer wg.Done()
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
