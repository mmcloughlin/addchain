package addchain

import (
	"fmt"
	"math/big"
)

// References:
//
// https://link.springer.com/content/pdf/10.1007/0-387-34805-0_37.pdf
// http://library.isical.ac.in:8080/jspui/bitstream/123456789/6441/1/DISS-285.pdf
// https://cr.yp.to/bib/2003/stam-thesis.pdf
// https://profs.info.uaic.ro/~tr/tr03-02.pdf
// https://koclab.cs.ucsb.edu/teaching/ecc/eccPapers/Doche-ch09.pdf
// https://github.com/kwantam/addchain

// BosCosterMakeSequence builds a sequence algorithm that applies a variant of the Bos-Coster heuristics.
func BosCosterMakeSequence() SequenceAlgorithm {
	return NewHeuristicSequenceAlgorithm(
		Approximation{Ratio: 8},
		Halving{},
	)
}

// Approximation is the Bos-Coster "Approximation" heuristic.
type Approximation struct {
	Ratio int64
}

func (h Approximation) String() string {
	return fmt.Sprintf("approximation(%d)", h.Ratio)
}

// Epsilon computes the epsilon value for a given target.
func (h Approximation) Epsilon(target *big.Int) *big.Int {
	r := big.NewInt(h.Ratio)
	return new(big.Int).Div(target, r)
}

// Suggest applies the "Approximation" heuristic. This heuristic looks for two
// elements a, b in the list that sum to something close to the target element
// f. That is, we look for f-(a+b) = epsilon where a ⩽ b and epsilon is a
// "small" positive value.
func (h Approximation) Suggest(s *SequenceState) []*Proposal {
	proposals := []*Proposal{}
	delta := new(big.Int)
	f, target := s.SplitTarget()
	eps := h.Epsilon(target)
	var min *big.Int
	for i, a := range f {
		for _, b := range f[i:] {
			delta.Add(a, b)
			delta.Sub(target, delta)
			if delta.Sign() < 0 {
				continue
			}

			insert := new(big.Int).Add(a, delta)

			if delta.Cmp(eps) <= 0 {
				// If within epsilon of the target, propose this as an insertion.
				proposals = append(proposals, ProposeInsert(insert))
			} else if min == nil || insert.Cmp(min) < 0 {
				// Otherwise keep track of it if it's the smallest we've seen so far.
				min = insert
			}
		}
	}

	if len(proposals) == 0 {
		proposals = []*Proposal{ProposeInsert(min)}
	}

	return proposals
}

/*
// division applies the "Division" heuristic.
func division(_ []*big.Int, target *big.Int) []*big.Int {
	// Small primes together with minimal addition chains for them.
	primes := []struct {
		P     int64
		Chain []int64 // excluding P
	}{
		{P: 19, Chain: []int64{1, 2, 4, 8, 16, 18}},
		{P: 17, Chain: []int64{1, 2, 4, 8, 9}},
		{P: 13, Chain: []int64{1, 2, 4, 8, 9}},
		{P: 11, Chain: []int64{1, 2, 3, 5, 10}},
		{P: 7, Chain: []int64{1, 2, 3, 5}},
		{P: 5, Chain: []int64{1, 2, 3}},
		{P: 3, Chain: []int64{1, 2}},
	}

	// Check if any of them divide the target.
	m, p := new(big.Int), new(big.Int)
	for _, prime := range primes {
		p.SetInt64(prime.P)
		if m.Mod(target, p).Sign() == 0 {
			d := new(big.Int).Div(target, p)
			insertions := []*big.Int{}
			for _, c := range prime.Chain {
				insert := new(big.Int).Mul(d, big.NewInt(c))
				insertions = append(insertions, insert)
			}
			return insertions
		}
	}

	return nil
}
*/

type Halving struct{}

func (h Halving) String() string {
	return "halving"
}

func (h Halving) Suggest(s *SequenceState) []*Proposal {
	f := s.Proto
	n := len(f)

	// Check the condition f / f_1 >= 2^u
	r := new(big.Int).Div(f[n-1], f[n-2])
	if r.BitLen() < 2 {
		return nil
	}
	u := r.BitLen() - 1

	// Compute k = floor(f / 2^u)
	k := new(big.Int).Rsh(f[n-1], uint(u))

	// Proposal to insert:
	// Delta d = f - k*2^u
	// Sequence k, 2*k, ..., k*2^u
	kshifts := []*big.Int{}
	for e := 0; e <= u; e++ {
		kshift := new(big.Int).Lsh(k, uint(e))
		kshifts = append(kshifts, kshift)
	}
	d := new(big.Int).Sub(f[n-1], kshifts[u])
	insertions := append([]*big.Int{d}, kshifts...)
	proposal := ProposeInsert(insertions...)

	return []*Proposal{proposal}
}
