package ast

import (
	"bytes"
	"mxshs/pyinterpreter/token"
	"strings"
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
    out.WriteString(as.TokenLiteral() + " ")
    
    if as.Value != nil {
        out.WriteString(as.Value.String())
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

type FloatLiteral struct {
    Token token.Token
    Value float64
}

func (fl *FloatLiteral) expressionNode() {}

func (fl *FloatLiteral) TokenLiteral() string {
    return fl.Token.Literal
}
func (fl *FloatLiteral) String() string {
    return fl.Token.Literal
}

type Boolean struct {
    Token token.Token
    Value bool
}

func (b *Boolean) expressionNode() {}

func (b *Boolean) TokenLiteral() string {
    return b.Token.Literal
}

func (b *Boolean) String() string {
    return b.Token.Literal
}

type StringLiteral struct {
    Token token.Token
    Value string
}

func (s *StringLiteral) expressionNode() {
}

func (s *StringLiteral) TokenLiteral() string {
    return s.Token.Literal
}

func (s *StringLiteral) String() string {
    return s.Token.Literal
}

type IfExpression struct {
    Token token.Token
    Condition Expression
    Consequence *BlockStatement
    Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

func (ie *IfExpression) TokenLiteral() string {
    return ie.Token.Literal
}

func (ie *IfExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(ie.Token.Literal + " ")
    out.WriteString(ie.Condition.String() + " ")
    out.WriteString(ie.Consequence.String())
    
    if ie.Alternative != nil {
        out.WriteString(" else ")
        out.WriteString(ie.Alternative.String())
    }
    
    out.WriteString(")")

    return out.String()
}

type BlockStatement struct {
    Token token.Token
    Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) TokenLiteral() string {
    return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
    var out bytes.Buffer

    for _, statement := range bs.Statements {
        out.WriteString(statement.String())
    }

    return out.String()
}

type FunctionStatement struct {
    Token token.Token
    Name *Name
    Arguments []*Name
    Body *BlockStatement
}

func (fs *FunctionStatement) statementNode() {}

func (fs *FunctionStatement) TokenLiteral() string {
    return fs.Token.Literal
}

func (fs *FunctionStatement) String() string {
    var out bytes.Buffer

    args := []string{}

    for _, arg := range fs.Arguments {
        args = append(args, arg.String())
    }

    out.WriteString(fs.TokenLiteral())
    out.WriteString("(" + strings.Join(args, ", ") + ")")
    out.WriteString(fs.Body.String())

    return out.String()
}

type PrefixExpression struct {
    Token token.Token
    Operator string
    Right Expression
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenLiteral() string {
    return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(pe.Operator)
    out.WriteString(pe.Right.String())
    out.WriteString(")")

    return out.String()
}

type InfixExpression struct {
    Token token.Token
    Operator string
    Left Expression
    Right Expression
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) TokenLiteral() string {
    return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(ie.Left.String())
    out.WriteString(" " + ie.Operator + " ")
    out.WriteString(ie.Right.String())
    out.WriteString(")")

    return out.String()
}

type CallExpression struct {
    Token token.Token
    Function Expression
    Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) TokenLiteral() string {
    return ce.Token.Literal
}

func (ce *CallExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(" + ce.Function.String())
    
    args := []string{}
    for _, arg := range ce.Arguments {
        args = append(args, arg.String())
    }

    out.WriteString("(" + strings.Join(args, ", ") + "))")
    
    return out.String()
}

type ListLiteral struct {
    Token token.Token
    Arr []Expression
}

func (ll *ListLiteral) expressionNode() {}

func (ll *ListLiteral) TokenLiteral() string {
    return ll.Token.Literal
}

func (ll *ListLiteral) String() string {
    var out bytes.Buffer

    out.WriteString("list([")

    for _, elem := range ll.Arr {
        out.WriteString(elem.String() + ", ")
    }

    out.WriteString("])")

    return out.String()
}

type IndexExpression struct {
    Token token.Token
    Struct Expression
    Value Expression
}

func (ie *IndexExpression) expressionNode() {}

func (ie *IndexExpression) TokenLiteral() string {
    return ie.Token.Literal
}

func (ie *IndexExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(" + ie.Struct.String())
    out.WriteString("[" + ie.Value.String() + "]")
    out.WriteString(")")

    return out.String()
}

