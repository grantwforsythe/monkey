package evaluator

import (
	"os"

	"fmt"

	"github.com/grantwforsythe/monkeylang/pkg/object"
)

var builtin = map[string]*object.Builtin{
	"len": {
		// Calculate the length of array or string.
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported. got=%s", args[0].Type())
			}
		},
	},
	"first": {
		// Get the first element of an array.
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			array, ok := args[0].(*object.Array)
			if !ok {
				return newError(
					"'first' only accepts an array as an argument. got=%s",
					args[0].Type(),
				)
			}

			if len(array.Elements) == 0 {
				return NULL
			}

			return array.Elements[0]
		},
	},
	"last": {
		// Get the last element of an array.
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			array, ok := args[0].(*object.Array)
			if !ok {
				return newError(
					"'last' only accepts an array as an argument. got=%s",
					args[0].Type(),
				)
			}

			length := len(array.Elements)
			if length == 0 {
				return NULL
			}

			return array.Elements[length-1]
		},
	},
	"rest": {
		// Return a copy of the array with the first element removed.
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			array, ok := args[0].(*object.Array)
			if !ok {
				return newError(
					"'rest' only accepts an array as an argument. got=%s",
					args[0].Type(),
				)
			}

			length := len(array.Elements)
			if length == 0 {
				return NULL
			}

			elements := make([]object.Object, length-1)

			if length == 1 {
				return &object.Array{Elements: elements}
			}

			copy(elements, array.Elements[1:])
			return &object.Array{Elements: elements}
		},
	},
	"push": {
		// Return a cloned array with a new value appended to it.
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return newError("wrong number of arguments. got=%d, want=>2", len(args))
			}

			array, ok := args[0].(*object.Array)
			if !ok {
				return newError(
					"the first argument needs to be of type ARRAY. got=%s",
					args[0].Type(),
				)
			}

			length := len(array.Elements)

			elements := make([]object.Object, length)
			copy(elements, array.Elements)
			elements = append(elements, args[1:]...)

			return &object.Array{Elements: elements}
		},
	},
	"quit": {
		Fn: func(args ...object.Object) object.Object {
			os.Exit(0)
			return &object.Null{}
		},
	},
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
