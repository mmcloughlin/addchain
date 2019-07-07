// Command addchain generates addition chains.
package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/google/subcommands"
)

func main() {
	base := command{
		log: log.New(os.Stderr, "addchain: ", 0),
	}

	subcommands.Register(&search{command: base}, "")
	subcommands.Register(&eval{command: base}, "")
	subcommands.Register(subcommands.HelpCommand(), "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}

// command is a base for all subcommands.
type command struct {
	log *log.Logger
}

func (c command) SetFlags(f *flag.FlagSet) {}

func (c command) UsageError(format string, args ...interface{}) subcommands.ExitStatus {
	c.log.Printf(format, args...)
	return subcommands.ExitUsageError
}

func (c command) Fail(format string, args ...interface{}) subcommands.ExitStatus {
	c.log.Printf(format, args...)
	return subcommands.ExitFailure
}

func (c command) Error(err error) subcommands.ExitStatus {
	return c.Fail(err.Error())
}
