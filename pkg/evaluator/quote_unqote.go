package evaluator

import (
	"github.com/grantwforsythe/monkeylang/pkg/ast"
	"github.com/grantwforsythe/monkeylang/pkg/object"
)

// Create a Quote object evaluating any calls to unquote
func quote(node ast.Node, env *object.Environment) object.Object {
	return &object.Quote{Node: evalUnquoteCalls(node, env)}
}

func evalUnquoteCalls(quote ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(quote, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}

		// We make this same assertion in isUnquoteCall so no need to do it again
		call, _ := node.(*ast.CallExpression)

		// Can only handle one expression
		if len(call.Arguments) != 1 {
			return node
		}

		eval := Eval(call.Arguments[0], env)

		convertible, ok := eval.(object.Convertible)
		if !ok {
			return node
		}

		return convertible.ToNode()
	})
}

// Check if a node is a call expression for unquote.
func isUnquoteCall(node ast.Node) bool {
	fn, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}
	return fn.Function.TokenLiteral() == "unquote"
}
