package stdlib

import (
	"bufio"
	"kat/value"
	"os"
)

var IoFuncs = map[string]value.Value{}

func init() {
	IoFuncs["input"] = &value.WrapperFunction{Name: "input", Fn: Input}
}

func Input(varargs ...value.Value) value.Value {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return &value.String{text}
}
