package gopool

import (
	"errors"
	"sync"
)

var (
	ErrNoWorkers = errors.New("Need at least one worker")
	ErrZeroQueue = errors.New("Queue capacity can't be zero")
)

type (
	// gopool can process anything that implements this interface
	Job interface {
		Execute() (Result, error)
	}

	Result struct {
		Response interface{}
	}

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

	Dispatcher struct {
		JobQueue   chan Job
		MaxWorkers int
		WaitGroup  *sync.WaitGroup
		// A pool of workers channels that are registered with the dispatcher
		WorkerPool chan chan Job
		// Collect errors
		ErrorCh  chan error
		ResultCh chan Result

		Errors  []error
		Results []Result
	}
)

func NewDispatcher(maxWorkers int, queueSize int) (*Dispatcher, error) {
	if maxWorkers < 1 {
		return nil, ErrNoWorkers
	}

	if queueSize < 1 {
		return nil, ErrZeroQueue
	}
	pool := make(chan chan Job, maxWorkers)
	jobq := make(chan Job, queueSize)
	errors := make(chan error)
	done := make(chan Result)
	return &Dispatcher{
		JobQueue:   jobq,
		MaxWorkers: maxWorkers,
		WorkerPool: pool,
		WaitGroup:  &sync.WaitGroup{},
		ErrorCh:    errors,
		ResultCh:   done,
	}, nil
}

// Enqueue one or many jobs to process
func (d *Dispatcher) Enqueue(joblist ...Job) {
	d.WaitGroup.Add(len(joblist))
	for _, job := range joblist {
		d.JobQueue <- job
	}
}

// Wait blocks until workers are done with their magic
func (d *Dispatcher) Wait() {
	d.WaitGroup.Wait()
}

// Run gets the workers ready to work and listens to what they have to say at the end of their job
func (d *Dispatcher) Run() {
	// Worker initialization
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, d.ErrorCh, d.ResultCh)
		worker.Start()
	}

	// Get ready to assign tasks
	go d.dispatch()

	// Listen for results or errors
	go func() {
		for {
			select {
			case err := <-d.ErrorCh:
				d.Errors = append(d.Errors, err)
			case res := <-d.ResultCh:
				d.Results = append(d.Results, res)
				d.WaitGroup.Done()
			}
		}
	}()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool
				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}

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
			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}
