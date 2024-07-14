package value

import (
	"fmt"
	"kat/ast"
	"strings"
)

type Type string

const (
	TYPE_INT          Type = "int"
	TYPE_FLOAT        Type = "float"
	TYPE_BOOL         Type = "bool"
	TYPE_STRING       Type = "string"
	TYPE_ARRAY        Type = "array"
	TYPE_MAP          Type = "map"
	TYPE_STRUCT       Type = "struct"
	TYPE_FUNCTION     Type = "function"
	TYPE_MODULE       Type = "module"
	TYPE_NULL         Type = "null"
	TYPE_SELF         Type = "self"
	TYPE_KEYVAL       Type = "keyval"
	TYPE_RETURN       Type = "return"
	TYPE_ERROR        Type = "error"
	TYPE_STD_FUNCTION Type = "std_function"
	TYPE_ENVIRONMENT  Type = "environment"
)

type Value interface {
	String() string
	Type() Type
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

func (f *Float) Type() Type {
	return TYPE_FLOAT
}

type Int struct {
	Value int64
}

func (i *Int) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Int) Type() Type {
	return TYPE_INT
}

type Bool struct {
	Value bool
}

func (b *Bool) String() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Bool) Type() Type {
	return TYPE_BOOL
}

type String struct {
	Value string
}

func (s *String) String() string {
	return s.Value
}

func (s *String) Type() Type {
	return TYPE_STRING
}

type Self struct {
	Value string
}

func (s *Self) String() string {
	return s.Value
}

func (s *Self) Type() Type {
	return TYPE_SELF
}

type Null struct{}

func (n *Null) String() string {
	return "null"
}

func (n *Null) Type() Type {
	return TYPE_NULL
}

type Function struct {
	Args []Value
	Body ast.Stmt
}

func (f *Function) String() string {
	return "fn"
}

func (f *Function) Type() Type {
	return TYPE_FUNCTION
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

func (s *Struct[T]) Type() Type {
	return TYPE_STRUCT
}

type Map[T any] struct {
	*KeyVal[T]
}

func (m *Map[T]) String() string {
	return "valuemap"
}

func (m *Map[T]) Type() Type {
	return TYPE_MAP
}

type KeyVal[T any] struct {
	Map map[string]T
}

func (kv *KeyVal[T]) String() string {
	return "valuekeyval"
}

func (kv *KeyVal[T]) Type() Type {
	return TYPE_KEYVAL
}

type Return struct {
	Value Value
}

func (r *Return) String() string {
	return fmt.Sprintf("%v", r.Value)
}

func (r *Return) Type() Type {
	return TYPE_RETURN
}

type Module struct {
	Value Value
}

func (m *Module) String() string {
	return fmt.Sprintf("%v", m.Value)
}

func (m *Module) Type() Type {
	return TYPE_MODULE
}

type Array struct {
	Value []Value
}

func (a *Array) String() string {
	return fmt.Sprintf("%v", a.Value)
}

func (a *Array) Type() Type {
	return TYPE_ARRAY
}

type WrapperFunction struct {
	Name string
	Fn   func(varargs ...Value) Value
}

func (wf WrapperFunction) String() string {
	return "wrapperfunction"
}

func (wf WrapperFunction) Type() Type {
	return TYPE_STD_FUNCTION
}

type Error struct {
	Value string
}

func (e *Error) String() string {
	return e.Value
}

func (e *Error) Type() Type {
	return TYPE_ERROR
}
