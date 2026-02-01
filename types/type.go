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

type AtomicType struct {
	Type string
}

func (t AtomicType) String() string {
	return t.Type
}

type ArrayType struct {
	Type ast.Type
}

func (t ArrayType) String() string {
	return fmt.Sprintf("[%s]", t.Type)
}

type MapType struct {
	KeyType   ast.Type
	ValueType ast.Type
}

func (t MapType) String() string {
	return fmt.Sprintf("[%s]%s", t.KeyType, t.ValueType)
}

var (
	INT    = AtomicType{Type: "int"}
	FLOAT  = AtomicType{Type: "float"}
	BOOL   = AtomicType{Type: "bool"}
	STRING = AtomicType{Type: "string"}
)

func IsInt(t ast.Type) bool {
	return t.(AtomicType).Type == "int"
}

func IsFloat(t ast.Type) bool {
	return t.(AtomicType).Type == "float"
}

func IsBool(t ast.Type) bool {
	return t.(AtomicType).Type == "bool"
}

func IsString(t ast.Type) bool {
	return t.(AtomicType).Type == "string"
}
