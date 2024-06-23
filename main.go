package main

import (
	"fmt"
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
	res := e.Eval(program)
	fmt.Println("result:", res)
}
