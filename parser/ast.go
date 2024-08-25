package parser

import (
	"bytes"

	"github.com/maiksch/best-lang/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p Program) String() string {
	var out bytes.Buffer

	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

// Identifier Expression

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Token.Literal }

// Integer Literal Expression

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return i.Token.Literal }

// Prefixed Expression
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Operand  Expression
}

func (p *PrefixExpression) expressionNode()      {}
func (p *PrefixExpression) TokenLiteral() string { return p.Token.Literal }
func (p *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Operand.String())
	out.WriteString(")")

	return out.String()
}

// Declare Statement

type DeclareStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (d *DeclareStatement) statementNode()       {}
func (d *DeclareStatement) TokenLiteral() string { return d.Token.Literal }
func (d *DeclareStatement) String() string {
	var out bytes.Buffer

	out.WriteString(d.Name.String())
	out.WriteString(" := ")

	if d.Value != nil {
		out.WriteString(d.Value.String())
	}

	return out.String()
}

// Return Statement

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (r *ReturnStatement) statementNode()       {}
func (r *ReturnStatement) TokenLiteral() string { return r.Token.Literal }
func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral())
	out.WriteString(" ")

	if r.Value != nil {
		out.WriteString(r.Value.String())
	}

	return out.String()
}

// Expression Statement

type ExpressionStatement struct {
	Token token.Token
	Value Expression
}

func (e *ExpressionStatement) statementNode()       {}
func (e *ExpressionStatement) TokenLiteral() string { return e.Token.Literal }
func (e *ExpressionStatement) String() string {
	if e.Value != nil {
		return e.Value.String()
	}
	return ""
}
