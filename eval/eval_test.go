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

func TestIfStatements(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    } {
        {"if (true): (69)", 69},
        {"if (false): (420)", nil},
        {"if (true): (69) else: (420)", 69},
        {"if (false): (69) else: (420)", 420},
        {"if (1 < 2): (69)", 69},
        {"if (1 > 2): (420)", nil}, 
        {"if (1 < 2): (69) else: (420)", 69},
        {"if (1 > 2): (69) else: (420)", 420},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        integer, ok := tt.expected.(int)
        if ok {
            testIntegerObject(t, evaluated, int64(integer))
        } else {
            testNullObject(t, evaluated)
        }
    }
}

func TestReturnStatement(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    } {
        {"return 69", 69},
        {"return 420 \n 420", 420},
        {"return 69 * 420 \n 1337", 28980},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testIntegerObject(t, evaluated, tt.expected)
    }
}

func testNullObject(t *testing.T, obj object.Object) bool {
    if obj.Type() == object.NULL_OBJ {
        return true
    } else {
        t.Errorf("Expected object type: %s, got: %T (%+v)", object.NULL_OBJ,
            obj, obj)
        return false
    }
}

func TestAssignStatements(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    } {
        {"a = 69 \n a", 69},
        {"a = 69 * 420 \n a", 28980},
        {"a = 69 \n b = 420 \n a + b", 489},
    }

    for _, tt := range tests {
        value := testEval(tt.input)
        testIntegerObject(t, value, tt.expected)
    }
}

func TestFunctionDef(t *testing.T) {
    tests := []struct {
        input string
        args []string
        body string
    } {
        {
            "def sum(a, b): (return a + b)",
            []string{"a", "b"},
            "return (a + b)"},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)

        function, ok := evaluated.(*object.Function)
        if !ok {
            t.Errorf("Expected object type: Function, got: %T (%+v)",
                evaluated, evaluated)
        }

        for i, arg := range tt.args {
            act_arg := function.Arguments[i].String()  
            if act_arg != arg {
                t.Errorf("Expected function arg at %d: %s, got: %s",
                    i, arg, act_arg)
            }
        }

        if function.Body.String() != tt.body {
            t.Errorf("Expected function body: %+v, got: %+v", 
                tt.body, function.Body.String())
        }
    }
}

func TestFunctions(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    }{
        {"def sum(a, b): return (a + b) \n sum(3, 5)", 8},
    }

    for _, tt := range tests {
        testIntegerObject(t, testEval(tt.input), tt.expected)
    }
}

func testEval(inp string) object.Object {
    l := lexer.GetLexer(inp)
    p := parser.GetParser(l)

    program := p.ParseProgram()

    env := object.NewEnv()

    return Eval(program, env)
}

