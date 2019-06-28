package addchain

import (
	"math/big"
	"reflect"
	"testing"
)

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
