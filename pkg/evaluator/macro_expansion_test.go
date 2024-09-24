package evaluator

import (
	"testing"

	"github.com/grantwforsythe/monkeylang/pkg/ast"
	"github.com/grantwforsythe/monkeylang/pkg/lexer"
	"github.com/grantwforsythe/monkeylang/pkg/object"
	"github.com/grantwforsythe/monkeylang/pkg/parser"
)

func TestDefineMacros(t *testing.T) {
	input := `
	let number = 1;
	let function = fn(x, y) { x + y };
	let mymacro = macro(x, y) { x + y };
	`

	program := testParseProgram(input)
	env := object.NewEnvironment()
	DefineMacros(program, env)

	if len(program.Statements) != 2 {
		t.Fatalf("the number of statements if not equal to 2. got=%d", len(program.Statements))
	}

	if _, ok := env.Get("number"); ok {
		t.Fatalf("'number' should be undefined")
	}

	if _, ok := env.Get("function"); ok {
		t.Fatalf("'function' should be undefined")
	}

	obj, ok := env.Get("mymacro")
	if !ok {
		t.Fatalf("'mymacro' should not be undefined")
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		t.Fatalf("obj is not of type *object.Macro. got=%T", macro)
	}

	if len(macro.Parameters) != 2 {
		t.Fatalf("macro is expected to have 2 parameters. got=%d", len(macro.Parameters))
	}

	if macro.Parameters[0].String() != "x" {
		t.Fatalf(
			"the first parameter of the macro is not equal to 'x'. got=%s",
			macro.Parameters[0].String(),
		)
	}

	if macro.Parameters[1].String() != "y" {
		t.Fatalf(
			"the first parameter of the macro is not equal to 'y'. got=%s",
			macro.Parameters[1].String(),
		)
	}

	if macro.Body.String() != "(x + y)" {
		t.Fatalf("the body of the macro is not equal to 'x + y'. got=%s", macro.Body.String())
	}
}

func TestExpandMacros(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`
			let infixExpression = macro() { quote(1 + 2); };
			infixExpression();
			`,
			`(1 + 2)`,
		},
		// TODO: Check why unquote doesn't evaluate the parameters
		{
			`
			let reverse = macro(a, b) { quote(unquote(b) - unquote(a)); };
			reverse(2 + 2, 10 - 5);
			`,
			`(10 - 5) - (2 + 2)`,
		},
	}

	for _, test := range tests {
		expected := testParseProgram(test.expected)
		program := testParseProgram(test.input)

		env := object.NewEnvironment()
		DefineMacros(program, env)
		expanded := ExpandMacros(program, env)

		if expanded.String() != expected.String() {
			t.Errorf("not equal. want=%q, got=%q",
				expected.String(), expanded.String())
		}
	}
}

func testParseProgram(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
