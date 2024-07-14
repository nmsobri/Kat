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

var Pkgs = &value.Map[value.Value]{KeyVal: &value.KeyVal[value.Value]{Map: make(map[string]value.Value)}}

func init() {
	// Register the standard library

	// fmt package
	Pkgs.Map["fmt"] = &value.Map[value.Value]{&value.KeyVal[value.Value]{stdlib.FmtFuncs}}

	// fmt io
	Pkgs.Map["io"] = &value.Map[value.Value]{&value.KeyVal[value.Value]{stdlib.IoFuncs}}
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
		return e.EvalProgram(stmt, env)

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
		return e.EvaluateBinaryExpr(stmt, env)

	case *ast.NodeConstStmt:
		return e.EvaluateConstStmt(stmt, env)

	case *ast.NodeLetStmt:
		return e.EvaluateLetStmt(stmt, env)

	case *ast.NodeIdentifier:
		return e.EvaluateIdentifier(stmt, env)

	case *ast.NodeFunctionStmt:
		return e.EvaluateFunctionStmt(stmt, env)

	case *ast.NodeFunctionCall:
		return e.EvalFunctionCall(stmt, env)

	case *ast.NodeBlockStmt:
		return e.EvalBlockStmt(stmt, env)

	case *ast.NodeReturnStmt:
		return e.EvalReturnStmt(result, stmt, env)

	case *ast.NodeConditionalStmt:
		return e.EvalConditionalStmt(stmt, env)

	case *ast.NodeStructStmt:
		return e.EvalStructStmt(stmt, env)

	case *ast.NodeStructExpr:
		return e.EvalStructExpr(stmt, env)

	case *ast.NodeMapExpr:
		return e.EvalMapExpr(stmt, env)

	case *ast.NodeModernForStmt:
		return e.EvalModernForStmt(stmt, env)

	case *ast.NodeClassicForStmt:
		return e.EvalClassicForStmt(stmt, env)

	case *ast.NodePostfixExpr:
		return e.EvalPostfixExpr(stmt, env)

	case *ast.NodePrefixExpr:
		return e.EvalPrefixExpr(stmt, env)

	case *ast.NodeImportExpr:
		return e.EvalImportExpr(stmt, env)

	case *ast.NodeArrayExpr:
		return e.EvalArrayExpr(stmt, env)

	case *ast.NodeIndexExpr:
		return e.EvalIndexExpr(stmt, env)

	case *ast.NodeSelf:
		return e.EvalSelf(stmt, env)

	default:
		msg := fmt.Sprintf("Unrecognized statement type: %T", stmt)
		return &value.Error{msg}
	}
}

func (e *Evaluator) EvalIndexExpr(stmt *ast.NodeIndexExpr, env *environment.Environment) value.Value {
	identifier := e.Eval(stmt.Identifier, env)

	if e.Error(identifier) {
		return identifier
	}

	switch node := identifier.(type) {
	case *value.Array:
		idx := e.Eval(stmt.Index, env)

		if e.Error(idx) {
			return idx
		}

		index, ok := idx.(*value.Int)

		if !ok {
			msg := "Array index is not an int"
			return &value.Error{msg}
		}

		if index.Value > int64(len(node.Value)) {
			msg := "Array index out of range"
			return &value.Error{msg}
		}

		return node.Value[index.Value]

	case *value.Map[value.Value]:
		idx := e.Eval(stmt.Index, env)

		if e.Error(idx) {
			return idx
		}

		index, ok := idx.(*value.String)

		if !ok {
			msg := "Map index is not a string"
			return &value.Error{msg}
		}

		val, ok := node.Map[index.Value]

		if !ok {
			msg := fmt.Sprintf("Map index %s is not found", index.Value)
			return &value.Error{msg}
		}

		return val

	default:
		msg := fmt.Sprintf("Unsupported index access on type %s", util.TypeOf(identifier))
		return &value.Error{msg}
	}
}

func (e *Evaluator) EvalSelf(stmt *ast.NodeSelf, env *environment.Environment) value.Value {
	self, ok := env.Get(stmt.Name)

	if !ok {
		msg := fmt.Sprintf("Symbol %s is not found", stmt.Name)
		return &value.Error{msg}
	}

	return self
}

func (e *Evaluator) EvalArrayExpr(stmt *ast.NodeArrayExpr, env *environment.Environment) value.Value {
	values := make([]value.Value, 0)

	for _, v := range stmt.Value {
		val := e.Eval(v, env)

		if e.Error(val) {
			return val
		}

		values = append(values, val)
	}

	return &value.Array{values}
}

func (e *Evaluator) EvalImportExpr(stmt *ast.NodeImportExpr, env *environment.Environment) value.Value {
	path := e.Eval(stmt.Path, env)

	if e.Error(path) {
		return path
	}

	pkg, ok := Pkgs.Map[path.(*value.String).Value]

	if !ok {
		msg := fmt.Sprintf("Package %s not found", path.(*value.String).Value)
		return &value.Error{msg}
	}

	return &value.Module{pkg}
}

func (e *Evaluator) EvalPrefixExpr(stmt *ast.NodePrefixExpr, env *environment.Environment) value.Value {
	right := e.Eval(stmt.Right, env)

	if e.Error(right) {
		return right
	}

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
			return &value.Error{msg}
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
			return &value.Error{msg}
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
			return &value.Error{msg}
		}

	default:
		msg := fmt.Sprintf("Unsupported operator: %s", stmt.Operator)
		return &value.Error{msg}
	}
}

func (e *Evaluator) EvalPostfixExpr(stmt *ast.NodePostfixExpr, env *environment.Environment) value.Value {
	left := e.Eval(stmt.Left, env)

	if e.Error(left) {
		return left
	}

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
			return &value.Error{msg}
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
			return &value.Error{msg}
		}

	default:
		msg := fmt.Sprintf("Unsupported operator: %s", stmt.Operator)
		return &value.Error{msg}
	}
}

func (e *Evaluator) EvalClassicForStmt(stmt *ast.NodeClassicForStmt, env *environment.Environment) value.Value {
	var result = value.NULL
	newEnv := environment.NewWithParent(env)

	e.Eval(stmt.PreExpr, newEnv) // pre expression

	condition := e.Eval(stmt.Condition, newEnv)

	if e.Error(condition) {
		return condition
	}

	for util.IsTruthy(condition) {
		e.Eval(stmt.Body, newEnv)
		e.Eval(stmt.PostExpr, newEnv) // post expression

		condition = e.Eval(stmt.Condition, newEnv)

		if e.Error(condition) {
			return condition
		}
	}

	return result
}

func (e *Evaluator) EvalModernForStmt(stmt *ast.NodeModernForStmt, env *environment.Environment) value.Value {
	var result = value.NULL
	condition := e.Eval(stmt.Condition, env)

	if e.Error(condition) {
		return condition
	}

	for util.IsTruthy(condition) {
		e.Eval(stmt.Body, env)

		condition = e.Eval(stmt.Condition, env)

		if e.Error(condition) {
			return condition
		}
	}
	return result
}

func (e *Evaluator) EvalMapExpr(stmt *ast.NodeMapExpr, env *environment.Environment) value.Value {
	keyVal := make(map[string]value.Value)

	for k, v := range stmt.Map {
		key, ok := k.(*ast.NodeIdentifier)

		if !ok {
			msg := fmt.Sprintf("Invalid struct key: %s", k)
			return &value.Error{msg}
		}

		keyVal[key.Name] = e.Eval(v, env)

		if e.Error(keyVal[key.Name]) {
			return keyVal[key.Name]
		}
	}

	return &value.Map[value.Value]{&value.KeyVal[value.Value]{keyVal}}
}

func (e *Evaluator) EvalStructExpr(stmt *ast.NodeStructExpr, env *environment.Environment) value.Value {
	ident, ok := stmt.Name.(*ast.NodeIdentifier)

	if !ok {
		msg := fmt.Sprintf("Invalid identifier: %s", stmt.Name)
		return &value.Error{msg}
	}

	structStmt, ok := env.Get(ident.Name)

	if !ok {
		msg := fmt.Sprintf("Struct %s is not found", ident.Name)
		return &value.Error{msg}
	}

	structStmtProp := structStmt.(*value.Struct[value.Value]).Prop
	keyMap := e.Eval(stmt.Values, env)

	if e.Error(keyMap) {
		return keyMap
	}

	props := keyMap.(*value.Map[value.Value])

	actualProps := make([]string, 0)
	for k := range props.Map {
		ok := util.InArray[string](structStmtProp, k)

		if !ok {
			msg := fmt.Sprintf("Unknown field %s on %s", k, ident.Name)
			return &value.Error{msg}
		}
		actualProps = append(actualProps, k)
	}

	return &value.Struct[value.Value]{ident.Name, actualProps, &value.KeyVal[value.Value]{props.Map}}
}

func (e *Evaluator) EvalStructStmt(stmt *ast.NodeStructStmt, env *environment.Environment) value.Value {
	var result value.Value = value.NULL
	identifier, ok := stmt.Identifier.(*ast.NodeIdentifier)

	if !ok {
		msg := fmt.Sprintf("Invalid identifier: %s", stmt.Identifier)
		return &value.Error{msg}
	}

	_, ok = env.Get(identifier.Name)

	if ok {
		msg := fmt.Sprintf("Symbol %s already exists", identifier.Name)
		return &value.Error{msg}
	}

	props := make([]string, 0)
	valKeyVal := make(map[string]value.Value)

	for _, p := range stmt.Properties {
		prop, ok := p.(*ast.NodeIdentifier)

		if !ok {
			msg := fmt.Sprintf("Invalid property: %s", p)
			return &value.Error{msg}
		}

		valKeyVal[prop.Name] = value.NULL
		props = append(props, prop.Name)
	}

	_struct := &value.Struct[value.Value]{identifier.Name, props, &value.KeyVal[value.Value]{Map: valKeyVal}}
	env.Set(identifier.Name, _struct)
	return result
}

func (e *Evaluator) EvalConditionalStmt(stmt *ast.NodeConditionalStmt, env *environment.Environment) value.Value {
	condition := e.Eval(stmt.Condition, env)

	if e.Error(condition) {
		return condition
	}

	if util.IsTruthy(condition) {
		return e.Eval(stmt.ThenArm, env)
	} else {
		return e.Eval(stmt.ElseArm, env)
	}
}

func (e *Evaluator) EvalReturnStmt(result value.Value, stmt ast.Node, env *environment.Environment) value.Value {
	result = e.Eval(stmt, env)

	if e.Error(result) {
		return result
	}

	return &value.Return{result}
}

func (e *Evaluator) EvalBlockStmt(stmt *ast.NodeBlockStmt, env *environment.Environment) value.Value {
	var result value.Value = value.NULL

	for _, stmt := range stmt.Body {
		result = e.Eval(stmt, env)

		if e.Error(result) {
			return result
		}

		if _, ok := result.(*value.Return); ok {
			return result.(*value.Return).Value
		}
	}
	return result
}

func (e *Evaluator) EvalFunctionCall(stmt *ast.NodeFunctionCall, env *environment.Environment) value.Value {
	var receiverInstance value.Value
	var identifier value.Value
	var identifierName string
	var fnEnv = environment.NewWithParent(env)
	var result value.Value = value.NULL

	switch node := stmt.Identifer.(type) {
	case *ast.NodeIdentifier:
		identifier = e.Eval(node, env)
		if e.Error(identifier) {
			return identifier
		}

		identifierName = node.Name

	case *ast.NodeBinaryExpr:
		receiverInstance = e.Eval(node.Left, env)
		if e.Error(receiverInstance) {
			return receiverInstance
		}

		ident, ok := node.Right.(*ast.NodeIdentifier)

		if !ok {
			msg := fmt.Sprintf("Invalid identifier: %s", node.Right)
			return &value.Error{msg}
		}

		identifier = &value.String{ident.Name}
		identifierName = ident.Name
	}

	// Params
	params := make([]value.Value, len(stmt.Parameters))
	for i, param := range stmt.Parameters {
		params[i] = e.Eval(param, env)
		if e.Error(params[i]) {
			return params[i]
		}
	}

	if receiverInstance != nil {
		switch receiveryType := receiverInstance.(type) {

		case *value.Struct[value.Value]:
			receiverName := receiveryType.Name
			_receiver, ok := env.Get(receiverName)

			if !ok {
				msg := fmt.Sprintf("Symbol %s is not found", receiverName)
				return &value.Error{msg}
			}

			receiver, ok := _receiver.(*value.Struct[value.Value])

			if !ok {
				msg := fmt.Sprintf("Symbol %s is not a valid receiver", _receiver)
				return &value.Error{msg}
			}

			_valFn, ok := receiver.Map[identifier.(*value.String).Value]

			if !ok {
				msg := fmt.Sprintf("Symbol %s is not found", identifier.(*value.String).Value)
				return &value.Error{msg}
			}

			valFn, ok := _valFn.(*value.Function)

			if !ok {
				msg := fmt.Sprintf("Symbol %s is not a function", identifier.(*value.String).Value)
				return &value.Error{msg}
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
				return &value.Error{msg}
			}

			// Args
			for i, _arg := range fnArgs {
				switch arg := _arg.(type) {
				case *value.String:
					fnEnv.Set(arg.Value, params[i])

				default:
					msg := fmt.Sprintf("Unrecognized arguement type: %s", util.TypeOf(arg))
					return &value.Error{msg}
				}
			}

			return e.Eval(valFn.Body, fnEnv)

		case *value.Module:
			valFn, ok := receiverInstance.(*value.Module).Value.(*value.Map[value.Value]).Map[identifierName]

			if !ok {
				msg := fmt.Sprintf("Symbol %s is not found", identifier.(*value.String).Value)
				return &value.Error{msg}
			}

			fn, ok := valFn.(*value.WrapperFunction)

			if !ok {
				return result
			}

			return fn.Fn(params...)

		default:
			msg := fmt.Sprintf("Unrecognized receiver type: %s", util.TypeOf(receiverInstance))
			return &value.Error{msg}
		}
	}

	valFn, ok := identifier.(*value.Function)

	if !ok {
		msg := fmt.Sprintf("Identifier %s is not a function", identifierName)
		return &value.Error{msg}
	}

	if len(valFn.Args) > len(params) {
		msg := fmt.Sprintf("Bad function arguments, expected %d, got %d", len(valFn.Args), len(params))
		return &value.Error{msg}
	}

	for i, _arg := range valFn.Args {
		switch arg := _arg.(type) {
		case *value.String:
			fnEnv.Set(arg.Value, params[i])

		default:
			msg := fmt.Sprintf("Unrecognized arguement type: %s", util.TypeOf(arg))
			return &value.Error{msg}
		}
	}

	return e.Eval(valFn.Body, fnEnv)
}

func (e *Evaluator) EvaluateFunctionStmt(stmt *ast.NodeFunctionStmt, env *environment.Environment) value.Value {
	var ident string
	var receiver string
	var result value.Value = value.NULL

	switch stmt.Identifier.(type) {
	case *ast.NodeIdentifier:
		ident = stmt.Identifier.(*ast.NodeIdentifier).Name

	case *ast.NodeBinaryExpr:
		receiver = stmt.Identifier.(*ast.NodeBinaryExpr).Left.(*ast.NodeIdentifier).Name
		ident = stmt.Identifier.(*ast.NodeBinaryExpr).Right.(*ast.NodeIdentifier).Name

	default:
		msg := fmt.Sprintf("Unrecognized function identifier type: %s", util.TypeOf(stmt.Identifier))
		return &value.Error{msg}
	}

	args := make([]value.Value, len(stmt.Arguements))

	for i, _arg := range stmt.Arguements {
		switch arg := _arg.(type) {
		case *ast.NodeIdentifier:
			args[i] = &value.String{arg.Name}

		case *ast.NodeSelf:
			if i != 0 {
				msg := fmt.Sprintf("self arguement should be at position 0, detected position: %d", i)
				return &value.Error{msg}
			}

			args[i] = &value.Self{arg.Name}

		default:
			msg := fmt.Sprintf("Unrecognized arguement type: %s", util.TypeOf(arg))
			return &value.Error{msg}
		}
	}

	valFn := &value.Function{Args: args, Body: stmt.Body}

	if receiver != "" {
		receiverVal, ok := env.Get(receiver)

		if !ok {
			msg := fmt.Sprintf("Symbol %s is not exists", ident)
			return &value.Error{msg}
		}

		_struct, ok := receiverVal.(*value.Struct[value.Value])

		if !ok {
			msg := fmt.Sprintf("Symbol %s is not a struct", ident)
			return &value.Error{msg}
		}

		if util.InArray[string](_struct.Prop, ident) {
			msg := fmt.Sprintf("Symbol %s already exists", ident)
			return &value.Error{msg}
		}

		// No need to set the value to the environment since
		// the struct is a pointer
		_struct.Prop = append(_struct.Prop, ident)
		_struct.KeyVal.Map[ident] = valFn
	} else {
		if _, ok := env.Get(ident); ok {
			msg := fmt.Sprintf("Symbol %s already exists", ident)
			return &value.Error{msg}
		}

		env.Set(ident, valFn)
	}

	return result
}

func (e *Evaluator) EvaluateIdentifier(stmt *ast.NodeIdentifier, env *environment.Environment) value.Value {
	val, ok := env.Get(stmt.Name)

	if !ok {
		msg := fmt.Sprintf("Symbol %s is not found", stmt.Name)
		return &value.Error{msg}
	}

	return val
}

func (e *Evaluator) EvaluateLetStmt(stmt *ast.NodeLetStmt, env *environment.Environment) value.Value {
	var result value.Value = value.NULL
	ident := stmt.Identifier.(*ast.NodeIdentifier).Name

	val := e.Eval(stmt.Value, env)
	if e.Error(val) {
		return val
	}

	if _, ok := env.Get(ident); ok {
		msg := fmt.Sprintf("Variable %s is already exists", ident)
		return &value.Error{msg}
	}

	env.Set(ident, val)
	return result
}

func (e *Evaluator) EvaluateConstStmt(stmt *ast.NodeConstStmt, env *environment.Environment) value.Value {
	var result value.Value = value.NULL
	ident := stmt.Identifier.(*ast.NodeIdentifier).Name

	val := e.Eval(stmt.Value, env)
	if e.Error(val) {
		return val
	}

	if _, ok := env.Get(ident); ok {
		msg := fmt.Sprintf("Constant %s already exists", ident)
		return &value.Error{msg}
	}

	env.Set(ident, val)
	return result
}

func (e *Evaluator) EvaluateBinaryExpr(stmt *ast.NodeBinaryExpr, env *environment.Environment) value.Value {
	var result value.Value = value.NULL

	switch stmt.Operator {

	case "+":
		left := e.Eval(stmt.Left, env)
		if e.Error(left) {
			return left
		}

		right := e.Eval(stmt.Right, env)
		if e.Error(right) {
			return right
		}

		val := left.(*value.Int).Value + right.(*value.Int).Value
		return &value.Int{val}

	case "-":
		left := e.Eval(stmt.Left, env)
		if e.Error(left) {
			return left
		}

		right := e.Eval(stmt.Right, env)
		if e.Error(right) {
			return right
		}

		val := left.(*value.Int).Value - right.(*value.Int).Value
		return &value.Int{val}

	case "*":
		left := e.Eval(stmt.Left, env)
		if e.Error(left) {
			return left
		}

		right := e.Eval(stmt.Right, env)
		if e.Error(right) {
			return right
		}

		val := left.(*value.Int).Value * right.(*value.Int).Value
		return &value.Int{val}

	case "/":
		left := e.Eval(stmt.Left, env)
		if e.Error(left) {
			return left
		}

		right := e.Eval(stmt.Right, env)
		if e.Error(right) {
			return right
		}

		val := left.(*value.Int).Value / right.(*value.Int).Value
		return &value.Int{val}

	case "%":
		left := e.Eval(stmt.Left, env)
		if e.Error(left) {
			return left
		}

		right := e.Eval(stmt.Right, env)
		if e.Error(right) {
			return right
		}

		val := left.(*value.Int).Value % right.(*value.Int).Value
		return &value.Int{val}

	case "=":
		var ident value.Value

		switch node := stmt.Left.(type) {
		case *ast.NodeIdentifier:
			ident = &value.String{node.Name}

		case *ast.NodeBinaryExpr:
			ident = e.Eval(node.Left, env)

			if e.Error(ident) {
				return ident
			}

		default:
			msg := fmt.Sprintf("Unrecognized assignment type: %s", util.TypeOf(stmt.Left))
			return &value.Error{msg}
		}

		switch node := stmt.Left.(type) {
		case *ast.NodeIdentifier:
			val := e.Eval(stmt.Right, env)

			if e.Error(val) {
				return val
			}

			realIdent := ident.(*value.String).Value

			if _, ok := env.Get(realIdent); !ok {
				msg := fmt.Sprintf("Variable %s is not found", realIdent)
				return &value.Error{msg}
			}

			env.Assign(realIdent, val)
			return val

		case *ast.NodeBinaryExpr:
			val := e.Eval(stmt.Right, env)

			if e.Error(val) {
				return val
			}

			receiver, ok := ident.(*value.Struct[value.Value])

			if !ok {
				msg := fmt.Sprintf("Invalid receiver type: %s", util.TypeOf(receiver))
				return &value.Error{msg}
			}

			right, ok := node.Right.(*ast.NodeIdentifier)

			if !ok {
				msg := fmt.Sprintf("Invalid identifier: %s", node.Right)
				return &value.Error{msg}
			}

			ident = &value.String{Value: right.Name}

			_, ok = receiver.Map[ident.(*value.String).Value]

			if !ok {
				msg := fmt.Sprintf("Symbol %s is not found", ident.(*value.String).Value)
				return &value.Error{msg}
			}

			receiver.Map[ident.(*value.String).Value] = val
			return result

		default:
			msg := fmt.Sprintf("Unrecognized assignment type: %s", util.TypeOf(stmt.Left))
			return &value.Error{msg}
		}

	case "<":
		left := e.Eval(stmt.Left, env)
		if e.Error(left) {
			return left
		}

		right := e.Eval(stmt.Right, env)
		if e.Error(right) {
			return right
		}

		if util.TypeOf(left) != util.TypeOf(right) {
			msg := fmt.Sprintf("Invalid operation %s < %s", left, right)
			return &value.Error{msg}
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
			return &value.Error{msg}
		}

	case ">":
		left := e.Eval(stmt.Left, env)
		if e.Error(left) {
			return left
		}

		right := e.Eval(stmt.Right, env)
		if e.Error(right) {
			return right
		}

		if util.TypeOf(left) != util.TypeOf(right) {
			msg := fmt.Sprintf("Invalid operation %s > %s", left, right)
			return &value.Error{msg}
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
			return &value.Error{msg}
		}

	case "<=":
		left := e.Eval(stmt.Left, env)
		if e.Error(left) {
			return left
		}

		right := e.Eval(stmt.Right, env)
		if e.Error(right) {
			return right
		}

		if util.TypeOf(left) != util.TypeOf(right) {
			msg := fmt.Sprintf("Invalid operation %s <= %s", left, right)
			return &value.Error{msg}
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
			if e.Error(_left) {
				return _left
			}

			_right := e.Eval(stmt.Right, env)
			if e.Error(_right) {
				return _right
			}

			left, _ := _left.(*value.Float)
			right, _ := _right.(*value.Float)

			if left.Value <= right.Value {
				return value.TRUE
			}

			return value.FALSE

		default:
			msg := fmt.Sprintf("Invalid operation %s <= %s", left, right)
			return &value.Error{msg}
		}

	case ">=":
		left := e.Eval(stmt.Left, env)
		if e.Error(left) {
			return left
		}

		right := e.Eval(stmt.Right, env)
		if e.Error(right) {
			return right
		}

		if util.TypeOf(left) != util.TypeOf(right) {
			msg := fmt.Sprintf("Invalid operation %s >= %s", left, right)
			return &value.Error{msg}
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
			return &value.Error{msg}
		}

	case "==":
		left := e.Eval(stmt.Left, env)
		if e.Error(left) {
			return left
		}

		right := e.Eval(stmt.Right, env)
		if e.Error(right) {
			return right
		}

		if left == right {
			return value.TRUE
		}

		return value.FALSE

	case "!=":
		left := e.Eval(stmt.Left, env)
		if e.Error(left) {
			return left
		}

		right := e.Eval(stmt.Right, env)
		if e.Error(right) {
			return right
		}

		if left != right {
			return value.TRUE
		}

		return value.FALSE

	case ".":
		receiver := e.Eval(stmt.Left, env)
		if e.Error(receiver) {
			return receiver
		}

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
				return &value.Error{msg}
			}

			return val

		case *value.Struct[value.Value]:
			val, ok := receiver.(*value.Struct[value.Value]).Map[right]

			if !ok {
				msg := fmt.Sprintf("Symbol %s is not found", right)
				return &value.Error{msg}
			}

			return val

		case *value.Self:
			self, _ := env.Get(receiver.(*value.Self).Value)
			return self

		default:
			msg := fmt.Sprintf("Unknown receiverInstance type %s for dot operator", util.TypeOf(receiver))
			return &value.Error{msg}
		}

	default:
		msg := fmt.Sprintf("Unrecognized operator: %s", stmt.Operator)
		return &value.Error{msg}
	}

}

func (e *Evaluator) EvalProgram(stmt *ast.NodeProgram, env *environment.Environment) value.Value {
	var result value.Value

	for _, stmt := range stmt.Body {
		result = e.Eval(stmt, env)

		if _, ok := result.(*value.Return); ok {
			return result
		}

		if _, ok := result.(*value.Error); ok {
			return result
		}
	}

	return result
}

func (e *Evaluator) Error(val value.Value) bool {
	return val.Type() == value.TYPE_ERROR
}
