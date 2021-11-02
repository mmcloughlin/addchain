package pass

import (
	"fmt"

	"github.com/mmcloughlin/addchain/acc/ir"
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

	// Initialize an allocation. This maintains a map from operand index to
	// variable, and a pool of free variables.
	allocation := newallocation()

	// Process instructions in reverse, similar to liveness analysis. Allocate
	// when variables are read, mark available for re-allocation on write.
	for i := len(p.Instructions) - 1; i >= 0; i-- {
		inst := p.Instructions[i]

		// The output operand variable now becomes free.
		v := allocation.Variable(inst.Output.Index)
		allocation.Free(v)

		// Inputs may need variables, if they are not already live.
		for _, input := range inst.Op.Inputs() {
			allocation.Allocate(input.Index)
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
	outv := allocation.Variable(p.Output().Index)
	for _, index := range p.Indexes {
		op := p.Operands[index]
		v := allocation.Variable(op.Index)
		_, ok := name[v]
		switch {
		// Operand index 0 is the input.
		case op.Index == 0:
			op.Identifier = a.Input
		// Use the output variable name after the last use of the input.
		case v == outv && op.Index >= lastinputread:
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

// allocation of a pool of variables to operands.
type allocation struct {
	// Allocation map from operand index to variable index.
	variable map[int]int

	// List of available variables.
	available []int

	// Total number of variables.
	n int
}

// newallocation initializes an empty allocation.
func newallocation() *allocation {
	return &allocation{
		variable:  map[int]int{},
		available: []int{},
		n:         0,
	}
}

// Allocate ensures a variable is allocated to operand index i.
func (a *allocation) Allocate(i int) {
	// Return if it already has an assignment.
	_, ok := a.variable[i]
	if ok {
		return
	}

	// If there's nothing available, we'll need one more temporary.
	if len(a.available) == 0 {
		a.available = append(a.available, a.n)
		a.n++
	}

	// Assign from the available list.
	last := len(a.available) - 1
	a.variable[i] = a.available[last]
	a.available = a.available[:last]
}

// Variable allocated to operand index i. Allocates one if it doesn't already
// have an allocation.
func (a *allocation) Variable(i int) int {
	a.Allocate(i)
	return a.variable[i]
}

// Free marks v as available to be allocated to another operand.
func (a *allocation) Free(v int) {
	a.available = append(a.available, v)
}
