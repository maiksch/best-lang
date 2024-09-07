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

	case *parser.PrefixExpression:
		return evalPrefixExpression(node)

	case *parser.InfixExpression:
		return evalInfixExpression(node)

	case *parser.BooleanLiteral:
		if node.Value {
			return TRUE
		} else {
			return FALSE
		}

	default:
		log.Panicf("eval for %T not implemented", node)
		return nil
	}
}

func evalInfixExpression(expr *parser.InfixExpression) Object {
	left := Eval(expr.Left)
	right := Eval(expr.Right)

	if left.Type() == INTEGER && right.Type() == INTEGER {
		l := left.(*Integer).Value
		r := right.(*Integer).Value

		switch expr.Operator {
		case "==":
			return toBooleanObject(l == r)

		case "!=":
			return toBooleanObject(l != r)

		case ">":
			return toBooleanObject(l > r)

		case "<":
			return toBooleanObject(l < r)

		case "+":
			return &Integer{Value: l + r}

		case "-":
			return &Integer{Value: l - r}

		case "/":
			return &Integer{Value: l / r}

		case "*":
			return &Integer{Value: l * r}
		}
	}

	if left.Type() == BOOLEAN && right.Type() == BOOLEAN {
		switch expr.Operator {
		case "==":
			return toBooleanObject(left == right)

		case "!=":
			return toBooleanObject(left != right)
		}
	}

	switch expr.Operator {
	case "==":
		return FALSE
	}

	return nil
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
			return toBooleanObject(!b.Value)
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

func toBooleanObject(v bool) *Boolean {
	if v {
		return TRUE
	} else {
		return FALSE
	}
}
