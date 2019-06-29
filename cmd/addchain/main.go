// Command addchain generates addition chains.
package main

import (
	"log"
	"os"

	"github.com/mmcloughlin/addchain"
	"github.com/mmcloughlin/addchain/internal/calc"
)

func main() {
	log.SetPrefix("addchain: ")
	log.SetFlags(0)

	// Parse command-line.
	if len(os.Args) < 2 {
		log.Fatal("usage: addchain expr")
	}
	expr := os.Args[1]

	// Evaluate expression.
	log.Printf("expr: %q", expr)

	n, err := calc.Eval(expr)
	if err != nil {
		log.Fatalf("failed to evaluate %q: %s", expr, err)
	}

	log.Printf("n: %s", n)

	// Initialize algorithm.
	a := addchain.NewDictAlgorithm(
		addchain.SlidingWindow{K: 4},
		addchain.NewContinuedFractions(addchain.DichotomicStrategy{}),
	)
	log.Printf("algorithm: %s", a)

	// Apply and report.
	c, err := a.FindChain(n)
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Produces(n); err != nil {
		log.Fatal(err)
	}

	for i, x := range c {
		log.Printf("%4d: %s", i, x)
	}

	// Analyze program.
	p, err := c.Program()
	if err != nil {
		log.Fatal(err)
	}

	doubles, adds := p.Count()
	log.Printf("total: %d", doubles+adds)
	log.Printf("doubles: %d", doubles)
	log.Printf("adds: %d", adds)
}
