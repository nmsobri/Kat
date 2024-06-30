package main

import (
	"fmt"
	"kat/environment"
	"kat/evaluator"
	"kat/lexer"
	"kat/parser"
)

func main() {
	source := util.ReadFile("./doc/function.kat")

	l := lexer.New(source)
	p := parser.New(l)

	program := p.ParseProgram()
	fmt.Println(program.String())

	e := evaluator.New(program)
	env := environment.New()
	res := e.Eval(program, env)

	fmt.Println()
	if e.IsError() {
		fmt.Println("Evaluation Errors:")
		for _, err := range e.Errors {
			fmt.Println(err)
		}
	}

	fmt.Println()
	fmt.Println("Result:", res)
}
