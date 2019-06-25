package addchain

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/mmcloughlin/addchain/bigint"
)

// References:
//
// https://link.springer.com/content/pdf/10.1007/0-387-34805-0_37.pdf
// http://library.isical.ac.in:8080/jspui/bitstream/123456789/6441/1/DISS-285.pdf
// https://cr.yp.to/bib/2003/stam-thesis.pdf
// https://profs.info.uaic.ro/~tr/tr03-02.pdf
// https://koclab.cs.ucsb.edu/teaching/ecc/eccPapers/Doche-ch09.pdf
// https://github.com/kwantam/addchain

// BosCosterMakeSequence applies a variant of the Bos-Coster MakeSequence
// algorithm to generate an addition sequence producing every element of targets.
func BosCosterMakeSequence(targets []*big.Int) ([]*big.Int, error) {
	// Initialize the protosequence.
	f := []*big.Int{big.NewInt(1), big.NewInt(2)}
	f = append(f, targets...)
	result := []*big.Int{}

	for len(f) > 2 {
		fmt.Println(f)

		// Pop the target element.
		top := len(f) - 1
		target := f[top]
		f = f[:top]
		result = append(result, target)

		// Heuristics.
		insert := approx(f, target)
		if insert == nil {
			insert = halving(f, target)
		}

		// Bail if we found nothing.
		if insert == nil {
			return nil, errors.New("failed to find addition sequence")
		}

		// Update protosequence.
		f = append(f, insert...)
		bigint.Sort(f)
	}

	return result, nil
}

// approx applies the "Approximation" heuristic.
//
// This heuristic looks for two elements a, b in the list that sum to something close to the top element f.
// That is, we look for f-(a+b) = epsilon where a <= b and epsilon is a "small" positive value.
func approx(f []*big.Int, target *big.Int) []*big.Int {
	// Look for the closest sum.
	delta := new(big.Int)
	var zero big.Int
	var mindelta *big.Int
	var besta *big.Int

	for i, a := range f {
		for _, b := range f[i:] {
			delta.Add(a, b)
			delta.Sub(target, delta)
			if delta.Cmp(&zero) < 0 {
				continue
			}
			if mindelta == nil || delta.Cmp(mindelta) < 0 {
				mindelta = new(big.Int).Set(delta)
				besta = a
			}
		}
	}

	// Exit if we didn't find anything at all.
	if mindelta == nil {
		return nil
	}

	// It it small enough? The literature is unclear on good choices for epsilon.
	// We follow Riad S. Wahby's implementation that uses epsilon approximately
	// log(target).
	//
	// Reference: https://github.com/kwantam/addchain/blob/abe1e1c254673e32ed923088b89080c14874e5b3/boscoster.go#L161-L164
	//
	//	func bc_approx_test(d []*big.Int) (int) {
	//	    var targ = d[len(d)-1]
	//	    var tmp = big.NewInt(0)
	//	    var eps = big.NewInt(int64(targ.BitLen() - 1))
	//
	// TODO(mbm): investigate choices of epsilon in Bos-Coster "Approximation" heuristic.
	epsilon := big.NewInt(int64(target.BitLen()))
	if mindelta.Cmp(epsilon) > 0 {
		return nil
	}

	// We have found a sum within epsilon of the target.
	// Return a + epsilon to be added to the protosequence.
	insert := new(big.Int).Add(besta, epsilon)
	return []*big.Int{insert}
}

// division applies the "Division" heuristic.
func division(f []*big.Int, target *big.Int) []*big.Int {
	return nil
}

// halving applies the "Halving" heuristic.
func halving(f []*big.Int, target *big.Int) []*big.Int {
	t := new(big.Int)

	// Look for target - f[i] = 2^u * k with maximal u.
	maxu := 0
	var s *big.Int
	for i := range f {
		t.Sub(target, f[i])
		u := bigint.TrailingZeros(t)
		if u >= maxu {
			maxu, s = u, f[i]
		}
	}

	// Bail if we didn't find anything even.
	if s == nil {
		return nil
	}

	// Otherwise we return the chain.
	insertions := []*big.Int{}
	k := t.Sub(target, s)
	for ; maxu >= 0; maxu-- {
		insertions = append(insertions, new(big.Int).Set(k))
		k.Rsh(k, 1)
	}

	return insertions
}
