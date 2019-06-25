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

	for _, x := range seq {
		t.Log(x)
	}
}
