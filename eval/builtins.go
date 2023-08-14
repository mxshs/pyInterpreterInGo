package eval

import (
//	"fmt"
	"mxshs/pyinterpreter/object"
)

var bltins = map[string]*object.Bltin{
    "len": &object.Bltin{
        Fn: pyLen,
    },
    "sum": &object.Bltin{
        Fn: pySum,
    },
}

func pyLen(args ...object.Object) object.Object {
    if len(args) != 1 {
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
    case *object.List:
        return &object.Integer{
            Value: int64(len(arg.Arr)),
        }
    default:
        return newError(
            "expected sequence-like argument, got %s argument",
            args[0].Type(),
        )
    }
}

func pySum(args ...object.Object) object.Object {
    switch args[0].Type() {
    default:
    //case object.INTEGER_OBJ:
        var s1 *float64
        var s2 *int64
        temp := int64(0)
        s2 = &temp
        // a := args[0].(*object.Integer).Value
        for _, obj := range args {
            if obj.Type() == object.INTEGER_OBJ {
                if s1 != nil {
                    *s1 += float64(obj.(*object.Integer).Value)
                } else {
                    *s2 += obj.(*object.Integer).Value
                }
            } else {
                if s1 == nil {
                    temp := float64(*s2)
                    s1 = &temp
                }
                *s1 += obj.(*object.Float).Value
            }
        }

        if s1 != nil {
            return &object.Float{Value: *s1}
        }

        return &object.Integer{
            Value: *s2,
        }
    }
}
