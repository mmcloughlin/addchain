package errutil

import "golang.org/x/xerrors"

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
