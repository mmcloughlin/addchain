package main

import (
	"context"
	"flag"
	"runtime"

	"github.com/google/subcommands"

	"github.com/mmcloughlin/addchain"
	"github.com/mmcloughlin/addchain/acc"
	"github.com/mmcloughlin/addchain/acc/printer"
	"github.com/mmcloughlin/addchain/internal/calc"
)

// search subcommand.
type search struct {
	command

	concurrency int
}

func (*search) Name() string     { return "search" }
func (*search) Synopsis() string { return "search for an addition chain." }
func (*search) Usage() string {
	return `Usage: search [-p <N>] <expr>

Search for an addition chain for <expr>.

 `
}

func (s *search) SetFlags(f *flag.FlagSet) {
	f.IntVar(&s.concurrency, "p", runtime.NumCPU(), "number of algorithms to run in parallel")
}

func (s *search) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() < 1 {
		return s.UsageError("missing expression")
	}
	expr := f.Arg(0)

	// Evaluate expression.
	s.log.Printf("expr: %q", expr)

	n, err := calc.Eval(expr)
	if err != nil {
		return s.Fail("failed to evaluate %q: %s", expr, err)
	}

	s.log.Printf("hex: %x", n)
	s.log.Printf("dec: %s", n)

	// Execute an ensemble of algorithms.
	ex := addchain.NewParallel()
	ex.SetLogger(s.log)
	ex.SetConcurrency(s.concurrency)

	as := addchain.Ensemble()
	rs := ex.Execute(n, as)

	// Report results.
	best := 0
	for i, r := range rs {
		s.log.Printf("algorithm: %s", r.Algorithm)
		if r.Err != nil {
			return s.Error(err)
		}
		doubles, adds := r.Program.Count()
		s.log.Printf("total: %d\tdoubles: \t%d adds: %d", doubles+adds, doubles, adds)
		if len(r.Program) < len(rs[best].Program) {
			best = i
		}
	}

	// Details for the best chain.
	b := rs[best]
	for n, op := range b.Program {
		s.log.Printf("[%3d] %3d+%3d\t%x", n+1, op.I, op.J, b.Chain[n+1])
	}
	s.log.Printf("best: %s", b.Algorithm)

	// Produce a program for it.
	p, err := acc.Decompile(b.Program)
	if err != nil {
		return s.Error(err)
	}

	syntax, err := acc.Build(p)
	if err != nil {
		return s.Error(err)
	}

	if err := printer.Print(syntax); err != nil {
		return s.Error(err)
	}

	return subcommands.ExitSuccess
}
