package bigopool

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParallelNoErrors(t *testing.T) {
	x := 0
	m := sync.Mutex{}
	ee := Parallel(
		func() error {
			m.Lock()
			defer m.Unlock()
			x++
			return nil
		},

		func() error {
			m.Lock()
			defer m.Unlock()
			x += 2
			return nil
		},

		func() error {
			m.Lock()
			defer m.Unlock()
			x += 3
			return nil
		},
	)
	assert.True(t, ee.IsEmpty())
	assert.Equal(t, 6, x)
}

func TestParallelWithErrors(t *testing.T) {
	y := 0
	ee := Parallel(
		func() error {
			return errors.New("error 1")
		},

		func() error {
			return errors.New("error 2")
		},

		func() error {
			y--
			return errors.New("error 3")
		},
	)
	assert.False(t, ee.IsEmpty())
	assert.Equal(t, 3, len(ee.All()))
	assert.NotNil(t, ee.ToError())
	assert.Equal(t, -1, y)
}

func TestCancelableParallelNoErrors(t *testing.T) {
	x := 0
	m := sync.Mutex{}
	ee := CancelableParallel(context.Background(),
		func(context.Context) error {
			m.Lock()
			defer m.Unlock()
			x++
			return nil
		},

		func(context.Context) error {
			m.Lock()
			defer m.Unlock()
			x += 2
			return nil
		},

		func(context.Context) error {
			m.Lock()
			defer m.Unlock()
			x += 3
			return nil
		},
	)
	assert.True(t, ee.IsEmpty())
	assert.Equal(t, 6, x)
}

func TestCancelableParallelWithError(t *testing.T) {
	ee := CancelableParallel(context.Background(),
		func(c context.Context) error {
			return errors.New("test")
		},
		func(c context.Context) error {
			// Set a timeout in case the context isn't canceled.
			timer := time.NewTimer(100 * time.Millisecond)

			select {
			case <-c.Done():
				return nil
			case <-timer.C:
				// Can't t.Fatal since we aren't in the main goroutine
				t.Log("Did not cancel context due to previous error")
				t.Fail()
				return nil
			}
		},
	)
	assert.False(t, ee.IsEmpty())
	assert.Len(t, ee.All(), 1)
}
