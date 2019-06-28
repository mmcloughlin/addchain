package addchain

import (
	"math/big"
	"testing"

	"github.com/mmcloughlin/addchain/internal/bigints"
)

// References:
//
//	[hehcc:exp]         Christophe Doche. Exponentiation. Handbook of Elliptic and Hyperelliptic Curve
//	                    Cryptography, chapter 9. 2006.
//	                    https://koclab.cs.ucsb.edu/teaching/ecc/eccPapers/Doche-ch09.pdf
//	[pairingsfinalexp]  Michael Scott, Naomi Benger, Manuel Charlemagne, Luis J. Dominguez Perez and
//	                    Ezekiel J. Kachisa. On the final exponentiation for calculating pairings on
//	                    ordinary elliptic curves. Cryptology ePrint Archive, Report 2008/490. 2008.
//	                    https://eprint.iacr.org/2008/490

func TestSequenceAlgorithms(t *testing.T) {
	as := []SequenceAlgorithm{
		// Continued fractions algorithms.
		NewContinuedFractions(BinaryStrategy{}),
		NewContinuedFractions(BinaryStrategy{Parity: 1}),
		NewContinuedFractions(DichotomicStrategy{}),

		// Heuristics algorithms.
		NewHeuristicAlgorithm(UseFirstHeuristic{
			Halving{},
			DeltaLargest{},
		}),
	}
	for _, a := range as {
		t.Run(a.String(), SequenceAlgorithmSuite(a))
	}
}

func SequenceAlgorithmSuite(a SequenceAlgorithm) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("as_chain_algorithm", ChainAlgorithmSuite(AsChainAlgorithm{a}))
		t.Run("known_sequences", CheckKnownSequences(a))
	}
}

func CheckKnownSequences(a SequenceAlgorithm) func(t *testing.T) {
	cases := []struct {
		Targets  []*big.Int
		Solution Chain
	}{
		// Example 9.39 in [hehcc:exp].
		{
			Targets:  bigints.Int64s(47, 117, 343, 499, 933, 5689),
			Solution: Int64s(1, 2, 4, 8, 10, 11, 18, 36, 47, 55, 91, 109, 117, 226, 343, 434, 489, 499, 933, 1422, 2844, 5688, 5689),
		},
		// [pairingsfinalexp] page 5.
		{
			Targets:  bigints.Int64s(6, 12, 18, 30, 36),
			Solution: Int64s(1, 2, 3, 6, 12, 18, 30, 36),
		},
		// [pairingsfinalexp] page 8.
		{
			Targets:  bigints.Int64s(4, 5, 7, 10, 15, 26, 28, 30, 55, 75, 80, 100, 108, 144),
			Solution: Int64s(1, 2, 4, 5, 7, 10, 15, 25, 26, 28, 30, 36, 50, 55, 75, 80, 100, 108, 144),
		},
		// [pairingsfinalexp] page 9.
		{
			Targets:  bigints.Int64s(3, 5, 7, 14, 15, 21, 25, 35, 49, 54, 62, 70, 87, 98, 112, 245, 273, 319, 343, 434, 450, 581, 609, 784, 931, 1407, 1911, 4802, 6517),
			Solution: Int64s(1, 2, 3, 4, 5, 7, 8, 14, 15, 16, 21, 25, 28, 35, 42, 49, 54, 62, 70, 87, 98, 112, 147, 245, 273, 294, 319, 343, 392, 434, 450, 581, 609, 784, 931, 1162, 1407, 1862, 1911, 3724, 4655, 4802, 6517),
		},
	}

	for _, c := range cases {
		if err := c.Solution.Superset(c.Targets); err != nil {
			panic(err)
		}
	}

	return func(t *testing.T) {
		for _, c := range cases {
			got := AssertSequenceAlgorithmProduces(t, a, c.Targets)
			t.Logf("     got: length=%d sequence=%v", len(got), got)
			t.Logf("solution: length=%d sequence=%v", len(c.Solution), c.Solution)
			if len(got) > len(c.Solution) {
				t.Logf("sub-optimal")
			}
		}
	}
}
