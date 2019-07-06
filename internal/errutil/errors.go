package errutil

import "fmt"

// AssertionFailure is used for an error resulting from the failure of an
// expected invariant.
func AssertionFailure(format string, args ...interface{}) error {
	return fmt.Errorf("assertion failure: "+format, args...)
}

// UnexpectedType builds an error for an unexpected type, typically in a type switch.
func UnexpectedType(t interface{}) error {
	return AssertionFailure("unexpected type %T", t)
}
