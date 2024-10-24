// Package compiler contains all of the logic for compiling an AST into bytecode.
// The Compiler is a stacked-based, single-pass compiler as objects are stored on a stacked and the AST is traversed only once.
package compiler

import (
	"fmt"

	"github.com/grantwforsythe/monkeylang/pkg/ast"
	"github.com/grantwforsythe/monkeylang/pkg/code"
	"github.com/grantwforsythe/monkeylang/pkg/object"
)

const DUMMY_OPERAND = 9999

// TODO: Prefix all errors with "Compiler error"

// EmittedInstruction represents an instruction emitted by the compiler.
type EmittedInstruction struct {
	Opcode   code.Opcode // Opcode represents the op of the emitted instruction.
	Position int         // Position represents the position of the emitted instruction in the slice of already emitted instructions.
}

type Compiler struct {
	instructions        code.Instructions   // instructions represent the instructions generated from source code.
	constants           []object.Object     // constants represents the constants pool the compiler.
	lastInstruction     *EmittedInstruction // lastInstruction represents the last emitted instruction by the compiler.
	previousInstruction *EmittedInstruction // previousInstruction represents the instruction emitted before the last instruction by the compiler.
}

// ByteCode represents a domain-specific language for a domain-specific virtual machine.
// It is called bytecode because the opcode in each instruction is one byte long.
type ByteCode struct {
	Instructions code.Instructions // Instructions represent the instructions generated by the compiler.
	Constants    []object.Object   // Constants represent the constants generated by the compiler.
}

// New initializes a new compiler.
func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    make([]object.Object, 0),
	}
}

// Compile traverses the nodes in the AST, converting it into bytecode.
func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, stmt := range node.Statements {
			err := c.Compile(stmt)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		// Expression statements emit a value but don't store it like an assignment statement
		c.emit(code.OpPop)

	case *ast.InfixExpression:
		// We are using one op for both greater than and less than, all that changes is the order in which values are emitted
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}

			c.emit(code.OpGT)
			return nil
		} else {
			err := c.Compile(node.Left)
			if err != nil {
				return err
			}

			err = c.Compile(node.Right)
			if err != nil {
				return err
			}
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case "==":
			c.emit(code.OpEQ)
		case "!=":
			c.emit(code.OpNEQ)
		case ">":
			c.emit(code.OpGT)
		default:
			return fmt.Errorf("unknown operator: %s", node.Operator)
		}

	case *ast.BooleanExpression:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		// The index of the newly added constant is used as an operand in the emitted instruction.
		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		jumpNotTruthyPosition := c.emit(code.OpJumpNotTruthy, DUMMY_OPERAND)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if c.lastInstructionIsPop() {
			c.removeLastInstruction()
		}

		// If there is no else statement, we just want to jump to the end of the consequence
		if node.Alternative == nil {
			// Replace the dummy operand with the position of the first instruction after the consequence of an if statement.
			c.updateOperand(jumpNotTruthyPosition, len(c.instructions))
		} else {
			// The jump statement is to jump over the else block
			// If the VM encounters the instruction in this context, it means the conditional evaluated to true
			jumpPosition := c.emit(code.OpJump, DUMMY_OPERAND)

			// Replace the dummy operand with the position after the jump operator which is the else block
			c.updateOperand(jumpNotTruthyPosition, len(c.instructions))

			err = c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIsPop() {
				c.removeLastInstruction()
			}

			// Replace the dummy operand with the position of the first instruction after the else block
			c.updateOperand(jumpPosition, len(c.instructions))
		}

	case *ast.BlockStatement:
		for _, stmt := range node.Statements {
			err := c.Compile(stmt)
			if err != nil {
				return err
			}
		}

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	}

	// Iterate over the instructions in memory, repeating the fetch-decode-execute cycle like in an actual machine.
	return nil
}

func (c *Compiler) ByteCode() *ByteCode {
	return &ByteCode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

// addConstant adds a constant to the constants pool.
// Returns the index of the newly added constant.
func (c *Compiler) addConstant(obj object.Object) int {
	// PERF: Unperformant way to add elements to a slice because the cap is 0 by default and will always be x2 the len by default
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// emit generates an instruction and add it to the results.
// Returns the position of the newly added instruction.
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	instruction := code.Make(op, operands...)
	// Starting position of the newly added instruction.
	position := len(c.instructions)
	// PERF: Unperformant way to add elements to a slice because the cap is 0 by default and will always be x2 the len by default
	c.instructions = append(c.instructions, instruction...)
	c.setLastInstruction(op, position)

	return position
}

func (c *Compiler) setLastInstruction(op code.Opcode, position int) {
	c.previousInstruction = c.lastInstruction
	c.lastInstruction = &EmittedInstruction{Opcode: op, Position: position}
}

// lastInstructionIsPop checks if the last emitted instruction was code.OpPop.
func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == code.OpPop
}

// removeLastInstruction removes the last emitted instruction from the slice of emitted instructions.
func (c *Compiler) removeLastInstruction() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
	// TODO: Set previous instruction to the previous previous instruction
	c.previousInstruction = nil
}

// updateOperand changes the operand for an instruction at a given position.
// It is assumed that the operand is valid for the instruction being updated.
func (c *Compiler) updateOperand(position, operand int) {
	op := code.Opcode(c.instructions[position])
	newInstruction := code.Make(op, operand)
	c.replaceInstruction(position, newInstruction)
}

// replaceInstruction replaces an instruction with a new instruction starting from a given posiution.
// It is assumed that the newInstruction is of the same type and length as the old instruction.
func (c *Compiler) replaceInstruction(position int, newInstruction code.Instructions) {
	for i := 0; i < len(newInstruction); i++ {
		c.instructions[position+i] = newInstruction[i]
	}
}
