package value

import (
	"fmt"
	"kat/ast"
	"strings"
)

type Value interface {
	String() string
}

var (
	TRUE  = &ValueBool{true}
	FALSE = &ValueBool{false}
	NULL  = &ValueNull{}
)

type ValueFloat struct {
	Value float64
}

func (vf *ValueFloat) String() string {
	return fmt.Sprintf("%.2f", vf.Value)
}

type ValueInt struct {
	Value int64
}

func (vi *ValueInt) String() string {
	return fmt.Sprintf("%d", vi.Value)
}

type ValueBool struct {
	Value bool
}

func (vb *ValueBool) String() string {
	return fmt.Sprintf("%t", vb.Value)
}

type ValueString struct {
	Value string
}

func (vs *ValueString) String() string {
	return vs.Value
}

type ValueSelf struct {
	Value string
}

func (vs *ValueSelf) String() string {
	return vs.Value
}

type ValueNull struct{}

func (vn *ValueNull) String() string {
	return "null"
}

type ValueFunction struct {
	Args []Value
	Body ast.Stmt
}

func (vn *ValueFunction) String() string {
	return "fn"
}

type ValueStruct[T any] struct {
	Name string
	Prop []string
	*ValueKeyVal[T]
}

func (vs *ValueStruct[T]) String() string {
	valStruct := make([]string, 0)

	for _, k := range vs.Prop {
		valStruct = append(valStruct, fmt.Sprintf("%s: %s", k, vs.Map[k]))
	}

	return fmt.Sprintf("%s{%s}", vs.Name, strings.Join(valStruct, ", "))
}

type ValueMap[T any] struct {
	*ValueKeyVal[T]
}

func (vs *ValueMap[T]) String() string {
	return "valuemap"
}

type ValueKeyVal[T any] struct {
	Map map[string]T
}

func (vs *ValueKeyVal[T]) String() string {
	return "valuekeyval"
}

type ValueReturn struct {
	Value Value
}

func (vn *ValueReturn) String() string {
	return fmt.Sprintf("%v", vn.Value)
}

type ValueModule struct {
	Value Value
}

func (vm *ValueModule) String() string {
	return fmt.Sprintf("%v", vm.Value)
}

type ValueArray struct {
	Value []Value
}

func (vn *ValueArray) String() string {
	return fmt.Sprintf("%v", vn.Value)
}

type WrapperFunction struct {
	Name string
	Fn   func(varargs ...Value) Value
}

func (wf WrapperFunction) String() string {
	return "wrapperfunction"
}
