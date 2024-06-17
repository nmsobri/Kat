package ast

import (
	"github.com/sanity-io/litter"
	"kat/token"
	"regexp"
)

// #######################################################
// ##################### Node Program ####################
// #######################################################
type NodeProgram struct {
	Body []Stmt // Statement
}

func (np NodeProgram) stmt() {}

func (np NodeProgram) String() string {
	litter.Config.FieldExclusions = regexp.MustCompile(`^(Token)$`)
	return litter.Sdump(np)
}

// #######################################################
// ##################### Node For Stmt ###################
// #######################################################
type NodeForStmt struct {
	Token     token.Token
	Condition Expr
	Body      BlockStmt
}

func (nf NodeForStmt) stmt() {}

// #######################################################
// #################### Node Const Stmt ##################
// #######################################################
type NodeConstStmt struct {
	Token      token.Token
	Identifier Expr
	Value      Expr
}

func (nc NodeConstStmt) stmt() {}

// #######################################################
// ################## Node Struct Stmt ###################
// #######################################################
type NodeStructStmt struct {
	Token      token.Token
	Identifier Expr
	Properties NodeStructProperties
}

func (ns NodeStructStmt) stmt() {}

// #######################################################
// ################## Node Struct Prop ###################
// #######################################################
type NodeStructProperties struct {
	Token      token.Token
	Properties []Expr
}

func (nsp NodeStructProperties) stmt() {}

// #######################################################
// ################ Node Function Stmt ###################
// #######################################################
type NodeFunctionStmt struct {
	Token      token.Token
	Identifier Expr
	Arguements []Expr
	Body       BlockStmt
}

func (nf NodeFunctionStmt) stmt() {}

// #######################################################
// ##################### Node Let Stmt ###################
// #######################################################
type NodeLetStmt struct {
	Token      token.Token
	Identifier Expr
	Value      Expr
}

func (nld NodeLetStmt) stmt() {}

// #######################################################
// #################### Node Expr Stmt ###################
// #######################################################
type NodeExprStmt struct {
	Expression Expr
}

func (ne NodeExprStmt) stmt() {}

// #######################################################
// ################# Node conditional stmt ###############
// #######################################################
type NodeConditionalStmt struct {
	Token     token.Token
	Condition Expr
	ThenArm   Stmt
	ElseArm   Stmt
}

func (nc NodeConditionalStmt) stmt() {}

// #######################################################
// #################### Node Block stmt ##################
// #######################################################
type BlockStmt struct {
	Body []Stmt
}

func (nb BlockStmt) stmt() {}
