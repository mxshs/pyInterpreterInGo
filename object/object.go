package object

import (
	"bytes"
	"fmt"
	"mxshs/pyinterpreter/ast"
	"strings"
)

const (
    INTEGER_OBJ = "INTEGER"
    BOOL_OBJ = "BOOLEAN"
    NULL_OBJ = "NULL"
    RETURN_VALUE = "RETURN_VALUE"
    ERROR_OBJ = "ERROR"
    FUNCTION_OBJ = "FUNCTION"
)

type ObjectType string

type Object interface {
    Type() ObjectType
    Inspect() string
}

type Integer struct {
    Value int64
}

func (i *Integer) Type() ObjectType {
    return INTEGER_OBJ
}

func (i *Integer) Inspect() string {
    return fmt.Sprintf("%d", i.Value)
}

type Boolean struct {
    Value bool
}

func (b *Boolean) Type() ObjectType {
    return BOOL_OBJ
}

func (b *Boolean) Inspect() string {
    return fmt.Sprintf("%t", b.Value)
}

type Null struct {
}

func (n *Null) Type() ObjectType {
    return NULL_OBJ
}

func (n *Null) Inspect() string {
    return "null"
}

type ReturnValue struct {
    Value Object
}

func (rv *ReturnValue) Type() ObjectType {
    return RETURN_VALUE
}

func (rv *ReturnValue) Inspect() string {
    return rv.Value.Inspect()
}

type Error struct {
    Message string
}

func (e *Error) Type() ObjectType {
    return ERROR_OBJ
}

func (e *Error) Inspect() string {
    return e.Message
}

type Function struct {
    Arguments []*ast.Name
    Body *ast.BlockStatement
    Env *Env
}

func (f *Function) Type() ObjectType {
    return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
    args := []string{}

    for _, arg := range f.Arguments {
        args = append(args, arg.String())
    }

    var out bytes.Buffer

    out.WriteString("function ")
    out.WriteString("(" + strings.Join(args, ", ") + ")")
    out.WriteString(":\n" + f.Body.String() + "\n")

    return out.String()
}

