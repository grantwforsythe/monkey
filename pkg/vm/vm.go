// Package vm contains the stack virutal machine.
// A stack virutal machine is one in which memory is allocated in The Stake which by convention is where the callstack is implemented.
// Like physical computers that execute machinecode, virutal machines execute bytecode.
package vm

import (
	"fmt"

	"github.com/grantwforsythe/monkeylang/pkg/code"
	"github.com/grantwforsythe/monkeylang/pkg/compiler"
	"github.com/grantwforsythe/monkeylang/pkg/object"
)

// StackSize represents the maximum number of elements in the stack.
const StackSize = 2048 // This number was abritarily choosen

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	// Instructions
	stack []object.Object
	// sp represents a stackpointer which always points to the next free space in the stack.
	sp int
}

// New creates a new virtual machine from bytecode.
func New(bytecode *compiler.ByteCode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

// StackTop gets the top element on the stack.
// Returns nil if the stack is empty.
func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

// Run is the fetch-decode-excute cycle for the virtual machine.
func (vm *VM) Run() error {
	// The fetch part.
	for ip := 0; ip < len(vm.instructions); ip++ {

		// The decode part.
		op := code.Opcode(vm.instructions[ip])

		// The execute part.
		switch op {
		case code.OpConstant:
			index := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.constants[index])
			if err != nil {
				return err
			}

		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()

			// BUG: Handle other object types
			result := &object.Integer{
				Value: left.(*object.Integer).Value + right.(*object.Integer).Value,
			}

			err := vm.push(result)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// pop removes the top object from the stack.
// Returns the pop object.
func (vm *VM) pop() object.Object {
	if vm.sp == 0 {
		return nil
	}

	obj := vm.stack[vm.sp-1]
	vm.sp -= 1
	return obj
}

// push adds an object to the top of the stack and increments the pointer.
// Returns an error if a stackover flow occurs.
func (vm *VM) push(obj object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stackover flow")
	}

	vm.stack[vm.sp] = obj
	vm.sp++

	return nil
}
