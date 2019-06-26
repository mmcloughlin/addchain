package addchain

import (
	"fmt"
	"math/big"

	"github.com/mmcloughlin/addchain/internal/bigint"
)

// ChainAlgorithm is a method of generating an addition chain for a target integer.
type ChainAlgorithm interface {
	fmt.Stringer
	FindChain(target *big.Int) (Program, error)
}

// NewChainAlgorithm is a convenience for building a chain algorithm from a function.
func NewChainAlgorithm(name string, search func(*big.Int) (Program, error)) ChainAlgorithm {
	return chainalgorithm{name: name, f: search}
}

type chainalgorithm struct {
	name string
	f    func(*big.Int) (Program, error)
}

func (a chainalgorithm) String() string                             { return a.name }
func (a chainalgorithm) FindChain(target *big.Int) (Program, error) { return a.f(target) }

// BinaryRightToLeft builds a chain algoirithm for the right-to-left binary method.
func BinaryRightToLeft() ChainAlgorithm {
	return NewChainAlgorithm("binary_left_to_right", func(n *big.Int) (Program, error) {
		p := Program{}
		s := -1
		d := 0
		b := new(big.Int).Set(n)
		for {
			if b.Bit(0) == 1 {
				if s >= 0 {
					s = p.Add(s, d)
				} else {
					s = d
				}
			}
			b.Rsh(b, 1)
			if bigint.IsZero(b) {
				break
			}
			d = p.Double(d)
		}
		return p, nil
	})
}
