package evaluator

// Boolean comparision is faster than Integer comparision because the former case is just doing a pointer comparision

// When finding a return value in a statement, we are not exiting but instead evaluating everything and only returning
// the returned value. A better approach would be to exit early (how would that work here?) or treeshake the AST after parsing

import (
	"fmt"

	"github.com/grantwforsythe/monkeylang/pkg/ast"
	"github.com/grantwforsythe/monkeylang/pkg/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}

		right := Eval(node.Right)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatement(node)

	case *ast.IfExpression:
		return evalIfExpression(node)

	case *ast.ReturnStatement:
		value := Eval(node.ReturnValue)
		if isError(value) {
			return value
		}

		return &object.ReturnValue{Value: value}

	case *ast.BooleanExpression:
		return evalBooleanExpression(node.Value)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}

	return nil
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt)

		switch result.(type) {
		case *object.ReturnValue:
			return result.(*object.ReturnValue).Value
		case *object.Error:
			return result
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

func evalBlockStatement(node *ast.BlockStatement) object.Object {
	var result object.Object

	for _, stmt := range node.Statements {
		result = Eval(stmt)

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

	// TODO: Figure out why we need to evaluate left before right
	// Because you read from left to right
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

func evalBooleanInfixExpresssion(operator string, left, right object.Object) object.Object {
	lValue := left.(*object.Boolean).Value
	rValue := right.(*object.Boolean).Value

	switch operator {
	case "==":
		return evalBooleanExpression(lValue == rValue)
	case "!=":
		return evalBooleanExpression(lValue != rValue)
	default:
		return NULL
	}
}

func evalIfExpression(node *ast.IfExpression) object.Object {
	condition := Eval(node.Condition)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(node.Consequence)
	} else if node.Alternative != nil {
		return Eval(node.Alternative)
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
