package addchain

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/mmcloughlin/addchain/internal/bigint"
	"github.com/mmcloughlin/addchain/internal/bigints"
)

// Chain is an addition chain.
type Chain []*big.Int

// New constructs the minimal chain {1}.
func New() Chain {
	return Chain{big.NewInt(1)}
}

// Int64s builds a chain from the given int64 values.
func Int64s(xs ...int64) Chain {
	return Chain(bigints.Int64s(xs...))
}

// Clone the chain.
func (c Chain) Clone() Chain {
	return bigints.Clone(c)
}

// AppendClone appends a copy of x to c.
func (c *Chain) AppendClone(x *big.Int) {
	*c = append(*c, bigint.Clone(x))
}

// End returns the last element of the chain.
func (c Chain) End() *big.Int {
	return c[len(c)-1]
}

// Op returns an Op that produces the kth position.
func (c Chain) Op(k int) (Op, error) {
	s := new(big.Int)
	for i := 0; i < k; i++ {
		for j := i; j < k; j++ {
			s.Add(c[i], c[j])
			if s.Cmp(c[k]) == 0 {
				return Op{i, j}, nil
			}
		}
	}
	return Op{}, fmt.Errorf("position %d is not the sum of previous entries", k)
}

// Program produces a program that generates the chain.
func (c Chain) Program() (Program, error) {
	if len(c) == 0 {
		return nil, errors.New("chain empty")
	}

	if c[0].Cmp(big.NewInt(1)) != 0 {
		return nil, errors.New("chain must start with 1")
	}

	p := Program{}
	for k := 1; k < len(c); k++ {
		op, err := c.Op(k)
		if err != nil {
			return nil, err
		}
		p = append(p, op)
	}

	return p, nil
}

// Validate checks that c is in fact an addition chain.
func (c Chain) Validate() error {
	_, err := c.Program()
	return err
}

// Produces checks that c is a valid chain ending with target.
func (c Chain) Produces(target *big.Int) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if c.End().Cmp(target) != 0 {
		return errors.New("chain does not end with target")
	}
	return nil
}

// Product computes the product of two addition chains. The is the "o times"
// operator defined in "Efficient computation of addition chains" by F.
// Bergeron, J. Berstel and S. Brlek.
func Product(a, b Chain) Chain {
	c := a.Clone()
	last := c.End()
	for _, x := range b[1:] {
		y := new(big.Int).Mul(last, x)
		c = append(c, y)
	}
	return c
}

// Plus adds x to the addition chain. This is the "o plus" operator defined in
// "Efficient computation of addition chains" by F. Bergeron, J. Berstel and S.
// Brlek.
func Plus(a Chain, x *big.Int) Chain {
	c := a.Clone()
	y := new(big.Int).Add(c.End(), x)
	return append(c, y)
}
