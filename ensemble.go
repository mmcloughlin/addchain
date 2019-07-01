package addchain

// Ensemble is a convenience for building an ensemble of chain algorithms intended for large integers.
func Ensemble() []ChainAlgorithm {
	seqalgs := []SequenceAlgorithm{
		NewContinuedFractions(BinaryStrategy{}),
		NewContinuedFractions(BinaryStrategy{Parity: 1}),
		NewContinuedFractions(DichotomicStrategy{}),

		NewHeuristicAlgorithm(UseFirstHeuristic{
			Halving{},
			DeltaLargest{},
		}),
		NewHeuristicAlgorithm(UseFirstHeuristic{
			Halving{},
			Approximation{},
		}),
	}

	// Build decomposers.
	decomposers := []Decomposer{}
	for k := uint(4); k <= 128; k *= 2 {
		decomposers = append(decomposers, SlidingWindow{K: k})
	}

	decomposers = append(decomposers, RunLength{T: 0})
	for t := uint(16); t <= 128; t *= 2 {
		decomposers = append(decomposers, RunLength{T: t})
	}

	for k := uint(2); k <= 8; k++ {
		decomposers = append(decomposers, Hybrid{K: k, T: 0})
		for t := uint(16); t <= 64; t *= 2 {
			decomposers = append(decomposers, Hybrid{K: k, T: t})
		}
	}

	// Build dictionary algorithms for every combination.
	as := []ChainAlgorithm{}
	for _, decomp := range decomposers {
		for _, seqalg := range seqalgs {
			a := NewDictAlgorithm(decomp, seqalg)
			as = append(as, a)
		}
	}

	// Wrap in an optimization layer.
	for i, a := range as {
		as[i] = Optimized{Algorithm: a}
	}

	return as
}
