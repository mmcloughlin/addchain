package acc

import (
	"testing"

	"github.com/mmcloughlin/addchain"
	"github.com/mmcloughlin/addchain/internal/test"
)

// CheckRandom runs the check function against randomly generated chain programs.
func CheckRandom(t *testing.T, check func(t *testing.T, p addchain.Program)) {
	gs := []addchain.Generator{
		addchain.RandomAddsGenerator{N: 10},
		addchain.NewRandomSolverGenerator(
			160,
			addchain.NewDictAlgorithm(
				addchain.SlidingWindow{K: 5},
				addchain.NewContinuedFractions(addchain.DichotomicStrategy{}),
			),
		),
	}

	for _, g := range gs {
		g := g // scopelint
		t.Run(g.String(), test.Trials(func(t *testing.T) bool {
			c, err := g.GenerateChain()
			if err != nil {
				t.Fatal(err)
			}

			p, err := c.Program()
			if err != nil {
				t.Fatal(err)
			}

			check(t, p)

			return true
		}))
	}
}
