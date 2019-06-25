package addchain

import "math/big"

// Binary uses the binary method to produce an addition chain for n.
func Binary(n *big.Int) Chain {
	var zero big.Int
	c := Chain{}
	s := -1
	d := 0
	b := new(big.Int).Set(n)
	for {
		if b.Bit(0) == 1 {
			if s >= 0 {
				s = c.Add(s, d)
			} else {
				s = d
			}
		}
		b.Rsh(b, 1)
		if b.Cmp(&zero) == 0 {
			break
		}
		d = c.Double(d)
	}
	return c
}
