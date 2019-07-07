package assert

import "testing"

// NoError fails the test on a non-nil error.
func NoError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
