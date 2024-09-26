// Package code contains all of the definitions for the bytecode format.
// The definition for Bytecode is defined in the compiler's package because of an import-cycle.
package code

import "fmt"

// Instruction represents instructions for the virtual machine.
type Instruction []byte

// Opcode represents the "operator" part of an instruction.
type Opcode byte

const (
	OpConstant Opcode = iota // OpConstant retrives the constant using the operand as an index and pushes it onto the stack.
)

// Definition represents the definition for an Opcode.
type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
}

// Lookup gets the Opcode definition for a given byte.
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d is not defined", op)
	}

	return def, nil
}
