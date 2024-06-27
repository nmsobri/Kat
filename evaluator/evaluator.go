package evaluator

import (
	"fmt"
	"kat/ast"
	"kat/environment"
	"kat/util"
	"kat/value"
	"log"
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
	var result value.Value = value.ValueNil{}

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

		case "=":
			ident := stmt.Left.(ast.NodeIdentifier).Name
			val := e.Eval(stmt.Right, env)

			if _, ok := env.Get(ident); !ok {
				msg := fmt.Sprintf("Variable %s not found", ident)
				panic(msg)
				return value.ValueNil{}
			}

			env.Set(ident, val)
			return val

		case "<":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)

			if util.TypeOf(left) != util.TypeOf(right) {
				msg := fmt.Sprintf("Invalid operation %s < %s", left, right)
				log.Fatal(msg)
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
				log.Fatal(msg)
			}

		case ">":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)

			if util.TypeOf(left) != util.TypeOf(right) {
				msg := fmt.Sprintf("Invalid operation %s > %s", left, right)
				log.Fatal(msg)
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
				log.Fatal(msg)

			}

		case "<=":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)

			if util.TypeOf(left) != util.TypeOf(right) {
				msg := fmt.Sprintf("Invalid operation %s <= %s", left, right)
				log.Fatal(msg)
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
				left, _ := left.(value.ValueFloat)
				right, _ := right.(value.ValueFloat)

				if left.Value <= right.Value {
					return value.TRUE
				}

				return value.FALSE
			default:
				msg := fmt.Sprintf("Invalid operation %s <= %s", left, right)
				log.Fatal(msg)

			}

		case ">=":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)

			if util.TypeOf(left) != util.TypeOf(right) {
				msg := fmt.Sprintf("Invalid operation %s >= %s", left, right)
				log.Fatal(msg)
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
				log.Fatal(msg)

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

		default:
			msg := fmt.Sprintf("Unrecognized operator: %s", stmt.Operator)
			log.Fatal(msg)
		}

	case ast.NodeConstStmt:
		ident := stmt.Identifier.(ast.NodeIdentifier).Name
		val := e.Eval(stmt.Value, env)

		if _, ok := env.Get(ident); ok {
			msg := fmt.Sprintf("Constant %s already exists", ident)
			log.Fatal(msg)
		}

		env.Set(ident, val)

	case ast.NodeLetStmt:
		ident := stmt.Identifier.(ast.NodeIdentifier).Name
		val := e.Eval(stmt.Value, env)

		if _, ok := env.Get(ident); ok {
			msg := fmt.Sprintf("Variable %s already exists", ident)
			log.Fatal(msg)
		}

		env.Set(ident, val)

	case ast.NodeIdentifier:
		val, ok := env.Get(stmt.Name)

		if !ok {
			msg := fmt.Sprintf("Variable %s not found", stmt.Name)
			log.Fatal(msg)
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
			log.Fatal(msg)
		}

		env.Set(ident, val)

	case ast.NodeFunctionCall:
		ident := stmt.Identifer.(ast.NodeIdentifier).Name

		val, ok := env.Get(ident)

		if !ok {
			msg := fmt.Sprintf("Function %s is not exists", ident)
			log.Fatal(msg)
		}

		params := make([]value.Value, len(stmt.Parameters))
		for i, param := range stmt.Parameters {
			params[i] = e.Eval(param, env)
		}

		fn := val.(value.ValueFunction)
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
			log.Fatal(msg)
		}

		_, ok = env.Get(identifier.Name)

		if ok {
			msg := fmt.Sprintf("Symbol %s already exists", identifier.Name)
			log.Fatal(msg)
		}

		props := make([]string, 0)
		valKeyVal := make(map[string]byte)

		for _, p := range stmt.Properties {
			prop, ok := p.(ast.NodeIdentifier)

			if !ok {
				msg := fmt.Sprintf("Invalid property: %s", p)
				log.Fatal(msg)
			}

			valKeyVal[prop.Name] = 1
			props = append(props, prop.Name)
		}

		env.Set(identifier.Name, value.ValueStruct[byte]{identifier.Name, props, value.ValueKeyVal[byte]{Map: valKeyVal}})

	case ast.NodeStructExpr:
		ident, ok := stmt.Name.(ast.NodeIdentifier)

		if !ok {
			msg := fmt.Sprintf("Invalid identifier: %s", stmt.Name)
			log.Fatal(msg)
		}

		structStmt, ok := env.Get(ident.Name)
		structStmtProp := structStmt.(value.ValueStruct[byte]).Prop

		if !ok {
			msg := fmt.Sprintf("Struct %s not found", ident.Name)
			log.Fatal(msg)
		}

		keyMap := e.Eval(stmt.Values, env)
		props := keyMap.(value.ValueKeyVal[value.Value])

		for k := range props.Map {
			ok := util.InArray[string](structStmtProp, k)

			if !ok {
				msg := fmt.Sprintf("Unknown field %s on %s", k, ident.Name)
				log.Fatal(msg)
			}
		}

		return value.ValueStruct[value.Value]{ident.Name, structStmtProp, value.ValueKeyVal[value.Value]{props.Map}}

	case ast.NodeMapExpr:
		keyVal := make(map[string]value.Value)

		for k, v := range stmt.Map {
			key, ok := k.(ast.NodeIdentifier)

			if !ok {
				msg := fmt.Sprintf("Invalid struct key: %s", k)
				log.Fatal(msg)
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

	default:
		msg := fmt.Sprintf("Unrecognized statement type: %T", stmt)
		log.Fatal(msg)
	}

	return result
}
