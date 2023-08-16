package repl

import (
	"fmt"
	"mxshs/pyinterpreter/eval"
	"mxshs/pyinterpreter/lexer"
	"mxshs/pyinterpreter/object"
	"mxshs/pyinterpreter/parser"
	"os"
)

func Read() {
    handle, err := os.ReadFile("test.py69")
    if err != nil {
        panic(err)
    }

    l := lexer.GetLexer(string(handle))
    p := parser.GetParser(l)
    program := p.ParseProgram()

    env := object.NewEnv()
    evaluated := eval.Eval(program, env)

    fmt.Print(evaluated.Inspect())
}

