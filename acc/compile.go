package acc

import (
	"fmt"

	"github.com/mmcloughlin/addchain"
	"github.com/mmcloughlin/addchain/acc/ir"
)

// Compile converts a program intermediate representation to a full addition
// chain program.
func Compile(p ir.Program) (addchain.Program, error) {
	c := addchain.Program{}
	for _, i := range p {
		var out int
		var err error

		switch op := i.Op.(type) {
		case ir.Add:
			out, err = c.Add(op.X.Index, op.Y.Index)
		case ir.Double:
			out, err = c.Double(op.X.Index)
		case ir.Shift:
			out, err = c.Shift(op.X.Index, op.S)
		default:
			return nil, fmt.Errorf("unexpected type %T", op)
		}

		if err != nil {
			return nil, err
		}
		if out != i.Output.Index {
			return nil, fmt.Errorf("incorrect output index")
		}
	}
	return c, nil
}

// Decompile an unrolled program into intermediate representation.
func Decompile(c addchain.Program) (ir.Program, error) {
	p := ir.Program{}
	for i := 0; i < len(c); i++ {
		op := c[i]

		// Regular addition.
		if !op.IsDouble() {
			p = append(p, ir.Instruction{
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
			p = append(p, ir.Instruction{
				Output: ir.Index(i + 1),
				Op: ir.Double{
					X: ir.Index(op.I),
				},
			})
			continue
		}

		i = j - 1
		p = append(p, ir.Instruction{
			Output: ir.Index(i + 1),
			Op: ir.Shift{
				X: ir.Index(op.I),
				S: uint(s),
			},
		})
	}
	return p, nil
}
