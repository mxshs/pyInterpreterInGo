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
    case *ast.FloatLiteral:
        return &object.Float{Value: node.Value}
    case *ast.Boolean:
        if node.Value {
            return TRUE
        } else {
            return FALSE
        }
    case *ast.StringLiteral:
        return &object.String{Value: node.Value}
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
    case *ast.ListLiteral:
        elements := []object.Object{}
        for _, elem := range node.Arr {
            elements = append(elements, Eval(elem, env))
        }
        return &object.List{Arr: elements}
    case *ast.IndexExpression:
        Struct := Eval(node.Struct, env)
        if isError(Struct) {
            return Struct
        }

        idx := Eval(node.Value, env)
        if isError(idx) {
            return idx
        }
        
        return evalIndexExpression(Struct, idx)
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
        return newError("unknown operator %s for type %s", op, operand.Type())
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
    if !IsNumeric(operand){
        return newError("unknown operator - for type %s",
            operand.Type(),
        )
    }

    if operand.Type() == object.FLOAT_OBJ {
        return &object.Float{Value: -operand.(*object.Float).Value}
    }

    return &object.Integer{Value: -operand.(*object.Integer).Value}
}

func evalInfixExpression(
    op string, left, right object.Object) object.Object {
    switch {
    case IsNumeric(left) && IsNumeric(right):
        if left.Type() == object.FLOAT_OBJ {
            if right.Type() == object.FLOAT_OBJ {
                return evalFloatInfixExpression(op, left, right)
            }

            return evalFloatInfixExpression(
                op,
                left,
                &object.Float{
                    Value: float64(right.(*object.Integer).Value),
                },
            )
        } else {
            if right.Type() == object.INTEGER_OBJ {
                return evalIntegerInfixExpression(op, left, right)
            }

            return evalFloatInfixExpression(
                op,
                &object.Float{
                    Value: float64(left.(*object.Integer).Value),
                },
                right,
            )
        }
    case left.Type() == object.BOOL_OBJ || right.Type() == object.BOOL_OBJ:
        return evalBoolInfixExpression(op, left, right)
    case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
        return evalStringInfixExpression(op, left, right)
    case left.Type() != right.Type():
        return newError("type mismatch in %s %s %s",
            op,
            left.Type(),
            right.Type(),
        )
    default:
        return newError("unknown operator %s for types %s and %s", 
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
            if rightVal < 0 {
                return &object.Float{
                    Value: powFloat(
                        float64(leftVal),
                        float64(rightVal),
                    ),
                }
            }
            return &object.Integer{Value: pow(leftVal, rightVal)}
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
            return newError("unknown operator %s for types %s and %s",
                op,
                left.Type(),
                right.Type(),
            )
        }
}

func evalFloatInfixExpression(
    op string, left, right object.Object) object.Object {
    leftVal := left.(*object.Float).Value
    rightVal := right.(*object.Float).Value

    switch op {
        case "+":
            return &object.Float{Value: leftVal + rightVal}
        case "-": 
            return &object.Float{Value: leftVal - rightVal}
        case "*":
            return &object.Float{Value: leftVal * rightVal}
        case "/":
            return &object.Float{Value: leftVal / rightVal}
        case "**":
            return &object.Float{Value: powFloat(leftVal, rightVal)} 
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
            return newError("unknown operator %s for types %s and %s",
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

func evalStringInfixExpression(op string,
    left, right object.Object) object.Object {
    switch op {
    case "+":
        return &object.String{
            Value: left.(*object.String).Value + right.(*object.String).Value}
    default:
        return newError("unknown operator %s for types %s and %s",
            op,
            left.Type(),
            right.Type(),
        )
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
    if ok {
        return val
    }

    val, ok = bltins[name.Value]
    if ok {
        return val
    }

    return newError("name is not declared: %s", name.Value)
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

    switch function := function.(type) {
    case *object.Bltin:
        return function.Fn(args...)
    case *object.Function:
        fnEnv := object.NewNestedEnv(function.Env)

        for i, arg := range function.Arguments {
            fnEnv.Set(arg.Value, args[i])
        }

        evaluated := Eval(function.Body, fnEnv)

        return convertFunctionReturn(evaluated)    
    default:
        return newError("expected type Function, got: %s", function.Type())
    }
}

func convertFunctionReturn(res object.Object) object.Object {
    if val, ok := res.(*object.ReturnValue); ok {
        return val.Value
    }

    return res
}

func evalIndexExpression(Struct, index object.Object) object.Object {
    switch {
    case Struct.Type() == object.LIST && index.Type() == object.INTEGER_OBJ:
        return evalListIndexExpression(Struct, index)
    default:
        return newError(
            "attempting to apply unsupported index: %s[%s]",
            Struct.Type(),
            index.Type(),
        )
    }
}

func evalListIndexExpression(list, index object.Object) object.Object {
    arr := list.(*object.List).Arr
    idx := index.(*object.Integer).Value

    if idx < 0 || idx > int64(len(arr) - 1) {
        return NULL
    }

    return arr[idx]
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

