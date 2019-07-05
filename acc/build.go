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
	prepare := pass.Concat(
		pass.Func(pass.ReadCounts),
		pass.NameByIndex("i"),
	)

	if err := prepare.Execute(p); err != nil {
		return nil, err
	}

	// Maintain an expression for each index.
	b := newbuilder(p)
	for _, i := range p.Instructions {
		if err := b.instruction(i); err != nil {
			return nil, err
		}
	}
	b.finalize()
	return b.chain, nil
}

type builder struct {
	chain *ast.Chain
	prog  *ir.Program
	expr  map[int]ast.Expr
	last  *ir.Operand
}

func newbuilder(p *ir.Program) *builder {
	return &builder{
		chain: &ast.Chain{},
		prog:  p,
		expr:  map[int]ast.Expr{},
	}
}

func (b *builder) finalize() {
	result := ast.Statement{
		Expr: b.operand(b.last),
	}
	b.chain.Statements = append(b.chain.Statements, result)
	return
}

func (b *builder) instruction(i ir.Instruction) error {
	e, err := b.operator(i.Op)
	if err != nil {
		return err
	}

	out := i.Output
	b.last = out
	idx := out.Index
	if b.prog.ReadCount[idx] <= 1 {
		b.expr[idx] = e
		return nil
	}

	name := ast.Identifier(out.Identifier)
	stmt := ast.Statement{
		Name: name,
		Expr: e,
	}
	b.chain.Statements = append(b.chain.Statements, stmt)
	b.expr[idx] = name

	return nil
}

func (b *builder) operator(op ir.Op) (ast.Expr, error) {
	switch o := op.(type) {
	case ir.Add:
		return ast.Add{
			X: b.operand(o.X),
			Y: b.operand(o.Y),
		}, nil
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

func (b *builder) operand(op *ir.Operand) ast.Expr {
	e, ok := b.expr[op.Index]
	if !ok {
		return ast.Operand(op.Index)
	}
	return e
}
