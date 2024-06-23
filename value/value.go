package value

import (
	"fmt"
	"kat/ast"
)

type Value interface {
	String() string
}

var (
	TRUE  = ValueBool{true}
	FALSE = ValueBool{false}
)

type ValueFloat struct {
	Value float64
}

func (vf ValueFloat) String() string {
	return fmt.Sprintf("%.2f", vf.Value)
}

type ValueInt struct {
	Value int64
}

func (vi ValueInt) String() string {
	return fmt.Sprintf("%d", vi.Value)
}

type ValueBool struct {
	Value bool
}

func (vb ValueBool) String() string {
	return fmt.Sprintf("%t", vb.Value)
}

type ValueString struct {
	Value string
}

func (vs ValueString) String() string {
	return vs.Value
}

type ValueNil struct{}

func (vn ValueNil) String() string {
	return "nil"
}

type ValueFunction struct {
	Args []Value
	Body ast.NodeBlockStmt
}

func (vn ValueFunction) String() string {
	return "fn"
}

type ValueReturn struct {
	Value Value
}

func (vn ValueReturn) String() string {
	return fmt.Sprintf("%v", vn.Value)
}
