package metavars

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"strconv"
	"strings"

	"github.com/mmcloughlin/addchain/internal/print"
)

// TODO: https://play.golang.org/p/sJ5QPrXB0eY

type Property struct {
	Name  string
	Doc   string
	Value string
}

type File struct {
	Package    string
	Properties []Property
}

func Write(w io.Writer, f *File) error {
	// Write source code.
	buf := bytes.NewBuffer(nil)
	p := print.New(buf)

	p.Linef("package %s", f.Package)
	p.NL()
	p.Linef("var (")
	p.Indent()
	for _, prop := range f.Properties {
		if prop.Doc != "" {
			p.Linef("// %s", prop.Doc)
		}
		p.Linef("%s = %q", prop.Name, prop.Value)
	}
	p.Dedent()
	p.Linef(")")

	if err := p.Error(); err != nil {
		return err
	}

	// Write to file.
	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = w.Write(src)
	return err
}

func Read(r io.Reader) (*File, error) {
	return parse("", r)
}

func ReadFile(filename string) (*File, error) {
	return parse(filename, nil)
}

func parse(filename string, src interface{}) (*File, error) {
	// Parse AST.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Build file from AST.
	b := &builder{f: &File{}}
	b.file(f)

	if b.err != nil {
		return nil, b.err
	}

	return b.f, nil
}

type builder struct {
	f   *File
	err error
}

func (b *builder) file(f *ast.File) {
	// Package name.
	b.f.Package = f.Name.Name

	// Declarations.
	for _, d := range f.Decls {
		b.decl(d)
	}
}

func (b *builder) decl(i ast.Decl) {
	switch d := i.(type) {
	case *ast.GenDecl:
		b.gendecl(d)
	default:
		b.seterror("unexpected declaration type %T", i)
	}
}

func (b *builder) gendecl(g *ast.GenDecl) {
	// Confirm this is a variable declaration.
	if g.Tok != token.VAR {
		b.seterror("expected variable declaration, got %s", g.Tok)
		return
	}

	// Process each spec.
	for _, spec := range g.Specs {
		v, ok := spec.(*ast.ValueSpec)
		if !ok {
			b.seterror("expected value spec")
		}
		b.valuespec(v)
	}
}

func (b *builder) valuespec(v *ast.ValueSpec) {
	p := Property{}

	// Extract name.
	if len(v.Names) != 1 {
		b.seterror("expected one name, got %d", len(v.Names))
		return
	}
	p.Name = v.Names[0].Name

	// Extract value.
	if len(v.Values) != 1 {
		b.seterror("expected one value, got %d", len(v.Values))
		return
	}

	expr := v.Values[0]
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		b.seterror("expected basic literal")
		return
	}

	if lit.Kind != token.STRING {
		b.seterror("expected string literal")
		return
	}

	s, err := strconv.Unquote(lit.Value)
	if err != nil {
		b.seterror("parse string literal: %w", err)
		return
	}

	p.Value = s

	// Optional documentation.
	if v.Doc != nil {
		const leader = "// "
		if len(v.Doc.List) != 1 {
			b.seterror("expect single line documentation")
			return
		}
		comment := v.Doc.List[0]
		if !strings.HasPrefix(comment.Text, leader) {
			b.seterror("expected documentation to have %q leader", leader)
			return
		}
		p.Doc = strings.TrimPrefix(comment.Text, leader)
	}

	// Add property to file.
	b.f.Properties = append(b.f.Properties, p)
}

func (b *builder) seterror(format string, args ...interface{}) {
	if b.err == nil {
		b.err = fmt.Errorf(format, args...)
	}
}

// *ast.File {
// 	1  .  Package: 1:1
// 	2  .  Name: *ast.Ident {
// 	3  .  .  NamePos: 1:9
// 	4  .  .  Name: "meta"
// 	5  .  }
// 	6  .  Decls: []ast.Decl (len = 1) {
// 	7  .  .  0: *ast.GenDecl {
// 	8  .  .  .  TokPos: 3:1
// 	9  .  .  .  Tok: var
// 10  .  .  .  Lparen: 3:5
// 11  .  .  .  Specs: []ast.Spec (len = 2) {
// 12  .  .  .  .  0: *ast.ValueSpec {
// 13  .  .  .  .  .  Doc: *ast.CommentGroup {
// 14  .  .  .  .  .  .  List: []*ast.Comment (len = 1) {
// 15  .  .  .  .  .  .  .  0: *ast.Comment {
// 16  .  .  .  .  .  .  .  .  Slash: 4:2
// 17  .  .  .  .  .  .  .  .  Text: "// ReleaseVersion is the version of the most recent release."
// 18  .  .  .  .  .  .  .  }
// 19  .  .  .  .  .  .  }
// 20  .  .  .  .  .  }
// 21  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
// 22  .  .  .  .  .  .  0: *ast.Ident {
// 23  .  .  .  .  .  .  .  NamePos: 5:2
// 24  .  .  .  .  .  .  .  Name: "ReleaseVersion"
// 25  .  .  .  .  .  .  .  Obj: *ast.Object {
// 26  .  .  .  .  .  .  .  .  Kind: var
// 27  .  .  .  .  .  .  .  .  Name: "ReleaseVersion"
// 28  .  .  .  .  .  .  .  .  Decl: *(obj @ 12)
// 29  .  .  .  .  .  .  .  .  Data: 0
// 30  .  .  .  .  .  .  .  }
// 31  .  .  .  .  .  .  }
// 32  .  .  .  .  .  }
// 33  .  .  .  .  .  Values: []ast.Expr (len = 1) {
// 34  .  .  .  .  .  .  0: *ast.BasicLit {
// 35  .  .  .  .  .  .  .  ValuePos: 5:19
// 36  .  .  .  .  .  .  .  Kind: STRING
// 37  .  .  .  .  .  .  .  Value: "\"0.2.0\""
// 38  .  .  .  .  .  .  }
// 39  .  .  .  .  .  }
// 40  .  .  .  .  }
