package addchain

import (
	"math/big"
	"testing"
)

func TestBosCosterMakeSequence(t *testing.T) {
	// Build targets list.
	nums := []int64{1, 2, 47, 117, 343, 499, 933, 5689}
	targets := []*big.Int{}
	for _, num := range nums {
		targets = append(targets, big.NewInt(num))
	}

	// According to Example 9.39 in
	// https://koclab.cs.ucsb.edu/teaching/ecc/eccPapers/Doche-ch09.pdf
	// we should get:
	//
	// 1, 2, 4, 8, 10, 11, 18, 36, 47 , 55, 91, 109, 117 , 226, 343 , 434, 489, 499 , 933 , 1422, 2844, 5688, 5689
	//
	// Our code is not the same, but the example shows it can be done with a chain of length 23.

	// Apply MakeSequence.
	seq, err := BosCosterMakeSequence(targets)
	if err != nil {
		t.Fatal(err)
	}

	for l, r := 0, len(seq)-1; l < r; l, r = l+1, r-1 {
		seq[l], seq[r] = seq[r], seq[l]
	}

	t.Logf("length=%d sequence=%v", len(seq), seq)
}
