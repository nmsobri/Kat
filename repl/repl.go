package repl

import (
	"bufio"
	"fmt"
	"io"
	"kat/compiler"
	"kat/lexer"
	"kat/parser"
	"kat/vm"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)

		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.New([]byte(line))
		p := parser.New(l)

		program := p.ParseProgram()

		c := compiler.New()

		if err := c.Compile(program); err != nil {
			fmt.Fprintf(out, "Woops! Compilaton failed:\n %s\n", err)
			continue
		}

		v := vm.New(c.Bytecode())

		if err := v.Run(); err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
			continue
		}

		stackTop := v.LastPoppedStackElem()
		io.WriteString(out, stackTop.String())
		io.WriteString(out, "\n")
	}
}
