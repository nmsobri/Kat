package main

import (
	"errors"
	"flag"
	"fmt"
	"kat/environment"
	"kat/evaluator"
	"kat/lexer"
	"kat/parser"
	"kat/types"
	"kat/util"
	"kat/value"
	"os"
)

var (
	versionNumber = "0.1"
	errStop       = errors.New("stop")
)

func main() {
	sourceFile, err := processFlag()

	if errors.Is(err, errStop) {
		os.Exit(1)
	}

	source := util.ReadFile(sourceFile)

	l := lexer.New(source)

	p := parser.New(l)
	program := p.ParseProgram()
	// fmt.Println(program.String())

	typeChecker := types.New()
	typeChecker.CheckProgram(program)

	e := evaluator.New(program)
	env := environment.New()
	res := e.Eval(program, env)

	if err, ok := res.(*value.Error); ok {
		fmt.Println(err)
	}
}

func processFlag() (string, error) {
	flag.Usage = func() {
		binary := os.Args[0]
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", binary)
		_, _ = fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	help := flag.Bool("help", false, "print this help")
	version := flag.Bool("version", false, "print version and exit")
	run := flag.Bool("run", false, "run kat file")

	flag.Parse()

	flagPresent := false
	flag.Visit(func(f *flag.Flag) {
		flagPresent = true
	})

	if *help || !flagPresent {
		flag.Usage()
		return "", errStop
	}

	if *version {
		_, _ = fmt.Fprintf(os.Stderr, "%s", versionNumber)
		return "", errStop
	}

	if *run {
		if flag.NArg() != 1 {
			binary := os.Args[0]
			_, _ = fmt.Fprintf(os.Stderr, "Usage: %s -run <file>\n\n", binary)
			os.Exit(1)
		}

		return flag.Arg(0), nil
	}

	panic("unreachable")
}
