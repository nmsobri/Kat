package repl

import (
	"bufio"
	"fmt"
	"io"
	"kat/compiler"
	"kat/lexer"
	"kat/parser"
	"kat/value"
	"kat/vm"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	symbolTable := compiler.NewSymbolTable()
	globals := make([]value.Value, vm.GLOBAL_SIZE)

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
		c := compiler.NewWithState(symbolTable)

		if err := c.Compile(program); err != nil {
			_, _ = fmt.Fprintf(out, "Woops! Compilaton failed: %s\n", err.Error())
			continue
		}

		symbolTable = c.SymbolTable()

		v := vm.NewWithState(c.Bytecode(), globals)

		if err := v.Run(); err != nil {
			_, _ = fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
			continue
		}

		globals = v.Globals()

		stackTop := v.LastPoppedStackElem()
		_, _ = io.WriteString(out, stackTop.String())
		_, _ = io.WriteString(out, "\n")
	}
}
