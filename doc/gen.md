# Output Generation

## Table of Contents

* [Template Reference](#template-reference)


## Template Reference

Templates are executed with the following data structure:

```go
type Data struct {
	Chain   addchain.Chain
	Ops     addchain.Program
	Script  *ast.Chain
	Program *ir.Program
}
```
