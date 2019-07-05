package errutil

import "fmt"

// UnexpectedType builds an error for an unexpected type, typically in a type switch.
func UnexpectedType(t interface{}) error {
	return fmt.Errorf("unexpected type %T", t)
}
