// Package code contains all of the definitions for the bytecode, which are instructions for the virtual machine, format and stack virtual machine.
// A stack virutal machine is one in which memory is allocated in The Stake which by convention is where the callstack is implemented.
// Like physical computers that execute machinecode, virutal machines execute bytecode.
// NOTE: The definition for Bytecode is defined in the compiler's package to avoid an import-cycle error.
package code

// The fetch-decode-execute cycle, aka instruction cycle, is the clock speed for a CPU.
// A computer's memory is segmented into "words" - smallest unit of memory - which is usually 32/64 bits.

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

// Instructions represents instructions for the virtual machine.
// An instruction is a small, basic command that tells the simulated processor what to do.
// It is made up of an Opcode and an operator.
type Instructions []byte

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		definition, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		// +1 because the ith position is the opcode
		operands, offset := ReadOperands(definition, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(definition, operands))

		i += 1 + offset
	}

	return strings.TrimRight(out.String(), "\n")
}

// TODO: Rename method
func (ins Instructions) fmtInstruction(defintion *Definition, operands []int) string {
	if len(operands) != len(defintion.OperandWidths) {
		return fmt.Sprintf(
			"ERROR: operand len %d does not match defined %d\n",
			len(operands),
			len(defintion.OperandWidths),
		)
	}

	switch len(defintion.OperandWidths) {
	case 1:
		// No newline is needed because it is include in ins.String()
		return fmt.Sprintf("%s %d", defintion.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operand count for %s\n", defintion.Name)
}

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

// Make creates an instruction from a given opcode and associated operands.
func Make(op Opcode, operands ...int) Instructions {
	definition, ok := definitions[op]
	if !ok {
		// BUG: nil or empty slice?
		return []byte{}
	}

	// +1 because we need to account for the length of the instruction
	offset := 1
	instruction := make([]byte, definition.OperandWidths[0]+offset)
	instruction[0] = byte(op)

	// Iterate over the defined OperandWidths, take the matching element from the given operands and put it into the instruction.
	for i, operand := range operands {
		switch definition.OperandWidths[i] {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(operand))
		}
		offset += definition.OperandWidths[i]
	}

	return instruction
}

// ReadOperands is the opposite of Make - converts a definition and instruction to respective opcode and operands.
// Returns the operands for the instruction and the offset which represents the index of the last operand in the instruction.
func ReadOperands(definition *Definition, instruction Instructions) ([]int, int) {
	operands := make([]int, len(definition.OperandWidths))
	offset := 0

	for i, width := range definition.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(instruction[offset:]))
		}

		offset += width
	}

	return operands, offset
}

// TODO: Figure out what is going on
// ReadUint16 reads an instruction which is stored in memory using big endian.
func ReadUint16(instruction Instructions) uint16 {
	return binary.BigEndian.Uint16(instruction)
}
