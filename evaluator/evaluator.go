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

func (e *Evaluator) Eval(astNode ast.Node, env *environment.Environment) value.Value {
	var result value.Value = value.NULL

	if e.IsError() {
		return result
	}

	switch stmt := astNode.(type) {

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
			var ident value.Value

			switch node := stmt.Left.(type) {
			case *ast.NodeIdentifier:
				ident = &value.ValueString{node.Name}

			case *ast.NodeBinaryExpr:
				ident = e.Eval(node.Left, env)

			default:
				msg := fmt.Sprintf("Unrecognized assignment type: %s", util.TypeOf(stmt.Left))
				e.Error(msg)
				return result
			}

			switch node := stmt.Left.(type) {
			case *ast.NodeIdentifier:
				val := e.Eval(stmt.Right, env)
				realIdent := ident.(*value.ValueString).Value

				if _, ok := env.Get(realIdent); !ok {
					msg := fmt.Sprintf("Variable %s is not found", realIdent)
					e.Error(msg)
					return result
				}

				env.Assign(realIdent, val)
				return val

			case *ast.NodeBinaryExpr:
				val := e.Eval(stmt.Right, env)

				receiver, ok := ident.(*value.ValueStruct[value.Value])

				if !ok {
					msg := fmt.Sprintf("Invalid receiver type: %s", util.TypeOf(receiver))
					e.Error(msg)
					return result
				}

				right, ok := node.Right.(*ast.NodeIdentifier)

				if !ok {
					msg := fmt.Sprintf("Invalid identifier: %s", node.Right)
					e.Error(msg)
					return result
				}

				ident = &value.ValueString{Value: right.Name}

				_, ok = receiver.Map[ident.(*value.ValueString).Value]

				if !ok {
					msg := fmt.Sprintf("Symbol %s is not found", ident.(*value.ValueString).Value)
					e.Error(msg)
					return result
				}

				receiver.Map[ident.(*value.ValueString).Value] = val
				return result

			default:
				msg := fmt.Sprintf("Unrecognized assignment type: %s", util.TypeOf(stmt.Left))
				e.Error(msg)
				return result
			}

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

			case *value.ValueSelf:
				self, _ := env.Get(receiver.(*value.ValueSelf).Value)
				return self

			default:
				e.Error(fmt.Sprintf("Unknown receiverInstance type %s for dot operator", util.TypeOf(receiver)))
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
			e.Error(msg)
			return result
		}

		args := make([]value.Value, len(stmt.Arguements))

		for i, _arg := range stmt.Arguements {
			switch arg := _arg.(type) {
			case *ast.NodeIdentifier:
				args[i] = &value.ValueString{arg.Name}

			case *ast.NodeSelf:
				if i != 0 {
					msg := fmt.Sprintf("self arguement should be at position 0, detected position: %d", i)
					e.Error(msg)
					return result
				}

				args[i] = &value.ValueSelf{arg.Name}

			default:
				msg := fmt.Sprintf("Unrecognized arguement type: %s", util.TypeOf(arg))
				e.Error(msg)
				return result
			}
		}

		valFn := &value.ValueFunction{Args: args, Body: stmt.Body}

		if receiver != "" {
			receiverVal, ok := env.Get(receiver)

			if !ok {
				msg := fmt.Sprintf("Symbol %s is not exists", ident)
				e.Error(msg)
				return result
			}

			_struct, ok := receiverVal.(*value.ValueStruct[value.Value])

			if !ok {
				msg := fmt.Sprintf("Symbol %s is not a struct", ident)
				e.Error(msg)
				return result
			}

			if util.InArray[string](_struct.Prop, ident) {
				msg := fmt.Sprintf("Symbol %s already exists", ident)
				e.Error(msg)
				return result
			}

			// No need to set the value to the environment since
			// the struct is a pointer
			_struct.Prop = append(_struct.Prop, ident)
			_struct.ValueKeyVal.Map[ident] = valFn
		} else {
			if _, ok := env.Get(ident); ok {
				msg := fmt.Sprintf("Symbol %s already exists", ident)
				e.Error(msg)
				return result
			}

			env.Set(ident, valFn)
		}

	case *ast.NodeFunctionCall:
		var identifierName string
		var receiverInstance value.Value
		var identifier value.Value
		var fnEnv = environment.NewWithParent(env)

		switch node := stmt.Identifer.(type) {
		case *ast.NodeIdentifier:
			identifier = e.Eval(node, env)
			identifierName = node.Name

		case *ast.NodeBinaryExpr:
			receiverInstance = e.Eval(node.Left, env)

			ident, ok := node.Right.(*ast.NodeIdentifier)

			if !ok {
				msg := fmt.Sprintf("Invalid identifier: %s", node.Right)
				e.Error(msg)
				return result
			}

			identifier = &value.ValueString{ident.Name}
			identifierName = ident.Name
		}

		// Params
		params := make([]value.Value, len(stmt.Parameters))
		for i, param := range stmt.Parameters {
			params[i] = e.Eval(param, env)
		}

		if receiverInstance != nil {
			switch receiveryType := receiverInstance.(type) {

			case *value.ValueStruct[value.Value]:
				receiverName := receiveryType.Name
				_receiver, ok := env.Get(receiverName)

				if !ok {
					msg := fmt.Sprintf("Symbol %s is not found", receiverName)
					e.Error(msg)
					return result
				}

				receiver, ok := _receiver.(*value.ValueStruct[value.Value])

				if !ok {
					msg := fmt.Sprintf("Symbol %s is not a valid receiver", _receiver)
					e.Error(msg)
					return result
				}

				_valFn, ok := receiver.Map[identifier.(*value.ValueString).Value]

				if !ok {
					msg := fmt.Sprintf("Symbol %s is not found", identifier.(*value.ValueString).Value)
					e.Error(msg)
					return result
				}

				valFn, ok := _valFn.(*value.ValueFunction)

				if !ok {
					msg := fmt.Sprintf("Symbol %s is not a function", identifier.(*value.ValueString).Value)
					e.Error(msg)
					return result
				}

				// Start bind params to args
				fnArgs := valFn.Args

				if len(valFn.Args) > 0 {
					self, ok := valFn.Args[0].(*value.ValueSelf)

					if ok {
						fnArgs = fnArgs[1:] // strip self
						fnEnv.Set(self.Value, receiverInstance)
					}

				}

				if len(fnArgs) > len(params) {
					msg := fmt.Sprintf("Bad function arguments, expected %d, got %d", len(fnArgs), len(params))
					e.Error(msg)
					return result
				}

				// Args
				for i, _arg := range fnArgs {
					switch arg := _arg.(type) {
					case *value.ValueString:
						fnEnv.Set(arg.Value, params[i])

					default:
						msg := fmt.Sprintf("Unrecognized arguement type: %s", util.TypeOf(arg))
						e.Error(msg)
						return result
					}
				}

				return e.Eval(valFn.Body, fnEnv)

			case *value.ValueModule:
				valFn, ok := receiverInstance.(*value.ValueModule).Value.(*environment.Environment).Get(identifier.(*value.ValueString).Value)

				if !ok {
					msg := fmt.Sprintf("Symbol %s is not found", identifier.(*value.ValueString).Value)
					e.Error(msg)
					return result
				}

				fn, ok := valFn.(*value.WrapperFunction)

				if !ok {
					return result
				}

				return fn.Fn(params...)

			default:
				msg := fmt.Sprintf("Unrecognized receiver type: %s", receiveryType)
				e.Error(msg)
				return result
			}
		}

		valFn, ok := identifier.(*value.ValueFunction)

		if !ok {
			msg := fmt.Sprintf("Identifier %s is not a function", identifierName)
			e.Error(msg)
			return result
		}

		if len(valFn.Args) > len(params) {
			msg := fmt.Sprintf("Bad function arguments, expected %d, got %d", len(valFn.Args), len(params))
			e.Error(msg)
			return result
		}

		for i, _arg := range valFn.Args {
			switch arg := _arg.(type) {
			case *value.ValueString:
				fnEnv.Set(arg.Value, params[i])

			default:
				msg := fmt.Sprintf("Unrecognized arguement type: %s", util.TypeOf(arg))
				e.Error(msg)
				return result
			}
		}

		result = e.Eval(valFn.Body, fnEnv)
		return result

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
		valKeyVal := make(map[string]value.Value)

		for _, p := range stmt.Properties {
			prop, ok := p.(*ast.NodeIdentifier)

			if !ok {
				msg := fmt.Sprintf("Invalid property: %s", p)
				e.Error(msg)
				return result
			}

			valKeyVal[prop.Name] = value.NULL
			props = append(props, prop.Name)
		}

		_struct := &value.ValueStruct[value.Value]{identifier.Name, props, &value.ValueKeyVal[value.Value]{Map: valKeyVal}}
		env.Set(identifier.Name, _struct)

	case *ast.NodeStructExpr:
		ident, ok := stmt.Name.(*ast.NodeIdentifier)

		if !ok {
			msg := fmt.Sprintf("Invalid identifier: %s", stmt.Name)
			e.Error(msg)
			return result
		}

		structStmt, ok := env.Get(ident.Name)

		if !ok {
			msg := fmt.Sprintf("Struct %s is not found", ident.Name)
			e.Error(msg)
			return result
		}

		structStmtProp := structStmt.(*value.ValueStruct[value.Value]).Prop
		keyMap := e.Eval(stmt.Values, env)
		props := keyMap.(*value.ValueKeyVal[value.Value])

		actualProps := make([]string, 0)
		for k := range props.Map {
			ok := util.InArray[string](structStmtProp, k)

			if !ok {
				msg := fmt.Sprintf("Unknown field %s on %s", k, ident.Name)
				e.Error(msg)
				return result
			}
			actualProps = append(actualProps, k)
		}

		return &value.ValueStruct[value.Value]{ident.Name, actualProps, &value.ValueKeyVal[value.Value]{props.Map}}

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
		right := e.Eval(stmt.Right, env)
		ident, isIdent := stmt.Right.(*ast.NodeIdentifier)

		switch stmt.Operator {
		case "++":
			switch right.(type) {
			case *value.ValueInt:
				if isIdent {
					env.Set(ident.Name, &value.ValueInt{right.(*value.ValueInt).Value + 1})
				}

				return &value.ValueInt{right.(*value.ValueInt).Value + 1}

			default:
				msg := fmt.Sprintf("Unsupported operator: %s for type %s", stmt.Operator, util.TypeOf(right))
				e.Error(msg)
				return result
			}

		case "--":
			switch right.(type) {
			case *value.ValueInt:
				if isIdent {
					env.Set(ident.Name, &value.ValueInt{right.(*value.ValueInt).Value - 1})
				}

				return &value.ValueInt{right.(*value.ValueInt).Value - 1}

			default:
				msg := fmt.Sprintf("Unsupported operator: %s for type %s", stmt.Operator, util.TypeOf(right))
				e.Error(msg)
				return result
			}

		case "-":
			switch right.(type) {
			case *value.ValueInt:
				if isIdent {
					env.Set(ident.Name, &value.ValueInt{-right.(*value.ValueInt).Value})
				}

				return &value.ValueInt{-right.(*value.ValueInt).Value}

			default:
				msg := fmt.Sprintf("Unsupported operator: %s for type %s", stmt.Operator, util.TypeOf(right))
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

	case *ast.NodeSelf:
		self, ok := env.Get(stmt.Name)

		if !ok {
			msg := fmt.Sprintf("Symbol %s is not found", stmt.Name)
			e.Error(msg)
			return result
		}

		return self

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
