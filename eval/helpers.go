package eval

import (
	"mxshs/pyinterpreter/object"
)

func pow(a int64, b int64) int64 {
    switch {
    case b == 0:
        return int64(1)
    default:
        base := int64(1)
        for b != 0 {
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
        res := 1.0
        enum := log(a) * b
        denom := 1.0

        for i := 1; i < 65; i += 1 {
            denom = denom * float64(i)
            member := floatBasePow(enum, int64(i)) / denom
            res += member
        }

        return res
    }
}

func log(a float64) float64 {
    res := 0.0
    member := (a - 1) / (a + 1)
    for i := 1; i < 259; i += 2 {
        res += 1.0 / float64(i) * floatBasePow(member, int64(i))
    }

    return res * 2
}

func floatBasePow(a float64, b int64) float64 {
    if b == 1 {
        return a
    }
    base := float64(1)
    for b != 0 {
        if b & 1 != 0 {
            base *= a
        }

        a = a * a
        b >>= 1
    }

    return base
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

