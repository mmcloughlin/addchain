package acc

import (
	"testing"

	"github.com/mmcloughlin/addchain/acc/eval"
	"github.com/mmcloughlin/addchain/acc/pass"
	"github.com/mmcloughlin/addchain/alg/ensemble"
	"github.com/mmcloughlin/addchain/alg/exec"
	"github.com/mmcloughlin/addchain/internal/bigint"
	"github.com/mmcloughlin/addchain/internal/results"
	"github.com/mmcloughlin/addchain/internal/test"
)

func TestResults(t *testing.T) {
	timer := test.Start()
	as := ensemble.Ensemble()
	for _, c := range results.Results {
		c := c // scopelint
		t.Run(c.Slug, func(t *testing.T) {
			timer.Check(t)

			// Execute.
			ex := exec.NewParallel()
			rs := ex.Execute(c.Target(), as)

			// Check each result.
			for _, r := range rs {
				if r.Err != nil {
					t.Fatalf("error with %s: %v", r.Algorithm, r.Err)
				}

				// Decompile into IR.
				p, err := Decompile(r.Program)
				if err != nil {
					t.Fatal(err)
				}

				// Allocate.
				a := pass.Allocator{
					Input:  "x",
					Output: "z",
					Format: "t%d",
				}
				if err := a.Execute(p); err != nil {
					t.Fatal(err)
				}

				// Evaluate.
				i := eval.NewInterpreter()
				x := bigint.One()
				i.Store(a.Input, x)
				i.Store(a.Output, x)

				if err := i.Execute(p); err != nil {
					t.Fatal(err)
				}

				// Verify output.
				output, ok := i.Load(a.Output)
				if !ok {
					t.Fatalf("missing output variable %q", a.Output)
				}

				expect := r.Chain.End()

				if !bigint.Equal(output, expect) {
					t.FailNow()
				}

			}
		})
	}
}
