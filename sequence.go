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
