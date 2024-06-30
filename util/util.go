package util

import (
	"fmt"
	"io"
	"kat/value"
	"log"
	"os"
)

func TypeOf(v any) string {
	return fmt.Sprintf("%T", v)
}

func IsTruthy(v value.Value) bool {
	switch v.(type) {
	case value.ValueInt:
		return v.(value.ValueInt).Value != 0

	case value.ValueFloat:
		return v.(value.ValueFloat).Value != 0

	case value.ValueBool:
		return v.(value.ValueBool).Value
	}

	return false
}

func InArray[T comparable](arr []T, key T) bool {
	for _, v := range arr {
		if v == key {
			return true
		}
	}

	return false
}

func ReadFile(fileName string) []byte {
	f, e := os.Open(fileName)

	if e != nil {
		log.Fatal(e)
	}

	source, e := io.ReadAll(f)

	if e != nil {
		log.Fatalln(e)
	}
	return source
}

func Repl() {
	fmt.Println("Welcome To Kat Repl")
}
