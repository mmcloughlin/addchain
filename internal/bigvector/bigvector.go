// Package bigvector implements operations on vectors of multi-precision integers.
package bigvector

import (
	"math/big"

	"github.com/mmcloughlin/addchain/internal/bigint"
)

var (
	zero = bigint.Zero()
	one  = bigint.One()
)

// Vector is a vector of big integers.
type Vector []*big.Int

// New constructs an n-dimensional zero vector.
func New(n int) Vector {
	v := make(Vector, n)
	for i := 0; i < n; i++ {
		v[i] = zero
	}
	return v
}

// NewBasis constructs an n-dimensional basis vector with a 1 in position i.
func NewBasis(n, i int) Vector {
	v := New(n)
	v[i] = one
	return v
}

// Add vectors.
func Add(u, v Vector) Vector {
	assertsamelen(u, v)
	n := len(u)
	w := make(Vector, n)
	for i := 0; i < n; i++ {
		w[i] = new(big.Int).Add(u[i], v[i])
	}
	return w
}

// Lsh left shifts every element of the vector v.
func Lsh(v Vector, s uint) Vector {
	n := len(v)
	w := make(Vector, n)
	for i := 0; i < n; i++ {
		w[i] = new(big.Int).Lsh(v[i], s)
	}
	return w
}

// assertsamelen panics if u and v are different lengths.
func assertsamelen(u, v Vector) {
	if len(u) != len(v) {
		panic("bigvector: length mismatch")
	}
}
