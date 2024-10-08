package parser

import (
	"log"
	"strconv"

	"github.com/maiksch/best-lang/lexer"
	"github.com/maiksch/best-lang/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x or !x
	CALL        // myFn()
)

var precedences = map[token.TokenType]int{
	token.EQUAL:     EQUALS,
	token.NOT_EQUAL: EQUALS,
	token.LT:        LESSGREATER,
	token.GT:        LESSGREATER,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.STAR:      PRODUCT,
	token.SLASH:     PRODUCT,
	token.LPAREN:    CALL,
}

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
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.EQUAL, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.STAR, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseFunctionCall)

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
	// log.Printf("%s %s\n", p.token.Type, p.token.Literal)
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
	for p.token.Type == token.NEWLINE {
		p.nextToken()
	}

	switch p.token.Type {
	case token.VARIABLE:
		return p.parseDeclarationStmt()
	case token.RETURN:
		return p.parseReturnStmt()
	default:
		return p.parseExpressionStmt()
	}
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	blockStmt := &BlockStatement{
		Token: p.token,
	}

	for !p.isPeekToken(token.EOF) &&
		!p.isPeekToken(token.RBRACE) &&
		!(p.token.Type == token.RBRACE && p.peekToken.Type == token.NEWLINE) {
		p.nextToken()

		stmt := p.parseStatement()

		blockStmt.Statements = append(blockStmt.Statements, stmt)
	}

	return blockStmt
}

func (p *Parser) parseDeclarationStmt() *DeclareStatement {
	s := &DeclareStatement{Token: p.token}

	p.assertNextToken(token.IDENTIFIER)

	s.Name = &Identifier{Token: p.token, Value: p.token.Literal}

	p.assertNextToken(token.ASSIGN)

	p.nextToken()

	s.Expression = p.parseExpression(LOWEST)

	p.assertEnd()

	return s
}

func (p *Parser) parseReturnStmt() *ReturnStatement {
	s := &ReturnStatement{Token: p.token}

	p.nextToken()

	s.Expression = p.parseExpression(LOWEST)

	p.assertEnd()

	return s
}

func (p *Parser) parseExpressionStmt() *ExpressionStatement {
	stmt := &ExpressionStatement{
		Token: p.token,
		Value: p.parseExpression(LOWEST),
	}

	if p.peekToken.Type == token.NEWLINE {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefix, ok := p.prefixParseFns[p.token.Type]
	if !ok {
		log.Printf("No prefix fn found for token %q", p.token.Literal)
		return nil
	}

	exp := prefix()

	for p.peekToken.Type != token.NEWLINE && precedence < p.peekPrecedence() {
		p.nextToken()

		infix, ok := p.infixParseFns[p.token.Type]
		if !ok {
			log.Panicf("No infix fn found for token %q", p.token.Literal)
		}

		exp = infix(exp)
	}

	return exp
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
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

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{
		Token: p.token,
		Value: p.token.Literal,
	}
}

func (p *Parser) parseBooleanLiteral() Expression {
	return &BooleanLiteral{
		Token: p.token,
		Value: p.token.Type == token.TRUE,
	}
}

func (p *Parser) parseIfExpression() Expression {
	expr := &IfExpression{
		Token: p.token,
	}

	p.nextToken()

	expr.Condition = p.parseExpression(LOWEST)

	if !p.isPeekToken(token.LBRACE) {
		log.Println("if expression missing opening {")
		return nil
	}

	expr.Consequence = p.parseBlockStatement()

	if p.isPeekToken(token.ELSE) {
		if !p.isPeekToken(token.LBRACE) {
			log.Println("else block is missing opening {")
			return nil
		}
		expr.Otherwise = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseFunctionExpression() Expression {
	expr := &FunctionLiteral{Token: p.token}

	if !p.isPeekToken(token.LPAREN) {
		log.Println("function is missing opening (")
		return nil
	}

	for !p.isPeekToken(token.RPAREN) {
		p.isPeekToken(token.KOMMA)

		if !p.isPeekToken(token.IDENTIFIER) {
			log.Println("function parameter found is not an identifier")
			return nil
		}
		expr.Parameters = append(expr.Parameters, &Identifier{
			Token: p.token,
			Value: p.token.Literal,
		})
	}

	if !p.isPeekToken(token.LBRACE) {
		log.Println("function body opening { missing")
		return nil
	}

	expr.Body = p.parseBlockStatement()

	return expr
}

func (p *Parser) parseFunctionCall(left Expression) Expression {
	expr := &FunctionCall{
		Token:    p.token,
		Function: left,
	}

	p.skipNewline()

	if p.isPeekToken(token.RPAREN) {
		return expr
	}

	for {
		p.skipNewline()
		p.nextToken()
		arg := p.parseExpression(LOWEST)
		expr.Arguments = append(expr.Arguments, arg)

		if !p.isPeekToken(token.KOMMA) {
			break
		}
	}

	p.skipNewline()

	if !p.isPeekToken(token.RPAREN) {
		log.Println("function call: closing ) missing")
		return nil
	}

	return expr
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()

	expr := p.parseExpression(LOWEST)

	if !p.isPeekToken(token.RPAREN) {
		log.Println("opened ( is missing closing )")
		return nil
	}

	return expr
}

func (p *Parser) parsePrefixExpression() Expression {
	prefixExp := &PrefixExpression{}

	prefixExp.Token = p.token
	prefixExp.Operator = p.token.Literal

	p.nextToken()

	prefixExp.Right = p.parseExpression(PREFIX)

	return prefixExp
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	infixExp := &InfixExpression{
		Token:    p.token,
		Operator: p.token.Literal,
		Left:     left,
	}

	precedence, ok := precedences[p.token.Type]
	if !ok {
		precedence = LOWEST
	}

	p.nextToken()

	infixExp.Right = p.parseExpression(precedence)

	return infixExp
}

func (p *Parser) assertEnd() {
	if p.peekToken.Type != token.NEWLINE && p.peekToken.Type != token.EOF && p.peekToken.Type != token.RBRACE {
		log.Panicf("invalid syntax. Expected end of statement but got %q", p.peekToken.Type)
	}
	p.nextToken()
}

func (p *Parser) assertNextToken(t token.TokenType) {
	if p.peekToken.Type != t {
		log.Panicf("invalid syntax. Expected %q but got %q", t, p.peekToken.Type)
	}
	p.nextToken()
}

func (p *Parser) skipNewline() {
	if p.peekToken.Type == token.NEWLINE {
		p.nextToken()
	}
}

func (p *Parser) isPeekToken(t token.TokenType) bool {
	ok := p.peekToken.Type == t
	if ok {
		p.nextToken()
	}
	return ok
}
