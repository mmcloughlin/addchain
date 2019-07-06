package pass

import (
	"testing"

	"github.com/mmcloughlin/addchain/acc/ir"
)

func TestNameBinaryValues(t *testing.T) {
	p := &ir.Program{
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
		},
	}

	// Pass to name 4-bit values.
	n := NameBinaryValues(4, "_%b")

	t.Logf("pre:\n%s", p)

	if err := n.Execute(p); err != nil {
		t.Fatal(err)
	}

	t.Logf("post:\n%s", p)

	// Expected names.
	expect := map[int]string{
		0: "_1",
		1: "_10",
		2: "_11",
		5: "", // should be larger than 4 bits
	}

	for idx, name := range expect {
		got := p.Operands[idx].Identifier
		if got != name {
			t.Errorf("operand %d has name %s; expected %s", idx, got, name)
		}
	}
}
