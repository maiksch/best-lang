package parser_test

import (
	"testing"

	"github.com/maiksch/best-lang/lexer"
	"github.com/maiksch/best-lang/parser"
	"github.com/maiksch/best-lang/token"
)

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    int64
	}{
		{"-1", "-", 1}, {"!4", "!", 4},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)

		program := p.ParseProgram()
		if program == nil {
			t.Fatalf("Program nil")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("Wrong number of statements.\n\tExpected: %d\n\tGot: %d", 1, len(program.Statements))
		}

		stmt := assertExpressionStatement(t, program.Statements[0])
		exp := assertPrefixExpression(t, stmt.Value)
		assertIntegerLiteral(t, exp.Operand, test.value)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `123
456`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("Program nil")
	}

	expect := []int64{123, 456}

	if len(program.Statements) != len(expect) {
		t.Fatalf("Wrong number of statements.\n\tExpected: %d\n\tGot: %d", len(expect), len(program.Statements))
	}

	for i, test := range expect {
		stmt := assertExpressionStatement(t, program.Statements[i])
		assertIntegerLiteral(t, stmt.Value, test)
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foo
bla`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("Program nil")
	}

	expect := []string{"foo", "bla"}

	if len(program.Statements) != len(expect) {
		t.Fatalf("Wrong number of statements.\n\tExpected: %d\n\tGot: %d", len(expect), len(program.Statements))
	}

	for i, test := range expect {
		stmt := assertExpressionStatement(t, program.Statements[i])
		if _, ok := stmt.Value.(*parser.Identifier); !ok {
			t.Fatalf("expresson is not of type Identifier. got %T", stmt.Value)
		}
		if stmt.String() != test {
			t.Fatalf("wrong expression.\n\texpected %q\n\tgot %q", test, stmt.String())
		}
	}
}

func TestToString(t *testing.T) {
	input := parser.Program{
		Statements: []parser.Statement{
			&parser.DeclareStatement{
				Token: token.Token{Type: token.VARIABLE, Literal: "var"},
				Name: &parser.Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "foo"},
					Value: "foo",
				},
				Value: &parser.Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "bar"},
					Value: "bar",
				},
			},
		},
	}

	result := input.String()
	expect := "foo := bar"

	if result != expect {
		t.Fatalf("String() failed.\n\texpected:%q\n\tgot:%q", expect, result)
	}
}

func TestReturnStatement(t *testing.T) {
	input := `return 5
return 10
return test(1, 2)`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("Program nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Wrong number of statements.\n\tExpected: 3\n\tGot: %d", len(program.Statements))
	}

	tests := []string{
		"return",
		"return",
		"return",
	}

	for i, test := range tests {
		statement := program.Statements[i].(*parser.ReturnStatement)
		if statement.TokenLiteral() != test {
			t.Fatalf("test failed: token value wrong.\n\texpected %q\n\tgot %q", test, statement.TokenLiteral())
		}
	}
}

func TestAssertNextToken(t *testing.T) {
	input := `var x 5`
	expect := "invalid syntax. Expected \"=\" but got \"INTEGER\""

	l := lexer.New(input)
	p := parser.New(l)

	defer func() {
		if r := recover(); r != expect {
			t.Fatalf("did not panic as expected\n\texpected: %q\n\tgot: %q", expect, r)
		}
	}()

	p.ParseProgram()
}

func TestDeclareStatement(t *testing.T) {
	input := `var x = 5
var y = 10
var foo = 1234`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("Program nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Wrong number of statements.\n\tExpected: 3\n\tGot: %d", len(program.Statements))
	}

	tests := []string{
		"x",
		"y",
		"foo",
	}

	for i, test := range tests {
		statement := program.Statements[i].(*parser.DeclareStatement)
		if statement.Name.Value != test {
			t.Fatalf("test failed: token value wrong.\n\texpected %q\n\tgot %q", test, statement.TokenLiteral())
		}
		if statement.Name.TokenLiteral() != test {
			t.Fatalf("test failed: token literal wrong.\n\texpected %q\n\tgot %q", test, statement.TokenLiteral())
		}
	}
}

func assertPrefixExpression(t *testing.T, expression parser.Expression) *parser.PrefixExpression {
	exp, ok := expression.(*parser.PrefixExpression)
	if !ok {
		t.Fatalf("expresson is not of type PrefixExpression. got %T", expression)
	}
	return exp
}

func assertExpressionStatement(t *testing.T, statement parser.Statement) *parser.ExpressionStatement {
	expressionStatement, ok := statement.(*parser.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not an ExpressionStatement. got %T", statement)
	}
	return expressionStatement
}

func assertIntegerLiteral(t *testing.T, exp parser.Expression, expect int64) {
	integerLiteralExp, ok := exp.(*parser.IntegerLiteral)
	if !ok {
		t.Fatalf("expresson is not of type IntegerLiteral. got %T", exp)
	}
	if integerLiteralExp.Value != expect {
		t.Fatalf("wrong value.\n\texpected %q\n\tgot %q", expect, integerLiteralExp.Value)
	}
}
