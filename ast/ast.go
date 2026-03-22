package ast

import (
	"bytes"
	"monkey/token"
)

// All ast node MUST implement Node interface (MUST have TokenLiteral() method)
type Node interface {
	TokenLiteral() string // Literal of the Token which this Node has (e.g. Identifier Node has IDENT `x` token literal)
	String() string       // debug
}

/*
type Node = Statement | Expression
type Program = Statement[]
type Statement = LetStatement | ...
*/

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// root node of ast
type Program struct {
	Statements []Statement
}

// impl Node
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Statemet
/*
Let Statement Node (`let x = 5`)
- token `=`
- field new identifier node
- right hand value expression
*/
type LetStatement struct {
	Token token.Token // token.LET token (`=`)
	Name  *Identifier // left-hand Identifier (`let x`, `let z`)
	Value Expression  // right-hand value experssion (`5 * 3`, `add(3, 4) + 1`)
}

// impl Statement
func (ls *LetStatement) statementNode() {}

// impl Node
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

/*
Identifier Node (`x`)
*/
type Identifier struct {
	Token token.Token // token.IDENT token
	Value string      // name (`x`)
}

// impl Experssion
func (i *Identifier) expressionNode() {}

// impl Node
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string {
	return i.Value
}

/*
return <expression>
*/
type ReturnStatement struct {
	Token       token.Token // RETURN token
	ReturnValue Expression
}

// implement Statement
func (rs *ReturnStatement) statementNode() {}

// implement Node
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ReturnStatement) String() string {
	// "return"+" "+"<ReturnValue>"+";"
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

/*
<Expression> + ";" (e.g. 5; add(3,4);
*/
type ExpressionStatement struct {
	Token      token.Token // first token of this expression
	Expression Expression
}

// impl Statement
func (es *ExpressionStatement) statementNode() {}

// impl Node
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) String() string {
	// TODO: remove nil-check later
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

// impl Expression
func (il *IntegerLiteral) expressionNode() {}

// impl Node
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token // prefix token e.g "!"
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
	// (!5), (!(-5))
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token // prefix token e.g "!"
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) String() string {
	// (!5), (!(-5))
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
