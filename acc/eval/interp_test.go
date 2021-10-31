package eval

import (
	"math/big"
	"testing"

	"github.com/mmcloughlin/addchain/acc/ir"
	"github.com/mmcloughlin/addchain/internal/bigint"
)

func TestInterpreter(t *testing.T) {
	// Construct a test input program with named operands. The result should be
	// 0b1111.
	a := ir.NewOperand("a", 0)
	b := ir.NewOperand("b", 1)
	c := ir.NewOperand("c", 2)
	d := ir.NewOperand("d", 4)
	e := ir.NewOperand("e", 5)
	p := &ir.Program{
		Instructions: []*ir.Instruction{
			{Output: b, Op: ir.Double{X: a}},
			{Output: c, Op: ir.Add{X: a, Y: b}},
			{Output: d, Op: ir.Shift{X: c, S: 2}},
			{Output: e, Op: ir.Add{X: d, Y: c}},
		},
	}

	t.Logf("program:\n%s", p)

	// Evaluate it.
	i := NewInterpreter()
	i.Store("a", big.NewInt(1))
	err := i.Execute(p)
	if err != nil {
		t.Fatal(err)
	}

	// Check variable values.
	expect := map[string]int64{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 12,
		"e": 15,
	}
	for v, x := range expect {
		got, ok := i.Load(v)
		if !ok {
			t.Fatalf("missing value for %q", v)
		}
		if !bigint.EqualInt64(got, x) {
			t.Errorf("got %s=%v; expect %v", v, got, x)
		}
	}
}
