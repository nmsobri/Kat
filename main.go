package main

import (
	"kat/repl"
	"os"
)

func main() {
	// source := util.ReadFile("./doc/stdlib.kat")

	// l := lexer.New(source)
	// p := parser.New(l)

	// program := p.ParseProgram()
	// //fmt.Println(program.String())

	// e := evaluator.New(program)
	// env := environment.New()
	// res := e.Eval(program, env)

	// if err, ok := res.(*value.Error); ok {
	// 	fmt.Println(err)
	// }
	repl.Start(os.Stdin, os.Stdout)
}
