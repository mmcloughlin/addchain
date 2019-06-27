package addchain

import (
	"math/big"
	"reflect"
	"testing"
)

func AssertChainAlgorithmGenerates(t *testing.T, a ChainAlgorithm, n *big.Int, expect Chain) {
	c, err := a.FindChain(n)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Validate(); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expect, c) {
		t.Fatalf("got %v; expect %v", c, expect)
	}
}

func AssertChainAlgorithmProduces(t *testing.T, a ChainAlgorithm, n *big.Int) {
	c, err := a.FindChain(n)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Produces(n)
	if err != nil {
		t.Log(c)
		t.Fatal(err)
	}
}
