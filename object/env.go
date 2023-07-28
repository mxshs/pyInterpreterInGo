package object

func NewNestedEnv(parent *Env) *Env {
    env := NewEnv()
    env.parent = parent
    return env
}

func NewEnv() *Env {
    env := make(map[string]Object)
    return &Env{store: env}
}

type Env struct {
    store map[string]Object
    parent *Env
}

func (e *Env) Get(name string) (Object, bool) {
    obj, ok := e.store[name]
    if !ok && e.parent != nil {
        obj, ok = e.parent.Get(name)
    }

    return obj, ok
}

func (e *Env) Set(name string, value Object) Object {
    _, ok := e.store[name]
    if !ok && e.parent != nil {
        e.parent.Set(name, value)
    } else {
        e.store[name] = value
    }
    
    return value
}

