package acc

import (
	"reflect"
	"testing"

	"github.com/mmcloughlin/addchain"
	"github.com/mmcloughlin/addchain/internal/test"
)

func TestDecompileRandomRoundTrip(t *testing.T) {
	gs := []addchain.Generator{
		addchain.RandomAddsGenerator{N: 10},
		addchain.NewRandomSolverGenerator(
			160,
			addchain.NewDictAlgorithm(
				addchain.SlidingWindow{K: 4},
				addchain.NewContinuedFractions(addchain.DichotomicStrategy{}),
			),
		),
	}
	for _, g := range gs {
		t.Run(g.String(), test.Trials(func(t *testing.T) bool {
			c, err := g.GenerateChain()
			if err != nil {
				t.Fatal(err)
			}

			expect, err := c.Program()
			if err != nil {
				t.Fatal(err)
			}

			p, err := Decompile(expect)
			if err != nil {
				t.Fatal(err)
			}

			got, err := Compile(p)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(expect, got) {
				t.Fatal("roundtrip failed")
			}

			return true
		}))
	}
}
