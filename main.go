package main

import (
	"fmt"
	"io"
	"kat/lexer"
	"kat/token"
	"log"
	"os"
)

func main() {
	source := ReadFile("./doc/main.kat")

	lex := lexer.New(source)

	currentToken := lex.Token()

	for currentToken.Type != token.EOF {
		fmt.Println(currentToken)
		currentToken = lex.Token()
	}

	fmt.Println(currentToken)
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
