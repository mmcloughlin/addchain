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

	// Execute an ensemble of algorithms.
	as := addchain.Ensemble()
	rs := addchain.Parallel(n, as)

	// Report results.
	best := 0
	for i, r := range rs {
		log.Printf("algorithm: %s", r.Algorithm)
		if r.Err != nil {
			log.Fatalf("error: %s", r.Err)
		}
		doubles, adds := r.Program.Count()
		log.Printf("total: %d\tdoubles: \t%d adds: %d", doubles+adds, doubles, adds)
		if len(r.Program) < len(rs[best].Program) {
			best = i
		}
	}

	// Details for the best chain.
	b := rs[best]
	for n, op := range b.Program {
		log.Printf("%3d:\t%d+%d\t%x", n+1, op.I, op.J, b.Chain[n+1])
	}
	log.Printf("best: %s", b.Algorithm)
}
