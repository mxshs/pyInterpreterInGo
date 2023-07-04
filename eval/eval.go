package eval

import (
    "mxshs/pyinterpreter/ast"
    "mxshs/pyinterpreter/object"
)

var (
    TRUE = &object.Boolean{Value: true}
    FALSE = &object.Boolean{Value: false}
    NULL = &object.Null{}
)

func Eval(node ast.Node) object.Object {
    switch node := node.(type) {
    case *ast.Program:
        return evalStatements(node.Statements)
    case *ast.ExpressionStatement:
        return Eval(node.Expression)
    case *ast.IntegerLiteral:
        return &object.Integer{Value: node.Value}
    case *ast.Boolean:
        if node.Value {
            return TRUE
        } else {
            return FALSE
        }
    case *ast.PrefixExpression:
        operand := Eval(node.Right)
        return evalPrefixExpression(node.Operator, operand)
    case *ast.InfixExpression:
        left := Eval(node.Left)
        right := Eval(node.Right)
        return evalInfixExpression(node.Operator, left, right)
    }

    return NULL
}

func evalStatements(statements []ast.Statement) object.Object {
    var res object.Object

    for _, statement := range statements {
        res = Eval(statement)
    }

    return res
}

func evalPrefixExpression(op string, operand object.Object) object.Object {
    switch op {
    case "!":
        return evalBangOperatorExpression(operand)
    case "-":
        return evalMinusPrefixOperatorExpression(operand)
    default:
        return NULL
    }
}

func evalBangOperatorExpression(operand object.Object) object.Object {
    switch operand {
    case TRUE:
        return FALSE
    case FALSE:
        return TRUE
    case NULL:
        return FALSE
    default:
        return FALSE
    }
}

func evalMinusPrefixOperatorExpression(operand object.Object) object.Object {
    if operand.Type() != object.INTEGER_OBJ {
        return NULL
    }

    return &object.Integer{Value: -operand.(*object.Integer).Value}
}

func evalInfixExpression(
    op string, left, right object.Object) object.Object {
    switch {
    case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
        return evalIntegerInfixExpression(op, left, right)
    default:
        return NULL
    }
}

func evalIntegerInfixExpression(
    op string, left, right object.Object) object.Object {
    leftVal := left.(*object.Integer).Value
    rightVal := right.(*object.Integer).Value

    switch op {
        case "+":
            return &object.Integer{Value: leftVal + rightVal}
        case "-": 
            return &object.Integer{Value: leftVal - rightVal}
        case "*":
            return &object.Integer{Value: leftVal * rightVal}
        case "/":
            return &object.Integer{Value: leftVal / rightVal}
        case "**":
            return &object.Integer{Value: pow(leftVal, rightVal).(int64)}
        default:
            return NULL
        }
}

