package repl

import (
	"bufio"
	"fmt"
	"io"

	"mxshs/pyinterpreter/lexer"
	"mxshs/pyinterpreter/token"
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

        for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
            fmt.Printf("%+v\n", tok)
        }
    }
}
