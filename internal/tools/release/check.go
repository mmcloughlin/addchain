package main

import (
	"context"
	"flag"
	"time"

	"github.com/google/subcommands"

	"github.com/mmcloughlin/addchain/internal/cli"
)

// check subcommand.
type check struct {
	cli.Command

	zenodo   Zenodo
	varsfile VarsFile
}

func (*check) Name() string     { return "check" }
func (*check) Synopsis() string { return "check release is ready to be tagged" }
func (*check) Usage() string {
	return `Usage: check

Perform some final checks to confirm release can be tagged.

`
}

func (cmd *check) SetFlags(f *flag.FlagSet) {
	cmd.zenodo.SetFlags(f)
	cmd.varsfile.SetFlags(f)
}

func (cmd *check) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Validate the version field.
	version, err := cmd.varsfile.Get("releaseversion")
	if err != nil {
		return cmd.Error(err)
	}

	if err := ValidateVersion(version); err != nil {
		return cmd.Error(err)
	}

	// Check the release date.
	releasedate, err := cmd.varsfile.Get("releasedate")
	if err != nil {
		return cmd.Error(err)
	}

	if _, err := time.Parse("2006-01-02", releasedate); err != nil {
		return cmd.Fail("release date should be in format YYYY-MM-DD")
	}

	// Check that a Zenodo deposit has been allocated.
	client, err := cmd.zenodo.Client()
	if err != nil {
		return cmd.Error(err)
	}

	id, err := cmd.varsfile.Get("zenodoid")
	if err != nil {
		return cmd.Error(err)
	}

	d, err := client.DepositionRetrieve(ctx, id)
	if err != nil {
		return cmd.Error(err)
	}

	if d.State != "inprogress" {
		return cmd.Fail("zenodo deposit in %q state", d.State)
	}

	if d.Submitted {
		return cmd.Fail("zenodo deposit has been published")
	}

	return subcommands.ExitSuccess
}
