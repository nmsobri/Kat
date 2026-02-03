package types

import (
	"fmt"
	"kat/ast"
)

/**
int
float
string
bool
*/

type Atomic struct {
	Type string
}

func (t Atomic) String() string {
	return t.Type
}

type Array struct {
	Type  string
	Value string
}

func (t Array) String() string {
	return fmt.Sprintf("[]%s", t.Value)
}

type Map struct {
	KeyType   ast.Type
	ValueType ast.Type
}

func (t Map) String() string {
	return fmt.Sprintf("[%s]%s", t.KeyType, t.ValueType)
}

var (
	INT    = Atomic{Type: "int"}
	FLOAT  = Atomic{Type: "float"}
	BOOL   = Atomic{Type: "bool"}
	STRING = Atomic{Type: "string"}
	ARRAY  = func(t string) ast.Type { return Array{Type: "array", Value: t} }
)

func IsInt(t ast.Type) bool {
	if _, ok := t.(Atomic); ok {
		return t.(Atomic).Type == "int"
	}

	return false
}

func IsFloat(t ast.Type) bool {
	if _, ok := t.(Atomic); ok {
		return t.(Atomic).Type == "float"
	}

	return false
}

func IsBool(t ast.Type) bool {
	if _, ok := t.(Atomic); ok {
		return t.(Atomic).Type == "bool"
	}

	return false
}

func IsString(t ast.Type) bool {
	if _, ok := t.(Atomic); ok {
		return t.(Atomic).Type == "string"
	}

	return false
}

func IsArray(t ast.Type) bool {
	if _, ok := t.(Array); ok {
		return t.(Array).Type == "array"
	}

	return false
}
