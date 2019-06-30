package bigint

import (
	"math/big"
	"testing"
)

func TestMask(t *testing.T) {
	if Mask(4, 16).Uint64() != 0xfff0 {
		t.Fail()
	}
}

func TestOnes(t *testing.T) {
	if Ones(8).Uint64() != 0xff {
		t.Fail()
	}
}

func TestExtract(t *testing.T) {
	x := big.NewInt(0xbeefcafe)
	if Extract(x, 4, 16).Uint64() != 0xcaf {
		t.Fail()
	}
}
