package alg

import (
	"fmt"
	"math/big"

	"github.com/mmcloughlin/addchain"
	"github.com/mmcloughlin/addchain/internal/bigint"
)

// ChainAlgorithm is a method of generating an addition chain for a target integer.
type ChainAlgorithm interface {
	fmt.Stringer
	FindChain(target *big.Int) (addchain.Chain, error)
}

// BinaryRightToLeft builds a chain algoirithm for the right-to-left binary method.
type BinaryRightToLeft struct{}

func (BinaryRightToLeft) String() string { return "binary_right_to_left" }

func (BinaryRightToLeft) FindChain(n *big.Int) (addchain.Chain, error) {
	c := addchain.Chain{}
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
