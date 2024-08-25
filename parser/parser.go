package parser

import (
	"log"
	"strconv"

	"github.com/maiksch/best-lang/lexer"
	"github.com/maiksch/best-lang/token"
)

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

type Parser struct {
	lexer *lexer.Lexer

	token     token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INTEGER, p.parseIntegerLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)

	// Read two tokens, so token and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(token token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[token] = fn
}

func (p *Parser) registerInfix(token token.TokenType, fn infixParseFn) {
	p.infixParseFns[token] = fn
}

func (p *Parser) nextToken() {
	p.token = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}

	for p.token.Type != token.EOF {
		statement := p.parseStatement()

		program.Statements = append(program.Statements, statement)

		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.token.Type {
	case token.VARIABLE:
		return p.parseDeclarationStmt()
	case token.RETURN:
		return p.parseReturnStmt()
	default:
		return p.parseExpressionStmt()
	}
}

func (p *Parser) parseDeclarationStmt() *DeclareStatement {
	s := &DeclareStatement{}

	s.Token = p.token

	p.assertNextToken(token.IDENTIFIER)

	s.Name = &Identifier{Token: p.token, Value: p.token.Literal}

	p.assertNextToken(token.ASSIGN)

	// todo: real expression parsing
	for p.token.Type != token.NEWLINE && p.token.Type != token.EOF {
		p.nextToken()
	}

	return s
}

func (p *Parser) parseReturnStmt() *ReturnStatement {
	s := &ReturnStatement{}

	s.Token = p.token

	// todo: real expression parsing
	for p.token.Type != token.NEWLINE && p.token.Type != token.EOF {
		p.nextToken()
	}

	return s
}

func (p *Parser) parseExpressionStmt() *ExpressionStatement {
	stmt := &ExpressionStatement{
		Token: p.token,
		Value: p.parseExpression(),
	}

	// todo: real expression parsing
	for p.token.Type != token.NEWLINE && p.token.Type != token.EOF {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression() Expression {
	prefix, ok := p.prefixParseFns[p.token.Type]
	if !ok {
		return nil
	}
	return prefix()
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{
		Token: p.token,
		Value: p.token.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() Expression {
	value, err := strconv.ParseInt(p.token.Literal, 0, 64)
	if err != nil {
		log.Fatal(err)
	}
	return &IntegerLiteral{
		Token: p.token,
		Value: value,
	}
}

func (p *Parser) parsePrefixExpression() Expression {
	prefixExp := &PrefixExpression{}

	prefixExp.Token = p.token
	prefixExp.Operator = p.token.Literal

	p.nextToken()

	prefixExp.Operand = p.parseExpression()

	return prefixExp
}

func (p *Parser) assertNextToken(t token.TokenType) {
	if p.peekToken.Type != t {
		log.Panicf("invalid syntax. Expected %q but got %q", t, p.peekToken.Type)
	}
	p.nextToken()
}
