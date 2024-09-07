package evaluator

import "github.com/maiksch/best-lang/parser"

func Eval(node parser.Node) Object {
	switch x := node.(type) {
	case *parser.Program:
		return evalStatements(x.Statements)

	case *parser.ExpressionStatement:
		return Eval(x.Value)

	case *parser.IntegerLiteral:
		return &Integer{Value: x.Value}

	default:
		panic("not implemented")
	}
}

func evalStatements(stmts []parser.Statement) Object {
	var result Object
	for _, stmt := range stmts {
		result = Eval(stmt)
	}
	return result
}
