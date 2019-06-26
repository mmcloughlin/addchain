package bigint

import "math/big"

// Zero returns 0.
func Zero() *big.Int {
	return big.NewInt(0)
}

// Equal returns whether x equals y.
func Equal(x, y *big.Int) bool {
	return x.Cmp(y) == 0
}

// IsZero returns true if x is zero.
func IsZero(x *big.Int) bool {
	return x.Sign() == 0
}

// IsNonZero returns true if x is non-zero.
func IsNonZero(x *big.Int) bool {
	return !IsZero(x)
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
