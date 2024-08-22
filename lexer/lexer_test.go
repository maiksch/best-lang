package lexer

import (
	"testing"

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

func TestBasicProgram(t *testing.T) {
	input := `five := 5
ten := 10

fn add(x, y) {
  return x + y
}

result := add(five, ten)`

	tests := []expectation{
		{token.IDENTIFIER, "five"},
		{token.DECLARE, ":="},
		{token.INTEGER, "5"},
		{token.NEWLINE, "\n"},
		{token.IDENTIFIER, "ten"},
		{token.DECLARE, ":="},
		{token.INTEGER, "10"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.FUNCTION, "fn"},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.KOMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.NEWLINE, "\n"},
		{token.RETURN, "return"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.NEWLINE, "\n"},
		{token.RBRACE, "}"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.IDENTIFIER, "result"},
		{token.DECLARE, ":="},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.KOMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	}

	runAndExpect(t, input, tests)
}

func TestAssignment(t *testing.T) {
	input := `five := 5`

	tests := []expectation{
		{token.IDENTIFIER, "five"},
		{token.DECLARE, ":="},
		{token.INTEGER, "5"},
		{token.EOF, ""},
	}

	runAndExpect(t, input, tests)
}

func TestWhitespace(t *testing.T) {
	input := `five     :=   	5`

	tests := []expectation{
		{token.IDENTIFIER, "five"},
		{token.DECLARE, ":="},
		{token.INTEGER, "5"},
		{token.EOF, ""},
	}

	runAndExpect(t, input, tests)
}

func TestSymbols(t *testing.T) {
	input := `+-/*():={},==<>!=
	`

	tests := []expectation{
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.STAR, "*"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.DECLARE, ":="},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.KOMMA, ","},
		{token.EQUAL, "=="},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.NOT_EQUAL, "!="},
		{token.NEWLINE, "\n"},
		{token.EOF, ""},
	}

	runAndExpect(t, input, tests)
}

func runAndExpect(t *testing.T, input string, tests []expectation) {
	l := NewLexer(input)

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
