package acc

import (
	"github.com/mmcloughlin/addchain/acc/ast"
	"github.com/mmcloughlin/addchain/acc/ir"
	"github.com/mmcloughlin/addchain/acc/pass"
	"github.com/mmcloughlin/addchain/internal/errutil"
)

// Build AST from a program in intermediate representation.
func Build(p *ir.Program) (*ast.Chain, error) {
	// Run some analysis passes first.
	err := pass.Exec(p,
		pass.Func(pass.ReadCounts),
		pass.NameByteValues,
		pass.NameXRuns,
		pass.NameByIndex("i"),
	)
	if err != nil {
		return nil, err
	}

	// Delegate to builder.
	b := newbuilder(p)
	if err := b.process(); err != nil {
		return nil, err
	}

	return b.chain, nil
}

type builder struct {
	chain *ast.Chain
	prog  *ir.Program
	expr  map[int]ast.Expr
}

func newbuilder(p *ir.Program) *builder {
	return &builder{
		chain: &ast.Chain{},
		prog:  p,
		expr:  map[int]ast.Expr{},
	}
}

func (b *builder) process() error {
	insts := b.prog.Instructions
	n := len(insts)
	for i := 0; i < n; i++ {
		inst := insts[i]
		out := inst.Output

		// Build expression for the result of this instruction.
		e, err := b.operator(inst.Op)
		if err != nil {
			return err
		}

		b.expr[out.Index] = e

		// If this output is read only by the following instruction, we don't need to
		// commit it to a variable.
		usedonce := b.prog.ReadCount[out.Index] == 1
		usednext := i+1 < n && ir.HasInput(insts[i+1].Op, out.Index)
		if usedonce && usednext {
			continue
		}

		// Otherwise write a statement for it.
		b.commit(inst.Output)
	}

	// Clear the name of the final statement.
	b.chain.Statements[len(b.chain.Statements)-1].Name = ""

	return nil
}

func (b *builder) operator(op ir.Op) (ast.Expr, error) {
	switch o := op.(type) {
	case ir.Add:
		return b.add(o)
	case ir.Double:
		return ast.Double{
			X: b.operand(o.X),
		}, nil
	case ir.Shift:
		return ast.Shift{
			X: b.operand(o.X),
			S: o.S,
		}, nil
	default:
		return nil, errutil.UnexpectedType(op)
	}
}

func (b *builder) add(a ir.Add) (ast.Expr, error) {
	// Addition operator construction is slightly delcate, since operand order
	// determines ordering of execution. By the design of instruction processing
	// above, the only way we can have multi-operator expressions is with a
	// sequence of operands that are used only once and in the following
	// instruction. This implies that only one of x and y can be an operator
	// expression. In order to preserve execution order, whichever one that is
	// needs to be the first operand.

	x := b.operand(a.X)
	y := b.operand(a.Y)

	switch {
	case ast.IsOp(x) && ast.IsOp(y):
		return nil, errutil.AssertionFailure("only one of x and y should be an operator expression")
	case ast.IsOp(y):
		x, y = y, x
	case ast.IsOp(x):
		// Nothing, it's already the first operand.
	}

	return ast.Add{
		X: x,
		Y: y,
	}, nil
}

func (b *builder) commit(op *ir.Operand) {
	name := ast.Identifier(op.Identifier)
	stmt := ast.Statement{
		Name: name,
		Expr: b.operand(op),
	}
	b.chain.Statements = append(b.chain.Statements, stmt)
	b.expr[op.Index] = name
}

func (b *builder) operand(op *ir.Operand) ast.Expr {
	e, ok := b.expr[op.Index]
	if !ok {
		return ast.Operand(op.Index)
	}
	return e
}
