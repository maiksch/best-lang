package main

import (
	"io"
	"os"

	"github.com/maiksch/best-lang/evaluator"
	"github.com/maiksch/best-lang/lexer"
	"github.com/maiksch/best-lang/parser"
	"github.com/maiksch/best-lang/repl"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		println("Welcome to the Best Lang REPL!")
		repl.Start(os.Stdin, os.Stdout)
	}

	prog := args[0]

	source, err := os.ReadFile(prog)
	if err != nil {
		panic(err)
	}

	env := evaluator.NewEnvrionment()
	lexer := lexer.New(string(source))
	parser := parser.New(lexer)

	ast := parser.ParseProgram()

	result := evaluator.Eval(ast, env)

	w := os.Stdout
	io.WriteString(w, result.Inspect())
	io.WriteString(w, "\n")
}
