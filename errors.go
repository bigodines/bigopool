package bigopool

import (
	"fmt"
	"sync"
)

type (
	// Errors is an interface for accessing a slice of errors.
	Errors interface {
		All() []error
		ToError() error
		IsEmpty() bool
		append(err error)
	}

	// errs is a thread safe struct for appending a slice of errors.
	errs struct {
		mutex sync.Mutex
		errs  []error
	}
)

// All returns the underlyings slice of errors.
func (ee *errs) All() []error {
	return ee.errs
}

// ToError returns all errors as a single error.
func (ee *errs) ToError() error {
	if len(ee.errs) == 0 {
		return nil
	}

	err := ee.errs[0]
	for _, otherErr := range ee.errs[1:] {
		err = fmt.Errorf("%v; %w", err, otherErr)
	}

	return err
}

// IsEmpty is true if there are no errors.
func (ee *errs) IsEmpty() bool {
	return len(ee.errs) == 0
}

// Error implements the error interface.
func (ee *errs) Error() string {
	if len(ee.errs) == 0 {
		return ""
	}

	return ee.ToError().Error()
}

// append safely appends to the error slice.
func (ee *errs) append(err error) {
	ee.mutex.Lock()
	ee.errs = append(ee.errs, err)
	ee.mutex.Unlock()
}
