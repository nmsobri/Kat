package ast

type Node interface {
	node()
	GetLocation() Location
}

type Stmt interface {
	Node
	stmt()
}

type Expr interface {
	Node
	expr()
}

type Location struct {
	Line   int
	Column int
}
