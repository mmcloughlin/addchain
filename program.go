package addchain

import "math/big"

// Op is an instruction to add positions I and J in a chain.
type Op struct{ I, J int }

// IsDouble returns whether this operation is a doubling.
func (o Op) IsDouble() bool { return o.I == o.J }

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
