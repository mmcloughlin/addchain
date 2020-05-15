package errutil

import (
	"errors"
	"fmt"
	"io"
)

// ErrNotImplemented is returned when feature is currently not implemented.
var ErrNotImplemented = errors.New("not implemented")

// AssertionFailure is used for an error resulting from the failure of an
// expected invariant.
func AssertionFailure(format string, args ...interface{}) error {
	return fmt.Errorf("assertion failure: "+format, args...)
}

// UnexpectedType builds an error for an unexpected type, typically in a type switch.
func UnexpectedType(t interface{}) error {
	return AssertionFailure("unexpected type %T", t)
}

// Errors is a collection of errors.
type Errors []error

// Add appends errors to the list.
func (e *Errors) Add(err ...error) {
	*e = append(*e, err...)
}

// Err returns an error equivalent to this error list.
// If the list is empty, Err returns nil.
func (e Errors) Err() error {
	if len(e) == 0 {
		return nil
	}
	return e
}

// Error implements the error interface.
func (e Errors) Error() string {
	switch len(e) {
	case 0:
		return "no errors"
	case 1:
		return e[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", e[0], len(e)-1)
}

// CheckClose closes c. If an error occurs it will be written to the error
// pointer errp, if it doesn't already reference an error. This is intended to
// allow you to properly check errors when defering a close call. In this case
// the error pointer should be the address of a named error return.
func CheckClose(errp *error, c io.Closer) {
	if err := c.Close(); err != nil && *errp == nil {
		*errp = err
	}
}
