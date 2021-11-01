package pass

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/mmcloughlin/addchain/acc/eval"
	"github.com/mmcloughlin/addchain/acc/ir"
	"github.com/mmcloughlin/addchain/acc/rand"
	"github.com/mmcloughlin/addchain/internal/bigint"
	"github.com/mmcloughlin/addchain/internal/test"
)

func TestAllocator(t *testing.T) {
	// Generate a program that will need n live variables.
	n := 13
	p := &ir.Program{}

	// Create n indexes that will be used later.
	for i := 1; i <= n; i++ {
		p.AddInstruction(&ir.Instruction{
			Output: ir.Index(i),
			Op: ir.Add{
				X: ir.One,
				Y: ir.Index(i - 1),
			},
		})
	}

	// Consume the n indicies. This ensures we'll have n live variables at some
	// point.
	for i := n + 1; i <= 2*n; i++ {
		p.AddInstruction(&ir.Instruction{
			Output: ir.Index(i),
			Op: ir.Add{
				X: ir.Index(i - n),
				Y: ir.Index(i - 1),
			},
		})
	}

	// Execute allocation pass.
	a := Allocator{
		Input:  "in",
		Output: "out",
		Format: "tmp%d",
	}

	Allocate(t, p, a)

	// Every operand should have a name.
	for _, operand := range p.Operands {
		if operand.Identifier == "" {
			t.Errorf("operand %s does not have a name", operand)
		}
	}

	// We expect to have n-1 temporaries. Note we do not expect n, since the output
	// variable should be used as a temporary.
	expect := make([]string, n-1)
	for i := 0; i < n-1; i++ {
		expect[i] = "tmp" + strconv.Itoa(i)
	}
	if !reflect.DeepEqual(expect, p.Temporaries) {
		t.Fatalf("got temporaries %v; expect %v", p.Temporaries, expect)
	}

	// Confirm input and output variables are used.
	if p.Operands[0].Identifier != a.Input {
		t.Errorf("unexpected input name: got %q expect %q", p.Operands[0].Identifier, a.Input)
	}

	if p.Output().Identifier != a.Output {
		t.Errorf("unexpected output name: got %q expect %q", p.Output().Identifier, a.Output)
	}
}

func TestAllocatorRandom(t *testing.T) {
	checks := []func(*testing.T, *ir.Program, Allocator){
		CheckInputOutputNotBothLive,
		CheckExecute,
	}

	r := rand.AddsGenerator{N: 64}
	test.Repeat(t, func(t *testing.T) bool {
		t.Helper()

		// Generate a random program.
		p, err := r.GenerateProgram()
		if err != nil {
			t.Fatal(err)
		}

		// Execute allocation pass.
		a := Allocator{
			Input:  "in",
			Output: "out",
			Format: "t%d",
		}

		Allocate(t, p, a)

		// Checks.
		for _, check := range checks {
			check(t, p, a)
		}

		return true
	})
}

func CheckInputOutputNotBothLive(t *testing.T, p *ir.Program, a Allocator) {
	t.Helper()

	// Verify the input and output are not live at the same time.
	live := map[string]bool{}
	for i := len(p.Instructions) - 1; i >= 0; i-- {
		inst := p.Instructions[i]

		// Update live set.
		delete(live, inst.Output.Identifier)
		for _, input := range inst.Op.Inputs() {
			live[input.Identifier] = true
		}

		if live[a.Input] && live[a.Output] {
			t.Fatalf("instruction %d: input %q and output %q both live", i, a.Input, a.Output)
		}
	}
}

func CheckExecute(t *testing.T, p *ir.Program, a Allocator) {
	t.Helper()

	// Execute with interpreter. Deliberately setup the input and output to
	// be aliased.
	i := eval.NewInterpreter()
	x := bigint.One()
	i.Store(a.Input, x)
	i.Store(a.Output, x)

	if err := i.Execute(p); err != nil {
		t.Fatal(err)
	}

	// Output should be the same as the target of the addition chain.
	output, ok := i.Load(a.Output)
	if !ok {
		t.Fatalf("missing output variable %q", a.Output)
	}

	if err := Eval(p); err != nil {
		t.Fatal(err)
	}

	expect := p.Chain.End()

	if !bigint.Equal(output, expect) {
		t.Logf("   got = %#x", output)
		t.Logf("expect = %#x", expect)
		t.Fatal("mismatch")
	}
}

func Allocate(t *testing.T, p *ir.Program, a Allocator) {
	t.Helper()

	t.Logf("pre alloc:\n%s", p)

	if err := a.Execute(p); err != nil {
		t.Fatal(err)
	}

	for _, op := range p.Operands {
		t.Logf("alloc %d: %s", op.Index, op.Identifier)
	}

	t.Logf("post alloc:\n%s", p)
}
