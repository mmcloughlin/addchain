package addchain

import (
	"reflect"
	"testing"
)

func TestProduct(t *testing.T) {
	a := Int64s(1, 2, 4, 6, 10)
	b := Int64s(1, 2, 4, 8)
	got := Product(a, b)
	expect := Int64s(1, 2, 4, 6, 10, 20, 40, 80)
	if !reflect.DeepEqual(expect, got) {
		t.Fatalf("Product(%v, %v) = %v; expect %v", a, b, got, expect)
	}
}
