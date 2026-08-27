package ast

import (
	"kat/token"
	"regexp"

	"github.com/sanity-io/litter"
)

// Simulate Tagged Union
// This mean following nodes can either be Statement or Node
type Statement struct{}

func (s Statement) stmt() {}
func (s Statement) node() {}

// #######################################################
// ##################### Node Program#####################
// #######################################################
type NodeProgram struct {
	Statement
	Body []Stmt // Statement
}

func (np *NodeProgram) String() string {
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
	Body      Stmt
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
	Body      Stmt
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
	Body       Stmt
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
// #################### Node Return stmt #################
// #######################################################
type NodeReturnStmt struct {
	Statement
	Token token.Token
	Value Expr
}

// #######################################################
// #################### Node Block stmt ##################
// #######################################################
type NodeBlockStmt struct {
	Statement
	Body []Stmt
}

// #######################################################
// ################### Node Import stmt ##################
// #######################################################
type NodeImportStmt struct {
	Statement
	Token token.Token
	Path  Expr
	Alias Expr
}
