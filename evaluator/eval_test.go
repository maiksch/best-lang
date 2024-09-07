package evaluator_test

import (
	"testing"

	"github.com/maiksch/best-lang/evaluator"
	"github.com/maiksch/best-lang/lexer"
	"github.com/maiksch/best-lang/parser"
)

func TestEvaluateIntegerLiteral(t *testing.T) {
	input := "1"

	actual := testEval(input)

	expectIntegerObject(t, actual, 1)
}

func expectIntegerObject(t *testing.T, v evaluator.Object, expect int64) {
	obj, ok := v.(*evaluator.Integer)
	if !ok {
		t.Fatalf("Expected integer object, got %T", v)
	}
	if obj.Value != 1 {
		t.Fatalf("Expected evaluated value to be %d. got %d", expect, obj.Value)
	}
}

func testEval(input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return evaluator.Eval(program)
}
