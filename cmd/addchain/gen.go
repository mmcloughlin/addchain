package main

import (
	"context"
	"flag"
	"io"
	"os"
	"text/template"

	"github.com/google/subcommands"

	"github.com/mmcloughlin/addchain/acc"
	"github.com/mmcloughlin/addchain/acc/ir"
	"github.com/mmcloughlin/addchain/acc/parse"
	"github.com/mmcloughlin/addchain/acc/pass"
	"github.com/mmcloughlin/addchain/internal/cli"
)

// generate subcommand.
type generate struct {
	cli.Command
}

func (*generate) Name() string     { return "gen" }
func (*generate) Synopsis() string { return "generate output from an addition chain program" }
func (*generate) Usage() string {
	return `Usage: gen [<filename>]

Generate output from an addition chain program.

`
}

func (cmd *generate) SetFlags(f *flag.FlagSet) {
}

func (cmd *generate) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) (status subcommands.ExitStatus) {
	// Read input.
	input, r, err := cli.OpenInput(f.Arg(0))
	if err != nil {
		return cmd.Error(err)
	}
	defer cmd.CheckClose(&status, r)

	// Parse to a syntax tree.
	s, err := parse.Reader(input, r)
	if err != nil {
		return cmd.Error(err)
	}

	// Translate to IR.
	p, err := acc.Translate(s)
	if err != nil {
		return cmd.Error(err)
	}

	// Perform temporary variable allocation.
	alloc := pass.Allocator{
		Input:  "x",
		Output: "z",
		Format: "t%d",
	}
	if err := alloc.Execute(p); err != nil {
		return cmd.Error(err)
	}

	// Generate output.
	// ast.Print(token.NewFileSet(), p)
	// 	tmpl := `
	// {{- range $i := .Program.Instructions -}}
	// {{- with add $i.Op -}}
	// {{ printf "add\t%s\t%s\t%s" $i.Output .X .Y }}
	// {{ end -}}
	//
	// {{- with double $i.Op -}}
	// {{ printf "double\t%s\t%s" $i.Output .X }}
	// {{ end -}}
	//
	// {{- with shift $i.Op -}}
	// {{ printf "shift\t%s\t%s\t%d" $i.Output .X .S }}
	// {{ end -}}
	// {{- end -}}
	// `

	// ec3:
	tmpl := `
{{- range $i := .Program.Instructions -}}
// Step {{ $i.Output.Index }}: {{ . }}
{{- with add $i.Op }}
Mul({{ $i.Output }}, {{ .X }}, {{ .Y }})
{{ end -}}

{{- with double $i.Op }}
Sqr({{ $i.Output }}, {{ .X }})
{{ end -}}

{{- with shift $i.Op -}}
{{- $first := 0 -}}
{{- if ne $i.Output.Identifier .X.Identifier }}
Sqr({{ $i.Output }}, {{ .X }})
{{- $first = 1 -}}
{{- end }}
for s := {{ $first }}; s < {{ .S }}; s++ {
	Sqr({{ $i.Output }}, {{ $i.Output }})
}
{{ end }}
{{ end -}}
`
	if err := Generate(os.Stdout, tmpl, p); err != nil {
		return cmd.Error(err)
	}

	return subcommands.ExitSuccess
}

func Generate(w io.Writer, tmpl string, p *ir.Program) error {
	// Custom template functions.
	funcs := template.FuncMap{
		"add": func(op ir.Op) ir.Op {
			if a, ok := op.(ir.Add); ok {
				return a
			}
			return nil
		},
		"double": func(op ir.Op) ir.Op {
			if d, ok := op.(ir.Double); ok {
				return d
			}
			return nil
		},
		"shift": func(op ir.Op) ir.Op {
			if s, ok := op.(ir.Shift); ok {
				return s
			}
			return nil
		},
	}

	// Parse template.
	t, err := template.New("").Funcs(funcs).Parse(tmpl)
	if err != nil {
		return err
	}

	// Prepare template data.
	type data struct {
		Program *ir.Program
	}
	d := data{Program: p}

	// Execute.
	return t.Execute(w, d)
}
