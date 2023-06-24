package token

type TokenType string

type Token struct {
    Type TokenType
    Literal string
}


const (
    EOF = "EOF"

    NAME = "NAME"
    INT = "INT"

    ASSIGN = "="
    PLUS = "+"

    COMMA = ","
    // SEMICOLON = ";"

    LPAR = "("
    RPAR = ")"
    LSQB = "{"
    RSQB = "}"

    FDEF = "def"
)



