package ast

type Chain struct {
	Intermediates []Intermediate
	Result        Expr
}

// Operand is an index into an addition chain.
type Operand int

// Identifier is a variable reference.
type Identifier string

type Expr interface{}

type Add struct {
	X, Y Expr
}

type Shift struct {
	X Expr
	S uint
}

type Intermediate struct {
	Name Identifier
	Expr Expr
}
