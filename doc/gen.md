# Output Generation

## Table of Contents

* [Template Reference](#template-reference)


## Template Reference

Templates use Go [`text/template`](https://pkg.go.dev/text/template) syntax. The data structure passed
to the template is:

```go
type Data struct {
	Chain   addchain.Chain
	Ops     addchain.Program
	Script  *ast.Chain
	Program *ir.Program
	Meta    *meta.Properties
}
```

In addition to the [builtin functions](https://pkg.go.dev/text/template#hdr-Functions),
templates may use:


**`add`** `func(ir.Op) ir.Op`: If the input operation is an `ir.Add` then return it, otherwise return `nil`.

**`double`** `func(ir.Op) ir.Op`: If the input operation is an `ir.Double` then return it, otherwise return `nil`.

**`shift`** `func(ir.Op) ir.Op`: If the input operation is an `ir.Shift` then return it, otherwise return `nil`.

**`inc`** `func(int) int`: Increment an integer.

**`format`** `func(interface {}) (string, error)`: Formats an addition chain script (`*ast.Chain`) as a string.

**`split`** `func(string, string) []string`: Calls `strings.Split`.

**`join`** `func([]string, string) string`: Calls `strings.Join`.

**`lines`** `func(string) []string`: Split input string into lines.

