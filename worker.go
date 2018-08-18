package gopool

type (
	Worker struct {
		// A pool of workers channels that are registered with the dispatcher
		WorkerPool chan chan Job
		// A channel for receiving a job that was dispatched
		jobCh chan Job
		// A channel for receiving a worker termination signal
		// (quits after processing)
		quit chan bool

		// reporting channels
		errCh    chan error
		resultCh chan Result
	}
)

// NewWorker creates a new worker that can be registered to a WorkerPool
// and receive jobs
func NewWorker(workerPool chan chan Job, errCh chan error, resultCh chan Result) Worker {
	return Worker{
		WorkerPool: workerPool,
		jobCh:      make(chan Job),
		quit:       make(chan bool),
		errCh:      errCh,
		resultCh:   resultCh,
	}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.jobCh

			select {
			case job := <-w.jobCh:
				result, err := job.Execute()
				if err != nil {
					w.errCh <- err
				}
				w.resultCh <- result
			}
		}
	}()
}
