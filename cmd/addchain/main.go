// Command addchain generates addition chains.
package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"

	"github.com/mmcloughlin/addchain/internal/cli"
	"github.com/mmcloughlin/addchain/internal/meta"
)

func main() {
	base := cli.NewBaseCommand("addchain")
	subcommands.Register(&search{Command: base}, "")
	subcommands.Register(&eval{Command: base}, "")
	subcommands.Register(&format{Command: base}, "")

	if meta.Meta.BuildVersion != "" {
		subcommands.Register(&version{version: meta.Meta.BuildVersion, Command: base}, "")
	}

	subcommands.Register(&cite{properties: meta.Meta, Command: base}, "")

	subcommands.Register(subcommands.HelpCommand(), "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
