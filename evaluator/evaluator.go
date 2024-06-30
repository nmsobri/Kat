package evaluator

import (
	"errors"
	"fmt"
	"kat/ast"
	"kat/environment"
	"kat/lexer"
	"kat/parser"
	"kat/util"
	"kat/value"
	"os"
)

type Evaluator struct {
	Errors []error
	Tree   ast.Stmt
}

func New(tree ast.Stmt) *Evaluator {
	return &Evaluator{
		Tree: tree,
	}
}

func (e *Evaluator) Eval(node ast.Node, env *environment.Environment) value.Value {
	var result value.Value = value.ValueNil{}

	if e.IsError() {
		return result
	}

	switch stmt := node.(type) {

	case ast.NodeProgram:
		for _, stmt := range stmt.Body {
			result = e.Eval(stmt, env)

			if _, ok := result.(value.ValueReturn); ok {
				return result
			}
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

		switch stmt.Operator {

		case "+":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)
			val := left.(value.ValueInt).Value + right.(value.ValueInt).Value
			return value.ValueInt{val}

		case "-":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)
			val := left.(value.ValueInt).Value - right.(value.ValueInt).Value
			return value.ValueInt{val}

		case "*":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)
			val := left.(value.ValueInt).Value * right.(value.ValueInt).Value
			return value.ValueInt{val}

		case "/":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)
			val := left.(value.ValueInt).Value / right.(value.ValueInt).Value
			return value.ValueInt{val}

		case "%":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)
			val := left.(value.ValueInt).Value % right.(value.ValueInt).Value
			return value.ValueInt{val}

		case "=":
			ident := stmt.Left.(ast.NodeIdentifier).Name
			val := e.Eval(stmt.Right, env)

			if _, ok := env.Get(ident); !ok {
				msg := fmt.Sprintf("Variable %s is not found", ident)
				e.Error(msg)

				return value.ValueNil{}
			}

			env.Assign(ident, val)
			return val

		case "<":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)

			if util.TypeOf(left) != util.TypeOf(right) {
				msg := fmt.Sprintf("Invalid operation %s < %s", left, right)
				e.Error(msg)
			}

			switch left.(type) {

			case value.ValueInt:
				left, _ := left.(value.ValueInt)
				right, _ := right.(value.ValueInt)

				if left.Value < right.Value {
					return value.TRUE
				}

				return value.FALSE

			case value.ValueFloat:
				left, _ := left.(value.ValueFloat)
				right, _ := right.(value.ValueFloat)

				if left.Value < right.Value {
					return value.TRUE
				}

				return value.FALSE

			default:
				msg := fmt.Sprintf("Invalid operation %s < %s", left, right)
				e.Error(msg)
			}

		case ">":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)

			if util.TypeOf(left) != util.TypeOf(right) {
				msg := fmt.Sprintf("Invalid operation %s > %s", left, right)
				e.Error(msg)
			}

			switch left.(type) {
			case value.ValueInt:
				left, _ := left.(value.ValueInt)
				right, _ := right.(value.ValueInt)

				if left.Value > right.Value {
					return value.TRUE
				}

				return value.FALSE

			case value.ValueFloat:
				left, _ := left.(value.ValueFloat)
				right, _ := right.(value.ValueFloat)

				if left.Value > right.Value {
					return value.TRUE
				}

				return value.FALSE

			default:
				msg := fmt.Sprintf("Invalid operation %s > %s", left, right)
				e.Error(msg)
			}

		case "<=":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)

			if util.TypeOf(left) != util.TypeOf(right) {
				msg := fmt.Sprintf("Invalid operation %s <= %s", left, right)
				e.Error(msg)
			}

			switch left.(type) {

			case value.ValueInt:
				left, _ := left.(value.ValueInt)
				right, _ := right.(value.ValueInt)

				if left.Value <= right.Value {
					return value.TRUE
				}

				return value.FALSE

			case value.ValueFloat:
				_left := e.Eval(stmt.Left, env)
				_right := e.Eval(stmt.Right, env)

				left, _ := _left.(value.ValueFloat)
				right, _ := _right.(value.ValueFloat)

				if left.Value <= right.Value {
					return value.TRUE
				}

				return value.FALSE

			default:
				msg := fmt.Sprintf("Invalid operation %s <= %s", left, right)
				e.Error(msg)
			}

		case ">=":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)

			if util.TypeOf(left) != util.TypeOf(right) {
				msg := fmt.Sprintf("Invalid operation %s >= %s", left, right)
				e.Error(msg)
			}

			switch left.(type) {

			case value.ValueInt:
				left, _ := left.(value.ValueInt)
				right, _ := right.(value.ValueInt)

				if left.Value >= right.Value {
					return value.TRUE
				}

				return value.FALSE

			case value.ValueFloat:
				left, _ := left.(value.ValueFloat)
				right, _ := right.(value.ValueFloat)

				if left.Value >= right.Value {
					return value.TRUE
				}

				return value.FALSE

			default:
				msg := fmt.Sprintf("Invalid operation %s >= %s", left, right)
				e.Error(msg)
			}

		case "==":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)

			if left == right {
				return value.TRUE
			}

			return value.FALSE

		case "!=":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)

			if left != right {
				return value.TRUE
			}

			return value.FALSE

		case ".":
			left := e.Eval(stmt.Left, env)
			right := stmt.Right.(ast.NodeIdentifier).Name

			env := left.(value.ValueEnv).Value.(*environment.Environment)
			val, ok := env.Get(right)

			if !ok {
				msg := fmt.Sprintf("Symbol %s is not found", right)
				e.Error(msg)
			}

			return val

		default:
			msg := fmt.Sprintf("Unrecognized operator: %s", stmt.Operator)
			e.Error(msg)
		}

	case ast.NodeConstStmt:
		ident := stmt.Identifier.(ast.NodeIdentifier).Name
		val := e.Eval(stmt.Value, env)

		if _, ok := env.Get(ident); ok {
			msg := fmt.Sprintf("Constant %s already exists", ident)
			e.Error(msg)
		}

		env.Set(ident, val)

	case ast.NodeLetStmt:
		ident := stmt.Identifier.(ast.NodeIdentifier).Name
		val := e.Eval(stmt.Value, env)

		if _, ok := env.Get(ident); ok {
			msg := fmt.Sprintf("Variable %s is already exists", ident)
			e.Error(msg)
		}

		env.Set(ident, val)

	case ast.NodeIdentifier:
		val, ok := env.Get(stmt.Name)

		if !ok {
			msg := fmt.Sprintf("Variable %s is not found", stmt.Name)
			e.Error(msg)
		}

		return val

	case ast.NodeFunctionStmt:
		ident := stmt.Identifier.(ast.NodeIdentifier).Name

		args := make([]value.Value, len(stmt.Arguements))
		for i, arg := range stmt.Arguements {
			args[i] = value.ValueString{arg.(ast.NodeIdentifier).Name}
		}

		val := value.ValueFunction{Args: args, Body: stmt.Body}

		if _, ok := env.Get(ident); ok {
			msg := fmt.Sprintf("Symbol %s already exists", ident)
			e.Error(msg)
		}

		env.Set(ident, val)

	case ast.NodeFunctionCall:
		val := e.Eval(stmt.Identifer, env)
		params := make([]value.Value, len(stmt.Parameters))
		for i, param := range stmt.Parameters {
			params[i] = e.Eval(param, env)
		}

		fn, ok := val.(value.ValueFunction)

		if !ok {
			return result
		}
		newEnv := environment.NewWithParent(env)

		for i, arg := range fn.Args {
			newEnv.Set(arg.String(), params[i])
		}

		result := e.Eval(fn.Body, newEnv)
		return result

	case ast.NodeBlockStmt:
		for _, stmt := range stmt.Body {
			result = e.Eval(stmt, env)

			if _, ok := result.(value.ValueReturn); ok {
				return result.(value.ValueReturn).Value
			}
		}

	case ast.NodeReturnStmt:
		return value.ValueReturn{e.Eval(stmt.Value, env)}

	case ast.NodeConditionalStmt:
		condition := e.Eval(stmt.Condition, env)

		if util.IsTruthy(condition) {
			return e.Eval(stmt.ThenArm, env)
		} else {
			return e.Eval(stmt.ElseArm, env)
		}

	case ast.NodeStructStmt:
		identifier, ok := stmt.Identifier.(ast.NodeIdentifier)

		if !ok {
			msg := fmt.Sprintf("Invalid identifier: %s", stmt.Identifier)
			e.Error(msg)
		}

		_, ok = env.Get(identifier.Name)

		if ok {
			msg := fmt.Sprintf("Symbol %s already exists", identifier.Name)
			e.Error(msg)
		}

		props := make([]string, 0)
		valKeyVal := make(map[string]byte)

		for _, p := range stmt.Properties {
			prop, ok := p.(ast.NodeIdentifier)

			if !ok {
				msg := fmt.Sprintf("Invalid property: %s", p)
				e.Error(msg)
			}

			valKeyVal[prop.Name] = 1
			props = append(props, prop.Name)
		}

		env.Set(identifier.Name, value.ValueStruct[byte]{identifier.Name, props, value.ValueKeyVal[byte]{Map: valKeyVal}})

	case ast.NodeStructExpr:
		ident, ok := stmt.Name.(ast.NodeIdentifier)

		if !ok {
			msg := fmt.Sprintf("Invalid identifier: %s", stmt.Name)
			e.Error(msg)
		}

		structStmt, ok := env.Get(ident.Name)
		structStmtProp := structStmt.(value.ValueStruct[byte]).Prop

		if !ok {
			msg := fmt.Sprintf("Struct %s is not found", ident.Name)
			e.Error(msg)
		}

		keyMap := e.Eval(stmt.Values, env)
		props := keyMap.(value.ValueKeyVal[value.Value])

		for k := range props.Map {
			ok := util.InArray[string](structStmtProp, k)

			if !ok {
				msg := fmt.Sprintf("Unknown field %s on %s", k, ident.Name)
				e.Error(msg)
			}
		}

		return value.ValueStruct[value.Value]{ident.Name, structStmtProp, value.ValueKeyVal[value.Value]{props.Map}}

	case ast.NodeMapExpr:
		keyVal := make(map[string]value.Value)

		for k, v := range stmt.Map {
			key, ok := k.(ast.NodeIdentifier)

			if !ok {
				msg := fmt.Sprintf("Invalid struct key: %s", k)
				e.Error(msg)
			}

			keyVal[key.Name] = e.Eval(v, env)
		}

		return value.ValueKeyVal[value.Value]{keyVal}

	case ast.NodeModernForStmt:
		condition := e.Eval(stmt.Condition, env)

		for util.IsTruthy(condition) {
			e.Eval(stmt.Body, env)
			condition = e.Eval(stmt.Condition, env)
		}

	case ast.NodeClassicForStmt:
		newEnv := environment.NewWithParent(env)

		e.Eval(stmt.PreExpr, newEnv) // pre expression
		condition := e.Eval(stmt.Condition, newEnv)

		for util.IsTruthy(condition) {
			e.Eval(stmt.Body, newEnv)
			e.Eval(stmt.PostExpr, newEnv) // post expression
			condition = e.Eval(stmt.Condition, newEnv)
		}

	case ast.NodePostfixExpr:
		left, ok := stmt.Left.(ast.NodeIdentifier)

		if !ok {
			msg := fmt.Sprintf("Invalid identifier: %s", stmt.Left)
			e.Error(msg)
		}

		leftVal, ok := env.Get(left.Name)

		if !ok {
			msg := fmt.Sprintf("Variable %s is not found", left.Name)
			e.Error(msg)
		}

		switch stmt.Operator {
		case "++":
			switch leftVal.(type) {
			case value.ValueInt:
				env.Set(left.Name, value.ValueInt{leftVal.(value.ValueInt).Value + 1})
				return value.ValueInt{leftVal.(value.ValueInt).Value}

			default:
				msg := fmt.Sprintf("Unsupported operator: %s for type %T", stmt.Operator, leftVal)
				e.Error(msg)
			}

		case "--":
			switch leftVal.(type) {
			case value.ValueInt:
				env.Set(left.Name, value.ValueInt{leftVal.(value.ValueInt).Value - 1})
				return value.ValueInt{leftVal.(value.ValueInt).Value}

			default:
				msg := fmt.Sprintf("Unsupported operator: %s for type %T", stmt.Operator, leftVal)
				e.Error(msg)
			}

		default:
			msg := fmt.Sprintf("Unsupported operator: %s", stmt.Operator)
			e.Error(msg)
		}

	case ast.NodePrefixExpr:
		right, ok := stmt.Right.(ast.NodeIdentifier)

		if !ok {
			msg := fmt.Sprintf("Invalid identifier: %s", stmt.Right)
			e.Error(msg)
		}

		rightVal, ok := env.Get(right.Name)

		if !ok {
			msg := fmt.Sprintf("Variable %s is not found", right.Name)
			e.Error(msg)
		}

		switch stmt.Operator {
		case "++":
			switch rightVal.(type) {
			case value.ValueInt:
				env.Set(right.Name, value.ValueInt{rightVal.(value.ValueInt).Value + 1})
				return value.ValueInt{rightVal.(value.ValueInt).Value + 1}

			default:
				msg := fmt.Sprintf("Unsupported operator: %s for type %T", stmt.Operator, rightVal)
				e.Error(msg)
			}

		case "--":
			switch rightVal.(type) {
			case value.ValueInt:
				env.Set(right.Name, value.ValueInt{rightVal.(value.ValueInt).Value - 1})
				return value.ValueInt{rightVal.(value.ValueInt).Value - 1}

			default:
				msg := fmt.Sprintf("Unsupported operator: %s for type %T", stmt.Operator, rightVal)
				e.Error(msg)
			}

		default:
			msg := fmt.Sprintf("Unsupported operator: %s", stmt.Operator)
			e.Error(msg)
		}

	case ast.NodeImportExpr:
		path := stmt.Path.(ast.NodeString).Value
		fileName := fmt.Sprintf("%s.kat", path)
		realPath := fmt.Sprintf("%s/%s/%s", os.Getenv("PWD"), "stdlib", fileName)

		source := util.ReadFile(realPath)
		l := lexer.New(source)
		p := parser.New(l)

		program := p.ParseProgram()

		e := New(program)
		env := environment.New()
		e.Eval(program, env)

		if e.IsError() {
			fmt.Println("Evaluation Errors:")
			for _, err := range e.Errors {
				fmt.Println(err)
			}
		}

		return value.ValueEnv{env}

	default:
		msg := fmt.Sprintf("Unrecognized statement type: %T", stmt)
		e.Error(msg)
	}

	return result
}

func (e *Evaluator) Error(msg string) {
	e.Errors = append(e.Errors, errors.New(msg))
}

func (e *Evaluator) IsError() bool {
	return len(e.Errors) > 0
}
