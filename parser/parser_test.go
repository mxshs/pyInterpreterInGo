package parser

import (
	"testing"

	"mxshs/pyinterpreter/ast"
	"mxshs/pyinterpreter/lexer"
)

func TestIntegerLiterals(t *testing.T) {
    input := "69"

    l := lexer.GetLexer(input)
    p := GetParser(l)
    program := p.ParseProgram()

    if len(program.Statements) != 1 {
        t.Fatalf(
            "expected number of statements: %d, got: %d",
            1,
            len(program.Statements),
        )
    }

    statement, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf(
            "expected statement of type ast.ExpressionStatement, got: %T",
            program.Statements[0],
        )
    }

    literal, ok := statement.Expression.(*ast.IntegerLiteral)
    if !ok {
        t.Fatalf(
            "expected expression of type ast.IntegerLiteral, got: %T",
            statement.Expression,
        )
    }

    if literal.Value != 69 {
        t.Fatalf(
            "expected expression to evaluate to %d, got %d",
            69,
            literal.Value,
        )
    }
}

func TestFloatLiterals(t *testing.T) {
    input := "420.69"

    l := lexer.GetLexer(input)
    p := GetParser(l)
    program := p.ParseProgram()

    if len(program.Statements) != 1 {
        t.Fatalf(
            "expected number of statements: %d, got: %d",
            1,
            len(program.Statements),
        )
    }

    statement, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf(
            "expected statement of type ast.ExpressionStatement, got: %T",
            program.Statements[0],
        )
    }

    literal, ok := statement.Expression.(*ast.FloatLiteral)
    if !ok {
        t.Fatalf(
            "expected expression of type ast.FloatLiteral, got: %T",
            statement.Expression,
        )
    }

    if literal.Value != 420.69 {
        t.Fatalf(
            "expected expression to evaluate to %f, got %f",
            420.69,
            literal.Value,
        )
    }
}

func TestStringLiterals(t *testing.T) {
    tests := []struct{
        input string
        expected string
    } {
        {"\"420.69\"", "420.69"},
        {"\"hello world\"", "hello world"},
    }

    for _, tt := range tests {
        l := lexer.GetLexer(tt.input)
        p := GetParser(l)
        program := p.ParseProgram()
   
        testStringExpression(program, tt.expected, t)
    }
}

func testStringExpression(program *ast.Program, value string, t *testing.T) {
    if len(program.Statements) != 1 {
        t.Fatalf(
            "expected number of statements: %d, got: %d",
            1,
            len(program.Statements),
        )
    }
    
    expression, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf(
            "expected statement of type ast.ExpressionStatement, got: %T",
            program.Statements[0],
        )
    }

    literal, ok := expression.Expression.(*ast.StringLiteral)
    if !ok {
        t.Fatalf(
            "expected expression of type ast.StringLiteral, got: %T",
            expression.Expression,
        )
    }

    if literal.Value != value {
        t.Fatalf(
            "expected expression to evaluate to %s, got: %s",
            value,
            literal.Value,
        )
    }
}

func TestBoolLiterals(t *testing.T) {
    tests := []struct{
        input string
        expected bool
    } {
        {"true", true},
        {"false", false},
    }

    for _, tt := range tests {
        l := lexer.GetLexer(tt.input)
        p := GetParser(l)
        program := p.ParseProgram()
   
        testBoolExpression(program, tt.expected, t)
    }
}

func testBoolExpression(program *ast.Program, value bool, t *testing.T) {
    if len(program.Statements) != 1 {
        t.Fatalf(
            "expected number of statements: %d, got: %d",
            1,
            len(program.Statements),
        )
    }

    statement, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf(
            "expected statement of type ast.ExpressionStatement, got: %T",
            program.Statements[0],
        )
    }

    expression, ok := statement.Expression.(*ast.Boolean)
    if !ok {
        t.Fatalf(
            "expected expression of type ast.Boolean, got: %T",
            statement.Expression,
        )
    }

    if expression.Value != value {
        t.Fatalf(
            "expected expression to evaluate to %v, got %v",
            value,
            expression.Value,
        )
    }
}

func TestAssignmentStatements(t *testing.T) {
    input := `a 3
    b =
    cd 3535
    `

    l := lexer.GetLexer(input)
    p := GetParser(l)
    program := p.ParseProgram()
    testParserErrors(t, p)    

    tests := []struct {
        expectedName string
    }{
        {"a"},
        {"b"},
        {"cd"},
    }

    for i, tt := range tests {
        statement := program.Statements[i]
        if !testAssignmentStatement(t, statement, tt.expectedName) {
            return
        }
    }
}

func TestPrefixExpressions(t *testing.T) {
    tests := []struct{
        input string
        prefix string
        value string
    } {
        {"!true", "!", "true"},
        {"!69", "!", "69"},
        {"-420", "-", "420"},
    }

    for _, tt := range tests {
        l := lexer.GetLexer(tt.input)
        p := GetParser(l)
        program := p.ParseProgram()
        testParserErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf(
                "expected number of statements: %d, got: %d",
                1,
                len(program.Statements),
            )
        }

        statement, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf(
                "expected statement of type ast.ExpressionStatement, got: %T",
                program.Statements[0],
            )
        }

        expression, ok := statement.Expression.(*ast.PrefixExpression)
        if !ok {
            t.Fatalf(
                "expected expression of type ast.PrefixExpression, got: %T",
                statement.Expression,
            )
        }

        if expression.Operator != tt.prefix {
            t.Fatalf(
                "expected prefix to be %s, got: %s",
                tt.prefix,
                expression.Operator,
            )
        }
        
        if expression.Right.String() != tt.value {
            t.Fatalf(
                "expected value to be %s, got: %s",
                tt.value,
                expression.Right.String(),
            )
        }
    }
}

func TestInfixExpressions(t *testing.T) {
    infixTests := []struct {
        input string
        leftValue int64
        operator string
        rightValue int64
    }{
        {"3 + 5", 3, "+", 5},
        {"69 - 420", 69, "-", 420},
        {"420 * 69", 420, "*", 69},
        {"6 / 9", 6, "/", 9},
        {"9 < 6", 9, "<", 6},
        {"6 > 9", 6, ">", 9},
        {"4 == 20", 4, "==", 20},
        {"20 != 4", 20, "!=", 4},
        {"4 ** 20", 4, "**", 20},
    }

    for _, tt := range infixTests {
        l := lexer.GetLexer(tt.input)
        p := GetParser(l)
        program := p.ParseProgram()
        testParserErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf(
                "expected number of statements: %d, got: %d",
                1,
                len(program.Statements),
            )
        }
        statement, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf(
                "expected statement of type ast.ExpressionStatement, got: %T",
                program.Statements[0],
            )
        }

        expression, ok := statement.Expression.(*ast.InfixExpression)
        if !ok {
            t.Fatalf(
                "expected expression of type ast.InfixExpression, got: %T",
                statement.Expression,
            )
        }

        if expression.Operator != tt.operator {
            t.Fatalf(
                "expected binary operator to be %s, got: %s",
                tt.operator,
                expression.Operator,
            )
        }
    }
}

func TestGroupExpressions(t *testing.T) {
    tests := []struct{
        input string
        expected string
    } {
        {"(3 + 5 * 6)", "(3 + (5 * 6))"},
    }

    for _, tt := range tests {
        l := lexer.GetLexer(tt.input)
        p := GetParser(l)
        program := p.ParseProgram()
        testParserErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf(
                "expected number of statements: %d, got: %d",
                1,
                len(program.Statements),
            )
        }

        statement, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf(
                "expected statement of type ast.ExpressionStatement, got: %T",
                program.Statements[0],
            )
        }

        if statement.Expression.String() != tt.expected {
            t.Fatalf(
                "expected %s as a grouped expression, got %s",
                tt.expected,
                statement.Expression.String(),
            )
        }
    }
}

func TestFunctionStatements(t *testing.T) {
    input := `def a(b, c):
        return a + b
    `

    l := lexer.GetLexer(input)
    p := GetParser(l)
    program := p.ParseProgram()
    testParserErrors(t, p)    

    tests := []struct {
        expectedName string
        expectedArgs []string
        expectedBody string
    }{
        {"a", []string{"b", "c"}, "return (a + b)"},
    }

    for _, tt := range tests {
        statement := program.Statements[0].(*ast.FunctionStatement)

        if statement.Name.String() != tt.expectedName {
            t.Fatalf(
                "expected function name: %s, got: %s",
                tt.expectedName,
                statement.Name.String(),
            )
        }

        for idx, val := range tt.expectedArgs {
            if statement.Arguments[idx].String() != val {
                t.Fatalf(
                    "expected %s at position %d, got: %s",
                    val,
                    idx,
                    statement.Arguments[idx].String(),
                )
            }
        }
        
        if statement.Body.String() != tt.expectedBody {
            t.Fatalf(
                "expected function name to be %s, got: %s",
                tt.expectedBody,
                statement.Body.String(),
            )
        }
    }
}

func TestCallExpression(t *testing.T) {
    input := `a(3, 5)`

    l := lexer.GetLexer(input)
    p := GetParser(l)
    program := p.ParseProgram()
    testParserErrors(t, p)    

    tests := []struct {
        expectedName string
        expectedArgs []string
    }{
        {"a", []string{"3", "5"}},
    }

    for _, tt := range tests {
        statement := program.Statements[0].(*ast.ExpressionStatement)
        call := statement.Expression.(*ast.CallExpression)

        if call.Function.String() != tt.expectedName {
            t.Fatalf(
                "expected function name to be %s, got: %s",
                tt.expectedName,
                call.Function.String(),
            )
        }

        for idx, val := range tt.expectedArgs {
            if call.Arguments[idx].String() != val {
                t.Fatalf(
                    "expected %s at position %d, got: %s",
                    val,
                    idx,
                    call.Arguments[idx].String(),
                )
            }
        }
    }
}

func TestIfStatements(t *testing.T) {
    tests := []struct{
        input string
        condition string
        consequence string
        alternative any 
    } {
        {
            "if 3 < 5: return 5",
            "(3 < 5)",
            "return 5",
            nil,
        },
        {
            "if 3 > 5:\n\treturn 5\nelse:\n\treturn 6",
            "(3 > 5)",
            "return 5",
            "return 6",
        },
        {
            "if 69 > 420: return 69\nelse: return 420",
            "(69 > 420)",
            "return 69",
            "return 420",
        },
    }

    for _, tt := range tests {

        l := lexer.GetLexer(tt.input)
        p := GetParser(l)
        program := p.ParseProgram()
        testParserErrors(t, p)    

        testIfStatement(
            program, tt.condition, tt.consequence, tt.alternative, t)
    }
}

func testIfStatement(
    program *ast.Program,
    condition, consequence string,
    alternative any,
    t *testing.T) {

        if len(program.Statements) != 1 {
            t.Fatalf(
                "expected number of statements: %d, got: %d",
                1,
                len(program.Statements),
            )
        }

        statement, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf(
                "expected statement of type ast.ExpressionStatement, got: %T",
                program.Statements[0],
            )
        }

        expression, ok := statement.Expression.(*ast.IfExpression)
        if !ok {
            t.Fatalf(
                "expected expression of type ast.IfExpression, got: %T",
                statement.Expression,
            )
        }
        
        if expression.Condition.String() != condition {
            t.Fatalf(
                "expected if condition to be %s, got %s",
                condition,
                expression.Condition.String(),
            )
        }

        if expression.Consequence.String() != consequence {
            t.Fatalf(
                "expected if consequence to be %s, got %s",
                consequence,
                expression.Consequence.String(),
            )
        }

        if alternative == nil {
            if expression.Alternative != nil {
                t.Fatalf(
                    "expected else clause to be %v, got: %s",
                    nil,
                    expression.Alternative.String(),
                )
            }
        } else if expression.Alternative.String() != alternative.(string) {
            t.Fatalf(
                "expected else clause to be %s, got: %s",
                alternative,
                expression.Alternative.String(),
            )
        }
}

func TestArrayExpression(t *testing.T) {
    tests := []struct{
        input string
        expected []string
    } {
        {"[]", []string{}},
        {"[1, \"2\", 3.0]", []string{"1", "2", "3.0"}},
    }

    for _, tt := range tests {
        l := lexer.GetLexer(tt.input)
        p := GetParser(l)
        program := p.ParseProgram()
        testParserErrors(t, p)    

        testArrayLiteral(program, tt.expected, t)
    }
}

func testArrayLiteral(program *ast.Program, expected []string, t *testing.T) {
        if len(program.Statements) != 1 {
            t.Fatalf(
                "expected number of statements: %d, got: %d",
                1,
                len(program.Statements),
            )
        }

        statement, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf(
                "expected statement of type ast.ExpressionStatement, got: %T",
                program.Statements[0],
            )
        }

        expression, ok := statement.Expression.(*ast.ListLiteral)
        if !ok {
            t.Fatalf(
                "expected expression of type ast.ListLiteral, got: %T",
                statement.Expression,
            )
        }

        for idx, elem := range expression.Arr {
            if elem.String() != expected[idx] {
                t.Fatalf(
                    "expected %s at position %d, got: %s",
                    expected[idx],
                    idx,
                    elem.String(),
                )
            }
        }
}

func TestIndexExpression(t *testing.T) {
    tests := []struct{
        input string
        expectedName string
        expectedIdx string
    } {
        {"a[1]", "a", "1"},
    }

    for _, tt := range tests {
        l := lexer.GetLexer(tt.input)
        p := GetParser(l)
        program := p.ParseProgram()
        testParserErrors(t, p)    

        if len(program.Statements) != 1 {
            t.Fatalf(
                "expected number of statements: %d, got: %d",
                1,
                len(program.Statements),
            )
        }

        statement, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf(
                "expected statement of type ast.ExpressionStatement, got: %T",
                program.Statements[0],
            )
        }

        expression, ok := statement.Expression.(*ast.IndexExpression)
        if !ok {
            t.Fatalf(
                "expected expression of type ast.IndexExpression, got: %T",
                statement.Expression,
            )
        }

        if expression.Struct.String() != tt.expectedName {
            t.Fatalf(
                "expected indexed struct name to be %s, got: %s",
                tt.expectedName,
                expression.Struct.String(),
            )
        }

        if expression.Value.String() != tt.expectedIdx {
            t.Fatalf(
                "expected index value to be %s, got: %s",
                tt.expectedIdx,
                expression.Value.String(),
            )
        }
    }
}

func testParserErrors(t *testing.T, p *Parser) {
    errors := p.Errors()

    if len(errors) == 0 {
        return
    }

    for _, msg := range errors {
        t.Logf("Parser error: %q:", msg)
    }
}

func testAssignmentStatement(t *testing.T, s ast.Statement, name string) bool {
    if s.TokenLiteral() != "=" {
        t.Logf(
            "expected assignment literal to be: =, got: %v",
            s.TokenLiteral(),
        )
        return false
    }

    assignment, ok := s.(*ast.AssignStatement)
    if !ok {
        t.Logf(
            "expected statement of type ast.AssignStatement, got: %T",
            s,
        )
        return false
    }

    if assignment.Name.Value != name {
        t.Logf(
            "expected %s as name(ident) in assignment, got: %s",
            name,
            assignment.Name.Value,
        )
        return false
    }

    if assignment.Name.TokenLiteral() != name {
        t.Logf(
            "expected %s as token for name in assignment, got: %s",
            name,
            assignment.Name,
        )
        return false
    }

    return true
}

