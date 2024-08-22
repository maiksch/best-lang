package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/maiksch/best-lang/lexer"
	"github.com/maiksch/best-lang/token"
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
		lexer := lexer.NewLexer(line)

		for t := lexer.NextToken(); t.Type != token.EOF; t = lexer.NextToken() {
			fmt.Fprintf(w, "%v\n", t)
		}
	}
}
