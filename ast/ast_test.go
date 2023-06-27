package ast

import (
	"testing"

	"mxshs/pyinterpreter/token"
)

func TestString(t *testing.T) {
    program := &Program{
        Statements: []Statement{
            &AssignStatement{
                Token: token.Token{Type: token.ASSIGN, Literal: "="},
                Name: &Name{
                    Token: token.Token{Type: token.NAME, Literal: "abc"},
                    Value: "abc",
                },
                Value: &Name{
                    Token: token.Token{Type: token.NAME, Literal: "def"},
                    Value: "def",
                },
            },
        },
    }
    
    if program.String() != "abc = def" {
        t.Errorf("Expected tree string: abc = def, got: %q", program.String())
    }
}

