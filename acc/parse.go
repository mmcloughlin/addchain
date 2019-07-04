package acc

import (
	"strings"

	"github.com/mmcloughlin/addchain/acc/ast"
	"github.com/mmcloughlin/addchain/acc/internal/parser"
)

//go:generate pigeon -o internal/parser/zparser.go acc.peg

// String parses s.
func String(s string) (*ast.Chain, error) {
	r := strings.NewReader(s)
	i, err := parser.ParseReader("string", r)
	if err != nil {
		return nil, err
	}
	return i.(*ast.Chain), nil
}
