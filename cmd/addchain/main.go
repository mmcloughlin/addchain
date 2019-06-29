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
	for _, r := range rs {
		log.Printf("algorithm: %s", r.Algorithm)
		if r.Err != nil {
			log.Printf("error: %s", r.Err)
			continue
		}
		doubles, adds := r.Program.Count()
		log.Printf("total: %d\tdoubles: \t%d adds: %d", doubles+adds, doubles, adds)
	}
}
