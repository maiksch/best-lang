package evaluator

import (
	"fmt"
	"log"

	"github.com/maiksch/best-lang/parser"
)

func Eval(node parser.Node) Object {
	switch node := node.(type) {
	case *parser.Program:
		return evalProgram(node)

	case *parser.ExpressionStatement:
		return Eval(node.Value)

	case *parser.BlockStatement:
		return evalBlockStatement(node)

	case *parser.ReturnStatement:
		return &ReturnValue{Value: Eval(node.Expression)}

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

func evalIfExpression(expr *parser.IfExpression) Object {
	condition := Eval(expr.Condition)
	if isError(condition) {
		return condition
	}

	if condition == TRUE {
		return Eval(expr.Consequence)
	}

	if expr.Otherwise != nil {
		return Eval(expr.Otherwise)
	}

	return &Nothing{}
}

func evalInfixExpression(expr *parser.InfixExpression) Object {
	left := Eval(expr.Left)
	if isError(left) {
		return left
	}
	right := Eval(expr.Right)
	if isError(right) {
		return right
	}

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

	return newError("operator type mismatch. %s %s %s", left.Type(), expr.Operator, right.Type())
}

func evalPrefixExpression(expr *parser.PrefixExpression) Object {
	value := Eval(expr.Right)
	if isError(value) {
		return value
	}

	if value.Type() == INTEGER && expr.Operator == "-" {
		i := value.(*Integer)
		i.Value = i.Value * -1
		return value
	}

	if value.Type() == BOOLEAN && expr.Operator == "!" {
		return toBooleanObject(!value.(*Boolean).Value)
	}

	return newError("invalid operator. %s%s", expr.Operator, value.Type())
}

func evalBlockStatement(block *parser.BlockStatement) Object {
	var result Object
	for _, stmt := range block.Statements {
		result = Eval(stmt)

		if result.Type() == ERROR || result.Type() == RETURN {
			return result
		}
	}
	return result
}

func evalProgram(prg *parser.Program) Object {
	var result Object
	for _, stmt := range prg.Statements {
		result = Eval(stmt)

		switch result := result.(type) {
		case *ReturnValue:
			return result.Value
		case *Error:
			return result
		}
	}
	return result
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR
	}
	return false
}

func newError(message string, a ...interface{}) *Error {
	return &Error{
		Message: fmt.Sprintf(message, a...),
	}
}

func toBooleanObject(v bool) *Boolean {
	if v {
		return TRUE
	} else {
		return FALSE
	}
}
