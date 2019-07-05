package ast

type Chain struct {
	Statements []Statement
}

type Statement struct {
	Name Identifier
	Expr Expr
}

type Expr interface{}

// Operand is an index into an addition chain.
type Operand int

// Identifier is a variable reference.
type Identifier string

type Add struct {
	X, Y Expr
}

type Shift struct {
	X Expr
	S uint
}

type Double struct {
	X Expr
}
