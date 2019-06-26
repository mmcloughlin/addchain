package addchain

import (
	"testing"

	"github.com/mmcloughlin/addchain/internal/bigints"
)

func TestBosCosterMakeSequence(t *testing.T) {
	a := BosCosterMakeSequence()
	t.Log(a)

	// According to Example 9.39 in
	// https://koclab.cs.ucsb.edu/teaching/ecc/eccPapers/Doche-ch09.pdf
	// we should get:
	//
	// 1, 2, 4, 8, 10, 11, 18, 36, 47 , 55, 91, 109, 117 , 226, 343 , 434, 489, 499 , 933 , 1422, 2844, 5688, 5689
	//
	// Our code is not the same, but the example shows it can be done with a chain of length 23.
	targets := bigints.Int64s([]int64{1, 2, 47, 117, 343, 499, 933, 5689})

	// Apply MakeSequence.
	c, err := a.FindSequence(targets)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("length=%d sequence=%v", len(c), c)
}
