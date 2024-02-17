package dict

import (
	"testing"

	"github.com/mmcloughlin/addchain/alg"
	"github.com/mmcloughlin/addchain/alg/algtest"
	"github.com/mmcloughlin/addchain/alg/contfrac"
)

func TestAlgorithms(t *testing.T) {
	as := []alg.ChainAlgorithm{
		// Dictionary-based algorithms.
		NewAlgorithm(
			FixedWindow{K: 7},
			contfrac.NewAlgorithm(contfrac.BinaryStrategy{}),
		),
		NewAlgorithm(
			SlidingWindow{K: 4},
			contfrac.NewAlgorithm(contfrac.DichotomicStrategy{}),
		),
		NewAlgorithm(
			SlidingWindowShortRTL{K: 4, Z: 2},
			contfrac.NewAlgorithm(contfrac.CoBinaryStrategy{}),
		),
		NewAlgorithm(
			Hybrid{TMin: 6, TMax: 20, Decomposer: SlidingWindowShortRTL{K: 5, Z: 2}},
			contfrac.NewAlgorithm(contfrac.CoBinaryStrategy{}),
		),

		// Runs algorithm.
		NewRunsAlgorithm(contfrac.NewAlgorithm(contfrac.DichotomicStrategy{})),
	}
	for _, a := range as {
		t.Run(a.String(), algtest.ChainAlgorithmSuite(a))
	}
}
