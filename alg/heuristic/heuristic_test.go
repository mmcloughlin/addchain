package heuristic

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/mmcloughlin/addchain/internal/bigints"
)

// References:
//
//	[hehcc:exp]  Christophe Doche. Exponentiation. Handbook of Elliptic and Hyperelliptic Curve
//	             Cryptography, chapter 9. 2006.
//	             https://koclab.cs.ucsb.edu/teaching/ecc/eccPapers/Doche-ch09.pdf

func TestHalving(t *testing.T) {
	cases := []struct {
		F      []*big.Int
		Target *big.Int
		Expect []*big.Int
	}{
		// Example from [hehcc:exp], page 163.
		{
			F:      bigints.Int64s(14),
			Target: big.NewInt(382),
			Expect: bigints.Int64s(14, 23, 46, 92, 184, 368),
		},
		// Simple powers of two case.
		{
			F:      bigints.Int64s(1, 2),
			Target: big.NewInt(8),
			Expect: bigints.Int64s(2, 4),
		},
	}
	h := Halving{}
	for _, c := range cases {
		if got := h.Suggest(c.F, c.Target); !reflect.DeepEqual(c.Expect, got) {
			t.Errorf("Suggest(%v, %v) = %v; expect %v", c.F, c.Target, got, c.Expect)
		}
	}
}
