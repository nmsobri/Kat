package vm

import (
	"fmt"
	"kat/ast"
	"kat/compiler"
	"kat/lexer"
	"kat/parser"
	"kat/value"
	"testing"
)

type vmTestCase struct {
	input    string
	expected any
}

func parse(input string) *ast.NodeProgram {
	l := lexer.New([]byte(input))
	p := parser.New(l)
	return p.ParseProgram()
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, test := range tests {
		program := parse(test.input)
		comp := compiler.New()
		err := comp.Compile(program)

		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())

		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.StackTop()
		testExpectedObject(t, test.expected, stackElem)
	}
}

func testExpectedObject(t *testing.T, expected any, actual value.Value) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)

		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	}
}

func testIntegerObject(expected int64, actual value.Value) error {
	actualValue, ok := actual.(*value.Int)

	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if actualValue.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", actualValue.Value, expected)
	}

	return nil
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3}, // FIXME
	}

	runVmTests(t, tests)
}
