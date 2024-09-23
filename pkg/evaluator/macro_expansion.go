package evaluator

import (
	"github.com/grantwforsythe/monkeylang/pkg/ast"
	"github.com/grantwforsythe/monkeylang/pkg/object"
)

// TODO: Handle nested macro definitions

// DefineMacro defines all top-level macros in the global environment and then removes all macro nodes from the AST.
func DefineMacros(program *ast.Program, env *object.Environment) {
	// Indexes for all macro definitions
	indexes := []int{}

	// Define all macro in the environment and store the indexes for where they are located.
	for i, statement := range program.Statements {
		// TODO: Have a better assertion
		if isMacroDefinition(statement) {
			// These conditions are already met in isMacroDefinition so no need to check again
			letStatement, _ := statement.(*ast.LetStatement)
			macroLiteral, _ := letStatement.Value.(*ast.MacroLiteral)

			// Add the macro to the environment
			env.Set(letStatement.Name.String(), &object.Macro{
				Parameters: macroLiteral.Parameters,
				Body:       macroLiteral.Body,
				Env:        env,
			})

			// Store the index for the macro
			indexes = append(indexes, i)
		}
	}

	// Remove all macro nodes from the AST
	for _, idx := range indexes {
		program.Statements = append(program.Statements[:idx], program.Statements[idx+1:]...)
	}
}

// ExpandMacros evaluate calls to macros and replaces the original call expression with the result of the evaluation in the AST.
func ExpandMacros(program ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node {
		exp, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		macro, ok := isMacroCall(exp, env)
		if !ok {
			return node
		}

		args := quoteArgs(exp.Arguments)
		extendedEnv := extendMacroEnv(macro, args)

		eval := Eval(macro.Body, extendedEnv)

		quote, ok := eval.(*object.Quote)
		if !ok {
			panic(
				"only AST nodes can be returned from macros, i.e. only Quote objects can be returned",
			)
		}

		return quote.Node
	})
}

// isMacroDefinition returns true if a statement is an *ast.MacroLiteral definition, else false.
func isMacroDefinition(statement ast.Statement) bool {
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		return false
	}

	_, ok = letStatement.Value.(*ast.MacroLiteral)
	return ok
}

// isMacroCall returns true if a call expression is on a macro, else false.
func isMacroCall(exp *ast.CallExpression, env *object.Environment) (*object.Macro, bool) {
	identifier, ok := exp.Function.(*ast.Identifier)
	if !ok {
		return nil, false
	}

	obj, ok := env.Get(identifier.Value)
	if !ok {
		return nil, false
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		return nil, false
	}

	return macro, true
}

// quoteArgs converts all arguments to a macro into *object.Quote objects.
func quoteArgs(args []ast.Expression) []*object.Quote {
	var quoteArgs []*object.Quote

	for _, arg := range args {
		quoteArgs = append(quoteArgs, &object.Quote{Node: arg})
	}

	return quoteArgs
}

// extendMacroEnv adds all arguments to a macro's scope.
// Returns the extended environment for the macro.
func extendMacroEnv(macro *object.Macro, args []*object.Quote) *object.Environment {
	extendedEnv := object.NewEnclosedEnvironment(macro.Env)

	for i, param := range macro.Parameters {
		extendedEnv.Set(param.Value, args[i])
	}

	return extendedEnv
}
