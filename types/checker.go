package types

import (
	"kat/ast"
	"log"
)

type Checker struct{}

func (c *Checker) CheckProgram(program *ast.NodeProgram) {
	for _, stmt := range program.Body {
		c.Check(stmt)
	}
}

func (c *Checker) Check(astNode ast.Node) {
	switch node := astNode.(type) {

	case *ast.NodeLetStmt:
		declaredType := node.Type
		expressionType := c.Infer(node.Value)
		location := node.Value.GetLocation()

		if !c.SubType(declaredType, expressionType) {
			log.Fatalf("expected type %s, but got %s at line:%d, column:%d",
				declaredType, expressionType, location.Line, location.Column,
			)
		}

	default:
		panic("Checker::Check todo")
	}

}

func (c *Checker) Infer(expr ast.Expr) Type {
	switch expr.(type) {

	case *ast.NodeInteger:
		return INT

	case *ast.NodeString:
		return STRING

	case *ast.NodeFloat:
		return FLOAT

	case *ast.NodeBoolean:
		return BOOL

	default:
		panic("unreachable")
	}

	return nil
}

func (c *Checker) SubType(a, b Type) bool {
	if IsBool(a) && IsBool(b) {
		return true
	}

	if IsInt(a) && IsInt(b) {
		return true
	}

	if IsFloat(a) && IsFloat(b) {
		return true
	}

	if IsString(a) && IsString(b) {
		return true
	}

	return false
}

func New() *Checker {
	return &Checker{}
}
