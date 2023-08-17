package repl

import (
	"bufio"
	"io"
	"strings"

	"mxshs/pyinterpreter/eval"
	"mxshs/pyinterpreter/lexer"
	"mxshs/pyinterpreter/object"
	"mxshs/pyinterpreter/parser"
)

const PROMPT = ">> "

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

        res := read(scanner, out)

        l := lexer.GetLexer(res)

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

func read(scanner *bufio.Scanner, out io.Writer) string {
    var block []string
    var indent int
    
    indentstr := " "

    for {
        io.WriteString(out, ">>" + indentstr)
        scanner.Scan()

        line := scanner.Text()

        if len(line) == 0 {
            if indent == 0 {
                return strings.Join(block, "\n")
            } else {
                indent -= 4
                indentstr = strings.Repeat(" ", indent + 1)
                continue
            }
        }

        block = append(block, strings.Repeat(" ", indent) + line)
        if line[len(line) - 1] == 58 {
            indent += 4
            indentstr = strings.Repeat(" ", indent + 1)
        }
    }
}

