package evaluator

import (
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

func (e *Evaluator) Eval(astNode ast.Node, env *environment.Environment) value.Value {
	var result value.Value = value.NULL

	switch stmt := astNode.(type) {

	case *ast.NodeProgram:
		for _, stmt := range stmt.Body {
			result = e.Eval(stmt, env)

			if _, ok := result.(*value.Return); ok {
				return result
			}

			if _, ok := result.(*value.Error); ok {
				return result
			}
		}

	case *ast.NodeExprStmt:
		return e.Eval(stmt.Expr, env)

	case *ast.NodeInteger:
		return &value.Int{stmt.Value}

	case *ast.NodeFloat:
		return &value.Float{stmt.Value}

	case *ast.NodeBoolean:
		return &value.Bool{stmt.Value}

	case *ast.NodeString:
		return &value.String{stmt.Value}

	case *ast.NodeBinaryExpr:
		switch stmt.Operator {

		case "+":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)
			val := left.(*value.Int).Value + right.(*value.Int).Value
			return &value.Int{val}

		case "-":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)
			val := left.(*value.Int).Value - right.(*value.Int).Value
			return &value.Int{val}

		case "*":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)
			val := left.(*value.Int).Value * right.(*value.Int).Value
			return &value.Int{val}

		case "/":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)
			val := left.(*value.Int).Value / right.(*value.Int).Value
			return &value.Int{val}

		case "%":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)
			val := left.(*value.Int).Value % right.(*value.Int).Value
			return &value.Int{val}

		case "=":
			var ident value.Value

			switch node := stmt.Left.(type) {
			case *ast.NodeIdentifier:
				ident = &value.String{node.Name}

			case *ast.NodeBinaryExpr:
				ident = e.Eval(node.Left, env)

			default:
				msg := fmt.Sprintf("Unrecognized assignment type: %s", util.TypeOf(stmt.Left))
				return e.Error(msg)
			}

			switch node := stmt.Left.(type) {
			case *ast.NodeIdentifier:
				val := e.Eval(stmt.Right, env)
				realIdent := ident.(*value.String).Value

				if _, ok := env.Get(realIdent); !ok {
					msg := fmt.Sprintf("Variable %s is not found", realIdent)
					return e.Error(msg)
				}

				env.Assign(realIdent, val)
				return val

			case *ast.NodeBinaryExpr:
				val := e.Eval(stmt.Right, env)

				receiver, ok := ident.(*value.Struct[value.Value])

				if !ok {
					msg := fmt.Sprintf("Invalid receiver type: %s", util.TypeOf(receiver))
					return e.Error(msg)
				}

				right, ok := node.Right.(*ast.NodeIdentifier)

				if !ok {
					msg := fmt.Sprintf("Invalid identifier: %s", node.Right)
					return e.Error(msg)
				}

				ident = &value.String{Value: right.Name}

				_, ok = receiver.Map[ident.(*value.String).Value]

				if !ok {
					msg := fmt.Sprintf("Symbol %s is not found", ident.(*value.String).Value)
					return e.Error(msg)
				}

				receiver.Map[ident.(*value.String).Value] = val
				return result

			default:
				msg := fmt.Sprintf("Unrecognized assignment type: %s", util.TypeOf(stmt.Left))
				return e.Error(msg)
			}

		case "<":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)

			if util.TypeOf(left) != util.TypeOf(right) {
				msg := fmt.Sprintf("Invalid operation %s < %s", left, right)
				return e.Error(msg)
			}

			switch left.(type) {

			case *value.Int:
				left, _ := left.(*value.Int)
				right, _ := right.(*value.Int)

				if left.Value < right.Value {
					return value.TRUE
				}

				return value.FALSE

			case *value.Float:
				left, _ := left.(*value.Float)
				right, _ := right.(*value.Float)

				if left.Value < right.Value {
					return value.TRUE
				}

				return value.FALSE

			default:
				msg := fmt.Sprintf("Invalid operation %s < %s", left, right)
				return e.Error(msg)
			}

		case ">":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)

			if util.TypeOf(left) != util.TypeOf(right) {
				msg := fmt.Sprintf("Invalid operation %s > %s", left, right)
				return e.Error(msg)
			}

			switch left.(type) {

			case *value.Int:
				left, _ := left.(*value.Int)
				right, _ := right.(*value.Int)

				if left.Value > right.Value {
					return value.TRUE
				}

				return value.FALSE

			case *value.Float:
				left, _ := left.(*value.Float)
				right, _ := right.(*value.Float)

				if left.Value > right.Value {
					return value.TRUE
				}

				return value.FALSE

			default:
				msg := fmt.Sprintf("Invalid operation %s > %s", left, right)
				return e.Error(msg)
			}

		case "<=":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)

			if util.TypeOf(left) != util.TypeOf(right) {
				msg := fmt.Sprintf("Invalid operation %s <= %s", left, right)
				return e.Error(msg)
			}

			switch left.(type) {

			case *value.Int:
				left, _ := left.(*value.Int)
				right, _ := right.(*value.Int)

				if left.Value <= right.Value {
					return value.TRUE
				}

				return value.FALSE

			case *value.Float:
				_left := e.Eval(stmt.Left, env)
				_right := e.Eval(stmt.Right, env)

				left, _ := _left.(*value.Float)
				right, _ := _right.(*value.Float)

				if left.Value <= right.Value {
					return value.TRUE
				}

				return value.FALSE

			default:
				msg := fmt.Sprintf("Invalid operation %s <= %s", left, right)
				return e.Error(msg)
			}

		case ">=":
			left := e.Eval(stmt.Left, env)
			right := e.Eval(stmt.Right, env)

			if util.TypeOf(left) != util.TypeOf(right) {
				msg := fmt.Sprintf("Invalid operation %s >= %s", left, right)
				return e.Error(msg)
			}

			switch left.(type) {

			case *value.Int:
				left, _ := left.(*value.Int)
				right, _ := right.(*value.Int)

				if left.Value >= right.Value {
					return value.TRUE
				}

				return value.FALSE

			case *value.Float:
				left, _ := left.(*value.Float)
				right, _ := right.(*value.Float)

				if left.Value >= right.Value {
					return value.TRUE
				}

				return value.FALSE

			default:
				msg := fmt.Sprintf("Invalid operation %s >= %s", left, right)
				return e.Error(msg)
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
			case *value.Null:
				return result

			case *value.Module:
				env = receiver.(*value.Module).Value.(*environment.Environment)
				val, ok := env.Get(right)

				if !ok {
					msg := fmt.Sprintf("Symbol %s is not found", right)
					return e.Error(msg)
				}

				return val

			case *value.Struct[value.Value]:
				val, ok := receiver.(*value.Struct[value.Value]).Map[right]

				if !ok {
					msg := fmt.Sprintf("Symbol %s is not found", right)
					return e.Error(msg)
				}

				return val

			case *value.Self:
				self, _ := env.Get(receiver.(*value.Self).Value)
				return self

			default:
				return e.Error(fmt.Sprintf("Unknown receiverInstance type %s for dot operator", util.TypeOf(receiver)))
			}

		default:
			msg := fmt.Sprintf("Unrecognized operator: %s", stmt.Operator)
			return e.Error(msg)
		}

	case *ast.NodeConstStmt:
		ident := stmt.Identifier.(*ast.NodeIdentifier).Name
		val := e.Eval(stmt.Value, env)

		if _, ok := env.Get(ident); ok {
			msg := fmt.Sprintf("Constant %s already exists", ident)
			return e.Error(msg)
		}

		env.Set(ident, val)

	case *ast.NodeLetStmt:
		ident := stmt.Identifier.(*ast.NodeIdentifier).Name
		val := e.Eval(stmt.Value, env)

		if _, ok := env.Get(ident); ok {
			msg := fmt.Sprintf("Variable %s is already exists", ident)
			return e.Error(msg)
		}

		env.Set(ident, val)

	case *ast.NodeIdentifier:
		val, ok := env.Get(stmt.Name)

		if !ok {
			msg := fmt.Sprintf("Symbol %s is not found", stmt.Name)
			return e.Error(msg)
		}

		return val

	case *ast.NodeFunctionStmt:
		var ident string
		var receiver string

		switch stmt.Identifier.(type) {
		case *ast.NodeIdentifier:
			ident = stmt.Identifier.(*ast.NodeIdentifier).Name

		case *ast.NodeBinaryExpr:
			receiver = stmt.Identifier.(*ast.NodeBinaryExpr).Left.(*ast.NodeIdentifier).Name
			ident = stmt.Identifier.(*ast.NodeBinaryExpr).Right.(*ast.NodeIdentifier).Name

		default:
			msg := fmt.Sprintf("Unrecognized function identifier type: %s", util.TypeOf(stmt.Identifier))
			return e.Error(msg)
		}

		args := make([]value.Value, len(stmt.Arguements))

		for i, _arg := range stmt.Arguements {
			switch arg := _arg.(type) {
			case *ast.NodeIdentifier:
				args[i] = &value.String{arg.Name}

			case *ast.NodeSelf:
				if i != 0 {
					msg := fmt.Sprintf("self arguement should be at position 0, detected position: %d", i)
					return e.Error(msg)
				}

				args[i] = &value.Self{arg.Name}

			default:
				msg := fmt.Sprintf("Unrecognized arguement type: %s", util.TypeOf(arg))
				return e.Error(msg)
			}
		}

		valFn := &value.Function{Args: args, Body: stmt.Body}

		if receiver != "" {
			receiverVal, ok := env.Get(receiver)

			if !ok {
				msg := fmt.Sprintf("Symbol %s is not exists", ident)
				return e.Error(msg)
			}

			_struct, ok := receiverVal.(*value.Struct[value.Value])

			if !ok {
				msg := fmt.Sprintf("Symbol %s is not a struct", ident)
				return e.Error(msg)
			}

			if util.InArray[string](_struct.Prop, ident) {
				msg := fmt.Sprintf("Symbol %s already exists", ident)
				return e.Error(msg)
			}

			// No need to set the value to the environment since
			// the struct is a pointer
			_struct.Prop = append(_struct.Prop, ident)
			_struct.KeyVal.Map[ident] = valFn
		} else {
			if _, ok := env.Get(ident); ok {
				msg := fmt.Sprintf("Symbol %s already exists", ident)
				return e.Error(msg)
			}

			env.Set(ident, valFn)
		}

	case *ast.NodeFunctionCall:
		var receiverInstance value.Value
		var receiverName string

		var identifier value.Value
		var identifierName string
		var fnEnv = environment.NewWithParent(env)

		switch node := stmt.Identifer.(type) {
		case *ast.NodeIdentifier:
			identifier = e.Eval(node, env)
			identifierName = node.Name

		case *ast.NodeBinaryExpr:
			receiverInstance = e.Eval(node.Left, env)
			receiverName = node.Left.(*ast.NodeIdentifier).Name

			ident, ok := node.Right.(*ast.NodeIdentifier)

			if !ok {
				msg := fmt.Sprintf("Invalid identifier: %s", node.Right)
				return e.Error(msg)
			}

			identifier = &value.String{ident.Name}
			identifierName = ident.Name
		}

		// Params
		params := make([]value.Value, len(stmt.Parameters))
		for i, param := range stmt.Parameters {
			params[i] = e.Eval(param, env)
		}

		if receiverInstance != nil {
			switch receiveryType := receiverInstance.(type) {

			case *value.Struct[value.Value]:
				receiverName := receiveryType.Name
				_receiver, ok := env.Get(receiverName)

				if !ok {
					msg := fmt.Sprintf("Symbol %s is not found", receiverName)
					return e.Error(msg)
				}

				receiver, ok := _receiver.(*value.Struct[value.Value])

				if !ok {
					msg := fmt.Sprintf("Symbol %s is not a valid receiver", _receiver)
					return e.Error(msg)
				}

				_valFn, ok := receiver.Map[identifier.(*value.String).Value]

				if !ok {
					msg := fmt.Sprintf("Symbol %s is not found", identifier.(*value.String).Value)
					return e.Error(msg)
				}

				valFn, ok := _valFn.(*value.Function)

				if !ok {
					msg := fmt.Sprintf("Symbol %s is not a function", identifier.(*value.String).Value)
					return e.Error(msg)
				}

				// Start bind params to args
				fnArgs := valFn.Args

				if len(valFn.Args) > 0 {
					self, ok := valFn.Args[0].(*value.Self)

					if ok {
						fnArgs = fnArgs[1:] // strip self
						fnEnv.Set(self.Value, receiverInstance)
					}

				}

				if len(fnArgs) > len(params) {
					msg := fmt.Sprintf("Bad function arguments, expected %d, got %d", len(fnArgs), len(params))
					return e.Error(msg)
				}

				// Args
				for i, _arg := range fnArgs {
					switch arg := _arg.(type) {
					case *value.String:
						fnEnv.Set(arg.Value, params[i])

					default:
						msg := fmt.Sprintf("Unrecognized arguement type: %s", util.TypeOf(arg))
						return e.Error(msg)
					}
				}

				return e.Eval(valFn.Body, fnEnv)

			case *value.Module:
				valFn, ok := receiverInstance.(*value.Module).Value.(*value.Map[value.Value]).Map[receiverName].(*value.Map[value.Value]).Map[identifierName]

				if !ok {
					msg := fmt.Sprintf("Symbol %s is not found", identifier.(*value.String).Value)
					return e.Error(msg)
				}

				fn, ok := valFn.(*value.WrapperFunction)

				if !ok {
					return result
				}

				return fn.Fn(params...)

			default:
				msg := fmt.Sprintf("Unrecognized receiver type: %s", receiveryType)
				return e.Error(msg)
			}
		}

		valFn, ok := identifier.(*value.Function)

		if !ok {
			msg := fmt.Sprintf("Identifier %s is not a function", identifierName)
			return e.Error(msg)
		}

		if len(valFn.Args) > len(params) {
			msg := fmt.Sprintf("Bad function arguments, expected %d, got %d", len(valFn.Args), len(params))
			return e.Error(msg)
		}

		for i, _arg := range valFn.Args {
			switch arg := _arg.(type) {
			case *value.String:
				fnEnv.Set(arg.Value, params[i])

			default:
				msg := fmt.Sprintf("Unrecognized arguement type: %s", util.TypeOf(arg))
				return e.Error(msg)
			}
		}

		result = e.Eval(valFn.Body, fnEnv)
		return result

	case *ast.NodeBlockStmt:
		for _, stmt := range stmt.Body {
			result = e.Eval(stmt, env)

			if _, ok := result.(*value.Return); ok {
				return result.(*value.Return).Value
			}
		}

	case *ast.NodeReturnStmt:
		return &value.Return{e.Eval(stmt.Value, env)}

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
			return e.Error(msg)
		}

		_, ok = env.Get(identifier.Name)

		if ok {
			msg := fmt.Sprintf("Symbol %s already exists", identifier.Name)
			return e.Error(msg)
		}

		props := make([]string, 0)
		valKeyVal := make(map[string]value.Value)

		for _, p := range stmt.Properties {
			prop, ok := p.(*ast.NodeIdentifier)

			if !ok {
				msg := fmt.Sprintf("Invalid property: %s", p)
				return e.Error(msg)
			}

			valKeyVal[prop.Name] = value.NULL
			props = append(props, prop.Name)
		}

		_struct := &value.Struct[value.Value]{identifier.Name, props, &value.KeyVal[value.Value]{Map: valKeyVal}}
		env.Set(identifier.Name, _struct)

	case *ast.NodeStructExpr:
		ident, ok := stmt.Name.(*ast.NodeIdentifier)

		if !ok {
			msg := fmt.Sprintf("Invalid identifier: %s", stmt.Name)
			return e.Error(msg)
		}

		structStmt, ok := env.Get(ident.Name)

		if !ok {
			msg := fmt.Sprintf("Struct %s is not found", ident.Name)
			return e.Error(msg)
		}

		structStmtProp := structStmt.(*value.Struct[value.Value]).Prop
		keyMap := e.Eval(stmt.Values, env)
		props := keyMap.(*value.Map[value.Value])

		actualProps := make([]string, 0)
		for k := range props.Map {
			ok := util.InArray[string](structStmtProp, k)

			if !ok {
				msg := fmt.Sprintf("Unknown field %s on %s", k, ident.Name)
				return e.Error(msg)
			}
			actualProps = append(actualProps, k)
		}

		return &value.Struct[value.Value]{ident.Name, actualProps, &value.KeyVal[value.Value]{props.Map}}

	case *ast.NodeMapExpr:
		keyVal := make(map[string]value.Value)

		for k, v := range stmt.Map {
			key, ok := k.(*ast.NodeIdentifier)

			if !ok {
				msg := fmt.Sprintf("Invalid struct key: %s", k)
				return e.Error(msg)
			}

			keyVal[key.Name] = e.Eval(v, env)
		}

		return &value.Map[value.Value]{&value.KeyVal[value.Value]{keyVal}}

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
		left := e.Eval(stmt.Left, env)
		ident, isIdent := stmt.Left.(*ast.NodeIdentifier)

		switch stmt.Operator {
		case "++":
			switch left.(type) {
			case *value.Int:
				if isIdent {
					env.Set(ident.Name, &value.Int{left.(*value.Int).Value + 1})
				}

				return &value.Int{left.(*value.Int).Value}

			default:
				msg := fmt.Sprintf("Unsupported operator: %s for type %s", stmt.Operator, util.TypeOf(left))
				return e.Error(msg)
			}

		case "--":
			switch left.(type) {
			case *value.Int:
				if isIdent {
					env.Set(ident.Name, &value.Int{left.(*value.Int).Value - 1})
				}

				return &value.Int{left.(*value.Int).Value}

			default:
				msg := fmt.Sprintf("Unsupported operator: %s for type %s", stmt.Operator, util.TypeOf(left))
				return e.Error(msg)
			}

		default:
			msg := fmt.Sprintf("Unsupported operator: %s", stmt.Operator)
			return e.Error(msg)
		}

	case *ast.NodePrefixExpr:
		right := e.Eval(stmt.Right, env)
		ident, isIdent := stmt.Right.(*ast.NodeIdentifier)

		switch stmt.Operator {
		case "++":
			switch right.(type) {
			case *value.Int:
				if isIdent {
					env.Set(ident.Name, &value.Int{right.(*value.Int).Value + 1})
				}

				return &value.Int{right.(*value.Int).Value + 1}

			default:
				msg := fmt.Sprintf("Unsupported operator: %s for type %s", stmt.Operator, util.TypeOf(right))
				return e.Error(msg)
			}

		case "--":
			switch right.(type) {
			case *value.Int:
				if isIdent {
					env.Set(ident.Name, &value.Int{right.(*value.Int).Value - 1})
				}

				return &value.Int{right.(*value.Int).Value - 1}

			default:
				msg := fmt.Sprintf("Unsupported operator: %s for type %s", stmt.Operator, util.TypeOf(right))
				return e.Error(msg)
			}

		case "-":
			switch right.(type) {
			case *value.Int:
				if isIdent {
					env.Set(ident.Name, &value.Int{-right.(*value.Int).Value})
				}

				return &value.Int{-right.(*value.Int).Value}

			default:
				msg := fmt.Sprintf("Unsupported operator: %s for type %s", stmt.Operator, util.TypeOf(right))
				return e.Error(msg)
			}

		default:
			msg := fmt.Sprintf("Unsupported operator: %s", stmt.Operator)
			return e.Error(msg)
		}

	case *ast.NodeImportExpr:
		pkgs := &value.Map[value.Value]{KeyVal: &value.KeyVal[value.Value]{Map: make(map[string]value.Value)}}

		// fmt package
		pkgs.Map["fmt"] = &value.Map[value.Value]{&value.KeyVal[value.Value]{stdlib.FmtFuncs}}

		// fmt io
		pkgs.Map["io"] = &value.Map[value.Value]{&value.KeyVal[value.Value]{stdlib.IoFuncs}}

		return &value.Module{pkgs}

	case *ast.NodeArrayExpr:
		values := make([]value.Value, 0)

		for _, v := range stmt.Value {
			values = append(values, e.Eval(v, env))
		}

		return &value.Array{values}

	case *ast.NodeIndexExpr:
		identifier := e.Eval(stmt.Identifier, env)

		switch node := identifier.(type) {
		case *value.Array:
			idx := e.Eval(stmt.Index, env)

			index, ok := idx.(*value.Int)

			if !ok {
				return e.Error("Array index is not an int")
			}

			if index.Value > int64(len(node.Value)) {
				return e.Error("Array index out of range")
			}

			return node.Value[index.Value]

		case *value.Map[value.Value]:
			idx := e.Eval(stmt.Index, env)

			index, ok := idx.(*value.String)

			if !ok {
				return e.Error("Map index is not a string")
			}

			val, ok := node.Map[index.Value]

			if !ok {
				msg := fmt.Sprintf("Map index %s is not found", index.Value)
				return e.Error(msg)
			}

			return val

		default:
			return e.Error(fmt.Sprintf("Unsupported index access on type %s", util.TypeOf(identifier)))
		}

	case *ast.NodeSelf:
		self, ok := env.Get(stmt.Name)

		if !ok {
			msg := fmt.Sprintf("Symbol %s is not found", stmt.Name)
			return e.Error(msg)
		}

		return self

	default:
		msg := fmt.Sprintf("Unrecognized statement type: %T", stmt)
		return e.Error(msg)
	}

	return result
}

func (e *Evaluator) Error(msg string) value.Value {
	return &value.Error{msg}
}
