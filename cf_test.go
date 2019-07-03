package addchain

import (
	"math/big"
	"testing"

	"github.com/mmcloughlin/addchain/internal/bigints"
)

func TestBinaryStrategy(t *testing.T) {
	a := AsChainAlgorithm{NewContinuedFractions(BinaryStrategy{})}
	n := big.NewInt(87)
	expect := bigints.Int64s(1, 2, 4, 5, 10, 20, 21, 42, 43, 86, 87)
	AssertChainAlgorithmGenerates(t, a, n, expect)
}

func TestCoBinaryStrategy(t *testing.T) {
	a := AsChainAlgorithm{NewContinuedFractions(CoBinaryStrategy{})}
	n := big.NewInt(87)
	expect := bigints.Int64s(1, 2, 3, 5, 10, 11, 21, 22, 43, 44, 87)
	AssertChainAlgorithmGenerates(t, a, n, expect)
}
