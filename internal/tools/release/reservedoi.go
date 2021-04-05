package main

import (
	"context"
	"flag"
	"strconv"

	"github.com/google/subcommands"

	"github.com/mmcloughlin/addchain/internal/cli"
)

// reservedoi subcommand.
type reservedoi struct {
	cli.Command

	httpclient HTTPClient
	zenodo     Zenodo
	varsfile   VarsFile
}

func (*reservedoi) Name() string     { return "reservedoi" }
func (*reservedoi) Synopsis() string { return "reserve zenodo doi for a new version" }
func (*reservedoi) Usage() string {
	return `Usage: reservedoi

Reserve DOI on Zenodo for a new release.

`
}

func (cmd *reservedoi) SetFlags(f *flag.FlagSet) {
	cmd.httpclient.SetFlags(f)
	cmd.zenodo.SetFlags(f)
	cmd.varsfile.SetFlags(f)
}

func (cmd *reservedoi) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Zenodo client.
	httpclient, err := cmd.httpclient.Client()
	if err != nil {
		return cmd.Error(err)
	}

	c, err := cmd.zenodo.Client(httpclient)
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

	// Fetch the new version.
	d, err := c.DepositionRetrieve(ctx, newid)
	if err != nil {
		return cmd.Error(err)
	}

	// Write back to variables file.
	if err := cmd.varsfile.Set("zenodoid", strconv.Itoa(d.ID)); err != nil {
		return cmd.Error(err)
	}

	if err := cmd.varsfile.Set("doi", d.DOI); err != nil {
		return cmd.Error(err)
	}

	return subcommands.ExitSuccess
}
