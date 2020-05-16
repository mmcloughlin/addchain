package alg_test

import (
	"testing"

	"github.com/mmcloughlin/addchain/alg"
	"github.com/mmcloughlin/addchain/alg/algtest"
)

func TestBinaryRightToLeft(t *testing.T) {
	algtest.ChainAlgorithm(t, alg.BinaryRightToLeft{})
}
