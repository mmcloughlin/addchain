package acc

import (
	"reflect"
	"testing"

	"github.com/mmcloughlin/addchain"
	"github.com/mmcloughlin/addchain/acc/ir"
)

func TestDecompileExample(t *testing.T) {
	p := addchain.Program{}
	p.Double(0)
	p.Add(0, 1)
	p.Shift(1, 3)
	p.Add(0, 5)

	t.Log(p)

	expect := &ir.Program{
		Instructions: []*ir.Instruction{
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

	got, err := Decompile(p)
	if err != nil {
		t.Fatal()
	}

	t.Logf("got:\n%s", got)

	if !reflect.DeepEqual(got, expect) {
		t.Logf("expect:\n%s", expect)
		t.Fatal("mismatch")
	}
}
