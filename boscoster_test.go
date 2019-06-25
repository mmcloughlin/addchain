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
