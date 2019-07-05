package ast

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mmcloughlin/addchain/internal/errutil"
)

// Print an AST node to standard out.
func Print(n interface{}) error {
	return Fprint(os.Stderr, n)
}

// Fprint writes the AST node n to w.
func Fprint(w io.Writer, n interface{}) error {
	p := &printer{
		out: w,
	}
	p.node(n)
	return p.err
}

type printer struct {
	out     io.Writer
	level   int   // current indentation level
	pending bool  // if there's a pending indentation
	err     error // saved error from printing
}

func (p *printer) node(n interface{}) {
	switch n := n.(type) {
	case *Chain:
		p.enter("chain")
		for _, stmt := range n.Statements {
			p.statement(stmt)
		}
		p.leave()
	case Statement:
		p.statement(n)
	case Operand:
		p.linef("operand(%d)", n)
	case Identifier:
		p.linef("identifier(%q)", n)
	case Add:
		p.add(n)
	case Double:
		p.double(n)
	case Shift:
		p.shift(n)
	default:
		p.seterror(errutil.UnexpectedType(n))
	}
}

func (p *printer) statement(stmt Statement) {
	p.enter("statement")
	p.printf("name = ")
	p.node(stmt.Name)
	p.printf("expr = ")
	p.node(stmt.Expr)
	p.leave()
}

func (p *printer) add(a Add) {
	p.enter("add")
	p.printf("x = ")
	p.node(a.X)
	p.printf("y = ")
	p.node(a.Y)
	p.leave()
}
func (p *printer) double(d Double) {
	p.enter("double")
	p.printf("x = ")
	p.node(d.X)
	p.leave()
}

func (p *printer) shift(s Shift) {
	p.enter("shift")
	p.linef("s = %d", s.S)
	p.printf("x = ")
	p.node(s.X)
	p.leave()
}

func (p *printer) enter(name string) {
	p.linef("%s {", name)
	p.level++
}

func (p *printer) leave() {
	p.level--
	p.linef("}")
}

func (p *printer) linef(format string, args ...interface{}) {
	p.printf(format, args...)
	p.nl()
}

func (p *printer) nl() {
	p.printf("\n")
	p.pending = true
}

func (p *printer) printf(format string, args ...interface{}) {
	if p.err != nil {
		return
	}
	if p.pending {
		indent := strings.Repeat("\t", p.level)
		format = indent + format
		p.pending = false
	}
	_, err := fmt.Fprintf(p.out, format, args...)
	p.seterror(err)
}

func (p *printer) seterror(err error) {
	if p.err == nil {
		p.err = err
	}
}
