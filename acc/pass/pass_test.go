package pass

import (
	"reflect"
	"testing"

	"github.com/mmcloughlin/addchain/acc/ir"
)

func TestReadCounts(t *testing.T) {
	p := &ir.Program{
		Instructions: []ir.Instruction{
			{
				Output: ir.Index(1),
				Op: ir.Double{
					X: ir.Index(0),
				},
			},
			{
				Output: ir.Index(2),
				Op: ir.Add{
					X: ir.Index(0),
					Y: ir.Index(1),
				},
			},
			{
				Output: ir.Index(5),
				Op: ir.Shift{
					X: ir.Index(1),
					S: 3,
				},
			},
			{
				Output: ir.Index(6),
				Op: ir.Add{
					X: ir.Index(0),
					Y: ir.Index(5),
				},
			},
		},
	}
	expect := map[int]int{
		0: 3,
		1: 2,
		5: 1,
	}

	if err := ReadCounts(p); err != nil {
		t.Fatal(p)
	}

	if !reflect.DeepEqual(expect, p.ReadCount) {
		t.FailNow()
	}
}
