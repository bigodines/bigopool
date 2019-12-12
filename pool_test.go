package bigopool

import (
	"fmt"
	"runtime"
	"testing"
	"time"

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

	assert.Equal(t, 1, len(d.Errors.All()))
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

	assert.Equal(t, 2, len(errs.All()))
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
	assert.Equal(t, 0, len(errors.All()))
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

func TestDispatcherCleanup(t *testing.T) {
	ngr := runtime.NumGoroutine()
	d, e := NewDispatcher(10, 20)
	if e != nil {
		t.Fail()
	}

	d.Run()
	d.Enqueue(EchoJob{}, EchoJob{}, EchoJob{})
	d.Wait()

	// sleep so goroutines have time to exit
	time.Sleep(1000 * time.Millisecond)
assert.Equal(t, ngr, runtime.NumGoroutine())
		t.Fail()
	}
}

// Benchmarks
func benchmarkEchoJob(w, q int, b *testing.B) {
	d, e := NewDispatcher(w, q)
	if e != nil {
		b.Fatal()
	}
	d.Run()

	for i := 0; i < b.N; i++ {
		d.Enqueue(EchoJob{})
	}

	d.Wait()
}

func Benchmark1Workers1Queue(b *testing.B)       { benchmarkEchoJob(1, 1, b) }
func Benchmark5Workers1000Queue(b *testing.B)    { benchmarkEchoJob(5, 1000, b) }
func Benchmark10Workers100Queue(b *testing.B)    { benchmarkEchoJob(10, 100, b) }
func Benchmark20Workers200Queue(b *testing.B)    { benchmarkEchoJob(20, 100, b) }
func Benchmark20Workers10000Queue(b *testing.B)  { benchmarkEchoJob(20, 10000, b) }
func Benchmark100Workers10000Queue(b *testing.B) { benchmarkEchoJob(100, 10000, b) }
