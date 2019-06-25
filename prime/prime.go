// Package prime provides representations of classes of prime numbers.
package prime

import (
	"fmt"
	"math/big"
)

// Crandall is a prime of the form 2ⁿ - c.
type Crandall struct {
	N int
	C int
}

// NewCrandall constructs a Crandall prime.
func NewCrandall(n, c int) Crandall {
	return Crandall{N: n, C: c}
}

func (p Crandall) String() string {
	return fmt.Sprintf("2^%d%+d", p.N, -p.C)
}

// Int returns the prime as a big integer.
func (p Crandall) Int() *big.Int {
	one := big.NewInt(1)
	c := big.NewInt(int64(p.C))
	e := new(big.Int).Lsh(one, uint(p.N))
	return new(big.Int).Sub(e, c)
}
