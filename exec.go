package addchain

import (
	"errors"
	"math/big"
	"sync"

	"github.com/mmcloughlin/addchain/internal/bigint"
)

// Result from applying an algorithm to a target.
type Result struct {
	Target    *big.Int
	Algorithm ChainAlgorithm
	Err       error
	Chain     Chain
	Program   Program
}

// Execute the algorithm on the target number n.
func Execute(n *big.Int, a ChainAlgorithm) Result {
	r := Result{
		Target:    n,
		Algorithm: a,
	}

	r.Chain, r.Err = a.FindChain(n)
	if r.Err != nil {
		return r
	}

	// Note this also performs validation.
	r.Program, r.Err = r.Chain.Program()
	if r.Err != nil {
		return r
	}

	// Still, verify that it produced what we wanted.
	if !bigint.Equal(r.Chain.End(), n) {
		r.Err = errors.New("did not produce the required value")
	}

	return r
}

// Parallel executes multiple algorithms in parallel.
func Parallel(n *big.Int, as []ChainAlgorithm) []Result {
	rs := make([]Result, len(as))
	var wg sync.WaitGroup
	for i, a := range as {
		wg.Add(1)
		go func(i int, a ChainAlgorithm) {
			defer wg.Done()
			rs[i] = Execute(n, a)
		}(i, a)
	}
	wg.Wait()
	return rs
}
