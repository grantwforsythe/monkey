package ast

type ModifierFunc func(Node) Node

// Apply a modifer to an AST node.
func Modify(node Node, modifier ModifierFunc) Node {
	// TODO: Replace underscores with proper error handling

	switch node := node.(type) {
	case *Program:
		for i, statement := range node.Statements {
			node.Statements[i], _ = Modify(statement, modifier).(Statement)
		}

	case *ExpressionStatement:
		node.Expression, _ = Modify(node.Expression, modifier).(Expression)

	case *CallExpression:
		for i, argument := range node.Arguments {
			node.Arguments[i], _ = Modify(argument, modifier).(Expression)
		}

	case *InfixExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Right, _ = Modify(node.Right, modifier).(Expression)

	case *PrefixExpression:
		node.Right, _ = Modify(node.Right, modifier).(Expression)

	case *IndexEpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Index, _ = Modify(node.Index, modifier).(Expression)

	case *IfExpression:
		node.Condition, _ = Modify(node.Condition, modifier).(Expression)
		node.Consequence, _ = Modify(node.Consequence, modifier).(*BlockStatement)
		if node.Alternative != nil {
			node.Alternative, _ = Modify(node.Alternative, modifier).(*BlockStatement)
		}

	case *BlockStatement:
		for i, statement := range node.Statements {
			node.Statements[i], _ = Modify(statement, modifier).(Statement)
		}

	case *ReturnStatement:
		node.ReturnValue, _ = Modify(node.ReturnValue, modifier).(Expression)

	case *LetStatement:
		node.Value, _ = Modify(node.Value, modifier).(Expression)

	case *FunctionLiteral:
		if node.Parameters != nil {
			for i, parameter := range node.Parameters {
				node.Parameters[i], _ = Modify(parameter, modifier).(*Identifier)
			}
		}
		node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)

	case *ArrayLiteral:
		for i, element := range node.Elements {
			node.Elements[i], _ = Modify(element, modifier).(Expression)
		}

	case *HashLiteral:
		modifiedPairs := make(map[Expression]Expression)

		for key, value := range node.Pairs {
			modifiedKey, _ := Modify(key, modifier).(Expression)
			modifiedValue, _ := Modify(value, modifier).(Expression)

			modifiedPairs[modifiedKey] = modifiedValue
		}

		node.Pairs = modifiedPairs
	}

	return modifier(node)
}
