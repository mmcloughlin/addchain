// Command addchain generates addition chains.
package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"

	"github.com/mmcloughlin/addchain/internal/cli"
)

func main() {
	base := cli.NewBaseCommand("addchain")
	subcommands.Register(&search{Command: base}, "")
	subcommands.Register(&eval{Command: base}, "")
	subcommands.Register(&format{Command: base}, "")
	subcommands.Register(subcommands.HelpCommand(), "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
