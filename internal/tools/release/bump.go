package main

import (
	"context"
	"flag"
	"time"

	"github.com/google/subcommands"

	"github.com/mmcloughlin/addchain/internal/cli"
)

// bump subcommand.
type bump struct {
	cli.Command

	varsfile    VarsFile
	releasedate string
}

func (*bump) Name() string     { return "bump" }
func (*bump) Synopsis() string { return "bump version" }
func (*bump) Usage() string {
	return `Usage: bump <version>

Bump version and update related files.

`
}

func (cmd *bump) SetFlags(f *flag.FlagSet) {
	cmd.varsfile.SetFlags(f)
	f.StringVar(&cmd.releasedate, "date", time.Now().UTC().Format("2006-01-02"), "release date")
}

func (cmd *bump) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Read arguments.
	if f.NArg() < 1 {
		return cmd.UsageError("missing version argument")
	}
	version := f.Arg(0)

	if err := ValidateVersion(version); err != nil {
		return cmd.Error(err)
	}

	// Set the version meta variable.
	if err := cmd.varsfile.Set("releaseversion", version); err != nil {
		return cmd.Error(err)
	}

	// Set the release date.
	if err := cmd.varsfile.Set("releasedate", cmd.releasedate); err != nil {
		return cmd.Error(err)
	}

	return subcommands.ExitSuccess
}
