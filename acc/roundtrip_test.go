package acc

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/mmcloughlin/addchain"
	"github.com/mmcloughlin/addchain/acc/parse"
	"github.com/mmcloughlin/addchain/acc/printer"
	"github.com/mmcloughlin/addchain/internal/test"
)

// The purpose of this test is to verify the full round trip from an addition
// chain program to source and back.
//
//            addchain.Chain		c
//                |   ^
//        Program |   | Evaluate
//                v   |
//           addchain.Program		p
//                |   ^
//      Decompile |   | Compile
//                v   |
//              ir.Program			r
//                |   ^
//          Build |   | Translate
//                v   |
//              ast.Chain			s
//                |   ^
//          Print |   | Parse
//                v   |
//              acc source			src

func TestRandomRoundTrip(t *testing.T) {
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
		t.Run(g.String(), test.Trials(func(t *testing.T) bool {
			// Generate a chain.
			c, err := g.GenerateChain()
			if err != nil {
				t.Fatal(err)
			}

			p, err := c.Program()
			if err != nil {
				t.Fatal(err)
			}

			// Decompile into IR.
			r, err := Decompile(p)
			if err != nil {
				t.Fatal(err)
			}

			// Build syntax tree.
			s, err := Build(r)
			if err != nil {
				t.Fatal(err)
			}

			// Print to source.
			var src bytes.Buffer
			err = printer.Fprint(&src, s)
			if err != nil {
				t.Fatal(err)
			}

			// Parse back to syntax.
			s2, err := parse.Reader("string", &src)
			if err != nil {
				t.Fatal(err)
			}

			// Translate back to IR.
			r2, err := Translate(s2)
			if err != nil {
				t.Fatal(err)
			}

			// Compile back to a program.
			p2, err := Compile(r2)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(p, p2) {
				t.Logf("p:\n%v", p)
				t.Logf("r:\n%v", r)
				t.Logf("src:\n%s", src.String())
				t.Logf("r2:\n%v", r2)
				t.Logf("p2:\n%v", p2)
				t.Fatal("roundtrip failure")
			}

			return true
		}))
	}
}
