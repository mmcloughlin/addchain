package ir

import (
	"fmt"
	"strings"

	"github.com/mmcloughlin/addchain"
)

type Program struct {
	Instructions []*Instruction

	// Pass/analysis results.
	Operands  map[int]*Operand
	ReadCount map[int]int
	Program   addchain.Program
	Chain     addchain.Chain
}

func (p *Program) AddInstruction(i *Instruction) {
	p.Instructions = append(p.Instructions, i)
}

func (p Program) String() string {
	var b strings.Builder
	for _, i := range p.Instructions {
		fmt.Fprintln(&b, i)
	}
	return b.String()
}

type Operand struct {
	Identifier string
	Index      int
}

func NewOperand(name string, i int) *Operand {
	return &Operand{
		Identifier: name,
		Index:      i,
	}
}

func Index(i int) *Operand {
	return NewOperand("", i)
}

var One = Index(0)

func (o Operand) String() string {
	if len(o.Identifier) > 0 {
		return o.Identifier
	}
	return fmt.Sprintf("[%d]", o.Index)
}

type Instruction struct {
	Output *Operand
	Op     Op
}

// Operands returns the input and output operands.
func (i Instruction) Operands() []*Operand {
	return append(i.Op.Inputs(), i.Output)
}

func (i Instruction) String() string {
	return fmt.Sprintf("%s \u2190 %s", i.Output, i.Op)
}

type Op interface {
	Inputs() []*Operand
	fmt.Stringer
}

type Add struct {
	X, Y *Operand
}

func (a Add) Inputs() []*Operand {
	return []*Operand{a.X, a.Y}
}

func (a Add) String() string {
	return fmt.Sprintf("%s + %s", a.X, a.Y)
}

type Double struct {
	X *Operand
}

func (d Double) Inputs() []*Operand {
	return []*Operand{d.X}
}

func (d Double) String() string {
	return fmt.Sprintf("2 * %s", d.X)
}

type Shift struct {
	X *Operand
	S uint
}

func (s Shift) Inputs() []*Operand {
	return []*Operand{s.X}
}

func (s Shift) String() string {
	return fmt.Sprintf("%s \u226a %d", s.X, s.S)
}

// HasInput reports whether the given operation takes idx as an input.
func HasInput(op Op, idx int) bool {
	for _, input := range op.Inputs() {
		if input.Index == idx {
			return true
		}
	}
	return false
}
