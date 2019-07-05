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
	for _, op := range c {
		inst := ir.Instruction{
			Output: ir.Index(len(p) + 1),
			Op: ir.Add{
				X: ir.Index(op.I),
				Y: ir.Index(op.J),
			},
		}
		p = append(p, inst)
	}
	return p, nil
}
