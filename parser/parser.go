package parser

import (
	"fmt"

	"mxshs/pyinterpreter/ast"
	"mxshs/pyinterpreter/lexer"
	"mxshs/pyinterpreter/token"
)

const (
    _ int = iota
    LOWEST
    EQUALS
    LESSGREATER
    SUM
    PRODUCT
    PREFIX
    CALL
)

type Parser struct {
    l *lexer.Lexer

    curToken token.Token
    peekToken token.Token
    errors []string

    prefixParsers map[token.TokenType]prefixParse
    infixParsers map[token.TokenType]infixParse
}

type (
    prefixParse func() ast.Expression
    infixParse func(ast.Expression) ast.Expression
)

func GetParser(l *lexer.Lexer) *Parser {
    p := &Parser{
        l: l,
        errors: []string{},
    }

    p.prefixParsers = make(map[token.TokenType]prefixParse)
    p.registerPrefix(token.NAME, p.parseName)

    p.nextToken()
    p.nextToken()

    return p
}

func (p *Parser) Errors() []string {
    return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
    msg := fmt.Sprintf("expected token of type: %s, got: %s", t, p.peekToken.Type)
    p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
    p.curToken = p.peekToken
    p.peekToken = p.l.NextToken()
}

func (p *Parser) parseStatement() ast.Statement {
    switch tok := p.curToken.Type; {
    case tok == token.NAME && p.peekTokenIs(token.ASSIGN):
        return p.parseAssignStatement()
    case tok == token.RETURN:
        return p.parseReturnStatement()
    default:
        return p.parseExpressionStatement()
    }
}

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
    tok, literal := p.curToken, p.curToken.Literal

    p.nextToken()

    statement := &ast.AssignStatement{Token: p.curToken}

    statement.Name = &ast.Name{Token: tok, Value: literal}

    for !p.tokenIs(token.NEWL) {
        p.nextToken()
    }

    return statement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
    statement := &ast.ReturnStatement{Token: p.curToken}

    p.nextToken()

    for !p.tokenIs(token.NEWL) {
        p.nextToken()
    }

    return statement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    statement := &ast.ExpressionStatement{Token: p.curToken}

    statement.Expression = p.parseExpression(LOWEST)

    if p.peekTokenIs(token.NEWL) {
        p.nextToken()
    }

    return statement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
    prefix := p.prefixParsers[p.curToken.Type]
    if prefix == nil {
        return nil
    }

    leftExp := prefix()

    return leftExp
}

func (p *Parser) parseName() ast.Expression {
    return &ast.Name{Token: p.curToken, Value: p.curToken.Literal}
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
        p.peekError(t)
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

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParse) {
    p.prefixParsers[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParse) {
    p.infixParsers[tokenType] = fn
}
