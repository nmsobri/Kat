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
	TRUE  = &Bool{true}
	FALSE = &Bool{false}
	NULL  = &Null{}
)

type Float struct {
	Value float64
}

func (f *Float) String() string {
	return fmt.Sprintf("%.2f", f.Value)
}

type Int struct {
	Value int64
}

func (i *Int) String() string {
	return fmt.Sprintf("%d", i.Value)
}

type Bool struct {
	Value bool
}

func (b *Bool) String() string {
	return fmt.Sprintf("%t", b.Value)
}

type String struct {
	Value string
}

func (s *String) String() string {
	return s.Value
}

type Self struct {
	Value string
}

func (s *Self) String() string {
	return s.Value
}

type Null struct{}

func (n *Null) String() string {
	return "null"
}

type Function struct {
	Args []Value
	Body ast.Stmt
}

func (f *Function) String() string {
	return "fn"
}

type Struct[T any] struct {
	Name string
	Prop []string
	*KeyVal[T]
}

func (s *Struct[T]) String() string {
	valStruct := make([]string, 0)

	for _, k := range s.Prop {
		valStruct = append(valStruct, fmt.Sprintf("%s: %s", k, s.Map[k]))
	}

	return fmt.Sprintf("%s{%s}", s.Name, strings.Join(valStruct, ", "))
}

type Map[T any] struct {
	*KeyVal[T]
}

func (m *Map[T]) String() string {
	return "valuemap"
}

type KeyVal[T any] struct {
	Map map[string]T
}

func (kv *KeyVal[T]) String() string {
	return "valuekeyval"
}

type Return struct {
	Value Value
}

func (r *Return) String() string {
	return fmt.Sprintf("%v", r.Value)
}

type Module struct {
	Value Value
}

func (m *Module) String() string {
	return fmt.Sprintf("%v", m.Value)
}

type Array struct {
	Value []Value
}

func (a *Array) String() string {
	return fmt.Sprintf("%v", a.Value)
}

type WrapperFunction struct {
	Name string
	Fn   func(varargs ...Value) Value
}

func (wf WrapperFunction) String() string {
	return "wrapperfunction"
}

type Error struct {
	Value string
}

func (e *Error) String() string {
	return e.Value
}
