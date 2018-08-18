package gopool

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type EchoJob struct{}

func (e EchoJob) Execute() (Result, error) {
	return Result{
		Body: "hello",
	}, nil
}

func TestBootstrap(t *testing.T) {
	d, e := NewDispatcher(2, 2)
	if e != nil {
		t.Fail()
	}
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
	d, e := NewDispatcher(2, 2)
	if e != nil {
		t.Fail()
	}
	d.Run()
	d.Enqueue(ErrorJob{})
	d.Wait()

	assert.Equal(t, 1, len(d.Errors))
}

func TestMixedErrors(t *testing.T) {
	d, e := NewDispatcher(2, 5)
	if e != nil {
		t.Fail()
	}
	d.Run()
	d.Enqueue(ErrorJob{})

	for i := 0; i < 10; i++ {
		d.Enqueue(EchoJob{})
	}

	d.Enqueue(ErrorJob{})
	d.Wait()

	assert.Equal(t, 2, len(d.Errors))
}

func TestAppendResults(t *testing.T) {
	d, e := NewDispatcher(2, 5)
	if e != nil {
		t.Fail()
	}
	d.Run()
	d.Enqueue(EchoJob{}, EchoJob{}, EchoJob{})

	d.Wait()
	assert.Equal(t, 3, len(d.Results))
	assert.Equal(t, 0, len(d.Errors))
}
