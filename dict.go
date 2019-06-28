package addchain

import (
	"fmt"
	"math/big"

	"github.com/mmcloughlin/addchain/internal/bigint"
)

// References:
//
//	[braueraddsubchains]  Martin Otto. Brauer addition-subtraction chains. PhD thesis, Universitat
//	                      Paderborn. 2001.
//	                      http://www.martin-otto.de/publications/docs/2001_MartinOtto_Diplom_BrauerAddition-SubtractionChains.pdf
//	[hehcc:exp]           Christophe Doche. Exponentiation. Handbook of Elliptic and Hyperelliptic Curve
//	                      Cryptography, chapter 9. 2006.
//	                      https://koclab.cs.ucsb.edu/teaching/ecc/eccPapers/Doche-ch09.pdf

// DictTerm represents the integer D * 2ᴱ.
type DictTerm struct {
	D uint64
	E uint
}

// Int converts the term to an integer.
func (t DictTerm) Int() *big.Int {
	x := big.NewInt(int64(t.D))
	x.Lsh(x, t.E)
	return x
}

// DictSum is the representation of an integer as a sum of dictionary terms.
// See [hehcc:exp] definition 9.34.
type DictSum []DictTerm

// Int computes the dictionary sum as an integer.
func (s DictSum) Int() *big.Int {
	x := bigint.Zero()
	for _, t := range s {
		x.Add(x, t.Int())
	}
	return x
}

// Decomposer is a method of breaking an integer into a dictionary sum.
type Decomposer interface {
	fmt.Stringer
	Decompose(x *big.Int) DictSum
}

// FixedWindow breaks integers into k-bit windows.
type FixedWindow struct {
	K uint
}

func (w FixedWindow) String() string { return fmt.Sprintf("fixed_window(%d)", w.K) }

// Decompose represents x in base 2ᵏ.
func (w FixedWindow) Decompose(x *big.Int) DictSum {
	sum := DictSum{}
	mask := bigint.Pow2(w.K)
	mask.Sub(mask, bigint.One())
	b := bigint.Clone(x)
	s := uint(0)
	for bigint.IsNonZero(b) {
		d := new(big.Int).And(b, mask)
		t := DictTerm{
			D: d.Uint64(),
			E: s,
		}
		sum = append(sum, t)
		b.Rsh(b, w.K)
		s += w.K
	}
	return sum
}

// SlidingWindow breaks integers into k-bit windows, skipping runs of zeros
// where possible. See [hehcc:exp] section 9.1.3 or [braueraddsubchains] section
// 1.2.3.
type SlidingWindow struct {
	K uint
}

func (w SlidingWindow) String() string { return fmt.Sprintf("sliding_window(%d)", w.K) }

// Decompose represents x in base 2ᵏ.
func (w SlidingWindow) Decompose(x *big.Int) DictSum {
	sum := DictSum{}
	mask := bigint.Pow2(w.K)
	mask.Sub(mask, bigint.One())
	b := bigint.Clone(x)
	s := uint(0)
	for bigint.IsNonZero(b) {
		// Find the next 1 bit.
		for b.Bit(0) == 0 {
			b.Rsh(b, 1)
			s++
		}

		// Extract the k-bit window here.
		d := new(big.Int).And(b, mask)
		t := DictTerm{
			D: d.Uint64(),
			E: s,
		}
		sum = append(sum, t)
		b.Rsh(b, w.K)
		s += w.K
	}
	return sum
}
