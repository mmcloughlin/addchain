package addchain

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/mmcloughlin/addchain/internal/bigint"
	"github.com/mmcloughlin/addchain/internal/test"
)

func TestDecomposersRandom(t *testing.T) {
	ds := []Decomposer{
		FixedWindow{K: 2},
		FixedWindow{K: 11},
		FixedWindow{K: 16},

		SlidingWindow{K: 2},
		SlidingWindow{K: 7},
		SlidingWindow{K: 12},
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, d := range ds {
		d := d
		t.Run(d.String(), test.Trials(func(t *testing.T) bool {
			n := bigint.RandBits(r, 256)
			got := d.Decompose(n)
			if !bigint.Equal(got.Int(), n) {
				t.Fatalf("got %v expect %v", got, n)
			}
			return true
		}))
	}
}

func TestFixedWindow(t *testing.T) {
	n := big.NewInt(0x10beef0)
	f := FixedWindow{K: 4}
	got := f.Decompose(n)
	nibbles := []int64{0, 0xf, 0xe, 0xe, 0xb, 0x0, 0x1}
	if len(got) != len(nibbles) {
		t.FailNow()
	}
	for i, n := range nibbles {
		if !bigint.EqualInt64(got[i].D, n) || got[i].E != uint(4*i) {
			t.FailNow()
		}
	}
}

func TestDictAlgorithm(t *testing.T) {
	a := NewDictAlgorithm(
		SlidingWindow{K: 4},
		NewContinuedFractions(DichotomicStrategy{}),
	)
	n := big.NewInt(587257)
	c := AssertChainAlgorithmProduces(t, a, n)
	t.Log(c)
}
