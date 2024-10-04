package evaluator_test

import (
	"testing"

	"github.com/maiksch/best-lang/evaluator"
	"github.com/maiksch/best-lang/lexer"
	"github.com/maiksch/best-lang/parser"
)

func TestEvalFunctionCall(t *testing.T) {
	input := `var foo = fn(x) {
		if x > 0 {
			return true
		}
		return false
	}
	foo(1)`
	actual := testEval(input)
	expectBooleanValue(t, actual, true)

	input = `
	fn(x) {
		if x > 0 {
			return true
		}
		return false
	}(1)`
	actual = testEval(input)
	expectBooleanValue(t, actual, true)

	input = `
	fn(x) {
		if x > 0 {
			return true
		}
		return false
	}(-1)`
	actual = testEval(input)
	expectBooleanValue(t, actual, false)

	input = `
	var foo = fn(x, y) { x + y}
	foo(foo(1, 2), 1 + 2)`
	actual = testEval(input)
	expectIntegerValue(t, actual, 6)

	input = `
	var x = 999999
	var foo = fn(x, y) { x + y}
	foo(foo(1, 2), 1 + 2)`
	actual = testEval(input)
	expectIntegerValue(t, actual, 6)

	input = `
	var x = 999999
	var foo = fn(y) { x + y }
	foo(1)`
	actual = testEval(input)
	expectIntegerValue(t, actual, 1000000)

	input = `
	var foo = fn(y) { x + y }
	foo(1)`
	actual = testEval(input)
	expectError(t, actual, "unknown identifier x")

	input = `
	var x = 0
	var foo = fn(y) { x + y }
	var bar = fn(y) {
		var x = 9999
		foo(y)
	}
	bar(1)`
	actual = testEval(input)
	expectIntegerValue(t, actual, 1)

	input = `
	var foo = fn(x) {
		if x > 0 {
			return true
		}
		return false
	}(1)
	if foo {
		return 123
	} else {
	 	return 666
	}`
	actual = testEval(input)
	expectIntegerValue(t, actual, 123)

	input = `
	var newAdder = fn(x) {
		return fn(y) { x + y }
	}
	var addTwo = newAdder(2)
	addTwo(2)`
	actual = testEval(input)
	expectIntegerValue(t, actual, 4)

	input = `
	var foo = fn(x, func) {
		return func(x) 
	}
	foo(
	 1,
	 fn(x) {
	 	return x + 4
	 }
	)`
	actual = testEval(input)
	expectIntegerValue(t, actual, 5)
}

func TestEvalFunctionLiteral(t *testing.T) {
	input := `fn(x) {
		if x > 0 {
			return true
		}
		return false
	}`
	fn := testEval(input)
	expectFunctionParameters(t, fn, "x")
	expectFunctionBody(t, fn, "if (x > 0) { return true } return false")
}

func expectFunctionBody(t *testing.T, actual evaluator.Object, expect string) {
	fn, ok := actual.(*evaluator.Function)
	if !ok {
		t.Fatalf("object is not a function. got %T", actual)
	}
	if fn.Body.String() != expect {
		t.Fatalf("function body wrong. expect\n%s\ngot\n%s", expect, fn.Body.String())
	}
}

func expectFunctionParameters(t *testing.T, actual evaluator.Object, expect ...string) {
	fn, ok := actual.(*evaluator.Function)
	if !ok {
		t.Fatalf("object is not a function. got %T", actual)
	}
	if len(fn.Parameters) != len(expect) {
		t.Fatalf("number of parameters wrong. expect: %d. got: %d", len(expect), len(fn.Parameters))
	}

	for i, param := range fn.Parameters {
		if param.Value != expect[i] {
			t.Fatalf("parameter %d expected to be %s. got %s", i, expect[i], param.Value)
		}
	}
}

func TestEvalDeclaration(t *testing.T) {
	input := `var x = 1
	x`
	actual := testEval(input)
	expectIntegerValue(t, actual, 1)
}

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

	input = `123 + "foo"`
	actual = testEval(input)
	expectError(t, actual, "operator type mismatch. INTEGER + STRING")

	input = `x`
	actual = testEval(input)
	expectError(t, actual, "unknown identifier x")
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

	input = `123 == "123"`
	actual = testEval(input)
	expectBooleanValue(t, actual, false)
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

func TestEvalStringLitearl(t *testing.T) {
	input := `"foo bar"`
	actual := testEval(input)
	expectStringValue(t, actual, "foo bar")

	input = `"foo" + " " + "bar"`
	actual = testEval(input)
	expectStringValue(t, actual, "foo bar")

	input = `"foo bar" == "foo" + " " + "bar"`
	actual = testEval(input)
	expectBooleanValue(t, actual, true)

	input = `"foo bar" != "baz"`
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
		if err, ok := v.(*evaluator.Error); ok {
			t.Fatalf("Expected boolean value, got error\n%s", err.Message)
		}
		t.Fatalf("Expected boolean value, got %T", v)
	}
	if obj.Value != expect {
		t.Fatalf("Expected evaluated value to be %t. got %t", expect, obj.Value)
	}
}

func expectStringValue(t *testing.T, v evaluator.Object, expect string) {
	obj, ok := v.(*evaluator.String)
	if !ok {
		t.Fatalf("Expected string value %s, got %T %s", expect, v, v.Inspect())
	}
	if obj.Value != expect {
		t.Fatalf("Expected evaluated value to be %s. got %s", expect, obj.Value)
	}
}

func expectIntegerValue(t *testing.T, v evaluator.Object, expect int64) {
	obj, ok := v.(*evaluator.Integer)
	if !ok {
		t.Fatalf("Expected integer value %d, got %T %s", expect, v, v.Inspect())
	}
	if obj.Value != expect {
		t.Fatalf("Expected evaluated value to be %d. got %d", expect, obj.Value)
	}
}

func testEval(input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvrionment()
	return evaluator.Eval(program, env)
}
