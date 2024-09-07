package evaluator

import "github.com/maiksch/best-lang/parser"

func Eval(node parser.Node) Object {
	switch node := node.(type) {
	case *parser.Program:
		return evalStatements(node.Statements)

	case *parser.ExpressionStatement:
		return Eval(node.Value)

	case *parser.IntegerLiteral:
		return &Integer{Value: node.Value}

	case *parser.BooleanLiteral:
		if node.Value {
			return TRUE
		} else {
			return FALSE
		}

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
