package repl

import (
	"bufio"
	"fmt"
	"io"

	"mxshs/pyinterpreter/eval"
	"mxshs/pyinterpreter/lexer"
	"mxshs/pyinterpreter/object"
	"mxshs/pyinterpreter/parser"
)

const PROMPT = "eval is not stable yet >> "

//const test = `def p(b, c):
//    if (b > c):
//        def g(i, e):
//            return i + e
//       return g(b, c)
//    else:
//        def g(i, e):
//            return i * e
//       return g(b, c)
//p(3,5)`

func StartREPL(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)
    env := object.NewEnv()

    for {
        fmt.Print(PROMPT)

        scanned := scanner.Scan()
        if !scanned {
            break
        }
    
        line := scanner.Text()
        l := lexer.GetLexer(line)
        
        p := parser.GetParser(l)
        program := p.ParseProgram()

        if len(p.Errors()) != 0 {
            for _, err := range p.Errors() {
                io.WriteString(out, "\t" + err + "\n")
            }
            continue
        }

        evaluated := eval.Eval(program, env)

        if evaluated != nil {
            io.WriteString(out, evaluated.Inspect())
            io.WriteString(out, "\n")
        }
    }
}

