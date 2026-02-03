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
		inferredType := c.Infer(node.Value)
		location := node.Value.GetLocation()

		if !c.SubType(declaredType, inferredType) {
			log.Fatalf("expected type %s, but got %s at line:%d, column:%d",
				declaredType, inferredType, location.Line, location.Column,
			)
		}

	default:
		panic("Checker::Check todo")
	}

}

func (c *Checker) Infer(expr ast.Expr) ast.Type {
	switch node := expr.(type) {

	case *ast.NodeInteger:
		return INT

	case *ast.NodeString:
		return STRING

	case *ast.NodeFloat:
		return FLOAT

	case *ast.NodeBoolean:
		return BOOL

	case *ast.NodeArrayExpr:
		arrayType := c.Infer(node.Value[0]) // get array type from first element

		for _, v := range node.Value {
			elementType := c.Infer(v)

			// todo: expression: [1,2, "hello"] gave: 026/02/03 14:04:32 expected type int, but got string at line:99, column:99999k
			// todo: should give proper error message like `expected type []int, got []int but contain string`
			// todo: plus location is incorrect
			if !c.SubType(arrayType, elementType) {
				log.Fatalf("expected type %s, but got %s at line:%d, column:%d",
					arrayType, elementType, 99, 99999,
				)
			}
		}

		return ARRAY(arrayType.String())

	default:
		panic("infer::unreachable")
	}

	return nil
}

func (c *Checker) SubType(a, b ast.Type) bool {
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

	if IsArray(a) && IsArray(b) && a.(Array).Value == b.(Array).Value {
		return true
	}

	return false
}

func New() *Checker {
	return &Checker{}
}
