package bigopool

import (
	"context"
	"sync"
)

// Parallel runs multiple functions in parallel and collects the errors safely.
func Parallel(ff ...func() error) Errors {
	var wg sync.WaitGroup
	ee := NewErrors()

	wg.Add(len(ff))

	for i := range ff {
		f := ff[i]

		go func() {
			defer wg.Done()

			if err := f(); err != nil {
				ee.Append(err)
			}
		}()
	}

	wg.Wait()

	return ee
}

// CancelableParallel runs multiple functions in parallel and collects the errors safely, while canceling the context
// passed to the remaining functions as soon as a function returns an error.
func CancelableParallel(ctx context.Context, ff ...func(context.Context) error) Errors {
	var wg sync.WaitGroup
	ee := NewErrors()

	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg.Add(len(ff))

	for i := range ff {
		f := ff[i]

		go func() {
			defer wg.Done()

			if err := f(cancelCtx); err != nil {
				cancel()
				ee.Append(err)
			}
		}()
	}

	wg.Wait()

	return ee
}
