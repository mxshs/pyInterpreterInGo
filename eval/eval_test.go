package eval

import (
    "mxshs/pyinterpreter/lexer"
    "mxshs/pyinterpreter/object"
    "mxshs/pyinterpreter/parser"

    "testing"
)

func TestEvalInteger(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    } {
        {"6", 6},
        {"69", 69},
        {"-6", -6},
        {"-69", -69},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testIntegerObject(t, evaluated, tt.expected)
    }
}

func testIntegerObject(t *testing.T, obj object.Object, exp int64) bool {
    result, ok := obj.(*object.Integer)
    if !ok {
        t.Errorf("Expected object type: %s, got: %T (%+v)", object.INTEGER_OBJ,
            obj, obj)
        return false
    }

    if result.Value != exp {
        t.Errorf("Expected object integer value: %d, got: %d",
            exp, result.Value)
        return false
    }

    return true
}

func TestEvalBool(t *testing.T) {
    tests := []struct {
        input string
        expected bool
    } {
        {"true", true},
        {"false", false},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testBoolObject(t, evaluated, tt.expected)
    }
}

func testBoolObject(t *testing.T, obj object.Object, exp bool) bool {
    result, ok := obj.(*object.Boolean)
    if !ok {
        t.Errorf("Expected object type: %s, got: %T (%+v)", object.BOOL_OBJ,
            obj, obj)
        return false
    }

    if result.Value != exp {
        t.Errorf("Expected object boolean value: %t, got: %t",
            exp, result.Value)
        return false
    }

    return true
}

func TestBangOperator(t * testing.T) {
    tests := []struct {
        input string
        expected bool
    } {
        {"!false", true},
        {"!true", false},
        {"!!false", false},
        {"!!true", true},
        {"!69", false},
        {"!!69", true},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testBoolObject(t, evaluated, tt.expected)
    }
}

func testEval(inp string) object.Object {
    l := lexer.GetLexer(inp)
    p := parser.GetParser(l)

    program := p.ParseProgram()

    return Eval(program)
}

