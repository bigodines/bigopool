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

	// Result can be anything defined by the worker
	Result interface{}

	Dispatcher struct {
		JobQueue   chan Job
		MaxWorkers int
		wg         *sync.WaitGroup
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
		wg:         &sync.WaitGroup{},
		ErrorCh:    errors,
		ResultCh:   done,
	}, nil
}

// Enqueue one or many jobs to process
func (d *Dispatcher) Enqueue(joblist ...Job) {
	d.wg.Add(len(joblist))
	for _, job := range joblist {
		d.JobQueue <- job
	}
}

// Wait blocks until workers are done with their magic
func (d *Dispatcher) Wait() {
	d.wg.Wait()
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
				d.wg.Done()
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
