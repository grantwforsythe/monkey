package evaluator

import (
	"github.com/grantwforsythe/monkeylang/pkg/object"
)

var builtin = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			str, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `len` not supported. got INTEGER")
			}

			return &object.Integer{Value: int64(len(str.Value))}
		},
	},
}
