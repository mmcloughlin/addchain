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

		RunLength{T: 0},
		RunLength{T: 1},
		RunLength{T: 3},
		RunLength{T: 7},
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
	n := big.NewInt(0xbeef0 << 3)
	f := FixedWindow{K: 4}
	got := f.Decompose(n)
	expect := DictSum{
		{D: big.NewInt(0xf), E: 7},
		{D: big.NewInt(0xe), E: 11},
		{D: big.NewInt(0xe), E: 15},
		{D: big.NewInt(0xb), E: 19},
	}
	if !DictSumEquals(got, expect) {
		t.Fatalf("got %v expect %v", got, expect)
	}
}

func TestRunLength(t *testing.T) {
	cases := []struct {
		T      uint
		X      int64
		Expect DictSum
	}{
		{
			T: 4,
			X: 0xff,
			Expect: DictSum{
				{D: big.NewInt(0xf), E: 0},
				{D: big.NewInt(0xf), E: 4},
			},
		},
		{
			T: 0,
			X: 0xff,
			Expect: DictSum{
				{D: big.NewInt(0xff), E: 0},
			},
		},
		{
			T: 4,
			X: 0xf0f0,
			Expect: DictSum{
				{D: big.NewInt(0xf), E: 4},
				{D: big.NewInt(0xf), E: 12},
			},
		},
		{
			T: 0,
			X: 0xf0f0,
			Expect: DictSum{
				{D: big.NewInt(0xf), E: 4},
				{D: big.NewInt(0xf), E: 12},
			},
		},
		{
			T: 3,
			X: 0xff,
			Expect: DictSum{
				{D: big.NewInt(0x3), E: 0},
				{D: big.NewInt(0x7), E: 2},
				{D: big.NewInt(0x7), E: 5},
			},
		},
	}
	for _, c := range cases {
		d := RunLength{T: c.T}
		if got := d.Decompose(big.NewInt(c.X)); !DictSumEquals(got, c.Expect) {
			t.Fatalf("Decompose(%#x) = %v; expect %v", c.X, got, c.Expect)
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

func DictSumEquals(a, b DictSum) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].E != b[i].E {
			return false
		}
		if !bigint.Equal(a[i].D, b[i].D) {
			return false
		}
	}
	return true
}
