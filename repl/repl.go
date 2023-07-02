package repl

import (
	"bufio"
	"fmt"
	"io"
	"reflect"

	"mxshs/pyinterpreter/lexer"
	"mxshs/pyinterpreter/parser"
	//"mxshs/pyinterpreter/token"
)

const PROMPT = "tokens only >> "

func StartREPL(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)

    for {
        fmt.Printf(PROMPT)

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

        for _, s := range program.Statements {
            io.WriteString(out, s.String() + "\n")
            io.WriteString(out, reflect.TypeOf(s).String() + "\n")
        }
//        for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
//            fmt.Printf("%+v\n", tok)
//        }
    }
}
