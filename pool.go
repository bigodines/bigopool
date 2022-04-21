package bigopool

import (
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
)

var (
	ErrNoWorkers = errors.New("Need at least one worker")
	ErrZeroQueue = errors.New("Queue capacity can't be zero")
)

type (
	// Job interface allows bigopool to process anything that implements Execute()
	Job interface {
		Execute() (Result, error)
	}

	// Result can be anything defined by the worker
	Result interface{}

	// Dispatcher is responsible for orchestrating jobs to workers and reporting results back
	Dispatcher struct {
		jobQueue   chan Job
		MaxWorkers int
		wg         *sync.WaitGroup
		// A pool of workers channels that are registered with the dispatcher
		workerPool chan chan Job
		quitCh     chan bool
		// Collect errors
		errorCh  chan error
		resultCh chan Result

		Errors  Errors
		Results []Result
	}
)

// NewDispatcher creates a new dispatcher
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
	quit := make(chan bool, 1)
	return &Dispatcher{
		jobQueue:   jobq,
		MaxWorkers: maxWorkers,
		workerPool: pool,
		wg:         &sync.WaitGroup{},
		errorCh:    errors,
		resultCh:   done,
		quitCh:     quit,
		Errors:     &errs{},
	}, nil
}

// Enqueue one or many jobs to process
func (d *Dispatcher) Enqueue(joblist ...Job) {
	d.wg.Add(len(joblist))
	for _, job := range joblist {
		d.jobQueue <- job
	}
}

// Wait blocks until workers are done with their magic
// return the results and errors
func (d *Dispatcher) Wait() ([]Result, Errors) {
	defer func() {
		// although it does not seem possible this is specifically for Dispatcher wg going negative
		//  if root cause is found and rectified this can be removed.
		if cause := recover(); cause != nil {
			fmt.Println(string(debug.Stack()), cause)
		}
	}()
	//defer d.cleanUp()
	d.wg.Wait()
	d.quitCh <- true
	return d.Results, d.Errors
}

func (d *Dispatcher) cleanUp() {
	close(d.errorCh)
	close(d.resultCh)
	close(d.quitCh)
}

// Run gets the workers ready to work and listens to what they have to say at the end of their job
func (d *Dispatcher) Run() {

	var workerCleanupWG sync.WaitGroup
	// Worker initialization
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.jobQueue, d.errorCh, d.resultCh)
		worker.Start(&workerCleanupWG)
	}

	// Listen for results or errors
	go func() {
		defer func() {
			close(d.jobQueue)
			workerCleanupWG.Wait()
			d.cleanUp()
		}()
		defer func() {
			// although it does not seem possible this is specifically for Dispatcher wg going negative
			//  if root cause is found and rectified this can be removed.
			if cause := recover(); cause != nil {
				fmt.Println(string(debug.Stack()), cause)
			}
		}()
		for {
			select {
			case err := <-d.errorCh:
				d.Errors.append(err)
			case res := <-d.resultCh:
				// If you are changing this code, please note this is not a thread safe append()
				d.Results = append(d.Results, res)
				d.wg.Done()
			case <-d.quitCh:
				return
			}
		}
	}()
}
