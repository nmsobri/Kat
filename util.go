package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type UtilT struct{}

var util UtilT

func (u UtilT) ReadFile(fileName string) []byte {
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

func (u UtilT) repl() {
	fmt.Println("Welcome To Kat Repl")
}
