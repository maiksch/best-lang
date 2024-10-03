package lexer

import "github.com/maiksch/best-lang/token"

type Lexer struct {
	input    string
	position int
}

func New(input string) *Lexer {
	return &Lexer{
		input:    input,
		position: -1,
	}
}

func (l *Lexer) readLiteral() string {
	position := l.position

	for {
		char := l.readChar()
		if !l.isLetter(char) {
			break
		}
	}

	word := l.input[position:l.position]

	l.position -= 1

	return word
}

func (l *Lexer) isLetter(char byte) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position

	for {
		char := l.readChar()
		if !l.isDigit(char) {
			break
		}
	}

	number := l.input[position:l.position]

	l.position -= 1

	return number
}

func (l *Lexer) isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func (l *Lexer) isWhitespace(char byte) bool {
	return char == ' ' || char == '\t'
}

func (l *Lexer) readChar() byte {
	l.position += 1
	if l.position >= len(l.input) {
		return 0
	}
	return l.input[l.position]
}

func (l *Lexer) peekChar() byte {
	position := l.position + 1
	if position >= len(l.input) {
		return 0
	}
	return l.input[position]
}

func (l *Lexer) NextToken() token.Token {
	ch := l.readChar()

	for l.isWhitespace(ch) {
		ch = l.readChar()
	}

	if ch == 0 {
		return token.Token{Type: token.EOF, Literal: "EOF"}
	}

	if t, ok := token.Symbols[ch]; ok {
		// Two symbol tokens
		switch t {
		case token.NEWLINE:
			return token.Token{Type: token.NEWLINE}
		case token.ASSIGN:
			if peek := l.peekChar(); peek == '=' {
				// Equal operator ==
				l.readChar()
				return token.Token{Type: token.EQUAL, Literal: token.EQUAL}
			}
		case token.BANG:
			if peek := l.peekChar(); peek == '=' {
				l.readChar()
				return token.Token{Type: token.NOT_EQUAL, Literal: token.NOT_EQUAL}
			}
		}

		// One symbol tokens
		return token.Token{Type: t, Literal: string(ch)}
	}

	if l.isLetter(ch) {
		literal := l.readLiteral()
		tokenType := token.GetWordTokenType(literal)
		return token.Token{Type: tokenType, Literal: literal}
	}

	if l.isDigit(ch) {
		number := l.readNumber()
		return token.Token{Type: token.INTEGER, Literal: number}
	}

	return token.Token{Type: token.ILLEGAL}
}
