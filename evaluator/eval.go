package evaluator

import (
	"fmt"
	"log"

	"github.com/maiksch/best-lang/parser"
)

func Eval(node parser.Node, env *Environment) Object {
	switch node := node.(type) {
	case *parser.Program:
		return evalProgram(node, env)

	case *parser.ExpressionStatement:
		return Eval(node.Value, env)

	case *parser.BlockStatement:
		return evalBlockStatement(node, env)

	case *parser.DeclareStatement:
		return evalDeclareStatement(node, env)

	case *parser.ReturnStatement:
		val := Eval(node.Expression, env)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}

	case *parser.FunctionLiteral:
		return &Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}

	case *parser.FunctionCall:
		return evalFunctionCall(node, env)

	case *parser.PrefixExpression:
		return evalPrefixExpression(node, env)

	case *parser.InfixExpression:
		return evalInfixExpression(node, env)

	case *parser.IfExpression:
		return evalIfExpression(node, env)

	case *parser.Identifier:
		return env.get(node)

	case *parser.IntegerLiteral:
		return &Integer{Value: node.Value}

	case *parser.BooleanLiteral:
		return toBooleanObject(node.Value)

	default:
		log.Panicf("eval for %T not implemented", node)
		return nil
	}
}

func evalFunctionCall(call *parser.FunctionCall, env *Environment) Object {
	fn := Eval(call.Function, env)
	if isError(fn) {
		return fn
	}

	if fn, ok := fn.(*Function); ok {
		closure := CloneEnvironment(fn.Env)

		for i, arg := range call.Arguments {
			val := Eval(arg, env)
			if isError(val) {
				return val
			}
			closure.set(fn.Parameters[i], val)
		}

		result := Eval(fn.Body, closure)

		if returnValue, ok := result.(*ReturnValue); ok {
			result = returnValue.Value
		}

		return result
	}

	return newError("unsuported function expression. %T", call.Function)
}

func evalIfExpression(expr *parser.IfExpression, env *Environment) Object {
	condition := Eval(expr.Condition, env)
	if isError(condition) {
		return condition
	}

	if condition == TRUE {
		return Eval(expr.Consequence, env)
	}

	if expr.Otherwise != nil {
		return Eval(expr.Otherwise, env)
	}

	return &Nothing{}
}

func evalInfixExpression(expr *parser.InfixExpression, env *Environment) Object {
	left := Eval(expr.Left, env)
	if isError(left) {
		return left
	}
	right := Eval(expr.Right, env)
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

func evalPrefixExpression(expr *parser.PrefixExpression, env *Environment) Object {
	value := Eval(expr.Right, env)
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

func evalDeclareStatement(stmt *parser.DeclareStatement, env *Environment) Object {
	value := Eval(stmt.Expression, env)
	if isError(value) {
		return value
	}

	return env.set(stmt.Name, value)
}

func evalBlockStatement(block *parser.BlockStatement, env *Environment) Object {
	var result Object
	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		if result.Type() == ERROR || result.Type() == RETURN {
			return result
		}
	}
	return result
}

func evalProgram(prg *parser.Program, env *Environment) Object {
	var result Object
	for _, stmt := range prg.Statements {
		result = Eval(stmt, env)

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
