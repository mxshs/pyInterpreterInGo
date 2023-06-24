package lexer

import (
    "mxshs/pyinterpreter/token"
)

type Lexer struct {
    input string
    position int
    readPosition int
    ch byte
    inputSize int
}

func GetLexer(input string) *Lexer {
    l := &Lexer{input: input, inputSize: len(input)}
    l.nextChar()
    return l
}

func (l *Lexer) nextChar() {
    if l.readPosition >= l.inputSize {
        l.ch = 0
    } else {
        l.ch = l.input[l.readPosition]
    }
    l.position = l.readPosition
    l.readPosition += 1
}

func (l *Lexer) nextToken() token.Token {
    var tok token.Token

    switch l.ch {
    case '=':
        tok = newToken(token.ASSIGN, l.ch)
    case '(':
        tok = newToken(token.LPAR, l.ch)
    case ')':
        tok = newToken(token.RPAR, l.ch)
    case ',':
        tok = newToken(token.COMMA, l.ch)
    case '+':
        tok = newToken(token.PLUS, l.ch)
    case '{':
        tok = newToken(token.LSQB, l.ch)
    case '}':
        tok = newToken(token.RSQB, l.ch)
    case 0:
        tok.Literal = ""
        tok.Type = token.EOF
    }

    l.nextChar()

    return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
    return token.Token{Type: tokenType, Literal: string(ch)}
}
