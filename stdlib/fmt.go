package stdlib

import (
	"fmt"
	"kat/value"
)

func Print(varargs ...value.Value) value.Value {
	args := buildArgs(varargs)
	fmt.Print(args...)
	return value.NULL
}

func Println(varargs ...value.Value) value.Value {
	args := buildArgs(varargs)
	fmt.Println(args...)
	return value.NULL
}

func Printf(varargs ...value.Value) value.Value {
	format := varargs[0].String()
	args := buildArgs(varargs[1:])
	fmt.Printf(format, args...)
	return value.NULL
}

func Sprintf(varargs ...value.Value) value.Value {
	format := varargs[0].String()
	args := buildArgs(varargs[1:])
	return &value.String{Value: fmt.Sprintf(format, args...)}
}

func buildArgs(varargs []value.Value) []any {
	args := make([]any, 0)

	for _, arg := range varargs {
		args = append(args, arg.String())
	}
	return args
}
