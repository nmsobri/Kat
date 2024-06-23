package main

import (
	"fmt"
	"kat/environment"
	"kat/evaluator"
	"kat/lexer"
	"kat/parser"
)

func main() {
	source := util.ReadFile("./doc/expression.kat")

	l := lexer.New(source)
	p := parser.New(l)

	program := p.ParseProgram()
	fmt.Println(program.String())

	e := evaluator.New(program)
	env := environment.New()
	res := e.Eval(program, env)

	fmt.Println()
	fmt.Println("Result:", res)
}
