package gen

import (
	"bufio"
	"io"
	"reflect"
	"strings"
	"text/template"

	"github.com/mmcloughlin/addchain"
	"github.com/mmcloughlin/addchain/acc"
	"github.com/mmcloughlin/addchain/acc/ast"
	"github.com/mmcloughlin/addchain/acc/ir"
	"github.com/mmcloughlin/addchain/acc/pass"
	"github.com/mmcloughlin/addchain/acc/printer"
	"github.com/mmcloughlin/addchain/meta"
)

type Data struct {
	Chain   addchain.Chain
	Ops     addchain.Program
	Script  *ast.Chain
	Program *ir.Program
	Meta    *meta.Properties
}

type Config struct {
	Allocator pass.Allocator
}

func PrepareData(cfg Config, s *ast.Chain) (*Data, error) {
	// Translate to IR.
	p, err := acc.Translate(s)
	if err != nil {
		return nil, err
	}

	// Apply processing passes: temporary variable allocation, and computing the
	// full addition chain sequence and operations.
	if err := pass.Exec(p, cfg.Allocator, pass.Func(pass.Eval)); err != nil {
		return nil, err
	}

	return &Data{
		Chain:   p.Chain,
		Ops:     p.Program,
		Script:  s,
		Program: p,
		Meta:    meta.Meta,
	}, nil
}

type Function struct {
	Name        string
	Description string
	Func        interface{}
}

func (f *Function) Signature() string {
	return reflect.ValueOf(f.Func).Type().String()
}

var Functions = []*Function{
	{
		Name:        "add",
		Description: "If the input operation is an `ir.Add` then return it, otherwise return `nil`.",
		Func: func(op ir.Op) ir.Op {
			if a, ok := op.(ir.Add); ok {
				return a
			}
			return nil
		},
	},

	{
		Name:        "double",
		Description: "If the input operation is an `ir.Double` then return it, otherwise return `nil`.",
		Func: func(op ir.Op) ir.Op {
			if d, ok := op.(ir.Double); ok {
				return d
			}
			return nil
		},
	},

	{
		Name:        "shift",
		Description: "If the input operation is an `ir.Shift` then return it, otherwise return `nil`.",
		Func: func(op ir.Op) ir.Op {
			if s, ok := op.(ir.Shift); ok {
				return s
			}
			return nil
		},
	},

	{
		Name:        "inc",
		Description: "Increment an integer.",
		Func:        func(n int) int { return n + 1 },
	},

	{
		Name:        "format",
		Description: "Formats an addition chain script (`*ast.Chain`) as a string.",
		Func:        printer.String,
	},

	{Name: "split", Description: "Calls `strings.Split`.", Func: strings.Split},
	{Name: "join", Description: "Calls `strings.Join`.", Func: strings.Join},

	{
		Name:        "lines",
		Description: "Split input string into lines.",
		Func: func(s string) []string {
			var lines []string
			scanner := bufio.NewScanner(strings.NewReader(s))
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			return lines
		},
	},
}

func Generate(w io.Writer, tmpl string, d *Data) error {
	// Custom template functions.
	funcs := template.FuncMap{}
	for _, f := range Functions {
		funcs[f.Name] = f.Func
	}

	// Parse template.
	t, err := template.New("").Funcs(funcs).Parse(tmpl)
	if err != nil {
		return err
	}

	// Execute.
	return t.Execute(w, d)
}
