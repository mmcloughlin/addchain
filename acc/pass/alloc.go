package pass

import (
	"fmt"

	"github.com/mmcloughlin/addchain/acc/ir"
	"github.com/mmcloughlin/addchain/internal/container/heap"
	"github.com/mmcloughlin/addchain/internal/errutil"
)

// Allocator pass assigns a minimal number of temporary variables to execute a program.
type Allocator struct {
	// Input is the name of the input variable. Note this is index 0, or the
	// identity element of the addition chain.
	Input string

	// Output is the name to give to the final output of the addition chain. This
	// variable may itself be used as a temporary during execution.
	Output string

	// Format defines how to format any temporary variables. This format string
	// must accept one integer value. For example "t%d" would be a reasonable
	// choice.
	Format string
}

// Execute performs temporary variable allocation.
func (a Allocator) Execute(p *ir.Program) error {
	// Canonicalize operands, collect unique indexes, and delete all names.
	if err := Exec(p, Func(CanonicalizeOperands), Func(Indexes), Func(ClearNames)); err != nil {
		return err
	}

	// Keep an allocation map from operand index to variable index.
	allocation := map[int]int{}

	// Keep a heap of available variables, and a total variable count.
	available := heap.NewMinInts()
	n := 0

	// Assign a variable for the output.
	out := p.Output()
	allocation[out.Index] = 0
	n = 1

	// Process instructions in reverse.
	for i := len(p.Instructions) - 1; i >= 0; i-- {
		inst := p.Instructions[i]

		// The output operand variable now becomes available.
		v, ok := allocation[inst.Output.Index]
		if !ok {
			return errutil.AssertionFailure("output operand %d missing allocation", inst.Output.Index)
		}
		available.Push(v)

		// Inputs may need variables, if they are not already live.
		for _, input := range inst.Op.Inputs() {
			_, ok := allocation[input.Index]
			if ok {
				continue
			}

			// If there's nothing available, we'll need one more temporary.
			if available.Empty() {
				available.Push(n)
				n++
			}

			allocation[input.Index] = available.Pop()
		}
	}

	// Assign names to the operands. Reuse of the output variable is handled
	// specially, since we have to account for the possibility that it could be
	// aliased with the input. Prior to the last use of the input, variable 0
	// will be a temporary, after it will be the output.
	lastinputread := 0
	for _, inst := range p.Instructions {
		for _, input := range inst.Op.Inputs() {
			if input.Index == 0 {
				lastinputread = inst.Output.Index
			}
		}
	}

	// Map from variable index to name.
	name := map[int]string{}
	for _, index := range p.Indexes {
		op := p.Operands[index]
		v := allocation[op.Index]
		_, ok := name[v]
		switch {
		// Operand index 0 is the input.
		case op.Index == 0:
			op.Identifier = a.Input
		// Variable index 0 is the output, as long as we're past the last use of
		// the input.
		case v == 0 && op.Index >= lastinputread:
			op.Identifier = a.Output
		// Unnamed variable: allocate a temporary.
		case !ok:
			name[v] = fmt.Sprintf(a.Format, len(p.Temporaries))
			p.Temporaries = append(p.Temporaries, name[v])
			fallthrough
		default:
			op.Identifier = name[v]
		}
	}

	return nil
}
