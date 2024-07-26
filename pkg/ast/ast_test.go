package ast

import (
	"testing"

	"github.com/grantwforsythe/monkeylang/pkg/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name:  &Identifier{Token: token.Token{Type: token.IDENT, Literal: "foo"}, Value: "foo"},
				Value: &Identifier{Token: token.Token{Type: token.IDENT, Literal: "5"}, Value: "5"},
			},
		},
	}

	if program.String() != "let foo = 5;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
