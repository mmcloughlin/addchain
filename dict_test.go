package addchain

import (
	"math/big"
	"math/rand"
	"reflect"
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
	expect := DictSum{
		{D: 0x0, E: 0},
		{D: 0xf, E: 4},
		{D: 0xe, E: 8},
		{D: 0xe, E: 12},
		{D: 0xb, E: 16},
		{D: 0x0, E: 20},
		{D: 0x1, E: 24},
	}
	if !reflect.DeepEqual(got, expect) {
		t.Fatalf("got %v expect %v", got, expect)
	}
}
