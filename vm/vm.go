package vm

import (
	"fmt"
	"kat/code"
	"kat/compiler"
	"kat/util"
	"kat/value"
)

const STACK_SIZE = 2048
const GLOBAL_SIZE = 65536

var (
	TRUE  = &value.Bool{Value: true}
	FALSE = &value.Bool{Value: false}
	NULL  = &value.Null{}
)

type VM struct {
	constants    []value.Value     // Refer bytecode constant pool
	instructions code.Instructions // Refer bytecode Instructions

	stack []value.Value
	sp    int // Always point to the next value. Top of the stack is stack[sp-1]

	globals []value.Value
}

func New(byteCode *compiler.Bytecode) *VM {
	return &VM{
		constants:    byteCode.Constants,
		instructions: byteCode.Instructions,

		stack:   make([]value.Value, STACK_SIZE),
		sp:      0,
		globals: make([]value.Value, GLOBAL_SIZE),
	}
}

func NewWithState(byteCode *compiler.Bytecode, globals []value.Value) *VM {
	vm := New(byteCode)
	vm.globals = globals
	return vm
}

func (v *VM) StackTop() value.Value {
	if v.sp == 0 {
		return nil
	}

	return v.stack[v.sp-1]
}

func (v *VM) Run() error {
	// Every index in VM.instruction is a `byte` ( 8 bit ), cause the type of VM.instructions is `byte`
	for ip := 0; ip < len(v.instructions); ip++ {
		op := code.Opcode(v.instructions[ip])

		switch op {

		case code.OpConstant:
			constIndex := code.ReadUint16(v.instructions[ip+1:])

			if err := v.push(v.constants[constIndex]); err != nil {
				return err
			}

			ip += 2

		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			if err := v.executeBinaryOperation(op); err != nil {
				return err
			}

		case code.OpGreaterThan, code.OpLessThan, code.OpEqual, code.OpNotEqual:
			if err := v.executeComparison(op); err != nil {
				return err
			}

		case code.OpTrue:
			if err := v.push(TRUE); err != nil {
				return err
			}

		case code.OpFalse:
			if err := v.push(FALSE); err != nil {
				return err
			}

		case code.OpPop:
			v.pop()

		case code.OpMinus:
			right := v.pop()
			if err := v.executeMinusOperator(right); err != nil {
				return err
			}

		case code.OpBang:
			right := v.pop()
			if err := v.executeBangOperator(right); err != nil {
				return err
			}

		case code.OpJumpIfFalse:
			condition := v.pop()
			pos := int(code.ReadUint16(v.instructions[ip+1:]))
			ip += 2 // skip jump operand

			if !util.IsTruthy(condition) {
				ip = pos - 1
			}

		case code.OpJump:
			pos := int(code.ReadUint16(v.instructions[ip+1:]))
			ip = pos - 1

		case code.OpNull:
			if err := v.push(NULL); err != nil {
				return err
			}

		case code.OpSetGlobal:
			index := code.ReadUint16(v.instructions[ip+1:])
			ip += 2

			value := v.pop()
			v.globals[index] = value

		case code.OpGetGlobal:
			index := code.ReadUint16(v.instructions[ip+1:])
			ip += 2

			value := v.globals[index]

			if err := v.push(value); err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *VM) push(value value.Value) error {
	if v.sp >= STACK_SIZE {
		return fmt.Errorf("stack overflow")
	}

	v.stack[v.sp] = value
	v.sp++
	return nil
}

func (v *VM) pop() value.Value {
	val := v.stack[v.sp-1]
	v.sp--
	return val
}

func (v *VM) LastPoppedStackElem() value.Value {
	return v.stack[v.sp]
}

func (v *VM) executeBinaryOperation(op code.Opcode) error {
	right := v.pop()
	left := v.pop()

	if right.Type() == value.TYPE_INT && left.Type() == value.TYPE_INT {
		return v.executeBinaryIntegerOperation(op, left, right)
	}

	return fmt.Errorf("Unsupported types for binary operation: %s %s", left.Type(), right.Type())
}

func (v *VM) executeComparison(op code.Opcode) error {
	right := v.pop()
	left := v.pop()

	if left.Type() == value.TYPE_INT && right.Type() == value.TYPE_INT {
		return v.executeIntegerComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		if err := v.push(nativeToBooleanObject(left == right)); err != nil {
			return err
		}

	case code.OpNotEqual:
		if err := v.push(nativeToBooleanObject(left != right)); err != nil {
			return err
		}

	default:
		return fmt.Errorf("Unknown operator %d (%s %s)", op, left.Type(), right.Type())
	}

	return nil
}

func (v *VM) executeBinaryIntegerOperation(op code.Opcode, _left value.Value, _right value.Value) error {
	left := _left.(*value.Int)
	right := _right.(*value.Int)

	var result *value.Int

	switch op {
	case code.OpAdd:
		result = &value.Int{Value: left.Value + right.Value}

	case code.OpSub:
		result = &value.Int{Value: left.Value - right.Value}

	case code.OpMul:
		result = &value.Int{Value: left.Value * right.Value}

	case code.OpDiv:
		result = &value.Int{Value: left.Value / right.Value}
	}

	return v.push(result)
}

func (v *VM) executeIntegerComparison(op code.Opcode, _left value.Value, _right value.Value) error {
	left := _left.(*value.Int).Value
	right := _right.(*value.Int).Value

	switch op {
	case code.OpEqual:
		if err := v.push(nativeToBooleanObject(left == right)); err != nil {
			return err
		}

	case code.OpNotEqual:
		if err := v.push(nativeToBooleanObject(left != right)); err != nil {
			return err
		}

	case code.OpGreaterThan:
		if err := v.push(nativeToBooleanObject(left > right)); err != nil {
			return err
		}

	case code.OpLessThan:
		if err := v.push(nativeToBooleanObject(left < right)); err != nil {
			return err
		}

	default:
		return fmt.Errorf("Unknown operator %d", op)
	}

	return nil
}

func nativeToBooleanObject(val bool) value.Value {
	if val {
		return TRUE
	}

	return FALSE
}

func (v *VM) executeMinusOperator(right value.Value) error {
	if right.Type() != value.TYPE_INT {
		return fmt.Errorf("Unsupported type for negation: %s", right.Type())
	}

	val := right.(*value.Int).Value
	return v.push(&value.Int{Value: -val})
}

func (v *VM) executeBangOperator(right value.Value) error {
	switch right {
	case TRUE:
		return v.push(FALSE)

	case FALSE:
		return v.push(TRUE)

	case NULL:
		return v.push(TRUE)

	default:
		return v.push(FALSE) // if its not a boolean value, consider everything else true, !true = false
	}
}

func (v *VM) Globals() []value.Value {
	return v.globals
}
