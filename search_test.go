package addchain

import (
	"math"
	"math/big"
	"math/rand"
	"testing"

	"github.com/mmcloughlin/addchain/prime"

	"github.com/mmcloughlin/addchain/internal/test"
)

func TestChainAlgorithms(t *testing.T) {
	as := []ChainAlgorithm{
		BinaryRightToLeft{},
		NewContinuedFractions(BinaryStrategy{}),
		NewContinuedFractions(BinaryStrategy{Parity: 1}),
		NewContinuedFractions(DichotomicStrategy{}),
	}
	for _, a := range as {
		t.Run(a.String(), ChainAlgorithmSuite(a))
	}
}

func ChainAlgorithmSuite(a ChainAlgorithm) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("powers_of_two", CheckPowersOfTwo(a, 100))
		t.Run("random_int64", test.Trials(CheckRandomInt64(a)))
		t.Run("primes", CheckPrimes(a))
	}
}

func CheckPowersOfTwo(a ChainAlgorithm, e int) func(t *testing.T) {
	return func(t *testing.T) {
		n := big.NewInt(1)
		for i := 0; i < e; i++ {
			AssertChainAlgorithmProduces(t, a, n)
			n.Lsh(n, 1)
		}
	}
}

func CheckRandomInt64(a ChainAlgorithm) func(t *testing.T) {
	return func(t *testing.T) {
		r := rand.Int63n(math.MaxInt64)
		n := big.NewInt(r)
		AssertChainAlgorithmProduces(t, a, n)
	}
}

func CheckPrimes(a ChainAlgorithm) func(t *testing.T) {
	return func(t *testing.T) {
		for _, p := range prime.Distinguished {
			AssertChainAlgorithmProduces(t, a, p.Int())
		}
	}
}
