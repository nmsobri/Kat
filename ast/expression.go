package ast

import (
	"kat/token"
)

// #######################################################
// ################### Node Boolean #####################
// #######################################################
type NodeBoolean struct {
	Token token.Token
	Value bool
}

func (nb NodeBoolean) expr() {}

// #######################################################
// ##################### Node Integer ####################
// #######################################################
type NodeInteger struct {
	Token token.Token
	Value int64
}

func (ni NodeInteger) expr() {}

// #######################################################
// ###################### Node Double ####################
// #######################################################
type NodeDouble struct {
	Token token.Token
	Value float64
}

func (nd NodeDouble) expr() {}

// #######################################################
// ##################### Node String #####################
// #######################################################
type NodeString struct {
	Token token.Token
	Value string
}

func (ns NodeString) expr() {}

// #######################################################
// #################### Node Index Expr ##################
// #######################################################
type NodeIndexExpr struct {
	Token      token.Token
	Identifier Expr
	Index      Expr
}

func (nie NodeIndexExpr) expr() {}

// #######################################################
// ################### Node Prefix Expr ##################
// #######################################################
type NodePrefixExpr struct {
	Token    token.Token
	Operator string
	Right    Expr
}

func (npe NodePrefixExpr) expr() {}

// #######################################################
// ################## Node Postfix Expr ##################
// #######################################################
type NodePostfixExpr struct {
	Token    token.Token
	Left     Expr
	Operator string
}

func (npe NodePostfixExpr) expr() {}

// #######################################################
// #################### Node BinaryExpr ##################
// #######################################################
type NodeBinaryExpr struct {
	Token    token.Token
	Left     Expr
	Right    Expr
	Operator string
}

func (nbe NodeBinaryExpr) expr() {}

// #######################################################
// ################### Node Function Call ################
// #######################################################
type NodeFunctionCall struct {
	Token token.Token
	Left  Expr
	Right []Expr
}

func (nfc NodeFunctionCall) expr() {}

// #######################################################
// ################### Node Struct Expr ##################
// #######################################################
type NodeStructExpr struct {
	Token      token.Token
	Identifier Expr
	Values     Expr
}

func (ns NodeStructExpr) expr() {}

// #######################################################
// ################### Node Conditional ##################
// #######################################################
type NodeTernaryExpr struct {
	Token     token.Token
	Condition Expr
	ThenArm   Expr
	ElseArm   Expr
}

func (nce NodeTernaryExpr) expr() {}

// #######################################################
// ##################### Node Map Expr ###################
// #######################################################
type NodeMapExpr struct {
	Token token.Token
	Map   map[Expr]Expr
}

func (nmd NodeMapExpr) expr() {}

// #######################################################
// #################### Node Array Expr ##################
// #######################################################
type NodeArrayExpr struct {
	Token token.Token
	Value []Expr
}

func (na NodeArrayExpr) expr() {}

// #######################################################
// #################### Node Identifier###################
// #######################################################
type NodeIdentifier struct {
	Token token.Token
	Name  string
}

func (ni NodeIdentifier) expr() {}

// #######################################################
// ################### Node Import Expr ##################
// #######################################################
type NodeImportExpr struct {
	Token token.Token
	Path  Expr
}

func (ni NodeImportExpr) expr() {}
