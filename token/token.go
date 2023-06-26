package token

type TokenType string

type Token struct {
    Type TokenType
    Literal string
}


const (
    EOF = "EOF"
    NEWL = "\n"

    NAME = "NAME"
    INT = "INT"

    ASSIGN = "="
    PLUS = "+"
    MINUS = "-"
    STAR = "*"
    SLASH = "/"
    BANG = "!"

    LT = "<"
    GT = ">"

    EQ = "=="
    NOT_EQ = "!="
    LESS_EQ = "<="
    GREATER_EQ = ">="


    PLUS_ASSIGN = "+="
    MINUS_ASSIGN = "-="
    TIMES_ASSIGN = "*="
    DIV_ASSIGN = "/="
 
    DOUBLE_STAR = "**"

    COMMA = ","
    // SEMICOLON = ";"
    COLON = ":"

    LPAR = "("
    RPAR = ")"
    LSQB = "{"
    RSQB = "}"

    FDEF = "def"
    BTRUE = "TRUE"
    BFALSE = "FALSE"
    IF = "IF"
    ELSE = "ELSE"
    FOR = "FOR"
    RETURN = "RETURN"
)

var keywords = map[string]TokenType{
    "def": FDEF,
    "true": BTRUE, 
    "false": BFALSE, 
    "if": IF,
    "else": ELSE,
    "for": FOR,
    "return": RETURN, 
}

func LookupKey(key string) TokenType{
    if tok, ok := keywords[key]; ok {
        return tok
    }

    return NAME
}
