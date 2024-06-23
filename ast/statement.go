package ast

import (
	"github.com/sanity-io/litter"
	"kat/token"
	"regexp"
)

// Simulate Tagged Union
// This mean following nodes can either be Statement or Node
type Statement struct{}

func (s Statement) stmt() {}
func (s Statement) node() {}

// #######################################################
// ##################### Node Tree ####################
// #######################################################
type NodeProgram struct {
	Statement
	Body []Stmt // Statement
}

func (np NodeProgram) String() string {
	litter.Config.FieldExclusions = regexp.MustCompile(`^(Token|Statement|Expression)$`)
	return litter.Sdump(np)
}

// #######################################################
// ################## Node Modern For Stmt ###############
// #######################################################
type NodeModernForStmt struct {
	Statement
	Token     token.Token
	Condition Expr
	Body      BlockStmt
}

// #######################################################
// ################# Node Classic For Stmt ###############
// #######################################################
type NodeClassicForStmt struct {
	Statement
	Token     token.Token
	Condition Expr
	PreExpr   Stmt
	PostExpr  Expr
	Body      BlockStmt
}

// #######################################################
// #################### Node Const Stmt ##################
// #######################################################
type NodeConstStmt struct {
	Statement
	Token      token.Token
	Identifier Expr
	Value      Expr
}

// #######################################################
// ################## Node Struct Stmt ###################
// #######################################################
type NodeStructStmt struct {
	Statement
	Token      token.Token
	Identifier Expr
	Properties NodeStructProperties
}

// #######################################################
// ################## Node Struct Prop ###################
// #######################################################
type NodeStructProperties struct {
	Statement
	Token      token.Token
	Properties []Expr
}

// #######################################################
// ################ Node Function Stmt ###################
// #######################################################
type NodeFunctionStmt struct {
	Statement
	Token      token.Token
	Identifier Expr
	Arguements []Expr
	Body       BlockStmt
}

// #######################################################
// ##################### Node Let Stmt ###################
// #######################################################
type NodeLetStmt struct {
	Statement
	Token      token.Token
	Identifier Expr
	Value      Expr
}

// #######################################################
// #################### Node Expr Stmt ###################
// #######################################################
type NodeExprStmt struct {
	Statement
	Expr Expr
}

// #######################################################
// ################# Node conditional stmt ###############
// #######################################################
type NodeConditionalStmt struct {
	Statement
	Token     token.Token
	Condition Expr
	ThenArm   Stmt
	ElseArm   Stmt
}

// #######################################################
// #################### Node Block stmt ##################
// #######################################################
type BlockStmt struct {
	Statement
	Body []Stmt
}
