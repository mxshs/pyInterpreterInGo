package lexer

import (
	"fmt"
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

func (l *Lexer) NextToken() token.Token {
    var tok token.Token

    l.omitSymbol()

    switch l.ch {
    case '=':
        if l.peekChar() == '=' {
            ch := l.ch
            l.nextChar()
            tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
        } else {
            tok = newToken(token.ASSIGN, l.ch)
        }
    case '+':
        if l.peekChar() == '=' {
            ch := l.ch
            l.nextChar()
            tok = token.Token{Type: token.PLUS_ASSIGN, Literal: string(ch) + string(l.ch)}
        } else {
            tok = newToken(token.PLUS, l.ch)
        }
    case '-':
        if l.peekChar() == '=' {
            ch := l.ch
            l.nextChar()
            tok = token.Token{Type: token.MINUS_ASSIGN, Literal: string(ch) + string(l.ch)}
        } else {
            tok = newToken(token.MINUS, l.ch)
        }
    case '*':
        peeked := l.peekChar()
        switch peeked {
        case '=':
            ch := l.ch
            l.nextChar()
            tok = token.Token{Type: token.TIMES_ASSIGN, Literal: string(ch) + string(l.ch)}
        case '*':
            ch := l.ch
            l.nextChar()
            tok = token.Token{Type: token.DOUBLE_STAR, Literal: string(ch) + string(l.ch)}
        default:
            tok = newToken(token.STAR, l.ch)
        }
    case '/':
        if l.peekChar() == '=' {
            ch := l.ch
            l.nextChar()
            tok = token.Token{Type: token.DIV_ASSIGN, Literal: string(ch) + string(l.ch)}
        } else {
            tok = newToken(token.SLASH, l.ch)
        }
    case '!':
        if l.peekChar() == '=' {
            ch := l.ch
            l.nextChar()
            tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
        } else {
            tok = newToken(token.BANG, l.ch)
        }
    case '<':
        if l.peekChar() == '=' {
            ch := l.ch
            l.nextChar()
            tok = token.Token{Type: token.LESS_EQ, Literal: string(ch) + string(l.ch)}
        } else {
            tok = newToken(token.LT, l.ch)
        }
    case '>':
        if l.peekChar() == '=' {
            ch := l.ch
            l.nextChar()
            tok = token.Token{Type: token.GREATER_EQ, Literal: string(ch) + string(l.ch)}
        } else {
            tok = newToken(token.GT, l.ch)
        }
    case '{':
        tok = newToken(token.LSQB, l.ch)
    case '}':
        tok = newToken(token.RSQB, l.ch)
    case '(':
        tok = newToken(token.LPAR, l.ch)
    case ')':
        tok = newToken(token.RPAR, l.ch)
    case '[':
        tok = newToken(token.LBR, l.ch)
    case ']':
        tok = newToken(token.RBR, l.ch)
    case ',':
        tok = newToken(token.COMMA, l.ch)
    case ':':
        tok = newToken(token.COLON, l.ch)
    case '"':
        tok.Type = token.STRING
        tok.Literal = l.readString()
    case '\n':
        tok = newToken(token.NEWL, l.ch)
    case 0:
        tok.Literal = ""
        tok.Type = token.EOF
    default:
        if isLetter(l.ch) {
            tok.Literal = l.readIdent()
            tok.Type = token.LookupKey(tok.Literal)
            return tok
        } else if isDigit(l.ch) {
            literal, flag := l.readNumber()
            tok.Literal = literal
            if flag {
                tok.Type = token.FLOAT
            } else {
                tok.Type = token.INT
            }
            return tok
        } else {
            panic(
                fmt.Sprintf(
                    "implement illegal chars, char %q caused error",
                    string(l.ch),
                ),
            )
        }
    }

    l.nextChar()

    return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
    return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdent() string {
    position := l.position
    for isLetter(l.ch) {
        l.nextChar()
    }

    return l.input[position:l.position]
}

func (l *Lexer) readNumber() (string, bool) {
    var flag bool
    position := l.position

    for {
        if isDigit(l.ch) {
            l.nextChar()
        } else if l.ch == '.' {
            l.nextChar()
            if flag {
                break
            }

            flag = true
        } else {
            break
        }
    }

    return l.input[position:l.position], flag
}

func (l *Lexer) readString() string {
    position := l.position + 1

    for {
        l.nextChar()
        if l.ch == '"' || l.ch == 0 || l.ch == '\n' {
            break
        }
    }

    return l.input[position:l.position]
}

func isLetter(ch byte) bool {
    return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
    return ch >= '0' && ch <= '9'
}

func (l *Lexer) omitSymbol() {
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
        l.nextChar()
    }
}

func (l *Lexer) peekChar() byte {
    if l.readPosition >= l.inputSize {
        return 0
    } else {
        return l.input[l.readPosition]
    }
}

