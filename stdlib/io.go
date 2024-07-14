package stdlib

import (
	"bufio"
	"kat/value"
	"log"
	"os"
)

var IoFuncs = map[string]value.Value{}

func init() {
	IoFuncs["read_line"] = &value.WrapperFunction{Name: "read_line", Fn: ReadLine}
	IoFuncs["append_to_file"] = &value.WrapperFunction{Name: "append_to_file", Fn: AppendtToFile}
	IoFuncs["write_to_file"] = &value.WrapperFunction{Name: "write_to_file", Fn: WriteToFile}
}

func ReadLine(varargs ...value.Value) value.Value {
	reader := bufio.NewReader(os.Stdin)
	text, _, _ := reader.ReadLine()
	return &value.String{string(text)}
}

func AppendtToFile(varargs ...value.Value) value.Value {
	if len(varargs) != 2 {
		log.Fatal("append_to_file requires 2 arguments")
	}

	filename := varargs[0].(*value.String).Value
	content := varargs[1].(*value.String).Value

	return writeFile(filename, content, os.O_WRONLY|os.O_CREATE|os.O_APPEND)

}

func WriteToFile(varargs ...value.Value) value.Value {
	if len(varargs) != 2 {
		log.Fatal("write_to_file requires 2 arguments")
	}

	filename := varargs[0].(*value.String).Value
	content := varargs[1].(*value.String).Value

	return writeFile(filename, content, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
}

func writeFile(filename string, content string, mode int) value.Value {
	file, err := os.OpenFile(filename, mode, 0644)

	if err != nil {
		return value.NULL
	}

	defer func() {
		_ = file.Close()
	}()

	// Write the content to the file.
	_, err = file.WriteString(content)

	if err != nil {
		return value.NULL
	}

	return value.NULL
}
