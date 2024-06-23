package evaluator

import (
	"kat/ast"
	"kat/value"
)

type Evaluator struct {
	Tree ast.Stmt
}

func New(tree ast.Stmt) *Evaluator {
	return &Evaluator{
		Tree: tree,
	}
}

func (e *Evaluator) Eval(node any) value.Value {
	var result value.Value = value.ValueInt{0}

	switch stmt := node.(type) {

	case ast.NodeProgram:
		for _, stmt := range stmt.Body {
			result = e.Eval(stmt)
		}

	case ast.NodeExprStmt:
		return e.Eval(stmt.Expression)

	case ast.NodeInteger:
		return value.ValueInt{stmt.Value}

	case ast.NodeDouble:
		return value.ValueDouble{stmt.Value}

	case ast.NodeBinaryExpr:
		left := e.Eval(stmt.Left)
		right := e.Eval(stmt.Right)

		switch stmt.Operator {
		case "+":
			val := left.(value.ValueInt).Value + right.(value.ValueInt).Value
			return value.ValueInt{val}

		case "-":
			val := left.(value.ValueInt).Value - right.(value.ValueInt).Value
			return value.ValueInt{val}

		case "*":
			val := left.(value.ValueInt).Value * right.(value.ValueInt).Value
			return value.ValueInt{val}

		case "/":
			val := left.(value.ValueInt).Value / right.(value.ValueInt).Value
			return value.ValueInt{val}
		}
	}

	return result
}
