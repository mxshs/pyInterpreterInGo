package eval

import (
	"mxshs/pyinterpreter/object"
)

type Num interface {
    ~float64 | ~int64
}

var builtin = map[string]*object.BltinFn {
    "len": &object.BltinFn{
        Fn: length,
    },
}

var aliases = map[string]any {
    "pow": pow,
}

func length(args ...object.Object) object.Object {
    switch {
    case len(args) > 1:
        return newError(
            "expected 1 sequence-like argument, got %d arguments",
            len(args),
        )
    }

    switch arg := args[0].(type) {
    case *object.String:
        return &object.Integer{
            Value: int64(len(arg.Value)),
        }
    default:
        return newError(
            "expected sequence-like argument, got %s type",
            args[0].Type(),
        )
    }
}

func pow(args ...object.Object) object.Object {
    if len(args) != 2 {
        return newError(
            "expected 2 numeric arguments, got %d arguments",
            len(args),
        )
    }
    
    var flag bool
    var aInt, bInt *object.Integer
    var aFloat, bFloat *object.Float

    if aFloat = CheckFloat(args[0]); aFloat != nil {
        flag = true
    } else if aInt = CheckInt(args[0]); aInt == nil {
        return newError(
            "expected only numeric arguments, got %s and %s",
            args[0].Type(),
            args[1].Type(),
        )
    }

    if bFloat = CheckFloat(args[1]); bFloat != nil {
        flag = true
        if aInt != nil {
            aFloat = args[0].(*object.Float)
        }
    } else if bInt = CheckInt(args[1]); bInt == nil {
        return newError(
            "expected only numeric arguments, got %s and %s",
            args[0].Type(),
            args[1].Type(),
        )
    }

    if flag {
        return &object.Float{
            Value: powT[float64](aFloat.Value, bFloat.Value),
        }
    } else {
        return &object.Integer{
            Value: powT[int64](aInt.Value, bInt.Value),
        }
    }
}

func powT[T Num](a any, b any) T {
    switch a := a.(type) {
    case int64:
        b, _ := b.(int64)
        switch {
            case b == 0:
                return 1
            case b < 0:
                return 0
            default:
                base := a
                for i := int64(1); i < b; i++ {
                    a = a * base
                }
            return T(a)
        }
    case float64:
        return T(0) 
    default:
        // Cannot happen
        return T(0)
    }
}
