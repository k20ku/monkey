package ast

import "monkey/token"

// All ast node MUST implement Node interface (MUST have TokenLiteral() method)
type Node interface {
	TokenLiteral() string // Literal of the Token which this Node has (e.g. Identifier Node has IDENT `x` token literal)
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
	expresstionNode()
}

// root node of ast
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

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

// implement Statement
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

/*
Identifier Node (`x`)
*/
type Identifier struct {
	Token token.Token // token.IDENT token
	Value string      // name (`x`)
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
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// implement Experssion
func (i *Identifier) expresstionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
