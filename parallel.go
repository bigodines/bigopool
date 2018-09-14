package bigopool

import (
	"sync"
)

// Parallel runs multiple functions in parallel and collects the errors safely.
func Parallel(ff ...func() error) Errors {
	var wg sync.WaitGroup
	var ee errs

	wg.Add(len(ff))

	for i := range ff {
		f := ff[i]

		go func() {
			defer wg.Done()

			if err := f(); err != nil {
				ee.append(err)
			}
		}()
	}

	wg.Wait()

	return &ee
}
