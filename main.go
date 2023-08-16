package main

import (
    "fmt"
    "os"
    "os/user"

    "mxshs/pyinterpreter/repl"
    //"mxshs/pyinterpreter/parser"
)

func main() {
//    parser.Run()
    repl.Read()
    user, err := user.Current()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Current user: %s\n", user.Username)
    repl.StartREPL(os.Stdin, os.Stdout)
}

