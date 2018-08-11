package gopool

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type EchoJob struct{}

func (e EchoJob) Execute() (Result, error) {
	return Result{}, nil
}

func TestBootstrap(t *testing.T) {
	d := NewDispatcher(2, 2)
	d.Run()
	for i := 0; i < 10; i++ {
		d.Enqueue(EchoJob{})
	}
	d.Wait()
}

type ErrorJob struct{}

func (e ErrorJob) Execute() (Result, error) {
	return Result{}, fmt.Errorf("Errored")
}
func TestErrors(t *testing.T) {
	d := NewDispatcher(2, 2)
	d.Run()
	d.Enqueue(ErrorJob{})
	d.Wait()

	assert.Equal(t, 1, len(d.Errors))
}

func TestMixedErrors(t *testing.T) {
	d := NewDispatcher(2, 2)
	d.Run()
	d.Enqueue(ErrorJob{})

	for i := 0; i < 10; i++ {
		d.Enqueue(EchoJob{})
	}

	d.Enqueue(ErrorJob{})
	d.Wait()

	assert.Equal(t, 2, len(d.Errors))

}
