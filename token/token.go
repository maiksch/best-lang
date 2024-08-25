package token

type TokenType string

const (
	ILLEGAL    = "ILLEGAL"
	EOF        = "EOF"
	IDENTIFIER = "IDENTIFIER"
	INTEGER    = "INTEGER"

	// Symbols
	ASSIGN  = "="
	LPAREN  = "("
	RPAREN  = ")"
	LBRACE  = "{"
	RBRACE  = "}"
	KOMMA   = ","
	PLUS    = "+"
	MINUS   = "-"
	SLASH   = "/"
	STAR    = "*"
	LT      = "<"
	GT      = ">"
	BANG    = "!"
	NEWLINE = "NEWLINE"

	// Two Symbol Tokens
	EQUAL     = "=="
	NOT_EQUAL = "!="

	// Keywords
	VARIABLE = "VAR"
	FUNCTION = "FN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

type Token struct {
	Type    TokenType
	Literal string
}

var Symbols = map[byte]TokenType{
	'=':  ASSIGN,
	'(':  LPAREN,
	')':  RPAREN,
	'{':  LBRACE,
	'}':  RBRACE,
	',':  KOMMA,
	'+':  PLUS,
	'-':  MINUS,
	'/':  SLASH,
	'*':  STAR,
	'<':  LT,
	'>':  GT,
	'!':  BANG,
	'\n': NEWLINE,
}

var keywords = map[string]TokenType{
	"var":    VARIABLE,
	"fn":     FUNCTION,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func GetWordTokenType(word string) TokenType {
	if token, ok := keywords[word]; ok {
		return token
	}
	return IDENTIFIER
}
