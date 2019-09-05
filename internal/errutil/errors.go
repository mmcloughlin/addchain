package errutil

import (
	"fmt"

	"golang.org/x/xerrors"
)

// ErrNotImplemented is returned when feature is currently not implemented.
var ErrNotImplemented = xerrors.New("not implemented")

// AssertionFailure is used for an error resulting from the failure of an
// expected invariant.
func AssertionFailure(format string, args ...interface{}) error {
	return xerrors.Errorf("assertion failure: "+format, args...)
}

// UnexpectedType builds an error for an unexpected type, typically in a type switch.
func UnexpectedType(t interface{}) error {
	return AssertionFailure("unexpected type %T", t)
}

// Errors is a collection of errors.
type Errors []error

// Add appends an error to the list.
func (e *Errors) Add(err error) {
	*e = append(*e, err)
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
