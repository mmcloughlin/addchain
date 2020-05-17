package results

import (
	"flag"
	"log"
	"os"
	"testing"

	"github.com/mmcloughlin/addchain/acc"
	"github.com/mmcloughlin/addchain/alg/ensemble"
	"github.com/mmcloughlin/addchain/alg/exec"
	"github.com/mmcloughlin/addchain/internal/test"
)

// verbose flag for customizing log output from algorithm executor.
var verbose = flag.Bool("verbose", false, "enable verbose logging")

func TestResults(t *testing.T) {
	as := ensemble.Ensemble()
	for _, c := range Results {
		c := c // scopelint
		t.Run(c.Slug, func(t *testing.T) {
			// Tests with a best known result are prioritized. Only run all tests in
			// stress test mode.
			if c.BestKnown == 0 {
				test.RequireStress(t)
			} else {
				test.RequireLong(t)
			}

			n := c.Target()
			t.Logf("n-%d=%x", c.D, n)

			// Execute.
			ex := exec.NewParallel()
			if *verbose {
				ex.SetLogger(log.New(os.Stderr, "", log.Ltime|log.Lmicroseconds))
			}
			rs := ex.Execute(n, as)

			// Process results.
			best := 0
			for i, r := range rs {
				if r.Err != nil {
					t.Logf("error with %s", r.Algorithm)
					t.Error(r.Err)
					continue
				}
				if len(r.Program) <= len(rs[best].Program) {
					best = i
				}
			}

			// Report the best.
			b := rs[best]
			t.Logf(" best: %s", b.Algorithm)
			t.Logf("total: %d", len(b.Program))

			if c.BestKnown > 0 {
				t.Logf("known: %d", c.BestKnown)
				t.Logf("delta: %+d", len(b.Program)-c.BestKnown)
			}

			// Ensure the recorded best length is correct.
			if c.Length != len(b.Program) {
				t.Errorf("incorrect best value %d; expect %d", c.Length, len(b.Program))
			}

			// Compare to golden file.
			golden, err := acc.LoadFile(test.GoldenName(c.Slug))
			if err != nil {
				t.Logf("failed to load golden file: %s", err)
			}

			save := test.Golden()
			switch {
			case golden == nil:
				t.Errorf("missing golden file")
				save = true
			case len(golden.Program) < len(b.Program):
				t.Errorf("regression from golden: %d to %d", len(golden.Program), len(b.Program))
			case len(golden.Program) > len(b.Program):
				t.Logf("improvement: %d to %d", len(golden.Program), len(b.Program))
				save = true
			}

			if save {
				t.Log("saving golden file")

				r, err := acc.Decompile(b.Program)
				if err != nil {
					t.Fatal(err)
				}

				if err := acc.Save(test.GoldenName(c.Slug), r); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
