package addchain

import (
	"fmt"
	"math/big"

	"github.com/mmcloughlin/addchain/internal/bigint"
)

// ChainAlgorithm is a method of generating an addition chain for a target integer.
type ChainAlgorithm interface {
	fmt.Stringer
	FindChain(target *big.Int) (Chain, error)
}

// BinaryRightToLeft builds a chain algoirithm for the right-to-left binary method.
type BinaryRightToLeft struct{}

func (BinaryRightToLeft) String() string { return "binary_right_to_left" }

func (BinaryRightToLeft) FindChain(n *big.Int) (Chain, error) {
	c := Chain{}
	b := new(big.Int).Set(n)
	d := bigint.One()
	var x *big.Int
	for bigint.IsNonZero(b) {
		c.AppendClone(d)
		if b.Bit(0) == 1 {
			if x == nil {
				x = bigint.Clone(d)
			} else {
				x.Add(x, d)
				c.AppendClone(x)
			}
		}
		b.Rsh(b, 1)
		d.Lsh(d, 1)
	}
	return c, nil
}

// Ensemble is a convenience for building an ensemble of chain algorithms intended for large integers.
func Ensemble() []ChainAlgorithm {
	seqalgs := []SequenceAlgorithm{
		NewContinuedFractions(BinaryStrategy{}),
		NewContinuedFractions(BinaryStrategy{Parity: 1}),
		NewContinuedFractions(DichotomicStrategy{}),
	}

	as := []ChainAlgorithm{}
	for k := 4; k <= 32; k *= 2 {
		for _, seqalg := range seqalgs {
			a := NewDictAlgorithm(
				SlidingWindow{K: uint(k)},
				seqalg,
			)
			as = append(as, a)
		}
	}

	return as
}
