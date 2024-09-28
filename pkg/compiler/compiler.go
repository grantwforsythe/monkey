package compiler

import (
	"go/ast"

	"github.com/grantwforsythe/monkeylang/pkg/code"
	"github.com/grantwforsythe/monkeylang/pkg/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

type ByteCode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

// New initializes a new compiler.
func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	return nil
}

func (c *Compiler) ByteCode() *ByteCode {
	return &ByteCode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}
