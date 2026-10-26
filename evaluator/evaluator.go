package evaluator

import (
	"MonkeyInterpreter/ast"
	"MonkeyInterpreter/object"
	"fmt"
	"reflect"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		ret := evalStatements(node.Statements)

		switch val := ret.(type) {
		case *object.FullReturnValue: //TODO: switch over to the new return system [target] <- [value]. <- can be called inside any block. Targets include prog, [function name], if, for, (no name to mean the next closest enclosing block), …
			ret = val.Value
		case *object.Error:
			ret = val
		}

		return ret
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		//right := Eval(node.Right)
		//return evalPrefixExpression(node.Operator, right) //TODO: just pass the node in and eval the right side in this function. That way node.Position() is always available

		return evalPrefixExpression(node)
	case *ast.InfixExpression:
		//left := Eval(node.Left)
		//right := Eval(node.Right)
		//return evalInfixExpression(node.Operator, left, right)

		return evalInfixExpression(node)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.DecimalLiteral:
		return &object.Decimal{Value: node.Value}
	case *ast.Boolean:
		if node.Value {
			return TRUE
		}
		return FALSE
	case *ast.BlockExpression:
		return evalStatements(node.Statements)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		return &object.ReturnValue{Value: val}
	case *ast.FullReturnStatement:
		val := Eval(node.ReturnValue)
		return &object.FullReturnValue{Value: val}
	default:
		fmt.Printf("Invalid node\n\tNode: %+v\n\t%s\n", node, reflect.TypeOf(node))
	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)

		switch val := result.(type) {
		case *object.ReturnValue:
			return val.Value
		case *object.FullReturnValue:
			return val
		case *object.Error:
			return val
		}
	}

	return result
}

func evalPrefixExpression(node *ast.PrefixExpression) object.Object {
	right := Eval(node.Right)

	if right == nil {
		return object.NewError(node.Right.Position(), "No value given for the expression %sx", node.Operator)
	}

	switch node.Operator {
	case "!":
		return evalBangOperatorExpression(node, right)
	case "-":
		return evalMinusOperatorExpression(node, right)
	default:
		return object.NewError(node.Position(), "unknown prefix operator %s %s", node.Operator, right.Type()) //for type !a
	}
}

func evalBangOperatorExpression(node *ast.PrefixExpression, right object.Object) object.Object {
	if right == nil {
		return object.NewError(node.Position(), "The value for this operation did not evaluate to a value: %s", node.Right.String())
	}

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
}

func evalMinusOperatorExpression(node *ast.PrefixExpression, right object.Object) object.Object {
	if right == nil {
		return object.NewError(node.Position(), "The value for this operation did not evaluate to a value: %s", node.Right.String())
	}

	switch r := right.(type) {
	case *object.Integer:
		result := &object.Integer{}
		result.Value = -r.Value
		return result
	case *object.Decimal:
		result := &object.Decimal{}
		result.Value = -r.Value
		return result
	default:
		return object.NewError(node.Position(), "invalid type %s for prefix operator %s", right.Type(), node.Operator)
	}
}

func evalInfixExpression(node *ast.InfixExpression) object.Object {
	left := Eval(node.Left)
	right := Eval(node.Right)

	if left == nil {
		return object.NewError(node.Position(), "The left side of this operation did not evaluate to a value: %s", node.Left.String())
	}

	if right == nil {
		return object.NewError(node.Position(), "The right side of this operation did not evaluate to a value: %s", node.Right.String())
	}

	operator := node.Operator

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(node, operator, left, right)
	case left.Type() == object.DECIMAL_OBJ && right.Type() == object.DECIMAL_OBJ:
		return evalDecimalInfixExpression(node, operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return object.NewError(node.Position(), "unknown infix operator %s %s %s", left.Type(), operator, right.Type()) //for type a + b
	}
}

func evalIntegerInfixExpression(node *ast.InfixExpression, operator string, left object.Object, right object.Object) object.Object {
	if left == nil {
		return object.NewError(node.Position(), "The left side of this operation did not evaluate to a value: %s", node.Left.String())
	}

	if right == nil {
		return object.NewError(node.Position(), "The right side of this operation did not evaluate to a value: %s", node.Right.String())
	}

	l, ok := left.(*object.Integer)
	if !ok {
		return object.NewError(node.Position(), "The left side of this operation is not an integer: %s", left.Inspect())
	}
	lv := l.Value

	r, ok := right.(*object.Integer)
	if !ok {
		return object.NewError(node.Position(), "The right side of this operation is not an integer: %s", right.Inspect())
	}
	rv := r.Value

	var res = &object.Integer{}
	switch operator {
	case "+":
		res.Value = lv + rv
	case "-":
		res.Value = lv - rv
	case "*":
		res.Value = lv * rv
	case "/":
		if rv == 0 {
			return object.NewError(node.Position(), "invalid division by zero")
		}
		res.Value = lv / rv
	case "%":
		if r.Value == 0 {
			return object.NewError(node.Position(), "invalid division by zero")
		}
		res.Value = lv % rv
	case "&":
		res.Value = lv & rv
	case "|":
		res.Value = lv | rv
	case "^":
		res.Value = lv ^ rv
	case "<<":
		res.Value = lv << rv
	case ">>":
		res.Value = lv >> rv
	case "<":
		return nativeBoolToBooleanObject(lv < rv)
	case ">":
		return nativeBoolToBooleanObject(lv > rv)
	case "==":
		return nativeBoolToBooleanObject(lv == rv)
	case "!=":
		return nativeBoolToBooleanObject(lv != rv)
	default:
		return object.NewError(node.Position(), "unknown infix operator %s %s %s", left.Type(), operator, right.Type()) //for type a + b
	}

	return res
}

func evalDecimalInfixExpression(node *ast.InfixExpression, operator string, left object.Object, right object.Object) object.Object {
	if left == nil {
		return object.NewError(node.Position(), "The left side of this operation did not evaluate to a value: %s", node.Left.String())
	}

	if right == nil {
		return object.NewError(node.Position(), "The right side of this operation did not evaluate to a value: %s", node.Right.String())
	}

	l, ok := left.(*object.Decimal)
	if !ok {
		return object.NewError(node.Position(), "The left side of this operation is not an integer: %s", left.Inspect())
	}
	lv := l.Value

	r, ok := right.(*object.Decimal)
	if !ok {
		return object.NewError(node.Position(), "The right side of this operation is not an integer: %s", right.Inspect())
	}
	rv := r.Value

	var res = &object.Decimal{}
	switch operator {
	case "+":
		res.Value = lv + rv
	case "-":
		res.Value = lv - rv
	case "*":
		res.Value = lv * rv
	case "/":
		if rv == 0 {
			return object.NewError(node.Position(), "invalid division by zero")
		}
		res.Value = lv / rv
	case "<":
		return nativeBoolToBooleanObject(lv < rv)
	case ">":
		return nativeBoolToBooleanObject(lv > rv)
	case "==":
		return nativeBoolToBooleanObject(lv == rv)
	case "!=":
		return nativeBoolToBooleanObject(lv != rv)
	default:
		return object.NewError(node.Position(), "unknown infix operator %s %s %s", left.Type(), operator, right.Type()) //for type a + b
	}

	return res
}

func nativeBoolToBooleanObject(b bool) object.Object {
	if b {
		return TRUE
	}
	return FALSE
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return object.NewError(ie.Condition.Position(), "Unknown boolean value for condition: %s", ie.Condition.String())
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
		return true
	}
}
