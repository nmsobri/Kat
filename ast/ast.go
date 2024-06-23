package ast

type Node interface {
	node()
}

type Stmt interface {
	Node
	stmt()
}

type Expr interface {
	Node
	expr()
}
