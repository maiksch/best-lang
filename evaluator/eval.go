package evaluator

import (
	"log"

	"github.com/maiksch/best-lang/parser"
)

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

	case *parser.PrefixExpression:
		return evalPrefixExpression(node)

	default:
		log.Panicf("eval for %T not implemented", node)
		return nil
	}
}

func evalPrefixExpression(expr *parser.PrefixExpression) Object {
	value := Eval(expr.Right)

	if i, ok := value.(*Integer); ok {
		switch expr.Operator {
		case "-":
			i.Value = i.Value * -1
		}
	}

	if b, ok := value.(*Boolean); ok {
		switch expr.Operator {
		case "!":
			if b == TRUE {
				value = FALSE
			} else {
				value = TRUE
			}
		}
	}

	return value
}

func evalStatements(stmts []parser.Statement) Object {
	var result Object
	for _, stmt := range stmts {
		result = Eval(stmt)
	}
	return result
}
