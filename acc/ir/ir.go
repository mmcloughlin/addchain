package ir

import (
	"fmt"
	"strings"
)

type Program struct {
	Instructions []Instruction
}

func (p *Program) AddInstruction(i Instruction) {
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
	Index      int
	Identifier string
}

func Index(i int) *Operand {
	return &Operand{Index: i}
}

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

func (i Instruction) String() string {
	return fmt.Sprintf("%s \u2190 %s", i.Output, i.Op)
}

type Op interface {
	fmt.Stringer
}

type Add struct {
	X, Y *Operand
}

func (a Add) String() string {
	return fmt.Sprintf("%s + %s", a.X, a.Y)
}

type Double struct {
	X *Operand
}

func (d Double) String() string {
	return fmt.Sprintf("2 * %s", d.X)
}

type Shift struct {
	X *Operand
	S uint
}

func (s Shift) String() string {
	return fmt.Sprintf("%s \u226a %d", s.X, s.S)
}
