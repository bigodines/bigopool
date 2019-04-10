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
	}

	// errs is a thread safe struct for appending a slice of errors.
	errs struct {
		mutex sync.Mutex
		all   []error
	}
)

// All returns the underlyings slice of errors.
func (ee *errs) All() []error {
	return ee.all
}

// ToError returns all errors as a single error.
func (ee *errs) ToError() error {
	if ee.IsEmpty() {
		return nil
	}

	return ee
}

// IsEmpty is true if there are no errors.
func (ee *errs) IsEmpty() bool {
	return len(ee.all) == 0
}

// Error implements the error interface.
func (ee *errs) Error() string {
	errorStr := ""
	for _, err := range ee.All() {
		errorStr = fmt.Sprintf("%s\n%s", errorStr, err.Error())
	}

	return errorStr
}

// append safely appends to the error slice.
func (ee *errs) append(err error) {
	ee.mutex.Lock()
	ee.all = append(ee.all, err)
	ee.mutex.Unlock()
}
