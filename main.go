package main

import (
	"github.com/maiksch/best-lang/repl"
)

func main() {
	repl.Debug(`
	var foo = fn(x, func) { return func(x) }
	foo(1, fn(){})
	`)

	// println("Welcome to the Best Lang REPL!")
	// repl.Start(os.Stdin, os.Stdout)
}
