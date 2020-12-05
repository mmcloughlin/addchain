// Package prime provides representations of classes of prime numbers.
package prime

import (
	"fmt"
	"math/big"

	"github.com/mmcloughlin/addchain/internal/bigint"
	"github.com/mmcloughlin/addchain/internal/polynomial"
)

// References:
//
//	[crandallprime]  Richard E. Crandall. Method and apparatus for public key exchange in a
//	                 cryptographic system. US Patent 5,159,632. 1992.
//	                 https://patents.google.com/patent/US5159632A
//	[isogenychains]  Koziel, Brian, Azarderakhsh, Reza, Jao, David and Mozaffari-Kermani, Mehran. On
//	                 Fast Calculation of Addition Chains for Isogeny-Based Cryptography. In
//	                 Information Security and Cryptology, pages 323--342. 2016.
//	                 http://faculty.eng.fau.edu/azarderakhsh/files/2016/11/Inscrypt2016.pdf
//	[solinasprime]   Jerome A. Solinas. Generalized Mersenne Primes. Technical Report CORR 99-39,
//	                 Centre for Applied Cryptographic Research (CACR) at the University of Waterloo.
//	                 1999. http://cacr.uwaterloo.ca/techreports/1999/corr99-39.pdf

// Prime is the interface for a prime number.
type Prime interface {
	Bits() int
	Int() *big.Int
	String() string
}

// Crandall represents a prime of the form 2ⁿ - c. Named after Richard E. Crandall [crandallprime].
type Crandall struct {
	N int
	C int
}

// NewCrandall constructs a Crandall prime.
func NewCrandall(n, c int) Crandall {
	return Crandall{N: n, C: c}
}

// Bits returns the number of bits required to represent p.
func (p Crandall) Bits() int {
	return p.N
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

// Solinas is a "Generalized Mersenne Prime", as introduced by Jerome Solinas
// [solinasprime]. Such primes are of the form f( 2ᵏ ) for a low-degree
// polynomial f.
type Solinas struct {
	F polynomial.Polynomial
	K uint
}

// NewSolinas constructs a Solinas prime.
func NewSolinas(f polynomial.Polynomial, k uint) Solinas {
	return Solinas{F: f, K: k}
}

// Bits returns the number of bits required to represent p.
func (p Solinas) Bits() int {
	return int(p.F.Degree() * p.K)
}

// Int returns p as an integer.
func (p Solinas) Int() *big.Int {
	return p.F.Evaluate(bigint.Pow2(p.K))
}

func (p Solinas) String() string {
	// Create another polynomial with all terms scaled by k.
	g := polynomial.Polynomial{}
	for _, t := range p.F {
		s := t
		s.N *= p.K
		g = append(g, s)
	}
	return g.Format("2")
}

// SmoothIsogeny is a prime of the form ℓ₁ᵃ * ℓ₂ᵇ * f ± 1 for small primes
// ℓ₁ and ℓ₂ and small cofactor f, as introduced in [isogenychains].
// These primes are useful for isogeny-based post-quantum cryptography.
type SmoothIsogeny struct {
	L1  uint
	L2  uint
	A   uint
	B   uint
	F   uint
	Add bool // True for +1, false for -1.
	p   *big.Int
}

// NewSmoothIsogeny constructs a smooth isogeny prime.
func NewSmoothIsogeny(l1 uint, l2 uint, a uint, b uint, f uint, add bool) SmoothIsogeny {
	t1 := new(big.Int).Exp(big.NewInt(int64(l1)), big.NewInt(int64(a)), nil)
	t2 := new(big.Int).Exp(big.NewInt(int64(l2)), big.NewInt(int64(b)), nil)
	p := big.NewInt(int64(f))
	p.Mul(p, t1)
	p.Mul(p, t2)
	if add {
		p.Add(p, bigint.One())
	} else {
		p.Sub(p, bigint.One())
	}
	return SmoothIsogeny{L1: l1, L2: l2, A: a, B: b, F: f, Add: add, p: p}
}

// Bits returns the bit length of p.
func (p SmoothIsogeny) Bits() int { return p.p.BitLen() }

// Int returns p as an integer.
func (p SmoothIsogeny) Int() *big.Int { return p.p }

func (p SmoothIsogeny) String() string {
	s := "-1"
	if p.Add {
		s = "+1"
	}
	return fmt.Sprintf("%d^%d*%d^%d*%d%s", p.L1, p.A, p.L2, p.B, p.F, s)
}

// Other is a prime whose structure does not match any of the other specific
// types in this package.
type Other struct {
	P *big.Int
}

// NewOther builds a prime from the provided integer.
func NewOther(p *big.Int) Other {
	return Other{P: p}
}

// MustHex parses a prime from the hex literal p. Panics on error.
func MustHex(p string) Other {
	return NewOther(bigint.MustHex(p))
}

// Bits returns the bit length of p.
func (p Other) Bits() int { return p.P.BitLen() }

// Int returns p as an integer.
func (p Other) Int() *big.Int { return p.P }

func (p Other) String() string { return fmt.Sprintf("%x", p.P) }
