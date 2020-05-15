package rand

import (
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/mmcloughlin/addchain"
	"github.com/mmcloughlin/addchain/alg"
	"github.com/mmcloughlin/addchain/internal/bigint"
	"github.com/mmcloughlin/addchain/internal/bigints"
)

// Generator can generate random chains.
type Generator interface {
	GenerateChain() (addchain.Chain, error)
	fmt.Stringer
}

// RandomAddsGenerator generates a random chain by making N random adds.
type RandomAddsGenerator struct {
	N int
}

func (r RandomAddsGenerator) String() string {
	return fmt.Sprintf("random_adds(%d)", r.N)
}

func (r RandomAddsGenerator) GenerateChain() (addchain.Chain, error) {
	c := addchain.New()
	for len(c) < r.N {
		i, j := rand.Intn(len(c)), rand.Intn(len(c))
		sum := new(big.Int).Add(c[i], c[j])
		c = bigints.InsertSortedUnique(c, sum)
	}
	return c, nil
}

// RandomSolverGenerator generates random N-bit values and uses an algorithm to
// build a chain for them.
type RandomSolverGenerator struct {
	N         uint
	Algorithm alg.ChainAlgorithm
	rand      *rand.Rand
}

func NewRandomSolverGenerator(n uint, a alg.ChainAlgorithm) RandomSolverGenerator {
	return RandomSolverGenerator{
		N:         n,
		Algorithm: a,
		rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r RandomSolverGenerator) String() string {
	return fmt.Sprintf("random_solver(%d,%s)", r.N, r.Algorithm)
}

func (r RandomSolverGenerator) GenerateChain() (addchain.Chain, error) {
	target := bigint.RandBits(r.rand, r.N)
	return r.Algorithm.FindChain(target)
}
