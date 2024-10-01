// Package code contains all of the definitions for the bytecode, which are instructions for the virtual machine, format and stack virtual machine.
// A stack virutal machine is one in which memory is allocated in The Stake which by convention is where the callstack is implemented.
// Like physical computers that execute machinecode, virutal machines execute bytecode.
// NOTE: The definition for Bytecode is defined in the compiler's package to avoid an import-cycle error.
package code

// The fetch-decode-execute cycle, aka instruction cycle, is the clock speed for a CPU.
// A computer's memory is segmented into "words" - smallest unit of memory - which is usually 32/64 bits.

import (
	"encoding/binary"
	"fmt"
)

// Instructions represents instructions for the virtual machine.
// An instruction is a small, basic command that tells the simulated processor what to do.
// It is made up of an Opcode and an operator.
type Instructions []byte

// Opcode represents the "operator" part of an instruction.
type Opcode byte

// We let iota generate the byte values because the actual values do not matter.
const (
	OpConstant Opcode = iota // OpConstant retrives the constant using the operand as an index and pushes it onto the stack.
)

// Definition represents the definition for an Opcode.
type Definition struct {
	Name          string // Name represents the name of the Opcode.
	OperandWidths []int  // OperandWidths represents the number of bytes each operand, the argument / parameter to an operator, uses.
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
}

// Lookup is a function used for debugging that gets the Opcode definition for a given byte.
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d is not defined", op)
	}

	return def, nil
}

// Make converts an instruction from a given opcode and associated operands.
func Make(op Opcode, operands []int) Instructions {
	definition, ok := definitions[op]
	if !ok {
		// BUG: nil or empty slice?
		return []byte{}
	}

	// TODO: Refactor this
	instruction := make([]byte, definition.OperandWidths[0]+1)
	instruction[0] = byte(op)

	// Iterate over the defined OperandWidths, take the matching element from the given operands and put it into the instruction.
	offset := 1
	for i, operand := range operands {
		switch definition.OperandWidths[i] {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(operand))
		}
		offset += definition.OperandWidths[i]
	}

	return instruction
}
