package addchain

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/mmcloughlin/addchain/internal/bigints"
)

// References:
//
//	[boscoster]                Bos, Jurjen and Coster, Matthijs. Addition Chain Heuristics. In Advances in
//	                           Cryptology --- CRYPTO' 89 Proceedings, pages 400--407. 1990.
//	[github:kwantam/addchain]  Riad S. Wahby. kwantam/addchain. Github Repository. Apache License, Version 2.0.
//	                           2018. https://github.com/kwantam/addchain
//	[hehcc:exp]                Christophe Doche. Exponentiation. Handbook of Elliptic and Hyperelliptic Curve
//	                           Cryptography, chapter 9. 2006.
//	                           https://koclab.cs.ucsb.edu/teaching/ecc/eccPapers/Doche-ch09.pdf
//	[modboscoster]             Ayan Nandy. Modifications of Bos and Coster’s Heuristics in search of a
//	                           shorter addition chain for faster exponentiation. Masters thesis, Indian
//	                           Statistical Institute Kolkata. 2011.
//	                           http://library.isical.ac.in:8080/jspui/bitstream/123456789/6441/1/DISS-285.pdf
//	[mpnt]                     F. L. Ţiplea, S. Iftene, C. Hriţcu, I. Goriac, R. Gordân and E. Erbiceanu.
//	                           MpNT: A Multi-Precision Number Theory Package, Number Theoretical Algorithms
//	                           (I). Technical Report TR03-02, Faculty of Computer Science, "Alexandru Ioan
//	                           Cuza" University, Iasi. 2003. https://profs.info.uaic.ro/~tr/tr03-02.pdf
//	[speedsubgroup]            Stam, Martijn. Speeding up subgroup cryptosystems. PhD thesis, Technische
//	                           Universiteit Eindhoven. 2003. https://cr.yp.to/bib/2003/stam-thesis.pdf

// Heuristic suggests insertions given a current protosequence.
type Heuristic interface {
	fmt.Stringer
	Suggest(f []*big.Int, target *big.Int) []*big.Int
}

// HeuristicAlgorithm searches for an addition sequence using a heuristic at
// each step. This implements the framework given in [mpnt], page 63, with the
// heuristic playing the role of the "newnumbers" function.
type HeuristicAlgorithm struct {
	heuristic Heuristic
}

// NewHeuristicAlgorithm builds a heuristic algorithm.
func NewHeuristicAlgorithm(h Heuristic) *HeuristicAlgorithm {
	return &HeuristicAlgorithm{
		heuristic: h,
	}
}

func (h HeuristicAlgorithm) String() string {
	return fmt.Sprintf("heuristic(%v)", h.heuristic)
}

// FindSequence searches for an addition sequence for the given targets.
func (h HeuristicAlgorithm) FindSequence(targets []*big.Int) (Chain, error) {
	// Initialize protosequence.
	leader := bigints.Int64s(1, 2)
	proto := append(leader, targets...)
	bigints.Sort(proto)
	proto = bigints.Unique(proto)
	c := []*big.Int{}

	for len(proto) > 2 {
		// Pop the target element.
		top := len(proto) - 1
		target := proto[top]
		proto = proto[:top]
		c = bigints.InsertSortedUnique(c, target)

		// Apply heuristic.
		insert := h.heuristic.Suggest(proto, target)
		if insert == nil {
			return nil, errors.New("failed to find sequence")
		}

		// Update protosequence.
		proto = bigints.MergeUnique(proto, insert)
	}

	// Prepare the chain to return.
	c = bigints.MergeUnique(leader, c)

	return Chain(c), nil
}

// DeltaLargest implements the simple heuristic of adding the delta between the
// largest two entries in the protosequence.
type DeltaLargest struct{}

func (DeltaLargest) String() string { return "delta_largest" }

// Suggest proposes inserting target-max(f).
func (DeltaLargest) Suggest(f []*big.Int, target *big.Int) []*big.Int {
	n := len(f)
	delta := new(big.Int).Sub(target, f[n-1])
	return []*big.Int{delta}
}
