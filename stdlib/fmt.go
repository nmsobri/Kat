package stdlib

import (
	"fmt"
	"kat/value"
)

var FmtFuncs = map[string]value.Value{}

func init() {
	FmtFuncs["print"] = &value.WrapperFunction{Name: "print", Fn: Print}
	FmtFuncs["println"] = &value.WrapperFunction{Name: "println", Fn: Println}
	FmtFuncs["printf"] = &value.WrapperFunction{Name: "printf", Fn: Printf}
	FmtFuncs["sprintf"] = &value.WrapperFunction{Name: "sprintf", Fn: Sprintf}
}

func init() {
	IoFuncs["input"] = &value.WrapperFunction{Name: "input", Fn: Input}
}

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
