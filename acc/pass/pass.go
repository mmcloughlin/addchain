package pass

import "github.com/mmcloughlin/addchain/acc/ir"

// ReadCounts computes how many times each index is read in the program. This
// populates the ReadCount field of the program.
func ReadCounts(p *ir.Program) error {
	p.ReadCount = map[int]int{}
	for _, i := range p.Instructions {
		for _, input := range i.Op.Inputs() {
			p.ReadCount[input.Index]++
		}
	}
	return nil
}
