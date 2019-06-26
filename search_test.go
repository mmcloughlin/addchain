package addchain

import (
	"math"
	"math/big"
	"math/rand"
	"testing"
)

func TestBinaryPowersOfTwo(t *testing.T) {
	n := big.NewInt(1)
	for i := 0; i < 100; i++ {
		p := BinaryRightToLeft(n)
		AssertProgramProduces(t, p, n)
		n.Lsh(n, 1)
	}
}

func TestBinaryRandomInt64(t *testing.T) {
	for trials := 0; trials < 1000; trials++ {
		r := rand.Int63n(math.MaxInt64)
		n := big.NewInt(r)
		p := BinaryRightToLeft(n)
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
