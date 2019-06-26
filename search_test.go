package addchain

import (
	"math"
	"math/big"
	"math/rand"
	"testing"

	"github.com/mmcloughlin/addchain/internal/test"
)

func TestChainAlgorithms(t *testing.T) {
	as := []ChainAlgorithm{
		BinaryRightToLeft(),
	}
	for _, a := range as {
		t.Run(a.String(), ChainAlgorithmSuite(a))
	}
}

func ChainAlgorithmSuite(a ChainAlgorithm) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("powers_of_two", CheckPowersOfTwo(a, 100))
		t.Run("random_int64", test.Trials(CheckRandomInt64(a)))
	}
}

func CheckPowersOfTwo(a ChainAlgorithm, e int) func(t *testing.T) {
	return func(t *testing.T) {
		n := big.NewInt(1)
		for i := 0; i < e; i++ {
			p, err := a.FindChain(n)
			if err != nil {
				t.Fatal(err)
			}
			AssertProgramProduces(t, p, n)
			n.Lsh(n, 1)
		}
	}
}

func CheckRandomInt64(a ChainAlgorithm) func(t *testing.T) {
	return func(t *testing.T) {
		r := rand.Int63n(math.MaxInt64)
		n := big.NewInt(r)
		p, err := a.FindChain(n)
		if err != nil {
			t.Fatal(err)
		}
		AssertProgramProduces(t, p, n)
	}
}

func AssertProgramProduces(t *testing.T, p Program, expect *big.Int) {
	c := p.Evaluate()
	err := c.Produces(expect)
	if err != nil {
		t.Fatal(err)
	}
}
