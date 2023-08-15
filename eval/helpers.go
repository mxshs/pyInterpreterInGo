package eval

import (
	"mxshs/pyinterpreter/object"
)

func pow(a int64, b int64) int64 {
    switch {
    case b == 0:
        return int64(1)
    case b < 0:
        panic("implement negative powers")
    default:
        base := int64(1)
        for b > 0 {
            if b & 1 != 0 {
                base *= a
            }

            a = a * a
            b >>= 1
        }
        
        return base
    }
}

func powFloat(a, b float64) float64 {
    switch {
    case a == 0 || b == 0:
        return 1
    default:
        return 0
    }
}

func IsNumeric(obj object.Object) bool {
    switch obj.Type() {
    case object.INTEGER_OBJ:
        return true
    case object.FLOAT_OBJ:
        return true
    default:
        return false
    }
}

