package main

import (
	"fmt"
	"io"
	"kat/lexer"
	"log"
	"os"
)

func main() {
	f, e := os.Open("./doc/main.kat")

	if e != nil {
		log.Fatal(e)
	}

	source, e := io.ReadAll(f)

	if e != nil {
		log.Fatalln(e)
	}

	input := string(source)
	lex := lexer.New(input)

	currentToken := lex.Next()
	fmt.Println(currentToken)
}

func repl() {
	fmt.Println("Welcome To Kat Repl")
}
