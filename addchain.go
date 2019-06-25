package addchain

import "math/big"

// Add is an instruction to add positions I and J in a chain.
type Add struct{ I, J int }

// IsDouble returns whether this addition is a doubling.
func (a Add) IsDouble() bool { return a.I == a.J }

// Chain is an addition chain.
type Chain []Add

func (c *Chain) Double(i int) int {
	return c.Add(i, i)
}

func (c *Chain) Add(i, j int) int {
	*c = append(*c, Add{i, j})
	return len(*c)
}

func (c Chain) Doubles() int {
	doubles, _ := c.Ops()
	return doubles
}

func (c Chain) Adds() int {
	_, adds := c.Ops()
	return adds
}

func (c Chain) Ops() (doubles, adds int) {
	for _, add := range c {
		if add.IsDouble() {
			doubles++
		} else {
			adds++
		}
	}
	return
}

// Evaluate executes the chain.
func (c Chain) Evaluate() []*big.Int {
	a := []*big.Int{big.NewInt(1)}
	for _, add := range c {
		sum := new(big.Int).Add(a[add.I], a[add.J])
		a = append(a, sum)
	}
	return a
}
