package assert

import "testing"

// NoError fails the test on a non-nil error.
func NoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

// Error fails the test on a nil error.
func Error(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error; got nil")
	}
}
