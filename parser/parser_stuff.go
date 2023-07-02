package parser

import (
	"fmt"
	//	"mxshs/pyinterpreter/ast"
	//	"mxshs/pyinterpreter/ast"
//	"mxshs/pyinterpreter/ast"
	"mxshs/pyinterpreter/lexer"

	// "mxshs/pyinterpreter/parser"
)

func Run() {
    input := `a = 3
    b = 5
    if (a > 3): (a = 5) else: (y = 5)
    print(5)
    def asd(a, b): (x = 5
    return a + b + x)`
    l := lexer.GetLexer(input)
    p := GetParser(l)
    program := p.ParseProgram()

    for _, s := range p.errors {
        fmt.Printf("%s", s)
    }
    for _, s := range program.Statements {
//        fmt.Printf("%s %d", s.String(), i)
        fmt.Printf("%s \n", s.String())
    }
//        stmt, ok := s.(*ast.ExpressionStatement)
//        if !ok {
//            fmt.Printf("wrong type, got: %T", s)
//        }
//
//        prefix, ok := stmt.Expression.(*ast.PrefixExpression)
//        if !ok {
//            fmt.Printf("wrong exp type, got: %T", stmt.Expression)
//        }
//        fmt.Printf("%s %s", prefix.Operator, prefix.String())
        // fmt.Printf("%s %s %s", s.(*ast.ExpressionStatement).Expression.String(), s.(*ast.ExpressionStatement).Expression.TokenLiteral(), program.Statements[0])
}

