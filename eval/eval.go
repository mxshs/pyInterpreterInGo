package eval

import (
	"fmt"
	"mxshs/pyinterpreter/ast"
	"mxshs/pyinterpreter/object"
)

var (
    TRUE = &object.Boolean{Value: true}
    FALSE = &object.Boolean{Value: false}
    NULL = &object.Null{}
)

func Eval(node ast.Node, env *object.Env) object.Object {
    switch node := node.(type) {
    case *ast.Program:
        return evalProgram(node.Statements, env)
    case *ast.ExpressionStatement:
        return Eval(node.Expression, env)
    case *ast.IntegerLiteral:
        return &object.Integer{Value: node.Value}
    case *ast.Boolean:
        if node.Value {
            return TRUE
        } else {
            return FALSE
        }
    case *ast.PrefixExpression:
        operand := Eval(node.Right, env)
        if isError(operand) {
            return operand
        } 
        return evalPrefixExpression(node.Operator, operand)
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
        return evalBlockStatement(node.Statements, env)
    case *ast.IfExpression:
        return evalIfExpression(node, env)
    case *ast.ReturnStatement:
        val := Eval(node.ReturnValue, env)
        if isError(val) {
            return val
        }

        return &object.ReturnValue{Value: val}
    case *ast.AssignStatement:
        val := Eval(node.Value, env)
        if isError(val) {
            return val
        }

        env.Set(node.Name.Value, val)
    case *ast.Name:
        return evalName(node, env) 
    case *ast.FunctionStatement:
        args := node.Arguments
        body := node.Body
        env.Set(node.Name.Value, &object.Function{Arguments: args, Env: env, Body: body})
    case *ast.CallExpression:
        // node.Function is literally a name (ident) of a func,
        // i.e. we address the current env to retrieve the actual function.
        // It means that anonymous functions are not supported yet.
        // Lambdas will be added later as a separate type, therefore no
        // further changes are expected in function call evaluation.
        function := Eval(node.Function, env)
        if isError(function) {
            return function
        }

        args := evalExpressions(node.Arguments, env)
        if len(args) == 1 && isError(args[0]) {
            return args[0]
        }
        
        return runFunction(function, args)
    }

    return NULL
}

func evalProgram(statements []ast.Statement, env *object.Env) object.Object {
    var res object.Object

    for _, statement := range statements {
        res = Eval(statement, env)

        switch res := res.(type) {
        case *object.ReturnValue:
            return res.Value
        case *object.Error:
            return res
        }
    }

    return res
}

func evalBlockStatement(
    statements []ast.Statement, env *object.Env) object.Object {
    
    var res object.Object

    for _, statement := range statements {
        res = Eval(statement, env)

        if res != nil {
            if (res.Type() == object.
                RETURN_VALUE || res.Type() == object.ERROR_OBJ) {
                return res        
            }
        }
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
        return newError("Unknown operator %s for type %s", op, operand.Type())
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
        return newError("Unknown operator - for type %s",
            operand.Type(),
        )
    }

    return &object.Integer{Value: -operand.(*object.Integer).Value}
}

func evalInfixExpression(
    op string, left, right object.Object) object.Object {
    switch {
    case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
        return evalIntegerInfixExpression(op, left, right)
    case left.Type() == object.BOOL_OBJ || right.Type() == object.BOOL_OBJ:
        return evalBoolInfixExpression(op, left, right)
    case left.Type() != right.Type():
        return newError("Type mismatch in %s %s %s",
            op,
            left.Type(),
            right.Type(),
        )
    default:
        return newError("Unknown operator %s for types %s and %s", 
            op,
            left.Type(),
            right.Type(),
        )
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
        case "<":
            if leftVal < rightVal {
                return TRUE
            } else {
                return FALSE
            }
        case ">":
            if leftVal > rightVal {
                return TRUE
            } else {
                return FALSE
            }
        case "<=":
            if leftVal <= rightVal {
                return TRUE
            } else {
                return FALSE
            }
        case ">=":
            if leftVal >= rightVal {
                return TRUE
            } else {
                return FALSE
            }
        case "==":
            if leftVal == rightVal {
                return TRUE
            } else {
                return FALSE
            }
        case "!=":
            if leftVal == rightVal {
                return FALSE
            } else {
                return TRUE
            }
        default:
            return newError("Unknown operator %s for types %s and %s",
                op,
                left.Type(),
                right.Type(),
            )
        }
}

func evalBoolInfixExpression(
    op string, left, right object.Object) object.Object {
        switch op {
        case "==":
            if left == right {
                return TRUE
            } else {
                return FALSE
            }
        case "!=":
            if left == right {
                return FALSE
            } else {
                return TRUE
            }
        default:
            return NULL
        }
}

func evalIfExpression(ie *ast.IfExpression, env *object.Env) object.Object {
    condition := Eval(ie.Condition, env)
    if isError(condition) {
        return condition
    }

    if checkCondition(condition) {
        return Eval(ie.Consequence, env)
    } else if ie.Alternative != nil {
        return Eval(ie.Alternative, env)
    } else {
        return NULL
    }
}

func evalName(name *ast.Name, env *object.Env) object.Object {
    val, ok := env.Get(name.Value)
    if !ok {
        return newError("Name is not declared: %s", name.Value)
    }

    return val
}

func evalExpressions(
    expressions []ast.Expression, env *object.Env) []object.Object {

    var res []object.Object

    for _, expr := range expressions {
        evaluated := Eval(expr, env)
        if isError(evaluated) {
            return []object.Object{evaluated}
        }

        res = append(res, evaluated)
    }

    return res
}

func runFunction(
    function object.Object, args []object.Object) object.Object {

    fn, ok := function.(*object.Function)
    if !ok {
        return newError("expected type Function, got: %s", fn.Type())
    }

    fnEnv := object.NewNestedEnv(fn.Env)
    
    for i, arg := range fn.Arguments {
        fnEnv.Set(arg.Value, args[i])
    }

    evaluated := Eval(fn.Body, fnEnv)

    return evaluated.(*object.ReturnValue).Value
}

func checkCondition(obj object.Object) bool {
    switch obj {
        case TRUE:
            return true
        case FALSE:
            return false
        case NULL:
            return false
        default:
            return true
    }
}

func newError(fmtString string, args ...interface{}) *object.Error {
    return &object.Error{Message: fmt.Sprintf(fmtString, args...)}
}

func isError(obj object.Object) bool {
    if obj != nil {
        if obj.Type() == object.ERROR_OBJ {
            return true
        }
    }

    return false
}

