package code

import (
	"bytes"
	"fmt"
)

const (
	OpConstant Opcode = iota
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpPop
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpLessThan
	OpBang
	OpMinus
	OpJump
	OpJumpIfFalse
	OpNull
	OpSetGlobal
	OpGetGlobal
)

type Opcode byte

type Instructions []byte

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(Opcode(ins[i]))

		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			return out.String()
		}

		operands, read := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))
		i += 1 + read
	}
	return out.String()
}

func (i Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
			len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant:    {"OpConstant", []int{2}},   // meaning this OpCode take one operand that have 2 byte length
	OpAdd:         {"OpAdd", []int{}},         // meaning this OpCode did not take any operand
	OpSub:         {"OpSub", []int{}},         // meaning this OpCode did not take any operand
	OpMul:         {"OpMul", []int{}},         // meaning this OpCode did not take any operand
	OpDiv:         {"OpDiv", []int{}},         // meaning this OpCode did not take any operand
	OpPop:         {"OpPop", []int{}},         // meaning this OpCode did not take any operand
	OpTrue:        {"OpTrue", []int{}},        // meaning this OpCode did not take any operand
	OpFalse:       {"OpFalse", []int{}},       // meaning this OpCode did not take any operand
	OpEqual:       {"OpEqual", []int{}},       // meaning this OpCode did not take any operand
	OpNotEqual:    {"OpNotEqual", []int{}},    // meaning this OpCode did not take any operand
	OpGreaterThan: {"OpGreaterThan", []int{}}, // meaning this OpCode did not take any operand
	OpLessThan:    {"OpLessThan", []int{}},    // meaning this OpCode did not take any operand
	OpBang:        {"OpBang", []int{}},        // meaning this OpCode did not take any operand
	OpMinus:       {"OpMinus", []int{}},       // meaning this OpCode did not take any operand

	OpJump:        {"OpJump", []int{2}},        // meaning this OpCode take one operand that have 2 byte length
	OpJumpIfFalse: {"OpJumpIfFalse", []int{2}}, // meaning this OpCode take one operand that have 2 byte length
	OpNull:        {"OpNull", []int{}},         // meaning this OpCode did not take any operand
	OpSetGlobal:   {"OpSetGlobal", []int{2}},   // meaning this OpCode did not take any operand
	OpGetGlobal:   {"OpGetGlobal", []int{2}},   // meaning this OpCode did not take any operand
}

func Lookup(op Opcode) (*Definition, error) {
	def, ok := definitions[op]

	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]

	if !ok {
		return []byte{}
	}

	instructionLen := 1

	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]

		switch width {
		case 2: // (16 bit)
			// Encode in Big Endian ( Most Significant Byte first)
			instruction[offset] = byte(o >> 8) // first byte, take most significant byte
			instruction[offset+1] = byte(o)    // second byte, take least significant byte
			offset += width
		}
	}

	return instruction
}

func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}
		offset += width
	}

	return operands, offset
}
func ReadUint16(ins Instructions) uint16 {
	return uint16(ins[0])<<8 | uint16(ins[1])
}
