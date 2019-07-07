package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"

	"github.com/mmcloughlin/addchain/acc/parse"
	"github.com/mmcloughlin/addchain/acc/printer"
)

// format subcommand.
type format struct {
	command
}

func (*format) Name() string     { return "fmt" }
func (*format) Synopsis() string { return "format an addition chain script" }
func (*format) Usage() string {
	return `Usage: fmt [<filename>]

Format an addition chain script.

 `
}

func (cmd *format) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Read input.
	input, r, err := OpenInput(f.Arg(0))
	if err != nil {
		return cmd.Error(err)
	}
	defer r.Close()

	// Parse to a syntax tree.
	s, err := parse.Reader(input, r)
	if err != nil {
		return cmd.Error(err)
	}

	// Print.
	if err := printer.Print(s); err != nil {
		return cmd.Error(err)
	}

	return subcommands.ExitSuccess
}
