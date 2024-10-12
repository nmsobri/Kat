package vm

import (
	"fmt"
	"kat/code"
	"kat/compiler"
	"kat/value"
)

const STACK_SIZE = 2048

type VM struct {
	constants    []value.Value     // Refer bytecode constant pool
	instructions code.Instructions // Refer bytecode Instructions

	stack []value.Value
	sp    int // Always point to the next value. Top of the stack is stack[sp-1]
}

func New(byteCode *compiler.Bytecode) *VM {
	return &VM{
		constants:    byteCode.Constants,
		instructions: byteCode.Instructions,

		stack: make([]value.Value, STACK_SIZE),
		sp:    0,
	}
}

func (v *VM) StackTop() value.Value {
	if v.sp == 0 {
		return nil
	}

	return v.stack[v.sp-1]
}

func (v *VM) Run() error {
	// Every index in VM.instrction is a `byte` ( 8 bit ), cause the type of VM.instructions is `byte`
	for ip := 0; ip < len(v.instructions); ip++ {
		op := code.Opcode(v.instructions[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(v.instructions[ip+1:])
			v.push(v.constants[constIndex])
			ip += 2

		case code.OpAdd:
			right := v.pop()
			left := v.pop()

			val := left.(*value.Int).Value + right.(*value.Int).Value
			o := &value.Int{Value: val}
			v.push(o)
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
