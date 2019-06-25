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
		c := Binary(n)
		AssertChainProduces(t, c, n)
		n.Lsh(n, 1)
	}
}

func TestBinaryRandomInt64(t *testing.T) {
	for trials := 0; trials < 1000; trials++ {
		r := rand.Int63n(math.MaxInt64)
		n := big.NewInt(r)
		c := Binary(n)
		AssertChainProduces(t, c, n)
	}
}

func AssertChainProduces(t *testing.T, c Chain, expect *big.Int) {
	x := c.Evaluate()
	if len(x) == 0 {
		t.Fatal("empty chain result")
	}
	last := x[len(x)-1]
	if last.Cmp(expect) != 0 {
		t.Fatalf("chain produced %s expect %s", last, expect)
	}
}
