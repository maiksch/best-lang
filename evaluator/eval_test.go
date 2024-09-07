package evaluator_test

import (
	"testing"

	"github.com/maiksch/best-lang/evaluator"
	"github.com/maiksch/best-lang/lexer"
	"github.com/maiksch/best-lang/parser"
)

func TestEvalBooleanLiteral(t *testing.T) {
	input := "true"
	actual := testEval(input)
	expectBooleanObject(t, actual, true)

	input = "false"
	actual = testEval(input)
	expectBooleanObject(t, actual, false)

	input = "!true"
	actual = testEval(input)
	expectBooleanObject(t, actual, false)

	input = "!false"
	actual = testEval(input)
	expectBooleanObject(t, actual, true)
}

func TestEvalIntegerLiteral(t *testing.T) {
	input := "1"
	actual := testEval(input)
	expectIntegerObject(t, actual, 1)

	input = "-100"
	actual = testEval(input)
	expectIntegerObject(t, actual, -100)
}

func expectBooleanObject(t *testing.T, v evaluator.Object, expect bool) {
	obj, ok := v.(*evaluator.Boolean)
	if !ok {
		t.Fatalf("Expected integer object, got %T", v)
	}
	if obj.Value != expect {
		t.Fatalf("Expected evaluated value to be %t. got %t", expect, obj.Value)
	}
}

func expectIntegerObject(t *testing.T, v evaluator.Object, expect int64) {
	obj, ok := v.(*evaluator.Integer)
	if !ok {
		t.Fatalf("Expected integer object, got %T", v)
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
