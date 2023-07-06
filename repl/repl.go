package repl

import (
	"bufio"
	"fmt"
	"io"
	//"reflect"

	"mxshs/pyinterpreter/eval"
	"mxshs/pyinterpreter/lexer"
	"mxshs/pyinterpreter/parser"
	//"mxshs/pyinterpreter/token"
)

const PROMPT = "eval is not stable yet >> "

func StartREPL(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)

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

        evaluated := eval.Eval(program)

        if evaluated != nil {
            io.WriteString(out, evaluated.Inspect())
            io.WriteString(out, "\n")
        }
//        for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
//            fmt.Printf("%+v\n", tok)
//        }
    }
}
