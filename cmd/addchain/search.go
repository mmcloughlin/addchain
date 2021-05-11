package main

import (
	"context"
	"flag"
	"runtime"

	"github.com/google/subcommands"

	"github.com/mmcloughlin/addchain/acc"
	"github.com/mmcloughlin/addchain/acc/printer"
	"github.com/mmcloughlin/addchain/alg/ensemble"
	"github.com/mmcloughlin/addchain/alg/exec"
	"github.com/mmcloughlin/addchain/internal/calc"
	"github.com/mmcloughlin/addchain/internal/cli"
	"github.com/mmcloughlin/profile"
)

// search subcommand.
type search struct {
	cli.Command

	concurrency int
	verbose     bool
}

func (*search) Name() string     { return "search" }
func (*search) Synopsis() string { return "search for an addition chain." }
func (*search) Usage() string {
	return `Usage: search [-v] [-p <N>] [-cpuprofile <file>] <expr>

Search for an addition chain for <expr>.

`
}

func (cmd *search) SetFlags(f *flag.FlagSet) {
	f.IntVar(&cmd.concurrency, "p", runtime.NumCPU(), "run `N` algorithms in parallel")
	f.BoolVar(&cmd.verbose, "v", false, "verbose output")
}

func (cmd *search) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) (status subcommands.ExitStatus) {
	// Enable profiling.
	defer profile.Start(
		profile.AllProfiles,
		profile.ConfigEnvVar("ADDCHAIN_PROFILE"),
	).Stop()

	// Parse arguments.
	if f.NArg() < 1 {
		return cmd.UsageError("missing expression")
	}
	expr := f.Arg(0)

	// Evaluate expression.
	cmd.Log.Printf("expr: %q", expr)

	n, err := calc.Eval(expr)
	if err != nil {
		return cmd.Fail("failed to evaluate %q: %s", expr, err)
	}

	cmd.Log.Printf("hex: %x", n)
	cmd.Log.Printf("dec: %s", n)

	// Execute an ensemble of algorithms.
	ex := exec.NewParallel()
	if cmd.verbose {
		ex.SetLogger(cmd.Log)
	}
	ex.SetConcurrency(cmd.concurrency)

	as := ensemble.Ensemble()
	rs := ex.Execute(n, as)

	// Report results.
	best := 0
	for i, r := range rs {
		cmd.Debugf("algorithm: %s", r.Algorithm)
		if r.Err != nil {
			return cmd.Fail("algorithm error: %v", r.Err)
		}
		doubles, adds := r.Program.Count()
		cmd.Debugf("total: %d\tdoubles: \t%d adds: %d", doubles+adds, doubles, adds)
		if len(r.Program) < len(rs[best].Program) {
			best = i
		}
	}

	// Details for the best chain.
	b := rs[best]
	for n, op := range b.Program {
		cmd.Debugf("[%3d] %3d+%3d\t%x", n+1, op.I, op.J, b.Chain[n+1])
	}
	cmd.Log.Printf("best: %s", b.Algorithm)

	// Produce a program for it.
	p, err := acc.Decompile(b.Program)
	if err != nil {
		return cmd.Error(err)
	}

	syntax, err := acc.Build(p)
	if err != nil {
		return cmd.Error(err)
	}

	if err := printer.Print(syntax); err != nil {
		return cmd.Error(err)
	}

	return subcommands.ExitSuccess
}

// Debugf prints a message in verbose mode only.
func (cmd *search) Debugf(format string, args ...interface{}) {
	if cmd.verbose {
		cmd.Log.Printf(format, args...)
	}
}
