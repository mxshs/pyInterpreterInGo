package object

import (
	"bytes"
	"fmt"
	"mxshs/pyinterpreter/ast"
	"strconv"
	"strings"
)

const (
    INTEGER_OBJ = "INTEGER"
    FLOAT_OBJ = "FLOAT"
    BOOL_OBJ = "BOOLEAN"
    STRING_OBJ = "STIRNG"
    NULL_OBJ = "NULL"
    RETURN_VALUE = "RETURN_VALUE"
    ERROR_OBJ = "ERROR"
    FUNCTION_OBJ = "FUNCTION"
    BLTIN = "BLTIN_FN"
    LIST = "LIST"
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

type Float struct {
    Value float64
}

func (f *Float) Type() ObjectType {
    return FLOAT_OBJ 
}

func (f *Float) Inspect() string {
    return strconv.FormatFloat(f.Value, 'f', -1, 64)
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

type String struct {
    Value string
}

func (s *String) Type() ObjectType {
    return STRING_OBJ
}

func (s *String) Inspect() string {
    return s.Value
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

type BuiltinFunction func(args ...Object) Object

type Bltin struct {
    Fn BuiltinFunction
}

func (b *Bltin) Type() ObjectType {
    return BLTIN
}

func (b *Bltin) Inspect() string {
    return "builtin function"
}

type List struct {
    Arr []Object
}

func (l *List) Type() ObjectType {
    return LIST
}

func (l *List) Inspect() string {
    var out bytes.Buffer

    out.WriteString("list([")

    elems := []string{}

    for _, elem := range l.Arr {
        elems = append(elems, elem.Inspect())
    }

    out.WriteString(strings.Join(elems, ", ") + "])")

    return out.String()
}

