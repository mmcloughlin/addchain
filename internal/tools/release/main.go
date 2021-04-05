package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"

	"github.com/mmcloughlin/addchain/internal/cli"
)

func main() {
	base := cli.NewBaseCommand("release")
	subcommands.Register(&bump{Command: base}, "")
	subcommands.Register(&reservedoi{Command: base}, "")
	subcommands.Register(&check{Command: base}, "")
	subcommands.Register(&upload{Command: base}, "")
	subcommands.Register(subcommands.HelpCommand(), "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
