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

func TestEvalFloat(t *testing.T) {
    tests := []struct {
        input string
        expected float64
    } {
        {"6.9", 6.9},
        {"42.0", 42.0},
        {"-6.9", -6.9},
        {"-42.0", -42.0},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testFloatObject(t, evaluated, tt.expected)
    }
}

func testFloatObject(t *testing.T, obj object.Object, exp float64) bool {
    result, ok := obj.(*object.Float)
    if !ok {
        t.Errorf("Expected object type: %s, got: %T (%+v)", object.FLOAT_OBJ,
            obj, obj)
        return false
    }

    if result.Value != exp {
        t.Errorf("Expected object integer value: %f, got: %f",
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

func TestEvalString(t *testing.T) {
    tests := []struct {
        input string
        expected string
    } {
        {"\"69420\"", "69420",},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testStringObject(t, evaluated, tt.expected)
    }
}

func testStringObject(t *testing.T, obj object.Object, exp string) bool {
    result, ok := obj.(*object.String)
    if !ok {
        t.Errorf(
            "expected object type: %s, got: %T",
            object.STRING_OBJ,
            obj,
        )

        return false
    }

    if result.Value != exp {
        t.Errorf(
            "expected string value to be: %s, got: %s",
            exp,
            result.Value,
        )

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
        {"if true: \n\t 69", 69},
        {"if false: 420", nil},
        {"if true: \n\t 69 \n else: 420", 69},
        {"if false: 69 \n else: \n\t 420", 420},
        {"if 1 < 2: \n\t 69", 69},
        {"if 1 > 2: 420", nil}, 
        {"if 1 < 2: \n\t 69 \n else: 420", 69},
        {"if 1 > 2: 69 \n else: \n\t 420", 420},
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
        {"return 420 \n 69", 420},
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
            "def c(a, b):\n\treturn a + b\nc",
            []string{"a", "b"},
            "return (a + b)"},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)

        if isError(evaluated) {
            t.Errorf(evaluated.(*object.Error).Message)
        }

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
        {"def c(a, b): \n\t return a + b \n c(3, 5)", 8},
    }

    for _, tt := range tests {
        testIntegerObject(t, testEval(tt.input), tt.expected)
    }
}

func TestBltins(t *testing.T) {
    tests := []struct {
        input string
        expected string
    } {
        {"len([1, 2.0, \"3\"])", "3"},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)

        if isError(evaluated) {
            t.Errorf(
                "unexpected Error during evaluation: %s",
                evaluated.(*object.Error).Message,
            )
        }

        if val := evaluated.Inspect(); val != tt.expected {
            t.Errorf(
                "expected result of infix expression to be: %s, got: %s",
                tt.expected,
                val,
            )
        }
    }
}

func TestInfixExpressions(t *testing.T) {
    tests := []struct {
        input string
        expected string
    } {
        {"69 + 420\n", "489"},
        {"69 - 420\n", "-351"},
        {"69 * 420\n", "28980"},
        {"5 / 2\n", "2"},
        {"2 ** 6\n", "64"},
        {"6.9 + 0.42\n", "7.32"},
        {"6.9 - 0.42\n", "6.48"},
        {"5.5 / 2\n", "2.75"},
        {"5.5 * 2\n", "11"},
        {"\"hello\" + \"world\"", "helloworld"},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)

        if isError(evaluated) {
            t.Errorf(
                "unexpected Error during evaluation: %s",
                evaluated.(*object.Error).Message,
            )
        }

        if val := evaluated.Inspect(); val != tt.expected {
            t.Errorf(
                "expected result of infix expression to be: %s, got :%s",
                tt.expected,
                val,
            )
        }
    }
}

func TestLists(t *testing.T) {
    tests := []struct {
        input string
        expectedTypes []string
        expectedValues []string
    } {
        {
            "[69, 4.20, \"3\"]",
            []string{object.INTEGER_OBJ, object.FLOAT_OBJ, object.STRING_OBJ},
            []string{"69", "4.2", "3"},
        },
    }

    for _, tt := range tests {
        testListObject(t, testEval(tt.input), tt.expectedTypes, tt.expectedValues)
    }
}

func testListObject(t *testing.T, obj object.Object, expTypes, expValues []string) bool {
    list, ok := obj.(*object.List)
    if !ok {
        t.Errorf(
            "expected object type: %s, got: %T (%+v)",
            object.LIST,
            obj,
            obj,
        )

        return false
    }

    for idx, val := range list.Arr {
        if string(val.Type()) != expTypes[idx] {
            t.Errorf(
                "expected %d element to be of type: %s, got: %T (%+v)",
                idx,
                expTypes[idx],
                val,
                val,
            )
        }
        
        if val.Inspect() != expValues[idx] {
            t.Errorf(
                "Expected %d element to be: %s, got: %s",
                idx,
                expValues[idx],
                val.Inspect(),
            )

            return false
        }
    }

    return true
}

func testEval(inp string) object.Object {
    l := lexer.GetLexer(inp)
    p := parser.GetParser(l)

    program := p.ParseProgram()

    env := object.NewEnv()

    return Eval(program, env)
}

