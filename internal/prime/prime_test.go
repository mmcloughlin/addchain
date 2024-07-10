package prime

import (
	"testing"

	"github.com/mmcloughlin/addchain/internal/polynomial"
)

func TestSolinas(t *testing.T) {
	// The "Goldilocks" prime.
	p := NewSolinas(polynomial.Polynomial{{A: -1, N: 0}, {A: -1, N: 1}, {A: 1, N: 2}}, 224)

	if p.Bits() != 448 {
		t.FailNow()
	}

	x := p.Int()
	decimal := "726838724295606890549323807888004534353641360687318060281490199180612328166730772686396383698676545930088884461843637361053498018365439"
	if x.String() != decimal {
		t.FailNow()
	}

	if p.String() != "2^448-2^224-1" {
		t.FailNow()
	}
}

func TestSmoothIsogeny(t *testing.T) {
	// p₅₁₂ from [isogenychains].
	p := NewSmoothIsogeny(2, 3, 253, 161, 7, false)

	if p.Bits() != 511 {
		t.FailNow()
	}

	x := p.Int()
	decimal := "6640624951081187159942983469764469901416062130859495455216614392426065341738463661693533115419196273210738003796604179119423082390833875356421735665631231"
	if x.String() != decimal {
		t.FailNow()
	}

	if p.String() != "2^253*3^161*7-1" {
		t.FailNow()
	}
}
