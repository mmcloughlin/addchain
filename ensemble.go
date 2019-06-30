package addchain

// Ensemble is a convenience for building an ensemble of chain algorithms intended for large integers.
func Ensemble() []ChainAlgorithm {
	seqalgs := []SequenceAlgorithm{
		NewContinuedFractions(BinaryStrategy{}),
		NewContinuedFractions(BinaryStrategy{Parity: 1}),
		NewContinuedFractions(DichotomicStrategy{}),
	}

	as := []ChainAlgorithm{}
	for k := uint(4); k <= 128; k *= 2 {
		for _, seqalg := range seqalgs {
			a := NewDictAlgorithm(
				SlidingWindow{K: k},
				seqalg,
			)
			as = append(as, a)
		}
	}

	return as
}
