package parser

import (
//	"fmt"
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

func testParserErrors(t *testing.T, p *Parser) {
    errors := p.Errors()

    if len(errors) == 0 {
        return
    }

    for _, msg := range errors {
        t.Errorf("Parser error: %q:", msg)
    }

    t.FailNow()
}

func testAssignmentStatement(t *testing.T, s ast.Statement, name string) bool {
    if s.TokenLiteral() != "=" {
        t.Errorf("Assignment literal supposed to be: =, got: %T", s.TokenLiteral())
        return false
    }

    assignment, ok := s.(*ast.AssignStatement)
    if !ok {
        t.Errorf("Assignment statement is not present, got %T", s)
        return false
    }

    if assignment.Name.Value != name {
        t.Errorf("Assignment statement name value supposed to be: %s, got: %s", name, assignment.Name.Value)
        return false
    }

    if assignment.Name.TokenLiteral() != name {
        t.Errorf("Assignment statement name literal supposed to be: %s, got: %s", name, assignment.Name)
        return false
    }

    return true
}

