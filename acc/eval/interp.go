// Package eval provides an interpreter for acc programs.
package eval

import (
	"fmt"
	"math/big"

	"github.com/mmcloughlin/addchain/acc/ir"
	"github.com/mmcloughlin/addchain/internal/errutil"
)

// Interpreter for acc programs.  In contrast to evaluation using indexes, the
// interpreter executes the program using operand variable names, as if it was a
// block of source code. Internally it maintains the state of every variable,
// and program instructions update that state.
type Interpreter struct {
	state map[string]*big.Int
}

// NewInterpreter builds a new interpreter. Initially, all variables are unset.
func NewInterpreter() *Interpreter {
	return &Interpreter{
		state: map[string]*big.Int{},
	}
}

// Load the named variable.
func (i *Interpreter) Load(v string) (*big.Int, bool) {
	x, ok := i.state[v]
	return x, ok
}

// Store x into the named variable.
func (i *Interpreter) Store(v string, x *big.Int) {
	i.state[v] = x
}

// Initialize the variable v to x. Errors if v is already defined.
func (i *Interpreter) Initialize(v string, x *big.Int) error {
	if _, ok := i.Load(v); ok {
		return fmt.Errorf("variable %q is already defined", v)
	}
	i.Store(v, x)
	return nil
}

// Execute the program p.
func (i *Interpreter) Execute(p *ir.Program) error {
	for _, inst := range p.Instructions {
		if err := i.instruction(inst); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) instruction(inst *ir.Instruction) error {
	output := i.output(inst.Output)
	switch op := inst.Op.(type) {
	case ir.Add:
		x, err := i.operands(op.X, op.Y)
		if err != nil {
			return err
		}
		output.Add(x[0], x[1])
	case ir.Double:
		x, err := i.operand(op.X)
		if err != nil {
			return err
		}
		output.Add(x, x)
	case ir.Shift:
		x, err := i.operand(op.X)
		if err != nil {
			return err
		}
		output.Lsh(x, op.S)
	default:
		return errutil.UnexpectedType(op)
	}
	return nil
}

func (i *Interpreter) output(operand *ir.Operand) *big.Int {
	if x, ok := i.Load(operand.Identifier); ok {
		return x
	}
	x := new(big.Int)
	i.Store(operand.Identifier, x)
	return x
}

func (i *Interpreter) operands(operands ...*ir.Operand) ([]*big.Int, error) {
	xs := make([]*big.Int, 0, len(operands))
	for _, operand := range operands {
		x, err := i.operand(operand)
		if err != nil {
			return nil, err
		}
		xs = append(xs, x)
	}
	return xs, nil
}

func (i *Interpreter) operand(operand *ir.Operand) (*big.Int, error) {
	if operand.Identifier == "" {
		return nil, fmt.Errorf("operand %s missing identifier", operand)
	}
	x, ok := i.Load(operand.Identifier)
	if !ok {
		return nil, fmt.Errorf("operand %q is not defined", operand)
	}
	return x, nil
}
