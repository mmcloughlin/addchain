// Package eval provides an interpreter for acc programs.
package eval

import (
	"fmt"
	"math/big"

	"github.com/mmcloughlin/addchain/acc/ir"
)

type Interpreter struct {
	state map[string]*big.Int
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		state: map[string]*big.Int{},
	}
}

// Load the variable v.
func (i *Interpreter) Load(v string) (*big.Int, bool) {
	x, ok := i.state[v]
	return x, ok
}

// Store x into the variable v.
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
	return nil
}

//
// func (e *Evaluator) assignment(a ast.Assignment) error {
// 	lhs := e.dst(a.LHS)
// 	switch expr := a.RHS.(type) {
// 	case ast.Pow:
// 		x, err := e.operands(expr.X, expr.N)
// 		if err != nil {
// 			return err
// 		}
// 		lhs.Exp(x[0], x[1], e.m)
// 	case ast.Inv:
// 		x, err := e.operand(expr.X)
// 		if err != nil {
// 			return err
// 		}
// 		lhs.ModInverse(x, e.m)
// 	case ast.Mul:
// 		x, err := e.operands(expr.X, expr.Y)
// 		if err != nil {
// 			return err
// 		}
// 		lhs.Mul(x[0], x[1])
// 	case ast.Neg:
// 		x, err := e.operand(expr.X)
// 		if err != nil {
// 			return err
// 		}
// 		lhs.Neg(x)
// 	case ast.Add:
// 		x, err := e.operands(expr.X, expr.Y)
// 		if err != nil {
// 			return err
// 		}
// 		lhs.Add(x[0], x[1])
// 	case ast.Sub:
// 		x, err := e.operands(expr.X, expr.Y)
// 		if err != nil {
// 			return err
// 		}
// 		lhs.Sub(x[0], x[1])
// 	case ast.Cond:
// 		x, err := e.operands(expr.X, expr.C)
// 		if err != nil {
// 			return err
// 		}
// 		if x[1].Sign() != 0 {
// 			lhs.Set(x[0])
// 		}
// 	case ast.Variable, ast.Constant:
// 		x, err := e.operand(expr)
// 		if err != nil {
// 			return err
// 		}
// 		lhs.Set(x)
// 	default:
// 		return errutil.UnexpectedType(expr)
// 	}
// 	lhs.Mod(lhs, e.m)
// 	return nil
// }
//
// func (e Evaluator) dst(v ast.Variable) *big.Int {
// 	if x, ok := e.Load(v); ok {
// 		return x
// 	}
// 	x := new(big.Int)
// 	e.Store(v, x)
// 	return x
// }
//
// func (e *Evaluator) operands(operands ...ast.Operand) ([]*big.Int, error) {
// 	xs := make([]*big.Int, 0, len(operands))
// 	for _, operand := range operands {
// 		x, err := e.operand(operand)
// 		if err != nil {
// 			return nil, err
// 		}
// 		xs = append(xs, x)
// 	}
// 	return xs, nil
// }
//
// func (e *Evaluator) operand(operand ast.Operand) (*big.Int, error) {
// 	switch op := operand.(type) {
// 	case ast.Variable:
// 		x, ok := e.Load(op)
// 		if !ok {
// 			return nil, xerrors.Errorf("variable %q is not defined", op)
// 		}
// 		return x, nil
// 	case ast.Constant:
// 		return new(big.Int).SetUint64(uint64(op)), nil
// 	default:
// 		return nil, errutil.UnexpectedType(op)
// 	}
// }
