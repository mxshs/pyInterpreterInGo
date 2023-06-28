package ast

import (
	"bytes"
	"mxshs/pyinterpreter/token"
)

type Node interface {
    TokenLiteral() string
    String() string
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

func (p *Program) String() string {
    var out bytes.Buffer

    for _, s := range p.Statements {
        out.WriteString(s.String())
    }

    return out.String()
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

func (as *AssignStatement) String() string {
    var out bytes.Buffer

    out.WriteString(as.Name.String() + " ")
    out.WriteString(as.TokenLiteral())
    
    if as.Value != nil {
        out.WriteString(" " + as.Value.String())
    }

    return out.String()
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

func (n *Name) String() string {
    return n.Value
} 

type ReturnStatement struct {
    Token token.Token
    ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
    return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
    var out bytes.Buffer

    out.WriteString(rs.TokenLiteral() + " ")
    
    if rs.ReturnValue != nil {
        out.WriteString(rs.ReturnValue.String())
    }

    return out.String()
}

type ExpressionStatement struct {
    Token token.Token
    Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string { 
    return es.Token.Literal 
}

func (es *ExpressionStatement) String() string {

    if es.Expression != nil {
        return es.Expression.String()
    }

    return ""
}

type IntegerLiteral struct {
    Token token.Token
    Value int64
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) TokenLiteral() string {
    return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
    return il.Token.Literal
}

