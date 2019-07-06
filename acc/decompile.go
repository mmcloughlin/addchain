package acc

import (
	"github.com/mmcloughlin/addchain"
	"github.com/mmcloughlin/addchain/acc/ir"
)

// Decompile an unrolled program into concise intermediate representation.
func Decompile(c addchain.Program) (*ir.Program, error) {
	p := &ir.Program{}
	for i := 0; i < len(c); i++ {
		op := c[i]

		// Regular addition.
		if !op.IsDouble() {
			p.AddInstruction(ir.Instruction{
				Output: ir.Index(i + 1),
				Op: ir.Add{
					X: ir.Index(op.I),
					Y: ir.Index(op.J),
				},
			})
			continue
		}

		// We have a double. Look ahead to see if this is a chain of doublings, which
		// can be encoded as a shift.
		j := i + 1
		for ; j < len(c) && c[j].I == j && c[j].J == j; j++ {
		}

		s := j - i

		// Shift size 1 encoded as a double.
		if s == 1 {
			p.AddInstruction(ir.Instruction{
				Output: ir.Index(i + 1),
				Op: ir.Double{
					X: ir.Index(op.I),
				},
			})
			continue
		}

		i = j - 1
		p.AddInstruction(ir.Instruction{
			Output: ir.Index(i + 1),
			Op: ir.Shift{
				X: ir.Index(op.I),
				S: uint(s),
			},
		})
	}
	return p, nil
}
