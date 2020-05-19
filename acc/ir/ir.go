// Package ir declares an intermediate representation for acc programs.
package ir

import (
	"fmt"
	"strings"

	"github.com/mmcloughlin/addchain"
)

// Program is a sequence of acc instructions.
type Program struct {
	Instructions []*Instruction

	// Pass/analysis results.
	Operands    map[int]*Operand
	ReadCount   map[int]int
	Program     addchain.Program
	Chain       addchain.Chain
	Temporaries []string
}

// AddInstruction appends an instruction to the program.
func (p *Program) AddInstruction(i *Instruction) {
	p.Instructions = append(p.Instructions, i)
}

// Output returns the output of the last instruction.
func (p Program) Output() *Operand {
	last := len(p.Instructions) - 1
	return p.Instructions[last].Output
}

func (p Program) String() string {
	var b strings.Builder
	for _, i := range p.Instructions {
		fmt.Fprintln(&b, i)
	}
	return b.String()
}

// Operand represents an element of an addition chain, with an optional
// identifier.
type Operand struct {
	Identifier string
	Index      int
}

// NewOperand builds a named operand for index i.
func NewOperand(name string, i int) *Operand {
	return &Operand{
		Identifier: name,
		Index:      i,
	}
}

// Index builds an unnamed operand for index i.
func Index(i int) *Operand {
	return NewOperand("", i)
}

// One is the first element in the addition chain, which by definition always
// has the value 1.
var One = Index(0)

func (o Operand) String() string {
	if len(o.Identifier) > 0 {
		return o.Identifier
	}
	return fmt.Sprintf("[%d]", o.Index)
}

// Instruction assigns the result of an operation to an operand.
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

// Op is an operation.
type Op interface {
	Inputs() []*Operand
	String() string
}

// Add is an addition operation.
type Add struct {
	X, Y *Operand
}

// Inputs returns the addends.
func (a Add) Inputs() []*Operand {
	return []*Operand{a.X, a.Y}
}

func (a Add) String() string {
	return fmt.Sprintf("%s + %s", a.X, a.Y)
}

// Double is a double operation.
type Double struct {
	X *Operand
}

// Inputs returns the operand.
func (d Double) Inputs() []*Operand {
	return []*Operand{d.X}
}

func (d Double) String() string {
	return fmt.Sprintf("2 * %s", d.X)
}

// Shift represents a shift-left operation, equivalent to repeat doubling.
type Shift struct {
	X *Operand
	S uint
}

// Inputs returns the operand to be shifted.
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
