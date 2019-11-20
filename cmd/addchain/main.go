// Command addchain generates addition chains.
package main

import (
	"context"
	"flag"
	"io"
	"io/ioutil"
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

// OpenInput is a convenience for possibly opening an input file, or otherwise returning standard in.
func OpenInput(filename string) (string, io.ReadCloser, error) {
	if filename == "" {
		return "<stdin>", ioutil.NopCloser(os.Stdin), nil
	}
	f, err := os.Open(filename)
	return filename, f, err
}
