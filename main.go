package main

import (
	"fmt"
	"io"
	"kat/lexer"
	"kat/parser"
	"log"
	"os"
)

func main() {
	source := ReadFile("./doc/main.kat")

	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	fmt.Println(program.Node(0))
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

func repl() {
	fmt.Println("Welcome To Kat Repl")
}
