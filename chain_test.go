package addchain

import (
	"reflect"
	"testing"
)

func TestChainOps(t *testing.T) {
	cases := []struct {
		Name   string
		Chain  Chain
		Expect [][]Op
	}{
		{
			Name:  "short",
			Chain: Int64s(1, 2),
			Expect: [][]Op{
				{{0, 0}}, // 2
			},
		},
		{
			Name:  "multiple_choices",
			Chain: Int64s(1, 2, 3, 4),
			Expect: [][]Op{
				{{0, 0}},         // 2
				{{0, 1}},         // 3
				{{0, 2}, {1, 1}}, // 4
			},
		},
		{
			Name:  "non_ascending",
			Chain: Int64s(1, 2, 3, 4, 7, 5, 6),
			Expect: [][]Op{
				{{0, 0}},                 // 2
				{{0, 1}},                 // 3
				{{0, 2}, {1, 1}},         // 4
				{{2, 3}},                 // 7
				{{0, 3}, {1, 2}},         // 5
				{{0, 5}, {1, 3}, {2, 2}}, // 6
			},
		},
	}
	for _, c := range cases {
		c := c // scopelint
		t.Run(c.Name, func(t *testing.T) {
			var got [][]Op
			for k := 1; k < len(c.Chain); k++ {
				got = append(got, c.Chain.Ops(k))
			}
			if !reflect.DeepEqual(got, c.Expect) {
				t.Logf("got    = %v", got)
				t.Logf("expect = %v", c.Expect)
				t.Fail()
			}
		})
	}
}

func TestProduct(t *testing.T) {
	a := Int64s(1, 2, 4, 6, 10)
	b := Int64s(1, 2, 4, 8)
	got := Product(a, b)
	expect := Int64s(1, 2, 4, 6, 10, 20, 40, 80)
	if !reflect.DeepEqual(expect, got) {
		t.Fatalf("Product(%v, %v) = %v; expect %v", a, b, got, expect)
	}
}
