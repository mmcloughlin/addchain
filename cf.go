package addchain

import (
	"fmt"
	"math/big"

	"github.com/mmcloughlin/addchain/internal/bigint"
	"github.com/mmcloughlin/addchain/internal/bigints"
)

// ContinuedFractionStrategy is a method of choosing the auxiliary integer k in
// the continued fraction method outlined by "Efficient computation of addition
// chains" by F. Bergeron, J. Berstel and S. Brlek.
type ContinuedFractionStrategy interface {
	fmt.Stringer
	K(n *big.Int) []*big.Int
}

// ContinuedFractions uses the continued fractions method for finding an
// addition chain.
type ContinuedFractions struct {
	strategy ContinuedFractionStrategy
}

func NewContinuedFractions(s ContinuedFractionStrategy) ContinuedFractions {
	return ContinuedFractions{
		strategy: s,
	}
}

func (a ContinuedFractions) String() string {
	return fmt.Sprintf("continued_fractions(%s)", a.strategy)
}

func (a ContinuedFractions) FindChain(target *big.Int) (Chain, error) {
	return a.FindSequence([]*big.Int{target})
}

func (a ContinuedFractions) FindSequence(targets []*big.Int) (Chain, error) {
	bigints.Sort(targets)
	return a.chain(targets), nil
}

func (a ContinuedFractions) minchain(n *big.Int) Chain {
	if bigint.IsPow2(n) {
		return bigint.Pow2UpTo(n)
	}

	if bigint.EqualInt64(n, 3) {
		return bigints.Int64s([]int64{1, 2, 3})
	}

	var min Chain
	for _, k := range a.strategy.K(n) {
		c := a.chain([]*big.Int{k, n})
		if min == nil || len(c) < len(min) {
			min = c
		}
	}

	return min
}

// chain produces a continued fraction chain for the given values. The slice ns
// must be in ascending order.
func (a ContinuedFractions) chain(ns []*big.Int) Chain {
	k := len(ns)
	if k == 1 || ns[k-2].Cmp(bigint.One()) <= 0 {
		return a.minchain(ns[k-1])
	}

	q, r := new(big.Int), new(big.Int)
	q.DivMod(ns[k-1], ns[k-2], r)

	cq := a.minchain(q)
	remaining := bigints.Clone(ns[:k-1])

	if bigint.IsZero(r) {
		return Product(a.chain(remaining), cq)
	}

	remaining = bigints.InsertSortedUnique(remaining, r)
	return Plus(Product(a.chain(remaining), cq), r)
}

type DichotomicStrategy struct{}

func (DichotomicStrategy) String() string { return "dichotomic" }

// K returns only one suggestion for k, namely floor(n / 2ʰ⁾ where h = log2(n)/2.
func (DichotomicStrategy) K(n *big.Int) []*big.Int {
	l := n.BitLen()
	h := uint(l) / 2
	k := new(big.Int).Div(n, bigint.Pow2(h))
	return []*big.Int{k}
}
