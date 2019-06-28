package addchain

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/mmcloughlin/addchain/internal/bigints"
)

// References
//
// https://link.springer.com/content/pdf/10.1007/0-387-34805-0_37.pdf
// http://library.isical.ac.in:8080/jspui/bitstream/123456789/6441/1/DISS-285.pdf
// https://cr.yp.to/bib/2003/stam-thesis.pdf
// https://profs.info.uaic.ro/~tr/tr03-02.pdf
// https://koclab.cs.ucsb.edu/teaching/ecc/eccPapers/Doche-ch09.pdf
// https://github.com/kwantam/addchain

// Heuristic suggests insertions given a current protosequence.
type Heuristic interface {
	fmt.Stringer
	Suggest(f []*big.Int, target *big.Int) []*big.Int
}

// HeuristicAlgorithm searches for an addition sequence using a heuristic at each step.
type HeuristicAlgorithm struct {
	heuristic Heuristic
}

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

func (DeltaLargest) Suggest(f []*big.Int, target *big.Int) []*big.Int {
	n := len(f)
	delta := new(big.Int).Sub(target, f[n-1])
	return []*big.Int{delta}
}
