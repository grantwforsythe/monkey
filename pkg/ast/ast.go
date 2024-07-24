package ast

import "github.com/grantwforsythe/monkeylang/pkg/token"

type Node interface {
	// The literal value of a token. This method will be used strictly for debugging and testing purposes
	TokenLiteral() string
}

type Statement interface {
	Node
	// Not strictly necessary but help the Go compiler
	statementNode()
}

type Expression interface {
	Node
	expresionNode()
}

// The root node of our AST
// Every valid Monkeylang program is a collection of statements
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

// The variable name
type Identifier struct {
	Token token.Token // The IDENT token
	Value string
}

func (i *Identifier) expressNode()         {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type LetStatement struct {
	Token token.Token // The LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
