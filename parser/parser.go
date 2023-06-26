package parser

import (
    "mxshs/pyinterpreter/token" 
    "mxshs/pyinterpreter/lexer"
    "mxshs/pyinterpreter/ast"
)

type Parser struct {
    l *lexer.Lexer

    curToken token.Token
    peekToken token.Token
}

func GetParser(l *lexer.Lexer) *Parser {
    p := &Parser{l: l}

    p.nextToken()
    p.nextToken()

    return p
}

func (p *Parser) nextToken() {
    p.curToken = p.peekToken
    p.peekToken = p.l.NextToken()
}

func (p *Parser) parseStatement() ast.Statement {
    switch p.curToken.Type {
    case token.NAME:
        return p.parseAssignStatement()
    default:
        return nil
    }
}

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
    tok, literal := p.curToken, p.curToken.Literal

    if !p.expectPeek(token.ASSIGN) {
        return nil
    }
    
    statement := &ast.AssignStatement{Token: p.curToken}

    statement.Name = &ast.Name{Token: tok, Value: literal}

    for !p.tokenIs(token.NEWL) {
        p.nextToken()
    }

    return statement
}

func (p *Parser) tokenIs (t token.TokenType) bool {
    return p.curToken.Type == t
}

func (p *Parser) peekTokenIs (t token.TokenType) bool {
    return p.peekToken.Type == t
}

func (p *Parser) expectPeek (t token.TokenType) bool {
    if p.peekTokenIs(t) {
        p.nextToken()
        return true
    } else {
        return false
    }
}

func (p *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{}
    program.Statements = []ast.Statement{}

    for p.curToken.Type != token.EOF {
        statement := p.parseStatement()
        if statement != nil {
            program.Statements = append(program.Statements, statement)
        }
        p.nextToken()
    }

    return program
}
