package addchain

import (
	"math/big"

	"github.com/mmcloughlin/addchain/internal/bigint"
)

// BinaryRightToLeft uses the right-to-left binary method to produce a program for n.
func BinaryRightToLeft(n *big.Int) Program {
	p := Program{}
	s := -1
	d := 0
	b := new(big.Int).Set(n)
	for {
		if b.Bit(0) == 1 {
			if s >= 0 {
				s = p.Add(s, d)
			} else {
				s = d
			}
		}
		b.Rsh(b, 1)
		if bigint.IsZero(b) {
			break
		}
		d = p.Double(d)
	}
	return p
}
