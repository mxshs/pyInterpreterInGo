package parser

import (
	"fmt"
	//	"mxshs/pyinterpreter/ast"
	"mxshs/pyinterpreter/ast"
	"mxshs/pyinterpreter/lexer"
	// "mxshs/pyinterpreter/parser"
)

func Run() {
    input := "35"
    l := lexer.GetLexer(input)
    p := GetParser(l)
    program := p.ParseProgram()

    for _, s := range program.Statements {
        stmt, ok := s.(*ast.ExpressionStatement)
        if !ok {
            fmt.Printf("wrong type, got: %T", s)
        }

        ident, ok := stmt.Expression.(*ast.IntegerLiteral)
        if !ok {
            fmt.Printf("wrong exp type, got: %T", stmt.Expression)
        }
        fmt.Printf("%d %s", ident.Value, ident.TokenLiteral())
        // fmt.Printf("%s %s %s", s.String(), s.TokenLiteral(), program.Statements[0])
    }
}

