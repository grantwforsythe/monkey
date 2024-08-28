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
		exp            *ExpressionStatement
		expectedToken  string
		expectedString string
	}{
		{exp: &ExpressionStatement{
			Token: token.Token{Type: token.IDENT, Literal: "5"},
			Expression: &Identifier{
				Token: token.Token{Type: token.IDENT, Literal: "5"},
				Value: "5",
			},
		},
			expectedToken: "5", expectedString: "5"},
		{exp: &ExpressionStatement{
			Token:      token.Token{Type: token.IDENT, Literal: "5"},
			Expression: nil,
		},
			expectedToken: "5", expectedString: ""},
	}

	for _, tt := range tests {
		if tt.exp.TokenLiteral() != tt.expectedToken {
			t.Errorf(
				"exp.TokenLiteral() is not equal to '%s'. got=%s",
				tt.expectedToken,
				tt.exp.TokenLiteral(),
			)
		}

		if tt.exp.String() != tt.expectedString {
			t.Errorf(
				"exp.String() is not equal to '%s'. got=%s",
				tt.expectedString,
				tt.exp.String(),
			)
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

func TestPrefixExpression(t *testing.T) {
	exp := &PrefixExpression{
		Token:    token.Token{Type: token.BANG, Literal: "!"},
		Operator: "!",
		Right: &Identifier{
			Token: token.Token{Type: token.IDENT, Literal: "foo"},
			Value: "foo",
		},
	}

	if exp.TokenLiteral() != "!" {
		t.Errorf("exp.TokenLiteral() is not equal to 'foo'. got=%s", exp.TokenLiteral())
	}

	if exp.String() != "(!foo)" {
		t.Errorf("exp.String() is not equal to '(!foo)'. got=%s", exp.String())
	}
}

func TestInfixExpression(t *testing.T) {
	exp := &InfixExpression{
		Token: token.Token{Type: token.PLUS, Literal: "+"},
		Left: &IntegerLiteral{
			Token: token.Token{Type: token.INT, Literal: "5"},
			Value: 5,
		},
		Operator: "+",
		Right: &IntegerLiteral{
			Token: token.Token{Type: token.INT, Literal: "5"},
			Value: 5,
		},
	}

	if exp.TokenLiteral() != "+" {
		t.Errorf("exp.TokenLiteral() is not equal to '+'. got=%s", exp.TokenLiteral())
	}

	if exp.String() != "(5 + 5)" {
		t.Errorf("exp.String() is not equal to '(!foo)'. got=%s", exp.String())
	}
}

func TestBooleanExpression(t *testing.T) {
	exp := &BooleanExpression{
		Token: token.Token{Type: token.TRUE, Literal: "true"},
		Value: true,
	}

	if exp.TokenLiteral() != "true" {
		t.Errorf("exp.TokenLiteral() is not equal to 'true'. got=%s", exp.TokenLiteral())
	}

	if exp.String() != "true" {
		t.Errorf("exp.String() is not equal to 'true'. got=%s", exp.String())
	}
}

func TestBlockStatement(t *testing.T) {
	stmt := &BlockStatement{
		Token: token.Token{Type: token.LPAREN, Literal: "("},
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
	}

	if stmt.TokenLiteral() != "(" {
		t.Errorf("stmt.TokenLiteral() is not equal to '('. got=%s", stmt.TokenLiteral())
	}

	if stmt.String() != "let foo = 5;" {
		t.Errorf("stmt.String() is not equal to 'let foo = 5;'. got=%s", stmt.String())
	}
}

func TestIfExpression(t *testing.T) {
	tests := []struct {
		exp            *IfExpression
		expectedToken  string
		expectedString string
	}{
		{exp: &IfExpression{
			Token: token.Token{Type: token.IF, Literal: "if"},
			Condition: &InfixExpression{
				Token: token.Token{Type: token.PLUS, Literal: "<"},
				Left: &IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Value: 5,
				},
				Operator: "<",
				Right: &IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "10"},
					Value: 10,
				},
			},
			Consequence: &BlockStatement{
				Token: token.Token{Type: token.LPAREN, Literal: "("},
				Statements: []Statement{
					&ExpressionStatement{
						Token: token.Token{Type: token.IDENT, Literal: "5"},
						Expression: &Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "5"},
							Value: "5",
						},
					},
				},
			},
		},
			expectedToken: "if", expectedString: "if(5 < 10) 5"},
		{exp: &IfExpression{
			Token: token.Token{Type: token.IF, Literal: "if"},
			Condition: &InfixExpression{
				Token: token.Token{Type: token.PLUS, Literal: "<"},
				Left: &IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Value: 5,
				},
				Operator: "<",
				Right: &IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "10"},
					Value: 10,
				},
			},
			Consequence: &BlockStatement{
				Token: token.Token{Type: token.LPAREN, Literal: "("},
				Statements: []Statement{
					&ExpressionStatement{
						Token: token.Token{Type: token.IDENT, Literal: "5"},
						Expression: &Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "5"},
							Value: "5",
						},
					},
				},
			},
			Alternative: &BlockStatement{
				Token: token.Token{Type: token.LPAREN, Literal: "("},
				Statements: []Statement{
					&ExpressionStatement{
						Token: token.Token{Type: token.IDENT, Literal: "10"},
						Expression: &Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "10"},
							Value: "10",
						},
					},
				},
			},
		},
			expectedToken: "if", expectedString: "if(5 < 10) 5 else 10"},
	}

	for _, tt := range tests {
		if tt.exp.TokenLiteral() != "if" {
			t.Errorf("tt.exp.TokenLiteral() is not equal to 'if'. got=%s", tt.exp.TokenLiteral())
		}

		if tt.exp.String() != tt.expectedString {
			t.Errorf(
				"tt.exp.String() is not equal to '%s'. got=%s",
				tt.expectedString,
				tt.exp.String(),
			)
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	exp := &StringLiteral{
		Token: token.Token{Type: token.STRING, Literal: "foobar"},
		Value: "foobar",
	}

	if exp.TokenLiteral() != "foobar" {
		t.Errorf("exp.TokenLiteral() wrong. got=%s", exp.TokenLiteral())
	}

	if exp.String() != "foobar" {
		t.Errorf("exp.String() wrong. got=%s", exp.String())
	}
}
