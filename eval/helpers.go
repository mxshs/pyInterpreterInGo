package eval

func pow(a int64, b int64) any {
    switch {
    case b == 0:
        return int64(1)
    case b < 0:
        return NULL
    default:
        base := a
        for i := int64(1); i < b; i++ {
            a = a * base
        }
        return a
    }
}

