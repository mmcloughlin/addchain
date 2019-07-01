package addchain

import (
	"fmt"
	"math/big"
)

// Optimized applies chain optimization to the result of a wrapped algorithm.
type Optimized struct {
	Algorithm ChainAlgorithm
}

func (o Optimized) String() string {
	return fmt.Sprintf("opt(%s)", o.Algorithm)
}

// FindChain delegates to the wrapped algorithm, then runs Optimize on the result.
func (o Optimized) FindChain(n *big.Int) (Chain, error) {
	c, err := o.Algorithm.FindChain(n)
	if err != nil {
		return nil, err
	}

	opt, err := Optimize(c)
	if err != nil {
		return nil, err
	}

	return opt, nil
}

// Optimize aims to remove redundancy from an addition chain.
func Optimize(c Chain) (Chain, error) {
	// Build program for c with all possible options at each step.
	ops := make([][]Op, len(c))
	for k := 1; k < len(c); k++ {
		ops[k] = c.Ops(k)
	}

	// Count how many times each index is used where it is the only available Op.
	counts := make([]int, len(c))
	for k := 1; k < len(c); k++ {
		if len(ops[k]) != 1 {
			continue
		}
		for _, i := range ops[k][0].Operands() {
			counts[i]++
		}
	}

	// Now, try to remove the positions which are never the only available op.
	remove := []int{}
	for k := 1; k < len(c)-1; k++ {
		if counts[k] > 0 {
			continue
		}

		// Prune places k is used.
		for l := k + 1; l < len(c); l++ {
			ops[l] = pruneuses(ops[l], k)

			// If this list now only has one element, the operands in it are now
			// indispensable.
			if len(ops[l]) == 1 {
				for _, i := range ops[l][0].Operands() {
					counts[i]++
				}
			}
		}

		// Mark k for deletion.
		remove = append(remove, k)
	}

	// Perform removals.
	pruned := Chain{}
	for i, x := range c {
		if len(remove) > 0 && remove[0] == i {
			remove = remove[1:]
			continue
		}
		pruned = append(pruned, x)
	}

	return pruned, nil
}

// pruneuses removes any uses of i from the list of operations.
func pruneuses(ops []Op, i int) []Op {
	filtered := ops[:0]
	for _, op := range ops {
		if !op.Uses(i) {
			filtered = append(filtered, op)
		}
	}
	return filtered
}
