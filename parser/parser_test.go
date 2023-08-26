package parser

import (
	//	"fmt"
	"fmt"
	"testing"

	"mxshs/pyinterpreter/ast"
	"mxshs/pyinterpreter/lexer"
)

func TestAssignmentStatements(t *testing.T) {
    input := `
    a 3
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
        // fmt.Printf(statement.TokenLiteral())
        if !testAssignmentStatement(t, statement, tt.expectedName) {
            return
        }
    }
}

func TestInfixExpressions(t * testing.T) {
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
            t.Fatalf("statements length expected to be: %d, got: %d",
                        1,
                        len(program.Statements))
        }
        statement, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("statement expected to be of type: ExpressionStatement, got: %T",
                        program.Statements[0])
        }

        expression, ok := statement.Expression.(*ast.InfixExpression)
        if !ok {
            t.Fatalf("expression expected to be of type: InfixExpression, got: %T",
                        statement.Expression)
        }

        if expression.Operator != tt.operator {
            t.Fatalf("expression operator expected to be: %s, got: %s",
                        tt.operator,
                        expression.Operator)
        }

        fmt.Printf("expression: %s", expression)
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
                "function name expected to be %s, got: %s",
                tt.expectedName,
                statement.Name.String(),
            )
        }

        for i, val := range tt.expectedArgs {
            if statement.Arguments[i].String() != val {
                t.Fatalf(
                    "function args expected to be %s, got: %s",
                    val,
                    statement.Arguments[i].String(),
                )
            }
        }
        
        if statement.Body.String() != tt.expectedBody {
            t.Fatalf(
                "function name expected to be %s, got: %s",
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
                "function name expected to be %s, got: %s",
                tt.expectedName,
                call.Function.String(),
            )
        }

        for i, val := range tt.expectedArgs {
            if call.Arguments[i].String() != val {
                t.Fatalf(
                    "function args expected to be %s, got: %s",
                    val,
                    call.Arguments[i].String(),
                )
            }
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
        t.Logf("Assignment literal supposed to be: =, got: %T", s.TokenLiteral())
        return false
    }

    assignment, ok := s.(*ast.AssignStatement)
    if !ok {
        t.Logf("Assignment statement is not present, got %T", s)
        return false
    }

    if assignment.Name.Value != name {
        t.Logf("Assignment statement name value supposed to be: %s, got: %s", name, assignment.Name.Value)
        return false
    }

    if assignment.Name.TokenLiteral() != name {
        t.Logf("Assignment statement name literal supposed to be: %s, got: %s", name, assignment.Name)
        return false
    }

    return true
}

