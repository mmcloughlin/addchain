package pass

import (
	"strconv"

	"github.com/mmcloughlin/addchain/acc/ir"
)

// Interface for a processing pass.
type Interface interface {
	Execute(*ir.Program) error
}

// Func adapts a function to the pass Interface.
type Func func(*ir.Program) error

// Execute calls p.
func (f Func) Execute(p *ir.Program) error {
	return f(p)
}

// Concat returns a pass that executes the given passes in order, stopping on
// the first error.
func Concat(passes ...Interface) Interface {
	return Func(func(p *ir.Program) error {
		for _, pass := range passes {
			if err := pass.Execute(p); err != nil {
				return err
			}
		}
		return nil
	})
}

// ReadCounts computes how many times each index is read in the program. This
// populates the ReadCount field of the program.
func ReadCounts(p *ir.Program) error {
	p.ReadCount = map[int]int{}
	for _, i := range p.Instructions {
		for _, input := range i.Op.Inputs() {
			p.ReadCount[input.Index]++
		}
	}
	return nil
}

// NameByIndex builds a pass that sets any unnamed operands to have name prefix
// + index.
func NameByIndex(prefix string) Interface {
	return Func(func(p *ir.Program) error {
		for _, i := range p.Instructions {
			for _, operand := range i.Operands() {
				if operand.Identifier != "" {
					continue
				}
				operand.Identifier = prefix + strconv.Itoa(operand.Index)
			}
		}
		return nil
	})
}
