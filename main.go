package main

import (
    "fmt"
    "os"
    "os/user"

    "mxshs/pyinterpreter/repl"
)

func main() {
    user, err := user.Current()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Yo %s, thats repl u can use\n", user.Username)
    repl.StartREPL(os.Stdin, os.Stdout)
}

