package pass

import (
	"fmt"
	"strconv"

	"github.com/mmcloughlin/addchain/acc/ir"
	"github.com/mmcloughlin/addchain/internal/errutil"
)

// References:
//
//	[curvechains]  Brian Smith. The Most Efficient Known Addition Chains for Field Element and
//	               Scalar Inversion for the Most Popular and Most Unpopular Elliptic Curves. 2017.
//	               https://briansmith.org/ecc-inversion-addition-chains-01 (accessed June 30, 2019)

// Naming schemes described in [curvechains].
var NameByteValues = NameBinaryValues(8, "_%b")

// NameBinaryValues assigns variable names to operands with values less than 2ᵏ.
// The identifier is determined from the format string, which should expect to
// take one *big.Int argument.
func NameBinaryValues(k int, format string) Interface {
	return Func(func(p *ir.Program) error {
		// We need canonical operands, and we need to know the chain values.
		if err := Exec(p, Func(CanonicalizeOperands), Func(Eval)); err != nil {
			return err
		}

		for _, operand := range p.Operands {
			// Skip if it already has a name.
			if operand.Identifier != "" {
				continue
			}

			// Fetch referenced value.
			idx := operand.Index
			if idx >= len(p.Chain) {
				return errutil.AssertionFailure("operand index %d out of bounds", idx)
			}
			x := p.Chain[idx]

			// Skip if too large.
			if x.BitLen() > k {
				continue
			}

			operand.Identifier = fmt.Sprintf(format, x)
		}

		return nil
	})
}

// NameByIndex builds a pass that sets any unnamed operands to have name prefix
// + index.
func NameByIndex(prefix string) Interface {
	return Func(func(p *ir.Program) error {
		if err := CanonicalizeOperands(p); err != nil {
			return err
		}
		for _, operand := range p.Operands {
			if operand.Identifier != "" {
				continue
			}
			operand.Identifier = prefix + strconv.Itoa(operand.Index)
		}
		return nil
	})
}
