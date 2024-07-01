package main

import (
	"fmt"
	"kat/environment"
	"kat/evaluator"
	"kat/lexer"
	"kat/parser"
	"kat/util"
)

func main() {
	source := util.ReadFile("./doc/import.kat")

	l := lexer.New(source)
	p := parser.New(l)

	program := p.ParseProgram()
	//fmt.Println(program.String())

	e := evaluator.New(program)
	env := environment.New()
	e.Eval(program, env)

	if e.IsError() {
		fmt.Println("Evaluation Errors:")
		for _, err := range e.Errors {
			fmt.Println(err)
		}
	}
}
