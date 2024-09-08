package evaluator_test

import (
	"testing"

	"github.com/maiksch/best-lang/evaluator"
	"github.com/maiksch/best-lang/lexer"
	"github.com/maiksch/best-lang/parser"
)

func TestEvalErrorHandling(t *testing.T) {
	input := "1 + true"
	actual := testEval(input)
	expectError(t, actual, "operator type mismatch. INTEGER + BOOLEAN")

	input = "-true"
	actual = testEval(input)
	expectError(t, actual, "invalid operator. -BOOLEAN")

	input = "!1"
	actual = testEval(input)
	expectError(t, actual, "invalid operator. !INTEGER")

	input = `1 + true
	2`
	actual = testEval(input)
	expectError(t, actual, "operator type mismatch. INTEGER + BOOLEAN")

	input = `true + true`
	actual = testEval(input)
	expectError(t, actual, "operator type mismatch. BOOLEAN + BOOLEAN")

	input = `
	if true {
		1 + true
		return 2
	}
	return 1`
	actual = testEval(input)
	expectError(t, actual, "operator type mismatch. INTEGER + BOOLEAN")

	input = `
	if true + 1 {
		return 2
	}`
	actual = testEval(input)
	expectError(t, actual, "operator type mismatch. BOOLEAN + INTEGER")

	input = `
	if 2 == 2 + true {
		return 2
	}`
	actual = testEval(input)
	expectError(t, actual, "operator type mismatch. INTEGER + BOOLEAN")

	input = `
	if 2 == true {
		return 2
	}`
	actual = testEval(input)
	expectError(t, actual, "operator type mismatch. INTEGER == BOOLEAN")

	input = `!(1 == true)`
	actual = testEval(input)
	expectError(t, actual, "operator type mismatch. INTEGER == BOOLEAN")
}

func TestEvalReturn(t *testing.T) {
	input := `1 + 1
	return true
	5 + 5`
	actual := testEval(input)
	expectBooleanValue(t, actual, true)

	input = `if true {
		return 1
	}
	1 + 1`
	actual = testEval(input)
	expectIntegerValue(t, actual, 1)

	input = `if true {
		if true {
			return true
		}
		return false
	}`
	actual = testEval(input)
	expectBooleanValue(t, actual, true)

	input = `if true {
		if true {
			false
		}
		return true
	}`
	actual = testEval(input)
	expectBooleanValue(t, actual, true)
}

func TestEvalIfExpression(t *testing.T) {
	input := "if 1 == 2 { 1 } else { 2 }"
	actual := testEval(input)
	expectIntegerValue(t, actual, 2)

	input = "if 2 == 2 { 1 } else { 2 }"
	actual = testEval(input)
	expectIntegerValue(t, actual, 1)

	input = "if false { 1 }"
	actual = testEval(input)
	expectNothingValue(t, actual)
}

func TestEvalInfixExpression(t *testing.T) {
	input := "1 == 1"
	actual := testEval(input)
	expectBooleanValue(t, actual, true)

	input = "true == true"
	actual = testEval(input)
	expectBooleanValue(t, actual, true)

	input = "true == !false"
	actual = testEval(input)
	expectBooleanValue(t, actual, true)

	input = "1 == 1123"
	actual = testEval(input)
	expectBooleanValue(t, actual, false)

	input = "true == 123"
	actual = testEval(input)
	expectBooleanValue(t, actual, false)

	input = "100 != 1000"
	actual = testEval(input)
	expectBooleanValue(t, actual, true)

	input = "1000 != 1000"
	actual = testEval(input)
	expectBooleanValue(t, actual, false)

	input = "true != false"
	actual = testEval(input)
	expectBooleanValue(t, actual, true)

	input = "100 < 1000"
	actual = testEval(input)
	expectBooleanValue(t, actual, true)

	input = "1 > 123"
	actual = testEval(input)
	expectBooleanValue(t, actual, false)

	input = "1 * 123 == 123"
	actual = testEval(input)
	expectBooleanValue(t, actual, true)
}

func TestEvalBooleanLiteral(t *testing.T) {
	input := "true"
	actual := testEval(input)
	expectBooleanValue(t, actual, true)

	input = "false"
	actual = testEval(input)
	expectBooleanValue(t, actual, false)

	input = "!true"
	actual = testEval(input)
	expectBooleanValue(t, actual, false)

	input = "!false"
	actual = testEval(input)
	expectBooleanValue(t, actual, true)
}

func TestEvalIntegerLiteral(t *testing.T) {
	input := "1"
	actual := testEval(input)
	expectIntegerValue(t, actual, 1)

	input = "-100"
	actual = testEval(input)
	expectIntegerValue(t, actual, -100)

	input = "-100 + 50"
	actual = testEval(input)
	expectIntegerValue(t, actual, -50)

	input = "2 * 34 / 2 - 10"
	actual = testEval(input)
	expectIntegerValue(t, actual, 24)

	input = "3 * (3 * 3) + 10"
	actual = testEval(input)
	expectIntegerValue(t, actual, 37)

	input = "(5 + 10 * 2 + 15 / 3) * 2 + -10"
	actual = testEval(input)
	expectIntegerValue(t, actual, 50)
}

func expectError(t *testing.T, actual evaluator.Object, expect string) {
	errorValue, ok := actual.(*evaluator.Error)
	if !ok {
		t.Fatalf("expected error value. got %T", actual)
	}
	if errorValue.Message != expect {
		t.Fatalf("wrong error\n\texpected: %s\n\tgot:      %s", expect, errorValue.Message)
	}
}

func expectNothingValue(t *testing.T, v evaluator.Object) {
	if _, ok := v.(*evaluator.Nothing); !ok {
		t.Fatalf("Expected nothing value, got %T", v)
	}
}

func expectBooleanValue(t *testing.T, v evaluator.Object, expect bool) {
	obj, ok := v.(*evaluator.Boolean)
	if !ok {
		t.Fatalf("Expected boolean value, got %T", v)
	}
	if obj.Value != expect {
		t.Fatalf("Expected evaluated value to be %t. got %t", expect, obj.Value)
	}
}

func expectIntegerValue(t *testing.T, v evaluator.Object, expect int64) {
	obj, ok := v.(*evaluator.Integer)
	if !ok {
		t.Fatalf("Expected integer value, got %T", v)
	}
	if obj.Value != expect {
		t.Fatalf("Expected evaluated value to be %d. got %d", expect, obj.Value)
	}
}

func testEval(input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return evaluator.Eval(program)
}
