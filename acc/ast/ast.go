package ast

type Chain struct {
	Statements []Statement
}

type Statement struct {
	Name Identifier
	Expr Expr
}

const (
	LowestPrec  = 0
	HighestPrec = 4
)

type Expr interface {
	Precedence() int
}

// Operand is an index into an addition chain.
type Operand int

func (Operand) Precedence() int { return HighestPrec }

// Identifier is a variable reference.
type Identifier string

func (Identifier) Precedence() int { return HighestPrec }

type Add struct {
	X, Y Expr
}

func (Add) Precedence() int { return 1 }

type Shift struct {
	X Expr
	S uint
}

func (Shift) Precedence() int { return 2 }

type Double struct {
	X Expr
}

func (Double) Precedence() int { return 3 }
