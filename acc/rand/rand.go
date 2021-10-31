package rand

import (
	"fmt"
	"math/rand"

	"github.com/mmcloughlin/addchain/acc/ir"
)

// Generator can generate random addition chain programs.
type Generator interface {
	GenerateProgram() (*ir.Program, error)
	String() string
}

// AddsGenerator generates a random program with N adds, in such a way that
// every index is used.
type AddsGenerator struct {
	N int
}

func (a AddsGenerator) String() string {
	return fmt.Sprintf("random_adds(%d)", a.N)
}

func (a AddsGenerator) GenerateProgram() (*ir.Program, error) {
	p := &ir.Program{}
	for i := 1; i <= a.N; i++ {
		p.AddInstruction(&ir.Instruction{
			Output: ir.Index(i),
			Op: ir.Add{
				X: ir.Index(i - 1), // ensure every index is used
				Y: ir.Index(rand.Intn(i)),
			},
		})
	}
	return p, nil
}
