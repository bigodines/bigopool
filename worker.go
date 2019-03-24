package bigopool

type (
	Worker struct {
		// A pool of workers channels that are registered with the dispatcher
		workerPool chan chan Job
		// A channel for receiving a job that was dispatched
		jobCh chan Job

		// reporting channels
		errCh    chan error
		resultCh chan Result

		quitCh chan bool
	}
)

// NewWorker creates a new worker that can be registered to a workerPool
// and receive jobs
func NewWorker(workerPool chan chan Job, errCh chan error, resultCh chan Result, quitCh chan bool) Worker {
	return Worker{
		workerPool: workerPool,
		jobCh:      make(chan Job),
		errCh:      errCh,
		resultCh:   resultCh,
		quitCh:     quitCh,
	}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.workerPool <- w.jobCh

			select {
			case job := <-w.jobCh:
				result, err := job.Execute()
				if err != nil {
					w.errCh <- err
				}
				w.resultCh <- result
				w.quitCh <- true
			}
		}
	}()
}
