package lexer_test

import (
	"testing"

	"github.com/maiksch/best-lang/lexer"
	"github.com/maiksch/best-lang/token"
)

type expectation struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func TestKeywords(t *testing.T) {
	input := `if else true false fn return`

	tests := []expectation{
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.FUNCTION, "fn"},
		{token.RETURN, "return"},
	}

	runAndExpect(t, input, tests)
}

func TestStringLiteral(t *testing.T) {
	input := `"test"
	"foo bar"`

	tests := []expectation{
		{token.STRING, "test"},
		{token.NEWLINE, ""},
		{token.STRING, "foo bar"},
	}

	runAndExpect(t, input, tests)
}

func TestBasicProgram(t *testing.T) {
	input := `var five = 5
var ten = 10

fn add(x, y) {
  return x + y
}

var result = add(five, ten)
"test"`

	tests := []expectation{
		{token.VARIABLE, "var"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INTEGER, "5"},
		{token.NEWLINE, ""},
		{token.VARIABLE, "var"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INTEGER, "10"},
		{token.NEWLINE, ""},
		{token.NEWLINE, ""},
		{token.FUNCTION, "fn"},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.KOMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.NEWLINE, ""},
		{token.RETURN, "return"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.NEWLINE, ""},
		{token.RBRACE, "}"},
		{token.NEWLINE, ""},
		{token.NEWLINE, ""},
		{token.VARIABLE, "var"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.KOMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RPAREN, ")"},
		{token.NEWLINE, ""},
		{token.EOF, ""},
	}

	runAndExpect(t, input, tests)
}

func TestWhitespace(t *testing.T) {
	input := `var   	five     =   	5`

	tests := []expectation{
		{token.VARIABLE, "var"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INTEGER, "5"},
		{token.EOF, ""},
	}

	runAndExpect(t, input, tests)
}

func TestAssignment(t *testing.T) {
	input := `var five = 5`

	tests := []expectation{
		{token.VARIABLE, "var"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INTEGER, "5"},
		{token.EOF, ""},
	}

	runAndExpect(t, input, tests)
}

func TestSymbols(t *testing.T) {
	input := `+-/*()={},==<>!=
	`

	tests := []expectation{
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.STAR, "*"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.ASSIGN, "="},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.KOMMA, ","},
		{token.EQUAL, "=="},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.NOT_EQUAL, "!="},
		{token.NEWLINE, ""},
		{token.EOF, ""},
	}

	runAndExpect(t, input, tests)
}

func runAndExpect(t *testing.T, input string, tests []expectation) {
	l := lexer.New(input)

	for _, test := range tests {
		token := l.NextToken()

		if token.Type != test.expectedType {
			t.Fatalf("test failed: token type wrong.\n\texpected %q\n\tgot %q", test.expectedType, token.Type)
		}

		if token.Literal != test.expectedLiteral {
			t.Fatalf("test failed: token literal wrong.\n\texpected %q\n\tgot %q", test.expectedLiteral, token.Literal)
		}
	}
}
