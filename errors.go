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
		Append(err error)
	}

	// errs is a thread safe struct for appending a slice of errors.
	errs struct {
		mutex sync.RWMutex
		all   []error
	}
)

// Creates an empty error struct
func NewErrors() *errs {
	return &errs{
		all:   []error{},
		mutex: sync.RWMutex{},
	}
}

// All returns the underlyings slice of errors.
func (ee *errs) All() []error {
	ee.mutex.RLock()
	defer ee.mutex.RUnlock()
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
	return len(ee.All()) == 0
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
func (ee *errs) Append(err error) {
	ee.mutex.Lock()
	defer ee.mutex.Unlock()
	ee.all = append(ee.all, err)

}
