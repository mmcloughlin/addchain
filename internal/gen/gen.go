package gen

import (
	"io"
	"text/template"

	"github.com/mmcloughlin/addchain"
	"github.com/mmcloughlin/addchain/acc"
	"github.com/mmcloughlin/addchain/acc/ast"
	"github.com/mmcloughlin/addchain/acc/ir"
	"github.com/mmcloughlin/addchain/acc/pass"
	"github.com/mmcloughlin/addchain/meta"
)

// Data provided to templates.
type Data struct {
	Chain   addchain.Chain
	Ops     addchain.Program
	Script  *ast.Chain
	Program *ir.Program
	Meta    *meta.Properties
}

// Config for template input generation.
type Config struct {
	// Allocator for temporary variables. This configuration determines variable
	// naming.
	Allocator pass.Allocator
}

// PrepareData builds input template data for the given addition chain script.
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

// Generate
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
