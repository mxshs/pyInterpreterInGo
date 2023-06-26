package ast

import "mxshs/pyinterpreter/token"

type Node interface {
    TokenLiteral() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}

type Program struct {
    Statements []Statement
}

func (p *Program) TokenLiteral() string {
    if len(p.Statements) > 0 {
        return p.Statements[0].TokenLiteral()
    } else {
        return ""
    }
}

type AssignStatement struct {
    Token token.Token
    Name *Name
    Value Expression
}

func (as *AssignStatement) statementNode() {
}

func (as *AssignStatement) TokenLiteral() string {
    return as.Token.Literal
}

type Name struct {
    Token token.Token
    Value string
}

func (n *Name) expressionNode() {
}

func (n *Name) TokenLiteral() string {
    return n.Token.Literal
}
