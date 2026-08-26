package main

import (
	"fmt"
	"kat/environment"
	"kat/evaluator"
	"kat/lexer"
	"kat/parser"
	"kat/util"
	"kat/value"
	"os"
	"path/filepath"
)

func main() {
	if err := checkUsage(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	source := util.ReadFile(os.Args[1])

	l := lexer.New(source)
	p := parser.New(l)

	program := p.ParseProgram()
	//fmt.Println(program.String())

	e := evaluator.New(program)
	env := environment.New()
	res := e.Eval(program, env)

	if err, ok := res.(*value.Error); ok {
		fmt.Println(err)
	}
}

func checkUsage() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("usage: %s <file>", filepath.Base(os.Args[0]))
	}

	return nil
}
