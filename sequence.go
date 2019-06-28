package addchain

import (
	"fmt"
	"math/big"
)

// SequenceAlgorithm is a method of generating an addition sequence for a set of
// target values.
type SequenceAlgorithm interface {
	fmt.Stringer
	FindSequence(targets []*big.Int) (Chain, error)
}

// AsChainAlgorithm adapts a sequence algorithm to a chain algorithm.
type AsChainAlgorithm struct {
	SequenceAlgorithm
}

// FindChain calls FindSequence with a singleton list containing the target.
func (a AsChainAlgorithm) FindChain(target *big.Int) (Chain, error) {
	return a.FindSequence([]*big.Int{target})
}
