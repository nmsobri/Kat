package evaluator

import (
	"errors"
	"fmt"
	"kat/ast"
	"kat/environment"
	"kat/stdlib"
	"kat/util"
	"kat/value"
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
	var result value.Value = value.NULL

	if e.IsError() {
		return result
	}

	switch stmt := node.(type) {

	case *ast.NodeProgram:
		for _, stmt := range stmt.Body {
			result = e.Eval(stmt, env)

			if _, ok := result.(*value.ValueReturn); ok {
				return result
			}
		}

	case *ast.NodeExprStmt:
		return e.Eval(stmt.Expr, env)

	case *ast.NodeInteger:
		return &value.ValueInt{stmt.Value}

	case *ast.NodeFloat:
		return &value.ValueFloat{stmt.Value}

	case *ast.NodeBoolean:
		return &value.ValueBool{stmt.Value}

	case *ast.NodeString:
		return &value.ValueString{stmt.Value}

	case *ast.NodeBinaryExpr:
		switch stmt.Operator {

		case "+":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)
			val := left.(*value.ValueInt).Value + right.(*value.ValueInt).Value
			return &value.ValueInt{val}

		case "-":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)
			val := left.(*value.ValueInt).Value - right.(*value.ValueInt).Value
			return &value.ValueInt{val}

		case "*":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)
			val := left.(*value.ValueInt).Value * right.(*value.ValueInt).Value
			return &value.ValueInt{val}

		case "/":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)
			val := left.(*value.ValueInt).Value / right.(*value.ValueInt).Value
			return &value.ValueInt{val}

		case "%":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)
			val := left.(*value.ValueInt).Value % right.(*value.ValueInt).Value
			return &value.ValueInt{val}

		case "=":
			ident := stmt.Left.(*ast.NodeIdentifier).Name
			val := e.Eval(stmt.Right, env)

			if _, ok := env.Get(ident); !ok {
				msg := fmt.Sprintf("Variable %s is not found", ident)
				e.Error(msg)
				return result
			}

			env.Assign(ident, val)
			return val

		case "<":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)

			if util.TypeOf(left) != util.TypeOf(right) {
				msg := fmt.Sprintf("Invalid operation %s < %s", left, right)
				e.Error(msg)
				return result
			}

			switch left.(type) {

			case *value.ValueInt:
				left, _ := left.(*value.ValueInt)
				right, _ := right.(*value.ValueInt)

				if left.Value < right.Value {
					return value.TRUE
				}

				return value.FALSE

			case *value.ValueFloat:
				left, _ := left.(*value.ValueFloat)
				right, _ := right.(*value.ValueFloat)

				if left.Value < right.Value {
					return value.TRUE
				}

				return value.FALSE

			default:
				msg := fmt.Sprintf("Invalid operation %s < %s", left, right)
				e.Error(msg)
				return result
			}

		case ">":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)

			if util.TypeOf(left) != util.TypeOf(right) {
				msg := fmt.Sprintf("Invalid operation %s > %s", left, right)
				e.Error(msg)
				return result
			}

			switch left.(type) {

			case *value.ValueInt:
				left, _ := left.(*value.ValueInt)
				right, _ := right.(*value.ValueInt)

				if left.Value > right.Value {
					return value.TRUE
				}

				return value.FALSE

			case *value.ValueFloat:
				left, _ := left.(*value.ValueFloat)
				right, _ := right.(*value.ValueFloat)

				if left.Value > right.Value {
					return value.TRUE
				}

				return value.FALSE

			default:
				msg := fmt.Sprintf("Invalid operation %s > %s", left, right)
				e.Error(msg)
				return result
			}

		case "<=":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)

			if util.TypeOf(left) != util.TypeOf(right) {
				msg := fmt.Sprintf("Invalid operation %s <= %s", left, right)
				e.Error(msg)
				return result
			}

			switch left.(type) {

			case *value.ValueInt:
				left, _ := left.(*value.ValueInt)
				right, _ := right.(*value.ValueInt)

				if left.Value <= right.Value {
					return value.TRUE
				}

				return value.FALSE

			case *value.ValueFloat:
				_left := e.Eval(stmt.Left, env)
				_right := e.Eval(stmt.Right, env)

				left, _ := _left.(*value.ValueFloat)
				right, _ := _right.(*value.ValueFloat)

				if left.Value <= right.Value {
					return value.TRUE
				}

				return value.FALSE

			default:
				msg := fmt.Sprintf("Invalid operation %s <= %s", left, right)
				e.Error(msg)
				return result
			}

		case ">=":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)

			if util.TypeOf(left) != util.TypeOf(right) {
				msg := fmt.Sprintf("Invalid operation %s >= %s", left, right)
				e.Error(msg)
				return result
			}

			switch left.(type) {

			case *value.ValueInt:
				left, _ := left.(*value.ValueInt)
				right, _ := right.(*value.ValueInt)

				if left.Value >= right.Value {
					return value.TRUE
				}

				return value.FALSE

			case *value.ValueFloat:
				left, _ := left.(*value.ValueFloat)
				right, _ := right.(*value.ValueFloat)

				if left.Value >= right.Value {
					return value.TRUE
				}

				return value.FALSE

			default:
				msg := fmt.Sprintf("Invalid operation %s >= %s", left, right)
				e.Error(msg)
				return result
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
			receiver := e.Eval(stmt.Left, env)
			right := stmt.Right.(*ast.NodeIdentifier).Name
			var env *environment.Environment

			switch receiver.(type) {
			case *value.ValueNull:
				return result

			case *value.ValueModule:
				env = receiver.(*value.ValueModule).Value.(*environment.Environment)

			case *value.ValueStruct[value.Value]:
				val, ok := receiver.(*value.ValueStruct[value.Value]).Map[right]

				if !ok {
					msg := fmt.Sprintf("Symbol %s is not found", right)
					e.Error(msg)
					return result
				}

				return val

			default:
				e.Error(fmt.Sprintf("Unknown receiver %s type for dot operator", util.TypeOf(receiver)))
				return result
			}

			val, ok := env.Get(right)

			if !ok {
				msg := fmt.Sprintf("Symbol %s is not found", right)
				e.Error(msg)
				return result
			}

			return val

		default:
			msg := fmt.Sprintf("Unrecognized operator: %s", stmt.Operator)
			e.Error(msg)
			return result
		}

	case *ast.NodeConstStmt:
		ident := stmt.Identifier.(*ast.NodeIdentifier).Name
		val := e.Eval(stmt.Value, env)

		if _, ok := env.Get(ident); ok {
			msg := fmt.Sprintf("Constant %s already exists", ident)
			e.Error(msg)
			return result
		}

		env.Set(ident, val)

	case *ast.NodeLetStmt:
		ident := stmt.Identifier.(*ast.NodeIdentifier).Name
		val := e.Eval(stmt.Value, env)

		if _, ok := env.Get(ident); ok {
			msg := fmt.Sprintf("Variable %s is already exists", ident)
			e.Error(msg)
			return result
		}

		env.Set(ident, val)

	case *ast.NodeIdentifier:
		val, ok := env.Get(stmt.Name)

		if !ok {
			msg := fmt.Sprintf("Symbol %s is not found", stmt.Name)
			e.Error(msg)
			return result
		}

		return val

	case *ast.NodeFunctionStmt:
		ident := stmt.Identifier.(*ast.NodeIdentifier).Name

		args := make([]value.Value, len(stmt.Arguements))
		for i, arg := range stmt.Arguements {
			args[i] = &value.ValueString{arg.(*ast.NodeIdentifier).Name}
		}

		val := &value.ValueFunction{Args: args, Body: stmt.Body}

		if _, ok := env.Get(ident); ok {
			msg := fmt.Sprintf("Symbol %s already exists", ident)
			e.Error(msg)
			return result
		}

		env.Set(ident, val)

	case *ast.NodeFunctionCall:
		val := e.Eval(stmt.Identifer, env)
		params := make([]value.Value, len(stmt.Parameters))

		for i, param := range stmt.Parameters {
			params[i] = e.Eval(param, env)
		}

		//fn, ok := val.(*value.ValueFunction)
		fn, ok := val.(*value.WrapperFunction)

		if !ok {
			return result
		}

		return fn.Fn(params...)

		//newEnv := environment.NewWithParent(env)
		//
		//for i, arg := range fn.Args {
		//	newEnv.Set(arg.String(), params[i])
		//}
		//
		//result := e.Eval(fn.Body, newEnv)
		//return result

	case *ast.NodeBlockStmt:
		for _, stmt := range stmt.Body {
			result = e.Eval(stmt, env)

			if _, ok := result.(*value.ValueReturn); ok {
				return result.(*value.ValueReturn).Value
			}
		}

	case *ast.NodeReturnStmt:
		return &value.ValueReturn{e.Eval(stmt.Value, env)}

	case *ast.NodeConditionalStmt:
		condition := e.Eval(stmt.Condition, env)

		if util.IsTruthy(condition) {
			return e.Eval(stmt.ThenArm, env)
		} else {
			return e.Eval(stmt.ElseArm, env)
		}

	case *ast.NodeStructStmt:
		identifier, ok := stmt.Identifier.(*ast.NodeIdentifier)

		if !ok {
			msg := fmt.Sprintf("Invalid identifier: %s", stmt.Identifier)
			e.Error(msg)
			return result
		}

		_, ok = env.Get(identifier.Name)

		if ok {
			msg := fmt.Sprintf("Symbol %s already exists", identifier.Name)
			e.Error(msg)
			return result
		}

		props := make([]string, 0)
		valKeyVal := make(map[string]byte)

		for _, p := range stmt.Properties {
			prop, ok := p.(*ast.NodeIdentifier)

			if !ok {
				msg := fmt.Sprintf("Invalid property: %s", p)
				e.Error(msg)
				return result
			}

			valKeyVal[prop.Name] = 1
			props = append(props, prop.Name)
		}

		env.Set(identifier.Name, &value.ValueStruct[byte]{identifier.Name, props, &value.ValueKeyVal[byte]{Map: valKeyVal}})

	case *ast.NodeStructExpr:
		ident, ok := stmt.Name.(*ast.NodeIdentifier)

		if !ok {
			msg := fmt.Sprintf("Invalid identifier: %s", stmt.Name)
			e.Error(msg)
			return result
		}

		structStmt, ok := env.Get(ident.Name)
		structStmtProp := structStmt.(*value.ValueStruct[byte]).Prop

		if !ok {
			msg := fmt.Sprintf("Struct %s is not found", ident.Name)
			e.Error(msg)
			return result
		}

		keyMap := e.Eval(stmt.Values, env)
		props := keyMap.(*value.ValueKeyVal[value.Value])

		for k := range props.Map {
			ok := util.InArray[string](structStmtProp, k)

			if !ok {
				msg := fmt.Sprintf("Unknown field %s on %s", k, ident.Name)
				e.Error(msg)
				return result
			}
		}

		return &value.ValueStruct[value.Value]{ident.Name, structStmtProp, &value.ValueKeyVal[value.Value]{props.Map}}

	case *ast.NodeMapExpr:
		keyVal := make(map[string]value.Value)

		for k, v := range stmt.Map {
			key, ok := k.(*ast.NodeIdentifier)

			if !ok {
				msg := fmt.Sprintf("Invalid struct key: %s", k)
				e.Error(msg)
				return result
			}

			keyVal[key.Name] = e.Eval(v, env)
		}

		return &value.ValueKeyVal[value.Value]{keyVal}

	case *ast.NodeModernForStmt:
		condition := e.Eval(stmt.Condition, env)

		for util.IsTruthy(condition) {
			e.Eval(stmt.Body, env)
			condition = e.Eval(stmt.Condition, env)
		}

	case *ast.NodeClassicForStmt:
		newEnv := environment.NewWithParent(env)

		e.Eval(stmt.PreExpr, newEnv) // pre expression
		condition := e.Eval(stmt.Condition, newEnv)

		for util.IsTruthy(condition) {
			e.Eval(stmt.Body, newEnv)
			e.Eval(stmt.PostExpr, newEnv) // post expression
			condition = e.Eval(stmt.Condition, newEnv)
		}

	case *ast.NodePostfixExpr:
		left, ok := stmt.Left.(*ast.NodeIdentifier)

		if !ok {
			msg := fmt.Sprintf("Invalid identifier: %s", stmt.Left)
			e.Error(msg)
			return result
		}

		leftVal, ok := env.Get(left.Name)

		if !ok {
			msg := fmt.Sprintf("Variable %s is not found", left.Name)
			e.Error(msg)
			return result
		}

		switch stmt.Operator {
		case "++":
			switch leftVal.(type) {
			case *value.ValueInt:
				env.Set(left.Name, &value.ValueInt{leftVal.(*value.ValueInt).Value + 1})
				return &value.ValueInt{leftVal.(*value.ValueInt).Value}

			default:
				msg := fmt.Sprintf("Unsupported operator: %s for type %T", stmt.Operator, leftVal)
				e.Error(msg)
				return result
			}

		case "--":
			switch leftVal.(type) {
			case *value.ValueInt:
				env.Set(left.Name, &value.ValueInt{leftVal.(*value.ValueInt).Value - 1})
				return &value.ValueInt{leftVal.(*value.ValueInt).Value}

			default:
				msg := fmt.Sprintf("Unsupported operator: %s for type %T", stmt.Operator, leftVal)
				e.Error(msg)
				return result
			}

		default:
			msg := fmt.Sprintf("Unsupported operator: %s", stmt.Operator)
			e.Error(msg)
			return result
		}

	case *ast.NodePrefixExpr:
		right, ok := stmt.Right.(*ast.NodeIdentifier)

		if !ok {
			msg := fmt.Sprintf("Invalid identifier: %s", stmt.Right)
			e.Error(msg)
			return result
		}

		rightVal, ok := env.Get(right.Name)

		if !ok {
			msg := fmt.Sprintf("Variable %s is not found", right.Name)
			e.Error(msg)
			return result
		}

		switch stmt.Operator {
		case "++":
			switch rightVal.(type) {
			case *value.ValueInt:
				env.Set(right.Name, &value.ValueInt{rightVal.(*value.ValueInt).Value + 1})
				return &value.ValueInt{rightVal.(*value.ValueInt).Value + 1}

			default:
				msg := fmt.Sprintf("Unsupported operator: %s for type %T", stmt.Operator, rightVal)
				e.Error(msg)
				return result
			}

		case "--":
			switch rightVal.(type) {
			case *value.ValueInt:
				env.Set(right.Name, &value.ValueInt{rightVal.(*value.ValueInt).Value - 1})
				return &value.ValueInt{rightVal.(*value.ValueInt).Value - 1}

			default:
				msg := fmt.Sprintf("Unsupported operator: %s for type %T", stmt.Operator, rightVal)
				e.Error(msg)
				return result
			}

		default:
			msg := fmt.Sprintf("Unsupported operator: %s", stmt.Operator)
			e.Error(msg)
			return result
		}

	case *ast.NodeImportExpr:
		//module := stmt.Path.(*ast.NodeString).Value
		env := environment.New()
		env.Set("Print", &value.WrapperFunction{Name: "Print", Fn: stdlib.Print})
		env.Set("Println", &value.WrapperFunction{Name: "Print", Fn: stdlib.Println})
		env.Set("Printf", &value.WrapperFunction{Name: "Print", Fn: stdlib.Printf})
		env.Set("Sprintf", &value.WrapperFunction{Name: "Print", Fn: stdlib.Sprintf})

		return &value.ValueModule{env}

	case *ast.NodeArrayExpr:
		values := make([]value.Value, 0)

		for _, v := range stmt.Value {
			values = append(values, e.Eval(v, env))
		}

		return &value.ValueArray{values}

	case *ast.NodeIndexExpr:
		identifier := e.Eval(stmt.Identifier, env)
		array, ok := identifier.(*value.ValueArray)

		if !ok {
			ident := stmt.Identifier.(*ast.NodeIdentifier)
			e.Error(fmt.Sprintf("Identifier: %s is not an array", ident.Name))
			return result
		}

		idx := e.Eval(stmt.Index, env)

		index, ok := idx.(*value.ValueInt)

		if !ok {
			e.Error("Array index is not an int")
			return result
		}

		if index.Value > int64(len(array.Value)) {
			e.Error("Array index out of range")
			return result
		}

		return array.Value[index.Value]

	default:
		msg := fmt.Sprintf("Unrecognized statement type: %T", stmt)
		e.Error(msg)
		return result
	}

	return result
}

func (e *Evaluator) Error(msg string) {
	e.Errors = append(e.Errors, errors.New(msg))
}

func (e *Evaluator) IsError() bool {
	return len(e.Errors) > 0
}
