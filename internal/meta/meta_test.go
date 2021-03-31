package meta

import (
	"testing"
)

func TestReleaseTime(t *testing.T) {
	_, err := ReleaseTime()
	if err != nil {
		t.Fatal(err)
	}
}
