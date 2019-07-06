// Command addchain generates addition chains.
package main

import (
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/mmcloughlin/addchain/acc"
	"github.com/mmcloughlin/addchain/acc/printer"

	"github.com/mmcloughlin/addchain"
	"github.com/mmcloughlin/addchain/internal/calc"
)

// l is the global logger.
var l = log.New(os.Stderr, "addchain: ", 0)

var concurrency = flag.Int("concurrency", runtime.NumCPU(), "number of algorithms to run concurrently")

func main() {
	// Parse command-line.
	flag.Parse()

	if flag.NArg() < 1 {
		l.Fatal("usage: addchain expr")
	}
	expr := flag.Arg(0)

	// Evaluate expression.
	l.Printf("expr: %q", expr)

	n, err := calc.Eval(expr)
	if err != nil {
		l.Fatalf("failed to evaluate %q: %s", expr, err)
	}

	l.Printf("hex: %x", n)
	l.Printf("dec: %d", n)

	// Execute an ensemble of algorithms.
	ex := addchain.NewParallel()
	ex.SetLogger(l)
	ex.SetConcurrency(*concurrency)

	as := addchain.Ensemble()
	rs := ex.Execute(n, as)

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
		log.Printf("[%3d] %3d+%3d\t%x", n+1, op.I, op.J, b.Chain[n+1])
	}
	log.Printf("best: %s", b.Algorithm)

	// Produce a program for it.
	p, err := acc.Decompile(b.Program)
	if err != nil {
		log.Fatal(err)
	}

	syntax, err := acc.Build(p)
	if err != nil {
		log.Fatal(err)
	}

	if err := printer.Print(syntax); err != nil {
		log.Fatal(err)
	}
}
