package ast

import (
	"testing"

	"github.com/grantwforsythe/monkeylang/pkg/token"
)

func TestProgram(t *testing.T) {
	tests := []struct {
		program  *Program
		expected string
	}{
		{program: &Program{
			Statements: []Statement{
				&LetStatement{
					Token: token.Token{Type: token.LET, Literal: "let"},
					Name: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "foo"},
						Value: "foo",
					},
					Value: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "5"},
						Value: "5",
					},
				},
			},
		}, expected: "let"},
		{program: &Program{}, expected: ""},
	}

	for _, tt := range tests {
		if tt.program.TokenLiteral() != tt.expected {
			t.Errorf(
				"program.TokenLiteral() is not equal to '%s'. got=%s",
				tt.expected,
				tt.program.TokenLiteral(),
			)
		}

	}

}

func TestIdentifier(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "foo"},
					Value: "foo",
				},
				Value: &Identifier{Token: token.Token{Type: token.IDENT, Literal: "5"}, Value: "5"},
			},
		},
	}

	stmt := program.Statements[0].(*LetStatement)
	if stmt.Name.TokenLiteral() != "foo" {
		t.Errorf(
			"stmt.Name is not equal to 'foo'. got=%s",
			stmt.Name.TokenLiteral(),
		)
	}

	if stmt.Name.String() != "foo" {
		t.Errorf("stmt.Name.String() wrong. got=%s", stmt.Name.String())
	}
}

func TestLetStatement(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "foo"},
					Value: "foo",
				},
				Value: &Identifier{Token: token.Token{Type: token.IDENT, Literal: "5"}, Value: "5"},
			},
		},
	}

	if program.Statements[0].TokenLiteral() != "let" {
		t.Errorf(
			"program.Statements[0].TokenLiteral() is not equal to 'let'. got=%s",
			program.Statements[0].TokenLiteral(),
		)
	}

	if program.String() != "let foo = 5;" {
		t.Errorf("program.String() wrong. got=%s", program.String())
	}
}

func TestReturnStatment(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ReturnStatement{
				Token: token.Token{Type: token.RETURN, Literal: "return"},
				ReturnValue: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "foo"},
					Value: "foo",
				},
			},
		},
	}

	stmt := program.Statements[0].(*ReturnStatement)
	if stmt.TokenLiteral() != "return" {
		t.Errorf(
			"stmt is not equal to 'return'. got=%s",
			program.Statements[0].TokenLiteral(),
		)
	}

	if stmt.String() != "return foo;" {
		t.Errorf("stmt.String() wrong. got=%s", stmt.String())
	}
}

func TestExpressionStatement(t *testing.T) {
	tests := []struct {
		program  *Program
		expected string
	}{
		{program: &Program{
			Statements: []Statement{
				&ExpressionStatement{
					Token: token.Token{Type: token.IDENT, Literal: "5"},
					Expression: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "5"},
						Value: "5",
					},
				},
			},
		}, expected: "5"},
		{program: &Program{
			Statements: []Statement{
				&ExpressionStatement{
					Token:      token.Token{Type: token.IDENT, Literal: "5"},
					Expression: nil,
				},
			},
		}, expected: ""},
	}

	for _, tt := range tests {
		if tt.program.String() != tt.expected {
			t.Errorf("program.String() wrong. got=%s", tt.program.String())
		}

		if tt.program.String() != tt.expected {
			t.Errorf("program.String() wrong. got=%s", tt.program.String())
		}
	}

}

func TestIntegerLiteralExpression(t *testing.T) {
	exp := &IntegerLiteral{
		Token: token.Token{Type: token.INT, Literal: "5"},
		Value: 5,
	}

	if exp.TokenLiteral() != "5" {
		t.Errorf("exp.TokenLiteral() wrong. got=%s", exp.TokenLiteral())
	}

	if exp.String() != "5" {
		t.Errorf("exp.String() wrong. got=%s", exp.String())
	}
}
