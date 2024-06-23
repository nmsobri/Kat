package evaluator

import (
	"kat/ast"
	"kat/environment"
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

func (e *Evaluator) Eval(node ast.Node, env *environment.Environment) value.Value {
	var result value.Value = value.ValueInt{0}

	switch stmt := node.(type) {

	case ast.NodeProgram:
		for _, stmt := range stmt.Body {
			result = e.Eval(stmt, env)
		}

	case ast.NodeExprStmt:
		return e.Eval(stmt.Expr, env)

	case ast.NodeInteger:
		return value.ValueInt{stmt.Value}

	case ast.NodeFloat:
		return value.ValueFloat{stmt.Value}

	case ast.NodeBoolean:
		return value.ValueBool{stmt.Value}

	case ast.NodeString:
		return value.ValueString{stmt.Value}

	case ast.NodeBinaryExpr:
		left := e.Eval(stmt.Left, env)
		right := e.Eval(stmt.Right, env)

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

		case "%":
			val := left.(value.ValueInt).Value % right.(value.ValueInt).Value
			return value.ValueInt{val}
		}
	}

	return result
}
