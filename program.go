package addchain

import (
	"math/big"

	"github.com/mmcloughlin/addchain/internal/bigint"
)

// Op is an instruction to add positions I and J in a chain.
type Op struct{ I, J int }

// IsDouble returns whether this operation is a doubling.
func (o Op) IsDouble() bool { return o.I == o.J }

// Operands returns the indicies used in this operation. This will contain one
// or two entries depending on whether this is a doubling.
func (o Op) Operands() []int {
	if o.IsDouble() {
		return []int{o.I}
	}
	return []int{o.I, o.J}
}

// Uses reports whether the given index is one of the operands.
func (o Op) Uses(i int) bool {
	return o.I == i || o.J == i
}

// Program is a sequence of operations.
type Program []Op

func (p *Program) Double(i int) int {
	return p.Add(i, i)
}

func (p *Program) Add(i, j int) int {
	*p = append(*p, Op{i, j})
	return len(*p)
}

func (p Program) Doubles() int {
	doubles, _ := p.Count()
	return doubles
}

func (p Program) Adds() int {
	_, adds := p.Count()
	return adds
}

func (p Program) Count() (doubles, adds int) {
	for _, op := range p {
		if op.IsDouble() {
			doubles++
		} else {
			adds++
		}
	}
	return
}

// Evaluate executes the program and returns the resulting chain.
func (p Program) Evaluate() Chain {
	c := New()
	for _, op := range p {
		sum := new(big.Int).Add(c[op.I], c[op.J])
		c = append(c, sum)
	}
	return c
}

// Dependencies returns an array of bitsets where each bitset contains the set
// of indicies that contributed to that position.
func (p Program) Dependencies() []*big.Int {
	bitsets := []*big.Int{bigint.One()}
	for i, op := range p {
		bitset := new(big.Int).Or(bitsets[op.I], bitsets[op.J])
		bitset.SetBit(bitset, i+1, 1)
		bitsets = append(bitsets, bitset)
	}
	return bitsets
}
