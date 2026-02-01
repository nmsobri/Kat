package ast

import (
	"fmt"
	"kat/token"
	"regexp"

	"github.com/sanity-io/litter"
)

type Type interface {
	fmt.Stringer
}

// Simulate Tagged Union
// This mean following nodes can either be Statement or Node
type Statement struct {
	Line   int
	Column int
}

func (s Statement) stmt() {}
func (s Statement) node() {}

func (s Statement) GetLocation() Location {
	return Location{
		Line:   s.Line,
		Column: s.Column,
	}
}

// #######################################################
// ##################### Node Program#####################😀
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
// ################## Node Modern For Stmt ###############😀
// #######################################################
type NodeModernForStmt struct {
	Statement
	Token     token.Token
	Condition Expr
	Body      Stmt
}

// #######################################################
// ################# Node Classic For Stmt ###############😀
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
// #################### Node Const Stmt ##################😀
// #######################################################
type NodeConstStmt struct {
	Statement
	Token      token.Token
	Identifier token.Token
	Value      Expr
}

// #######################################################
// ################## Node Struct Stmt ###################😀
// #######################################################
type NodeStructStmt struct {
	Statement
	Token      token.Token
	Identifier Expr
	Properties []Expr
}

// #######################################################
// ################ Node Function Stmt ###################😀
// #######################################################
type NodeFunctionStmt struct {
	Statement
	Token      token.Token
	Identifier Expr
	Arguements []Expr
	Body       Stmt
}

// #######################################################
// ##################### Node Let Stmt ###################😀
// #######################################################
type NodeLetStmt struct {
	Statement
	Token      token.Token
	Identifier token.Token
	Type       Type
	Value      Expr
}

// #######################################################
// #################### Node Expr Stmt ###################😀
// #######################################################
type NodeExprStmt struct {
	Statement
	Expr Expr
}

// #######################################################
// ################# Node conditional stmt ###############😀
// #######################################################
type NodeConditionalStmt struct {
	Statement
	Token     token.Token
	Condition Expr
	ThenArm   Stmt
	ElseArm   Stmt
}

// #######################################################
// #################### Node Return stmt #################😀
// #######################################################
type NodeReturnStmt struct {
	Statement
	Token token.Token
	Value Expr
}

// #######################################################
// #################### Node Block stmt ##################😀
// #######################################################
type NodeBlockStmt struct {
	Statement
	Body []Stmt
}
