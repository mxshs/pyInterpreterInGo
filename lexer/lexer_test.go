package lexer

import (
    "testing"

    "mxshs/pyinterpreter/token"
)

func TestTokens(t *testing.T) {
    input := `
    a = 35
    b = 5

    def add(a, b):
        return a + b

    def compare(a, b):
        if a >= b:
            print(a)
        if a <= b:
            print(b)
        return a, b

    def do_stuff(a, b):
        a *= 3
        a += 3
        b /= 3
        b -= 3
        c = a - b
        c = a * b
        c = a / b
        c = !a
        c = a ** 2
        print(c != a)
        print(c == a)
        d = {hey: world}
        if a:
            print(a)
        else:
            print(wtf)
        a = true
        b = false
        
        for i in range(b):
            print(i)

        return c

    c = add(a, b)

    print(c)`

    tests := []struct {
        expectedType token.TokenType
        expectedLiteral string
    }{
        {token.NAME, "a"},
        {token.ASSIGN, "="},
        {token.INT, "35"},
        {token.NAME, "b"},
        {token.ASSIGN, "="},
        {token.INT, "5"},
        {token.FDEF, "def"},
        {token.NAME, "add"},
        {token.LPAR, "("},
        {token.NAME, "a"},
        {token.COMMA, ","},
        {token.NAME, "b"},
        {token.RPAR, ")"},
        {token.COLON, ":"},
        {token.RETURN, "return"},
        {token.NAME, "a"},
        {token.PLUS, "+"},
        {token.NAME, "b"},
        {token.FDEF, "def"},
        {token.NAME, "compare"},
        {token.LPAR, "("},
        {token.NAME, "a"},
        {token.COMMA, ","},
        {token.NAME, "b"},
        {token.RPAR, ")"},
        {token.COLON, ":"},
        {token.IF, "if"},
        {token.NAME, "a"},
        {token.GREATER_EQ, ">="},
        {token.NAME, "b"},
        {token.COLON, ":"},
        {token.NAME, "print"},
        {token.LPAR, "("},
        {token.NAME, "a"},
        {token.RPAR, ")"},
        {token.IF, "if"},
        {token.NAME, "a"},
        {token.LESS_EQ, "<="},
        {token.NAME, "b"},
        {token.COLON, ":"},
        {token.NAME, "print"},
        {token.LPAR, "("},
        {token.NAME, "b"},
        {token.RPAR, ")"}, 
        {token.RETURN, "return"},
        {token.NAME, "a"},
        {token.COMMA, ","},
        {token.NAME, "b"},

        {token.FDEF, "def"},
        {token.NAME, "do_stuff"},
        {token.LPAR, "("},
        {token.NAME, "a"},
        {token.COMMA, ","},
        {token.NAME, "b"},
        {token.RPAR, ")"},
        {token.COLON, ":"},
        {token.NAME, "a"},
        {token.TIMES_ASSIGN, "*="},
        {token.INT, "3"},
        {token.NAME, "a"},
        {token.PLUS_ASSIGN, "+="},
        {token.INT, "3"},
        {token.NAME, "b"},
        {token.DIV_ASSIGN, "/="},
        {token.INT, "3"},
        {token.NAME, "b"},
        {token.MINUS_ASSIGN, "-="},
        {token.INT, "3"},
        {token.NAME, "c"},
        {token.ASSIGN, "="},
        {token.NAME, "a"},
        {token.MINUS, "-"},
        {token.NAME, "b"},
        
        {token.NAME, "c"},
        {token.ASSIGN, "="},
        {token.NAME, "a"},
        {token.STAR, "*"},
        {token.NAME, "b"},

        {token.NAME, "c"},
        {token.ASSIGN, "="},
        {token.NAME, "a"},
        {token.SLASH, "/"},
        {token.NAME, "b"},

        {token.NAME, "c"},
        {token.ASSIGN, "="},
        {token.BANG, "!"},
        {token.NAME, "a"},

        {token.NAME, "c"},
        {token.ASSIGN, "="},
        {token.NAME, "a"},
        {token.DOUBLE_STAR, "**"},
        {token.INT, "2"},

        {token.NAME, "print"},
        {token.LPAR, "("},
        {token.NAME, "c"},
        {token.NOT_EQ, "!="},
        {token.NAME, "a"},
        {token.RPAR, ")"}, 

        {token.NAME, "print"},
        {token.LPAR, "("},
        {token.NAME, "c"},
        {token.EQ, "=="},
        {token.NAME, "a"},
        {token.RPAR, ")"}, 

        {token.NAME, "d"},
        {token.ASSIGN, "="},
        {token.LSQB, "{"},
        {token.NAME, "hey"},
        {token.COLON, ":"},
        {token.NAME, "world"},
        {token.RSQB, "}"},

        {token.IF, "if"},
        {token.NAME, "a"},
        {token.COLON, ":"},
        {token.NAME, "print"},
        {token.LPAR, "("},
        {token.NAME, "a"},
        {token.RPAR, ")"},

        {token.ELSE, "else"},
        {token.COLON, ":"},
        {token.NAME, "print"},
        {token.LPAR, "("},
        {token.NAME, "wtf"},
        {token.RPAR, ")"},

        {token.NAME, "a"},
        {token.ASSIGN, "="},
        {token.BTRUE, "true"},
        {token.NAME, "b"},
        {token.ASSIGN, "="},
        {token.BFALSE, "false"},

        {token.FOR, "for"},
        {token.NAME, "i"},
        {token.NAME, "in"},
        {token.NAME, "range"},
        {token.LPAR, "("},
        {token.NAME, "b"},
        {token.RPAR, ")"},
        {token.COLON, ":"},
        {token.NAME, "print"},
        {token.LPAR, "("},
        {token.NAME, "i"},
        {token.RPAR, ")"},

        {token.RETURN, "return"},
        {token.NAME, "c"},
        
        {token.NAME, "c"},
        {token.ASSIGN, "="},
        {token.NAME, "add"},
        {token.LPAR, "("},
        {token.NAME, "a"},
        {token.COMMA, ","},
        {token.NAME, "b"},
        {token.RPAR, ")"},
        {token.NAME, "print"},
        {token.LPAR, "("},
        {token.NAME, "c"},
        {token.RPAR, ")"},
    }

    l := GetLexer(input)

    for i, tt := range tests {
        tok := l.nextToken()

        if tok.Type != tt.expectedType {
            t.Fatalf("tests[%d] - wrong token type: expected %q, got %q",
                i, tt.expectedType, tok.Type)
        }

        if tok.Literal != tt.expectedLiteral {
            t.Fatalf("tests[%d] - wrong token literal: expected %q, got %q",
                i, tt.expectedLiteral, tok.Literal)
        }
    }
}
