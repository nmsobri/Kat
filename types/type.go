package types

import "fmt"

/**
int
float
string
bool
*/

type Type interface {
	fmt.Stringer
}

type AtomicType struct {
	Type string
}

func (t AtomicType) String() string {
	return t.Type
}

type ArrayType struct {
	Type Type
}

func (t ArrayType) String() string {
	return fmt.Sprintf("[%s]", t.Type)
}

type MapType struct {
	KeyType   Type
	ValueType Type
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

func IsInt(t Type) bool {
	return t.(AtomicType).Type == "int"
}

func IsFloat(t Type) bool {
	return t.(AtomicType).Type == "float"
}

func IsBool(t Type) bool {
	return t.(AtomicType).Type == "bool"
}

func IsString(t Type) bool {
	return t.(AtomicType).Type == "string"
}
