package addchain

import (
	"fmt"
	"math/big"
	"sort"

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

// SortByExponent sorts terms in ascending order of the exponent E.
func (s DictSum) SortByExponent() {
	sort.Slice(s, func(i, j int) bool { return s[i].E < s[j].E })
}

// Dictionary returns the distinct D values in the terms of this sum. The values
// are returned in ascending order.
func (s DictSum) Dictionary() []uint64 {
	set := map[uint64]bool{}
	for _, t := range s {
		set[t.D] = true
	}

	dict := make([]uint64, 0, len(set))
	for d := range set {
		dict = append(dict, d)
	}

	sort.Slice(dict, func(i, j int) bool { return dict[i] < dict[j] })

	return dict
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
	targets := []*big.Int{}
	for _, d := range dict {
		targets = append(targets, big.NewInt(int64(d)))
	}

	c, err := a.seqalg.FindSequence(targets)
	if err != nil {
		return nil, err
	}

	// Build chain for n out of the dictionary.
	k := len(sum) - 1
	cur := big.NewInt(int64(sum[k].D))
	for ; k > 0; k-- {
		// Shift until the next exponent.
		for i := sum[k].E; i > sum[k-1].E; i-- {
			cur.Lsh(cur, 1)
			c.AppendClone(cur)
		}

		// Add in the dictionary term at this position.
		cur.Add(cur, big.NewInt(int64(sum[k-1].D)))
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
