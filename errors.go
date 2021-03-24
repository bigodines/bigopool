package bigopool

import (
	"sync"

	"go.uber.org/multierr"
)

type (
	// Errors is an interface for accessing a slice of errors.
	Errors interface {
		All() []error
		ToError() error
		IsEmpty() bool
	}

	// errs is a thread safe struct for appending a slice of errors.
	errs struct {
		mutex sync.Mutex
		errs  error
	}
)

// All returns the underlyings slice of errors.
func (ee *errs) All() []error {
	return multierr.Errors(ee.errs)
}

// ToError returns all errors as a single error.
func (ee *errs) ToError() error {
	if ee.errs == nil {
		return nil
	}

	return ee.errs
}

// IsEmpty is true if there are no errors.
func (ee *errs) IsEmpty() bool {
	return len(ee.All()) == 0
}

// Error implements the error interface.
func (ee *errs) Error() string {
	if ee.errs == nil {
		return ""
	}

	return ee.errs.Error()
}

// append safely appends to the error slice.
func (ee *errs) append(err error) {
	ee.mutex.Lock()
	ee.errs = multierr.Append(ee.errs, err)
	ee.mutex.Unlock()
}
