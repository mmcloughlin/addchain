package bigint

import "math/big"

// Zero returns 0.
func Zero() *big.Int {
	return big.NewInt(0)
}

// One returns 1.
func One() *big.Int {
	return big.NewInt(1)
}

// Equal returns whether x equals y.
func Equal(x, y *big.Int) bool {
	return x.Cmp(y) == 0
}

// EqualInt64 is a convenience for checking if x equals the int64 value y.
func EqualInt64(x *big.Int, y int64) bool {
	return Equal(x, big.NewInt(y))
}

// IsZero returns true if x is zero.
func IsZero(x *big.Int) bool {
	return x.Sign() == 0
}

// IsNonZero returns true if x is non-zero.
func IsNonZero(x *big.Int) bool {
	return !IsZero(x)
}

// Clone returns a copy of x.
func Clone(x *big.Int) *big.Int {
	return new(big.Int).Set(x)
}

// Pow2 returns 2ᵉ.
func Pow2(e uint) *big.Int {
	return new(big.Int).Lsh(One(), e)
}

// IsPow2 returns whether x is a power of 2.
func IsPow2(x *big.Int) bool {
	e := x.BitLen()
	if e == 0 {
		return false
	}
	return Equal(x, Pow2(uint(e-1)))
}

// Pow2UpTo returns all powers of two ⩽ x.
func Pow2UpTo(x *big.Int) []*big.Int {
	p := One()
	ps := []*big.Int{}
	for p.Cmp(x) <= 0 {
		ps = append(ps, Clone(p))
		p.Lsh(p, 1)
	}
	return ps
}

// TrailingZeros returns the number of trailing zero bits in x. Returns 0 if x is 0.
func TrailingZeros(x *big.Int) int {
	if x.BitLen() == 0 {
		return 0
	}
	n := 0
	for ; x.Bit(n) == 0; n++ {
	}
	return n
}
