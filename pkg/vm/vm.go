// Package vm contains the stack virtual machine.
// A stack virtual machine is one in which memory is allocated in The Stake which by convention is where the callstack is implemented.
// Like physical computers that execute machinecode, virtual machines execute bytecode.
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

var TRUE = &object.Boolean{Value: true}
var FALSE = &object.Boolean{Value: false}

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

// lastPoppedStackElem gets the last element popped from the stack. Values are not zero out when they are popped from the stack, instead the stackpointer is decremented.
// Returns the object last popped from the stack.
func (vm *VM) lastPoppedStackElem() object.Object {
	// The stackpointer always points to the next free slot in memory.
	return vm.stack[vm.sp]
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

		case code.OpEQ:
			right := vm.pop()
			left := vm.pop()

			var val bool
			switch {
			case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
				val = left.(*object.Integer).Value == right.(*object.Integer).Value
			case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
				val = left.(*object.Boolean).Value == right.(*object.Boolean).Value
			default:
				return fmt.Errorf("type mismatch: %s == %s", left.Type(), right.Type())
			}

			err := vm.push(convertBooleanToObject(val))
			if err != nil {
				return err
			}

		case code.OpNEQ:
			right := vm.pop()
			left := vm.pop()

			var val bool
			switch {
			case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
				val = left.(*object.Integer).Value != right.(*object.Integer).Value
			case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
				val = left.(*object.Boolean).Value != right.(*object.Boolean).Value
			default:
				return fmt.Errorf("type mismatch: %s != %s", left.Type(), right.Type())
			}

			err := vm.push(convertBooleanToObject(val))
			if err != nil {
				return err
			}

		case code.OpGT:
			right := vm.pop()
			left := vm.pop()

			var val bool
			switch {
			case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
				val = left.(*object.Integer).Value > right.(*object.Integer).Value
			default:
				return fmt.Errorf("type mismatch: %s > %s", left.Type(), right.Type())
			}

			err := vm.push(convertBooleanToObject(val))
			if err != nil {
				return err
			}

		case code.OpSub:
			right := vm.pop()
			left := vm.pop()

			// BUG: Handle other object types
			err := vm.push(
				&object.Integer{
					Value: left.(*object.Integer).Value - right.(*object.Integer).Value,
				},
			)
			if err != nil {
				return err
			}

		case code.OpMul:
			right := vm.pop()
			left := vm.pop()

			// BUG: Handle other object types
			err := vm.push(
				&object.Integer{
					Value: left.(*object.Integer).Value * right.(*object.Integer).Value,
				},
			)
			if err != nil {
				return err
			}

		case code.OpDiv:
			right := vm.pop()
			left := vm.pop()

			// BUG: Handle other object types
			err := vm.push(
				&object.Integer{
					Value: left.(*object.Integer).Value / right.(*object.Integer).Value,
				},
			)
			if err != nil {
				return err
			}

		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()

			// BUG: Handle other object types
			err := vm.push(
				&object.Integer{
					Value: left.(*object.Integer).Value + right.(*object.Integer).Value,
				},
			)
			if err != nil {
				return err
			}

		case code.OpTrue:
			err := vm.push(TRUE)
			if err != nil {
				return err
			}

		case code.OpMinus:
			operand := vm.pop()

			if operand.Type() != object.INTEGER_OBJ {
				return fmt.Errorf("unsupported type for negation: %s", operand.Type())
			}

			err := vm.push(&object.Integer{Value: -operand.(*object.Integer).Value})
			if err != nil {
				return err
			}

		case code.OpBang:
			err := vm.executeBangOperator()
			if err != nil {
				return err
			}

		case code.OpFalse:
			err := vm.push(FALSE)
			if err != nil {
				return err
			}

		case code.OpPop:
			vm.pop()

		}
	}

	return nil
}

// TODO: Refactor stack into own struct

// pop removes the top object from the stack.
// Returns the pop object.
func (vm *VM) pop() object.Object {
	if vm.sp == 0 {
		return nil
	}

	// Does not actually zero out the value, instead we just change the stackpointer.
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

// convertBooleanToObject converts a boolean value into a *object.Boolean.
func convertBooleanToObject(val bool) *object.Boolean {
	if val {
		return TRUE
	}
	return FALSE
}

// executeBangOperator negates the last value pushed onto the stack.
// If the last value is not of type *object.Boolean, it will default to FALSE.
func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	switch operand {
	case TRUE:
		return vm.push(FALSE)
	case FALSE:
		return vm.push(TRUE)
	// BUG: Potential bug here as any object on the stack will have a falsely value prefixed with the bang operator
	default:
		return vm.push(FALSE)
	}
}
