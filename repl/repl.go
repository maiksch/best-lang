package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/maiksch/best-lang/lexer"
	"github.com/maiksch/best-lang/parser"
)

const PROMPT = ">> "

func Start(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)

	for {
		fmt.Fprintf(w, "%s", PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		lexer := lexer.New(line)
		parser := parser.New(lexer)

		programm := parser.ParseProgram()

		io.WriteString(w, programm.String())
		io.WriteString(w, "\n")
	}
}
