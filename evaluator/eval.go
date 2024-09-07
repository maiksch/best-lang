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

	case *parser.BlockStatement:
		return evalStatements(node.Statements)

	case *parser.PrefixExpression:
		return evalPrefixExpression(node)

	case *parser.InfixExpression:
		return evalInfixExpression(node)

	case *parser.IfExpression:
		return evalIfExpression(node)

	case *parser.IntegerLiteral:
		return &Integer{Value: node.Value}

	case *parser.BooleanLiteral:
		return toBooleanObject(node.Value)

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

	if value.Type() == INTEGER && expr.Operator == "-" {
		i := value.(*Integer)
		i.Value = i.Value * -1
	}

	if value.Type() == BOOLEAN && expr.Operator == "!" {
		return toBooleanObject(!value.(*Boolean).Value)
	}

	return value
}

func evalIfExpression(expr *parser.IfExpression) Object {
	condition := Eval(expr.Condition)

	if condition == TRUE {
		return Eval(expr.Consequence)
	}

	if expr.Otherwise != nil {
		return Eval(expr.Otherwise)
	}

	return &Nothing{}
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
