package evaluator

// Boolean comparision is faster than Integer comparision because the former case is just doing a pointer comparision

// When finding a return value in a statement, we are not exiting but instead evaluating everything and only returning
// the returned value. A better approach would be to exit early (how would that work here?) or treeshake the AST after parsing

import (
	"fmt"

	"github.com/grantwforsythe/monkeylang/pkg/ast"
	"github.com/grantwforsythe/monkeylang/pkg/object"
)

// NOTE: env can be refactored into the package scope so it does not need to be passed around

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.FunctionLiteral:
		return &object.Function{Body: node.Body, Env: env, Parameters: node.Parameters}

	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}

		// Evaluate the arguments
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(fn, args)

	case *ast.ReturnStatement:
		value := Eval(node.ReturnValue, env)
		if isError(value) {
			return value
		}

		return &object.ReturnValue{Value: value}

	case *ast.LetStatement:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}

		env.Set(node.Name.String(), value)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.BooleanExpression:
		return evalBooleanExpression(node.Value)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)

		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return &object.Array{Elements: elements}

	case *ast.IndexEpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		switch obj := result.(type) {
		case *object.ReturnValue:
			return obj.Value
		case *object.Error:
			return obj
		}
	}

	return result
}

func evalBooleanExpression(value bool) object.Object {
	if value {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		switch right {
		case TRUE:
			return FALSE
		case FALSE:
			return TRUE
		case NULL:
			return TRUE
		default:
			return FALSE
		}
	case "-":
		if right, ok := right.(*object.Integer); ok {
			return &object.Integer{Value: -1 * right.Value}
		}
		return newError("unknown operator: -%s", right.Type())
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBlockStatement(node *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range node.Statements {
		result = Eval(stmt, env)

		if result == nil {
			continue
		}

		if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
			return result
		}
	}

	return result
}

// TODO: Swap function signature so it is left than right
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return evalBooleanExpression(left == right)
	case operator == "!=":
		return evalBooleanExpression(left != right)
	// If left and right are boolean objects, they are evaluated to either TRUE or FALSE which are constant pointer
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	lValue := left.(*object.Integer).Value
	rValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: lValue + rValue}
	case "-":
		return &object.Integer{Value: lValue - rValue}
	case "*":
		return &object.Integer{Value: lValue * rValue}
	case "/":
		return &object.Integer{Value: lValue / rValue}
	case "<":
		return evalBooleanExpression(lValue < rValue)
	case ">":
		return evalBooleanExpression(lValue > rValue)
	case "==":
		return evalBooleanExpression(lValue == rValue)
	case "!=":
		return evalBooleanExpression(lValue != rValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	lValue := left.(*object.String).Value
	rValue := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: lValue + rValue}
	case "==":
		return evalBooleanExpression(lValue == rValue)
	case "!=":
		return evalBooleanExpression(lValue != rValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		// There are currently only 3 object types: Integer, Boolean, and Null
		// If obj is not of type Boolean or Null, then we know it has to be of type Integer
		// TODO: Handle the case for more object types
		return obj.(*object.Integer).Value > 0
	}
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if obj, ok := env.Get(node.Value); ok {
		return obj
	}

	if builtin, ok := builtin[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", node.Value)
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, expression := range expressions {
		eval := Eval(expression, env)
		if isError(eval) {
			return []object.Object{eval}
		}

		result = append(result, eval)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:

		// Assign the arguments to their corresponding parameter
		enclosedEnv := object.NewEnclosedEnvironment(fn.Env)
		for paramIdx, param := range fn.Parameters {
			enclosedEnv.Set(param.Value, args[paramIdx])
		}

		eval := Eval(fn.Body, enclosedEnv)

		if result, ok := eval.(*object.ReturnValue); ok {
			return result.Value
		}

		// TODO: Figure out what could be returned here
		return eval
	// TODO: Figure out why this works here
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		array := left.(*object.Array)
		idx := index.(*object.Integer).Value

		// Indexing out of bounds
		if idx < 0 || idx > int64(len(array.Elements)-1) {
			return NULL
		}

		return array.Elements[idx]
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}
