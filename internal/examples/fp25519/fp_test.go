package fp25519

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func Trials() int {
	if testing.Short() {
		return 1 << 4
	}
	return 1 << 12
}

func RandElt(t *testing.T) *Elt {
	t.Helper()
	one := new(big.Int).SetInt64(1)
	max := new(big.Int).Sub(p, one)
	x, err := rand.Int(rand.Reader, max)
	if err != nil {
		t.Fatal(err)
	}
	x.Add(x, one)
	return new(Elt).SetInt(x)
}

func TestInv(t *testing.T) {
	for trial := 0; trial < Trials(); trial++ {
		x := RandElt(t)
		got := new(Elt).Inv(x)
		expect := new(big.Int).ModInverse(x.Int(), p)
		if got.Int().Cmp(expect) != 0 {
			t.FailNow()
		}
	}
}

func TestInvAlias(t *testing.T) {
	for trial := 0; trial < Trials(); trial++ {
		x := RandElt(t)
		expect := new(Elt).Inv(x) // non-aliased
		x.Inv(x)                  // aliased
		if x.Int().Cmp(expect.Int()) != 0 {
			t.FailNow()
		}
	}
}
