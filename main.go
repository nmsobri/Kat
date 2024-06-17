package main

import (
	"fmt"
	"kat/lexer"
	"kat/parser"
)

func main() {
	source := util.ReadFile("./doc/main.kat")

	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	fmt.Println(program.String())
}
