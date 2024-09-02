package parser_test

import (
	"testing"

	"github.com/maiksch/best-lang/lexer"
	"github.com/maiksch/best-lang/parser"
	"github.com/maiksch/best-lang/token"
)

func TestFunctionExpression(t *testing.T) {
	input := "fn(a, b) { return a + b }"

	l := lexer.New(input)
	p := parser.New(l).ParseProgram()

	stmt := expectExpressionStatement(t, p.Statements[0])
	fnExpr, ok := stmt.Value.(*parser.FunctionExpression)
	if !ok {
		t.Fatalf("expression is not an IfExpression. got %T", stmt.Value)
	}
	if len(fnExpr.Parameters) != 2 {
		t.Fatalf("amount of parameters wrong. expected %v but got %v", 1, len(fnExpr.Parameters))
	}
	expectIdentifier(t, &fnExpr.Parameters[0], "a")
	expectIdentifier(t, &fnExpr.Parameters[1], "b")
	returnStmt := expectReturnStatement(t, fnExpr.Body.Statements[0])
	expectInfixExpression(t, returnStmt.Value, "a", "+", "b")
}

func TestIfExpressions(t *testing.T) {
	input := "if true == 1 { 1 } else { 2 }"

	l := lexer.New(input)
	p := parser.New(l).ParseProgram()
	expectStatements(t, p, 1)

	stmt := expectExpressionStatement(t, p.Statements[0])

	ifExpr, ok := stmt.Value.(*parser.IfExpression)
	if !ok {
		t.Fatalf("expression is not an IfExpression. got %T", stmt.Value)
	}
	expectInfixExpression(t, ifExpr.Condition, true, "==", 1)
	consequence := expectExpressionStatement(t, ifExpr.Consequence.Statements[0])
	expectLiteralExpression(t, consequence.Value, 1)
	otherwise := expectExpressionStatement(t, ifExpr.Otherwise.Statements[0])
	expectLiteralExpression(t, otherwise.Value, 2)
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"1 + 2", "(1 + 2)"},
		{"1 - 2 * 3", "(1 - (2 * 3))"},
		{"1 * 2 + 3", "((1 * 2) + 3)"},
		{"1 == 2 * 3", "(1 == (2 * 3))"},
		{"1 > false + 3 * 4 < 2", "((1 > (false + (3 * 4))) < 2)"},
		{"1 == 2 * 4 + !5 < 6 / true < 3", "(1 == ((((2 * 4) + (!5)) < (6 / true)) < 3))"},
		{"(1 + 2) * 3", "((1 + 2) * 3)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l).ParseProgram()

		expectProgram(t, p, test.expect)
	}
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"1 + 2", 1, "+", 2},
		{"1 - 2", 1, "-", 2},
		{"1 * 2", 1, "*", 2},
		{"1 / 2", 1, "/", 2},
		{"1 > 2", 1, ">", 2},
		{"1 < 2", 1, "<", 2},
		{"1 == 2", 1, "==", 2},
		{"true != false", true, "!=", false},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l).ParseProgram()

		expectStatements(t, p, 1)
		stmt := expectExpressionStatement(t, p.Statements[0])
		expectInfixExpression(t, stmt.Value, test.leftValue, test.operator, test.rightValue)
	}
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"-1", "-", 1}, {"!true", "!", true},
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

		stmt := expectExpressionStatement(t, program.Statements[0])
		exp := expectPrefixExpression(t, stmt.Value)
		expectLiteralExpression(t, exp.Right, test.value)
	}
}

func TestLiteralExpression(t *testing.T) {
	input := `123
456
false
true
identifier`

	l := lexer.New(input)
	p := parser.New(l).ParseProgram()

	expect := []interface{}{123, 456, false, true, "identifier"}

	expectStatements(t, p, len(expect))

	for i, test := range expect {
		stmt := expectExpressionStatement(t, p.Statements[i])
		expectLiteralExpression(t, stmt.Value, test)
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
return identifier
return true`

	l := lexer.New(input)
	p := parser.New(l).ParseProgram()

	expect := []interface{}{5, 10, "identifier", true}

	expectStatements(t, p, len(expect))

	for i, test := range expect {
		statement := p.Statements[i].(*parser.ReturnStatement)
		if statement.TokenLiteral() != "return" {
			t.Fatalf("test failed: token value wrong.\n\texpected: %q\n\tgot:     %q", test, statement.TokenLiteral())
		}
		expectLiteralExpression(t, statement.Value, test)
	}
}

func expectReturnStatement(t *testing.T, stmt parser.Statement) *parser.ReturnStatement {
	returnStmt := stmt.(*parser.ReturnStatement)
	if returnStmt.Token.Type != token.RETURN {
		t.Fatalf("test failed: return statement expected got: %q", returnStmt.Token.Type)
	}
	return returnStmt
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
var y = x
var foo = true`

	l := lexer.New(input)
	p := parser.New(l).ParseProgram()

	expect := []struct {
		identifier string
		value      interface{}
	}{
		{identifier: "x", value: 5},
		{identifier: "y", value: "x"},
		{identifier: "foo", value: true},
	}

	expectStatements(t, p, len(expect))

	for i, test := range expect {
		statement := p.Statements[i].(*parser.DeclareStatement)
		if statement.Name.String() != test.identifier {
			t.Fatalf("test failed: identifier wrong.\n\texpected %q\n\tgot %q", test, statement.TokenLiteral())
		}
		expectLiteralExpression(t, statement.Value, test.value)
	}
}

func expectProgram(t *testing.T, p *parser.Program, expect string) {
	actual := p.String()
	if actual != expect {
		t.Errorf("wrong program.\n\texpected: %v\n\tgot:      %v", expect, actual)
	}
}

func expectPrefixExpression(t *testing.T, expression parser.Expression) *parser.PrefixExpression {
	exp, ok := expression.(*parser.PrefixExpression)
	if !ok {
		t.Fatalf("expresson is not of type PrefixExpression. got %T", expression)
	}
	return exp
}

func expectExpressionStatement(t *testing.T, statement parser.Statement) *parser.ExpressionStatement {
	expressionStatement, ok := statement.(*parser.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not an ExpressionStatement. got %T", statement)
	}
	return expressionStatement
}

func expectLiteralExpression(t *testing.T, expr parser.Expression, expect interface{}) {
	switch v := expect.(type) {
	case int:
		expectIntegerLiteral(t, expr, int64(v))
	case int64:
		expectIntegerLiteral(t, expr, v)
	case bool:
		expectBooleanLiteral(t, expr, v)
	case string:
		expectIdentifier(t, expr, v)
	default:
		t.Fatalf("Check for undefined literal expession of type %T", v)
	}
}

func expectInfixExpression(t *testing.T, expr parser.Expression, left interface{}, operator string, right interface{}) {
	exp, ok := expr.(*parser.InfixExpression)
	if !ok {
		t.Fatalf("expresson is not of type InfixExpression. got %T", expr)
	}

	if exp.Operator != operator {
		t.Fatalf("wrong operator.\n\tGot:    %s\n\tExpect: %s", exp.Operator, operator)
	}

	expectLiteralExpression(t, exp.Left, left)
	expectLiteralExpression(t, exp.Right, right)
}

func expectIdentifier(t *testing.T, expr parser.Expression, expect string) {
	identifier, ok := expr.(*parser.Identifier)
	if !ok {
		t.Fatalf("expression is not of type Identifier. Got %T", expr)
	}

	if identifier.Value != expect {
		t.Fatalf("wrong identifier name.\n\tGot    %s\n\tExpect:%s", identifier.Value, expect)
	}
}

func expectBooleanLiteral(t *testing.T, expr parser.Expression, expect bool) {
	booleanLit, ok := expr.(*parser.BooleanLiteral)
	if !ok {
		t.Fatalf("expresson is not of type BooleanLiteral. got %T", expr)
	}
	if booleanLit.Value != expect {
		t.Fatalf("wrong value.\n\texpected: %v\n\tgot:      %v", expect, booleanLit.Value)
	}
}

func expectIntegerLiteral(t *testing.T, exp parser.Expression, expect int64) {
	integerLiteralExp, ok := exp.(*parser.IntegerLiteral)
	if !ok {
		t.Fatalf("expresson is not of type IntegerLiteral. got %T", exp)
	}
	println("Comparing", integerLiteralExp.Value, expect)
	if integerLiteralExp.Value != expect {
		t.Fatalf("wrong value.\n\texpected: %v\n\tgot:      %v", expect, integerLiteralExp.Value)
	}
}

func expectStatements(t *testing.T, p *parser.Program, expect int) {
	if len(p.Statements) != expect {
		t.Fatalf("Wrong number of statements.\n\tExpected: %d\n\tGot: %d", 1, len(p.Statements))
	}
}
