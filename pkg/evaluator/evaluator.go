package evaluator

// Boolean comparision is faster than Integer comparision because the former case is just doing a pointer comparision

import (
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
		return evalStatement(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		right := Eval(node.Right)
		left := Eval(node.Left)
		return evalInfixExpression(right, left, node.Operator)

	case *ast.BooleanExpression:
		return evalBooleanExpression(node.Value)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}

	return nil
}

func evalStatement(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
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
		return NULL
	default:
		// TODO: Improve as this is prone to cause some grief
		return NULL
	}
}

func evalInfixExpression(right, left object.Object, operator string) object.Object {
	switch {
	case right.Type() == object.INTEGER_OBJ && left.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(right, left, operator)
	// If left and right are boolean objects, they are evaluated to either TRUE or FALSE which are constant pointer
	case operator == "==":
		return evalBooleanExpression(left == right)
	case operator == "!=":
		return evalBooleanExpression(left != right)
	// TODO: Implement
	default:
		return NULL
	}
}

func evalIntegerInfixExpression(right, left object.Object, operator string) object.Object {
	rValue := right.(*object.Integer).Value
	lValue := left.(*object.Integer).Value

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
		return NULL
	}
}

func evalBooleanInfixExpresssion(right, left object.Object, operator string) object.Object {
	rValue := right.(*object.Boolean).Value
	lValue := left.(*object.Boolean).Value

	switch operator {
	case "==":
		return evalBooleanExpression(lValue == rValue)
	case "!=":
		return evalBooleanExpression(lValue != rValue)
	// TODO: Implement
	default:
		return NULL
	}
}
