package addchain

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/mmcloughlin/addchain/internal/ints"

	"github.com/mmcloughlin/addchain/internal/bigint"
	"github.com/mmcloughlin/addchain/internal/bigints"
)

// References:
//
//	[braueraddsubchains]  Martin Otto. Brauer addition-subtraction chains. PhD thesis, Universitat
//	                      Paderborn. 2001.
//	                      http://www.martin-otto.de/publications/docs/2001_MartinOtto_Diplom_BrauerAddition-SubtractionChains.pdf
//	[genshortchains]      Kunihiro, Noboru and Yamamoto, Hirosuke. New Methods for Generating Short
//	                      Addition Chains. IEICE Transactions on Fundamentals of Electronics
//	                      Communications and Computer Sciences. 2000.
//	                      https://pdfs.semanticscholar.org/b398/d10faca35af9ce5a6026458b251fd0a5640c.pdf
//	[hehcc:exp]           Christophe Doche. Exponentiation. Handbook of Elliptic and Hyperelliptic Curve
//	                      Cryptography, chapter 9. 2006.
//	                      https://koclab.cs.ucsb.edu/teaching/ecc/eccPapers/Doche-ch09.pdf

// DictTerm represents the integer D * 2ᴱ.
type DictTerm struct {
	D *big.Int
	E uint
}

// Int converts the term to an integer.
func (t DictTerm) Int() *big.Int {
	return new(big.Int).Lsh(t.D, t.E)
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

// SortByExponent sorts terms in ascending order of the exponent E.
func (s DictSum) SortByExponent() {
	sort.Slice(s, func(i, j int) bool { return s[i].E < s[j].E })
}

// Dictionary returns the distinct D values in the terms of this sum. The values
// are returned in ascending order.
func (s DictSum) Dictionary() []*big.Int {
	dict := make([]*big.Int, 0, len(s))
	for _, t := range s {
		dict = append(dict, t.D)
	}
	bigints.Sort(dict)
	return bigints.Unique(dict)
}

// Decomposer is a method of breaking an integer into a dictionary sum.
type Decomposer interface {
	fmt.Stringer
	Decompose(x *big.Int) DictSum
}

// FixedWindow breaks integers into k-bit windows.
type FixedWindow struct {
	K uint // Window size.
}

func (w FixedWindow) String() string { return fmt.Sprintf("fixed_window(%d)", w.K) }

// Decompose represents x in terms of k-bit windows from left to right.
func (w FixedWindow) Decompose(x *big.Int) DictSum {
	sum := DictSum{}
	h := x.BitLen()
	for h > 0 {
		l := ints.Max(h-int(w.K), 0)
		d := bigint.Extract(x, uint(l), uint(h))
		if bigint.IsNonZero(d) {
			sum = append(sum, DictTerm{D: d, E: uint(l)})
		}
		h = l
	}
	sum.SortByExponent()
	return sum
}

// SlidingWindow breaks integers into k-bit windows, skipping runs of zeros
// where possible. See [hehcc:exp] section 9.1.3 or [braueraddsubchains] section
// 1.2.3.
type SlidingWindow struct {
	K uint // Window size.
}

func (w SlidingWindow) String() string { return fmt.Sprintf("sliding_window(%d)", w.K) }

// Decompose represents x in base 2ᵏ.
func (w SlidingWindow) Decompose(x *big.Int) DictSum {
	sum := DictSum{}
	h := x.BitLen() - 1
	for h >= 0 {
		// Find first 1.
		for h >= 0 && x.Bit(h) == 0 {
			h--
		}

		if h < 0 {
			break
		}

		// Look down k positions.
		l := ints.Max(h-int(w.K)+1, 0)

		// Advance to the next 1.
		for x.Bit(l) == 0 {
			l++
		}

		sum = append(sum, DictTerm{
			D: bigint.Extract(x, uint(l), uint(h+1)),
			E: uint(l),
		})

		h = l - 1
	}
	sum.SortByExponent()
	return sum
}

// RunLength decomposes integers in to runs of 1s up to a maximal length. See
// [genshortchains] Section 3.1.
type RunLength struct {
	T uint // Maximal run length. Zero means no limit.
}

func (r RunLength) String() string { return fmt.Sprintf("run_length(%d)", r.T) }

// Decompose breaks x into runs of 1 bits.
func (r RunLength) Decompose(x *big.Int) DictSum {
	sum := DictSum{}
	i := x.BitLen() - 1
	for i >= 0 {
		// Find first 1.
		for i >= 0 && x.Bit(i) == 0 {
			i--
		}

		if i < 0 {
			break
		}

		// Look for the end of the run.
		s := i
		for i >= 0 && x.Bit(i) == 1 && (r.T == 0 || uint(s-i) < r.T) {
			i--
		}

		// We have a run from s to i+1.
		sum = append(sum, DictTerm{
			D: bigint.Ones(uint(s - i)),
			E: uint(i + 1),
		})
	}
	sum.SortByExponent()
	return sum
}

// Hybrid is a mix of the sliding window and run length decomposition methods,
// similar to the "Hybrid Method" of [genshortchains] Section 3.3.
type Hybrid struct {
	K uint // Window size.
	T uint // Maximal run length. Zero means no limit.
}

func (h Hybrid) String() string { return fmt.Sprintf("hybrid(%d,%d)", h.K, h.T) }

// Decompose breaks x into k-bit sliding windows or runs of 1s up to length T.
func (h Hybrid) Decompose(x *big.Int) DictSum {
	sum := DictSum{}

	// Clone since we'll be modifying it.
	y := bigint.Clone(x)

	// Process runs of length at least K.
	i := y.BitLen() - 1
	for i >= 0 {
		// Find first 1.
		for i >= 0 && y.Bit(i) == 0 {
			i--
		}

		if i < 0 {
			break
		}

		// Look for the end of the run.
		s := i
		for i >= 0 && y.Bit(i) == 1 && (h.T == 0 || uint(s-i) < h.T) {
			i--
		}

		// We have a run from s to i+1. Skip it if its short.
		n := uint(s - i)
		if n <= h.K {
			continue
		}

		// Add it to the sum and remove it from the integer.
		sum = append(sum, DictTerm{
			D: bigint.Ones(n),
			E: uint(i + 1),
		})

		y.Xor(y, bigint.Mask(uint(i+1), uint(s+1)))
	}

	// Process what remains with a sliding window.
	w := SlidingWindow{K: h.K}
	rem := w.Decompose(y)

	sum = append(sum, rem...)
	sum.SortByExponent()

	return sum
}

// DictAlgorithm implements a general dictionary-based chain construction
// algorithm, as in [braueraddsubchains] Algorithm 1.26. This operates in three
// stages: decompose the target into a sum of dictionray terms, use a sequence
// algorithm to generate the dictionary, then construct the target from the
// dictionary terms.
type DictAlgorithm struct {
	decomp Decomposer
	seqalg SequenceAlgorithm
}

// NewDictAlgorithm builds a dictionary algorithm that breaks up integers using
// the decomposer d and uses the sequence algorithm s to generate dictionary entries.
func NewDictAlgorithm(d Decomposer, a SequenceAlgorithm) *DictAlgorithm {
	return &DictAlgorithm{
		decomp: d,
		seqalg: a,
	}
}

func (a DictAlgorithm) String() string {
	return fmt.Sprintf("dictionary(%s,%s)", a.decomp, a.seqalg)
}

// FindChain builds an addition chain producing n. This works by using the
// configured Decomposer to represent n as a sum of dictionary terms, then
// delegating to the SequenceAlgorithm to build a chain producing the
// dictionary, and finally using the dictionary terms to construct n. See
// [genshortchains] Section 2 for a full description.
func (a DictAlgorithm) FindChain(n *big.Int) (Chain, error) {
	// Decompose the target.
	sum := a.decomp.Decompose(n)
	sum.SortByExponent()

	// Extract dictionary.
	dict := sum.Dictionary()

	// Use the sequence algorithm to produce a chain for each element of the dictionary.
	c, err := a.seqalg.FindSequence(dict)
	if err != nil {
		return nil, err
	}

	// Build chain for n out of the dictionary.
	k := len(sum) - 1
	cur := bigint.Clone(sum[k].D)
	for ; k > 0; k-- {
		// Shift until the next exponent.
		for i := sum[k].E; i > sum[k-1].E; i-- {
			cur.Lsh(cur, 1)
			c.AppendClone(cur)
		}

		// Add in the dictionary term at this position.
		cur.Add(cur, sum[k-1].D)
		c.AppendClone(cur)
	}

	for i := sum[0].E; i > 0; i-- {
		cur.Lsh(cur, 1)
		c.AppendClone(cur)
	}

	// Prepare chain for returning.
	bigints.Sort(c)
	c = Chain(bigints.Unique(c))

	return c, nil
}
