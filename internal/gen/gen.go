package gen

import (
	"io"
	"text/template"

	"github.com/mmcloughlin/addchain/acc/ir"
)

type Data struct {
	Program *ir.Program
}

func BuildData(p *ir.Program) *Data {
	return &Data{
		Program: p,
	}
}

func Generate(w io.Writer, tmpl string, d *Data) error {
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

	// Execute.
	return t.Execute(w, d)
}
