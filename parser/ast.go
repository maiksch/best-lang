package parser

import (
	"bytes"
	"fmt"

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

// String Literal Expression

type StringLiteral struct {
	Token token.Token
	Value string
}

func (i *StringLiteral) expressionNode()      {}
func (i *StringLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *StringLiteral) String() string       { return i.Token.Literal }

// Boolean Literal Expression

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode()      {}
func (b *BooleanLiteral) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanLiteral) String() string {
	return fmt.Sprint(b.Value)
}

// Prefix Expression
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode()      {}
func (p *PrefixExpression) TokenLiteral() string { return p.Token.Literal }
func (p *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

// Infix Expression

type InfixExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (i *InfixExpression) expressionNode()      {}
func (i *InfixExpression) TokenLiteral() string { return i.Token.Literal }
func (i *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
}

// If Expression

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Otherwise   *BlockStatement
}

func (i *IfExpression) expressionNode()      {}
func (i *IfExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(i.Condition.String())
	out.WriteString(" { ")
	out.WriteString(i.Consequence.String())
	out.WriteString(" } ")
	if i.Otherwise != nil {
		out.WriteString("else { ")
		out.WriteString(i.Otherwise.String())
		out.WriteString(" } ")
	}

	return out.String()
}

// Function Literal

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *FunctionLiteral) expressionNode()      {}
func (f *FunctionLiteral) TokenLiteral() string { return f.Token.Literal }
func (f *FunctionLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("fn(")
	for i, param := range f.Parameters {
		if i != 0 {
			out.WriteString(", ")
		}
		out.WriteString(param.String())
	}
	out.WriteString(")")

	out.WriteString("{")
	out.WriteString(f.Body.String())
	out.WriteString("}")

	return out.String()
}

type FunctionCall struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (f *FunctionCall) expressionNode()      {}
func (f *FunctionCall) TokenLiteral() string { return f.Token.Literal }
func (f *FunctionCall) String() string {
	var out bytes.Buffer

	out.WriteString(f.Function.String())
	out.WriteString("(")
	for i, arg := range f.Arguments {
		if i != 0 {
			out.WriteString(", ")
		}
		out.WriteString(arg.String())
	}
	out.WriteString(")")

	return out.String()
}

// Block Statement

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) statementNode()       {}
func (b *BlockStatement) TokenLiteral() string { return b.Token.Literal }
func (b *BlockStatement) String() string {
	var out bytes.Buffer
	for _, stmt := range b.Statements {
		out.WriteString(stmt.String())
	}
	return out.String()
}

// Declare Statement

type DeclareStatement struct {
	Token      token.Token
	Name       *Identifier
	Expression Expression
}

func (d *DeclareStatement) statementNode()       {}
func (d *DeclareStatement) TokenLiteral() string { return d.Token.Literal }
func (d *DeclareStatement) String() string {
	var out bytes.Buffer

	out.WriteString("var ")
	out.WriteString(d.Name.String())
	out.WriteString(" = ")

	if d.Expression != nil {
		out.WriteString(d.Expression.String())
	}

	return out.String()
}

// Return Statement

type ReturnStatement struct {
	Token      token.Token
	Expression Expression
}

func (r *ReturnStatement) statementNode()       {}
func (r *ReturnStatement) TokenLiteral() string { return r.Token.Literal }
func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral())
	out.WriteString(" ")

	if r.Expression != nil {
		out.WriteString(r.Expression.String())
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
