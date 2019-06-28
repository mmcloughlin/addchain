package addchain

import "testing"

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
	}
}
