package gopool

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type EchoJob struct{}

func (e EchoJob) Execute() (Result, error) {
	return "hello", nil
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
type ret struct{}

func (e ErrorJob) Execute() (Result, error) {
	return ret{}, fmt.Errorf("Errored")
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
	_, errs := d.Wait()

	assert.Equal(t, 2, len(errs))
}

func TestAppendResults(t *testing.T) {
	d, e := NewDispatcher(2, 5)
	if e != nil {
		t.Fail()
	}
	d.Run()
	d.Enqueue(EchoJob{}, EchoJob{}, EchoJob{})

	results, errors := d.Wait()
	assert.Equal(t, 3, len(results))
	assert.Equal(t, 0, len(errors))
}

func TestInvalid(t *testing.T) {
	_, e := NewDispatcher(0, 1000)
	if e == nil {
		t.Fail()
	}

	_, e = NewDispatcher(10, 0)
	if e == nil {
		t.Fail()
	}
}
