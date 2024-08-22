package main

import (
	"os"

	"github.com/maiksch/best-lang/repl"
)

func main() {
	println("Welcome to the Best Lang REPL!")

	repl.Start(os.Stdin, os.Stdout)
}
