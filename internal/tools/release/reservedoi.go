package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"

	"github.com/mmcloughlin/addchain/internal/cli"
)

// reservedoi subcommand.
type reservedoi struct {
	cli.Command

	zenodo   Zenodo
	varsfile VarsFile
}

func (*reservedoi) Name() string     { return "reservedoi" }
func (*reservedoi) Synopsis() string { return "reserve zenodo doi for a new version" }
func (*reservedoi) Usage() string {
	return `Usage: reservedoi

Reserve DOI on Zenodo for a new release.

`
}

func (cmd *reservedoi) SetFlags(f *flag.FlagSet) {
	cmd.zenodo.SetFlags(f)
	cmd.varsfile.SetFlags(f)
}

func (cmd *reservedoi) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Zenodo client.
	c, err := cmd.zenodo.Client()
	if err != nil {
		return cmd.Error(err)
	}

	// Fetch existing Zenodo ID.
	id, err := cmd.varsfile.Get("zenodoid")
	if err != nil {
		return cmd.Error(err)
	}

	cmd.Log.Printf("current zenodo id %s", id)

	// Start a new deposit.
	newid, err := c.DepositionNewVersion(ctx, id)
	if err != nil {
		return cmd.Error(err)
	}

	cmd.Log.Printf("new version id %s", newid)

	// Write back to variables file.
	if err := cmd.varsfile.Set("zenodoid", newid); err != nil {
		return cmd.Error(err)
	}

	return subcommands.ExitSuccess
}
