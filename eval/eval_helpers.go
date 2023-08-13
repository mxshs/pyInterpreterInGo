package eval

import (
	"mxshs/pyinterpreter/object"
)

func CheckInt(val object.Object) (*object.Integer) {
    if res, ok := val.(*object.Integer); ok {
        return res
    }

    return nil
}

func CheckFloat(val object.Object) (*object.Float) {
    if res, ok := val.(*object.Float); ok {
        return res
    }

    return nil
}

func IsNumeric(val object.Object) bool {
    if _, ok := val.(*object.Integer); ok {
        return true
    } else if _, ok := val.(*object.Float); ok {
        return true
    }

    return false
}

