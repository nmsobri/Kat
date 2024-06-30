package ast

import (
	"kat/token"
)

// Simulate Tagged Union
// This mean following nodes can either be Expression or Node
type Expression struct{}

func (e Expression) expr() {}
func (e Expression) node() {}

// #######################################################
// ################### Node Boolean ######################😀
// #######################################################
type NodeBoolean struct {
	Expression
	Token token.Token
	Value bool
}

// #######################################################
// ##################### Node Integer ####################😀
// #######################################################
type NodeInteger struct {
	Expression
	Token token.Token
	Value int64
}

// #######################################################
// ###################### Node Float #####################😀
// #######################################################
type NodeFloat struct {
	Expression
	Token token.Token
	Value float64
}

// #######################################################
// ##################### Node String #####################😀
// #######################################################
type NodeString struct {
	Expression
	Token token.Token
	Value string
}

// #######################################################
// #################### Node Index Expr ##################
// #######################################################
type NodeIndexExpr struct {
	Expression
	Token      token.Token
	Identifier Expr
	Index      Expr
}

// #######################################################
// ################### Node Prefix Expr ##################😀
// #######################################################
type NodePrefixExpr struct {
	Expression
	Token    token.Token
	Operator string
	Right    Expr
}

// #######################################################
// ################## Node Postfix Expr ##################😀
// #######################################################
type NodePostfixExpr struct {
	Expression
	Token    token.Token
	Left     Expr
	Operator string
}

// #######################################################
// #################### Node BinaryExpr ##################😀
// #######################################################
type NodeBinaryExpr struct {
	Expression
	Token    token.Token
	Left     Expr
	Right    Expr
	Operator string
}

// #######################################################
// ################### Node Function Call ################😀
// #######################################################
type NodeFunctionCall struct {
	Expression
	Token      token.Token
	Identifer  Expr
	Parameters []Expr
}

// #######################################################
// ################### Node Struct Expr ##################😀
// #######################################################
type NodeStructExpr struct {
	Expression
	Token  token.Token
	Name   Expr
	Values Expr
}

// #######################################################
// ################### Node Conditional ##################
// #######################################################
type NodeTernaryExpr struct {
	Expression
	Token     token.Token
	Condition Expr
	ThenArm   Expr
	ElseArm   Expr
}

// #######################################################
// ##################### Node Map Expr ###################
// #######################################################
type NodeMapExpr struct {
	Expression
	Token token.Token
	Map   map[Expr]Expr
}

// #######################################################
// #################### Node Array Expr ##################
// #######################################################
type NodeArrayExpr struct {
	Expression
	Token token.Token
	Value []Expr
}

// #######################################################
// ################### Node Identifier ###################😀
// #######################################################
type NodeIdentifier struct {
	Expression
	Token token.Token
	Name  string
}

// #######################################################
// ################### Node Import Expr ##################
// #######################################################
type NodeImportExpr struct {
	Expression
	Token token.Token
	Path  Expr
}
