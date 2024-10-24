package compiler

import (
	"fmt"
	"testing"

	"github.com/grantwforsythe/monkeylang/pkg/ast"
	"github.com/grantwforsythe/monkeylang/pkg/code"
	"github.com/grantwforsythe/monkeylang/pkg/lexer"
	"github.com/grantwforsythe/monkeylang/pkg/object"
	"github.com/grantwforsythe/monkeylang/pkg/parser"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []any
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			"1 + 2",
			[]any{1, 2},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			"1 - 2",
			[]any{1, 2},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			},
		},
		{
			"1 * 2",
			[]any{1, 2},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			},
		},
		{
			"1 / 2",
			[]any{1, 2},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			},
		},
		{
			"50 / 2 * 2 + 10 - 5",
			[]any{50, 2, 2, 10, 5},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpMul),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			},
		},
		{
			"-(50 / 2 * 2 + 10 - 5)",
			[]any{50, 2, 2, 10, 5},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpMul),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpSub),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},
		{
			"-1",
			[]any{1},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},
		{
			"1; 2",
			[]any{1, 2},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestBooleanExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			"true; false",
			[]any{},
			[]code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
		{
			"false",
			[]any{},
			[]code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
		{
			"1 > 2",
			[]any{1, 2},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGT),
				code.Make(code.OpPop),
			},
		},
		{
			"1 < 2",
			[]any{2, 1},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGT),
				code.Make(code.OpPop),
			},
		},
		{
			"1 == 2",
			[]any{1, 2},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpEQ),
				code.Make(code.OpPop),
			},
		},
		{
			"1 != 2",
			[]any{1, 2},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNEQ),
				code.Make(code.OpPop),
			},
		},
		{
			"true == true",
			[]any{},
			[]code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpTrue),
				code.Make(code.OpEQ),
				code.Make(code.OpPop),
			},
		},
		{
			"false == true",
			[]any{},
			[]code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpTrue),
				code.Make(code.OpEQ),
				code.Make(code.OpPop),
			},
		},
		{
			"true != true",
			[]any{},
			[]code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpTrue),
				code.Make(code.OpNEQ),
				code.Make(code.OpPop),
			},
		},
		{
			"!true",
			[]any{},
			[]code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
		},
		{
			"!true != true",
			[]any{},
			[]code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpBang),
				code.Make(code.OpTrue),
				code.Make(code.OpNEQ),
				code.Make(code.OpPop),
			},
		},
		{
			"!(1 != 2)",
			[]any{1, 2},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNEQ),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			"if (true) { 10 }; 3333",
			[]any{10, 3333},
			[]code.Instructions{
				// index: 000
				code.Make(code.OpTrue),
				// index: 001
				code.Make(code.OpJumpNotTruthy, 7),
				// index: 004
				code.Make(code.OpConstant, 0),
				// index: 007
				code.Make(code.OpPop),
				// index: 008
				code.Make(code.OpConstant, 1),
				// index: 011
				code.Make(code.OpPop),
			},
		},
		{
			"if (false) { 10 }; 3333",
			[]any{10, 3333},
			[]code.Instructions{
				// index: 000
				code.Make(code.OpFalse),
				// index: 001
				code.Make(code.OpJumpNotTruthy, 7),
				// index: 004
				code.Make(code.OpConstant, 0),
				// index: 007
				code.Make(code.OpPop),
				// index: 008
				code.Make(code.OpConstant, 1),
				// index: 011
				code.Make(code.OpPop),
			},
		},
		{
			"if (true) { 10 } else { 20 }; 3333",
			[]any{10, 20, 3333},
			[]code.Instructions{
				// index: 000
				code.Make(code.OpTrue),
				// index: 001
				code.Make(code.OpJumpNotTruthy, 10),
				// index: 004
				code.Make(code.OpConstant, 0),
				// index: 007
				code.Make(code.OpJump, 13),
				// index: 010
				code.Make(code.OpConstant, 1),
				// index: 013
				code.Make(code.OpPop),
				// index: 014
				code.Make(code.OpConstant, 2),
				// index: 017
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, test := range tests {
		program := parse(test.input)

		compiler := New()
		err := compiler.Compile(program)

		if err != nil {
			t.Errorf("compiler error: %s", err)
		}

		bytecode := compiler.ByteCode()
		err = testInstructions(test.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Errorf("testInstructions failed: %s", err)
		}

		err = testConstants(test.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Errorf("testConstants failed: %s", err)
		}
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testInstructions(expected []code.Instructions, actual code.Instructions) error {
	concatted := concatInstructions(expected)

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instruction length. expected=%s, got=%s", concatted, actual)
	}

	for i, instruction := range concatted {
		if actual[i] != instruction {
			return fmt.Errorf(
				"wrong instruction at %d. expected=%q, got=%q",
				i,
				instruction,
				actual[i],
			)
		}
	}

	return nil
}

func concatInstructions(instructions []code.Instructions) code.Instructions {
	// Getting the length for the first instuction so that we can efficiently append elements to the slice.
	out := code.Instructions{}

	for _, instruction := range instructions {
		out = append(out, instruction...)
	}

	return out
}

func testConstants(expected []any, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf(
			"wrong number of constants. expected=%d, got=%d",
			len(expected),
			len(actual),
		)
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("failed to create consant in position %d: %s", i, err)
			}
		}
	}

	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not of type *object.Integer. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("mismatched values. expected=%d, got=%d", expected, result.Value)
	}

	return nil
}
