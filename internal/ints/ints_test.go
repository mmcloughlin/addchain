package ints

import "testing"

func TestNextMultiple(t *testing.T) {
	for a := 65; a <= 128; a++ {
		if NextMultiple(a, 64) != 128 {
			t.FailNow()
		}
	}
}
