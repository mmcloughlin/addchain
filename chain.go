package addchain

import (
	"errors"
	"fmt"
	"math/big"
)

// Chain is an addition chain.
type Chain []*big.Int

// New constructs the minimal chain {1}.
func New() Chain {
	return Chain{big.NewInt(1)}
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
