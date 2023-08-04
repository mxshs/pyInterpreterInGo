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

func (e *Env) findInGlobalScope(name string) (*Env, bool) {
    _, ok := e.store[name]
    if !ok && e.parent != nil {
        e, ok = e.parent.findInGlobalScope(name)
    }

    return e, ok
}

func (e *Env) Set(name string, value Object) Object {
    parent, ok := e.findInGlobalScope(name)
    if ok {
        parent.store[name] = value
    } else {
        e.store[name] = value
    }

    return value
}

