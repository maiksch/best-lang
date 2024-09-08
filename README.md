# The goal of this project

I am following the book `Writing an Interpreter in Go` by Thorsten Ball, but the goal is to change the details of the language to fit what I want to have in a programming language.

This is what I want to implement in the future, independent of the book:

- Static Typing (starting of with dynamic typing)
- Structural Typing
- Enums
- No null (Optionals instead)
- No try catch (somehow something better)

And some nice to haves:

- Stack traces for errors
- Line numbers for errors

Statements
`variable declaration/assignment`
`return`

Expressions
`string literal`
`number literal`
`fn definition`

## Example Code:
```
var age = 1
var name = "Monkey"
var array = [1, 2, 3, 4, 5]
var map = {name: "thorsten", age: 28}

var add = fn(a, b) {
  var c = a + b
  return c
}

add(1, 2)
```