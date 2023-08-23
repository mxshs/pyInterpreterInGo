package parser

import (
	"fmt"
	"strconv"

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
    POWER
    PREFIX
    CALL
    INDEX
)

// TODO: test assignment as infix expression, i.e. ignore assign statement and
// construct assign node in infix parsing (should work, cause the only other
// case, where "=" can be encountered is block statement, which is parsed by a
// different parsing fn).

var precedenceMap = map[token.TokenType]int{
    token.EQ: EQUALS,
    token.NOT_EQ: EQUALS,
    token.LT: LESSGREATER,
    token.GT: LESSGREATER,
    token.GREATER_EQ: LESSGREATER,
    token.LESS_EQ: LESSGREATER,
    token.PLUS: SUM,
    token.MINUS: SUM,
    token.SLASH: PRODUCT,
    token.STAR: PRODUCT,
    token.DOUBLE_STAR: POWER,
    token.LPAR: CALL,
    token.LBR: INDEX,
}

type Parser struct {
    l *lexer.Lexer

    curToken token.Token
    peekToken token.Token
    Depth int
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

    p.curToken = p.l.NextToken()
    p.peekToken = p.l.NextToken()

    p.prefixParsers = make(map[token.TokenType]prefixParse)
    p.registerPrefix(token.NAME, p.parseName)
    p.registerPrefix(token.INT, p.parseIntegerLiteral)
    p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
    p.registerPrefix(token.BTRUE, p.parseBoolean)
    p.registerPrefix(token.BFALSE, p.parseBoolean)
    p.registerPrefix(token.STRING, p.parseString)
    p.registerPrefix(token.LPAR, p.parseGroupedExpression)
    p.registerPrefix(token.IF, p.parseIfExpression)
    p.registerPrefix(token.MINUS, p.parsePrefixExpression)
    p.registerPrefix(token.BANG, p.parsePrefixExpression)
    p.registerPrefix(token.STAR, p.parsePrefixExpression)
    p.registerPrefix(token.DOUBLE_STAR, p.parsePrefixExpression)
    p.registerPrefix(token.LBR, p.parseListExpression)

    p.infixParsers = make(map[token.TokenType]infixParse)
    p.registerInfix(token.LPAR, p.parseCallExpression)
    p.registerInfix(token.EQ, p.parseInfixExpression)
    p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
    p.registerInfix(token.LT, p.parseInfixExpression)
    p.registerInfix(token.GT, p.parseInfixExpression)
    p.registerInfix(token.PLUS, p.parseInfixExpression)
    p.registerInfix(token.MINUS, p.parseInfixExpression)
    p.registerInfix(token.SLASH, p.parseInfixExpression)
    p.registerInfix(token.STAR, p.parseInfixExpression)
    p.registerInfix(token.DOUBLE_STAR, p.parseInfixExpression)
    p.registerInfix(token.LBR, p.parseIndexExpression)

    //p.nextToken()
    //p.nextToken()

    return p
}

func (p *Parser) Errors() []string {
    return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
    msg := fmt.Sprintf(
        "expected token of type: %s, got: %s",
        t,
        p.peekToken.Type,
    )
    p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
    p.curToken = p.peekToken
    p.Depth = p.l.GetDepth()
    p.peekToken = p.l.NextToken()
}

func (p *Parser) parseStatement() ast.Statement {
    switch tok := p.curToken.Type; {
    case tok == token.NAME && p.peekTokenIs(token.ASSIGN):
        return p.parseAssignStatement()
    case tok == token.FDEF:
        return p.parseFunctionStatement()
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

    p.nextToken()

    statement.Value = p.parseExpression(LOWEST)

    if p.peekTokenIs(token.NEWL) {
        p.nextToken()
    }

    return statement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
    statement := &ast.ReturnStatement{Token: p.curToken}

    p.nextToken()

    statement.ReturnValue = p.parseExpression(LOWEST)

    if p.peekTokenIs(token.NEWL) {
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
        p.addNoPrefixParseError(p.curToken.Type)
        return nil
    }

    leftExp := prefix()

    for !p.peekTokenIs(token.NEWL) && precedence < p.peekPrecedence() {
        infix := p.infixParsers[p.peekToken.Type]
        if infix == nil {
            return leftExp
        }

        p.nextToken()

        leftExp = infix(leftExp)

        if p.curToken.Type == token.EOF {
            break
        }
    }

    return leftExp
}

func (p *Parser) addNoPrefixParseError(t token.TokenType) {
    msg := fmt.Sprintf("prefix parse function not found for type %s %s", t, p.peekToken.Literal)
    p.errors = append(p.errors, msg)
}

func (p *Parser) parseName() ast.Expression {
    return &ast.Name{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
    literal := &ast.IntegerLiteral{Token: p.curToken}

    value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
    if err != nil {

        msg := fmt.Sprintf(
            "error during parsing %q as integer",
            p.curToken.Literal,
        )
        p.errors = append(p.errors, msg)
        return nil
    }

    literal.Value = value

    return literal
}

func (p *Parser) parseFloatLiteral() ast.Expression {
    literal := &ast.FloatLiteral{Token: p.curToken}

    value, err := strconv.ParseFloat(p.curToken.Literal, 64)
    if err != nil {
        msg := fmt.Sprintf(
            "error during parsing %q as float",
            p.curToken.Literal,
        )
        p.errors = append(p.errors, msg)
        return nil
    }

    literal.Value = value

    return literal
}

func (p *Parser) parseBoolean() ast.Expression {
    literal := &ast.Boolean{Token: p.curToken}

    value, err := strconv.ParseBool(p.curToken.Literal)
    if err != nil {
        msg := fmt.Sprintf(
            "boolean parse error. Cannot convert %s to boolean",
            p.curToken.Literal,
        )
        p.errors = append(p.errors, msg)
        return nil
    }

    literal.Value = value

    return literal
}

func (p *Parser) parseString() ast.Expression {
    literal := &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}

    return literal
}

func (p *Parser) parsePrefixExpression() ast.Expression {
    expression := &ast.PrefixExpression{
        Token: p.curToken,
        Operator: p.curToken.Literal,
    }

    p.nextToken()

    expression.Right = p.parseExpression(PREFIX)

    return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
    expression := &ast.InfixExpression{
        Token: p.curToken,
        Operator: p.curToken.Literal,
        Left: left,
    }

    precedence := p.precedenceIs()
    p.nextToken()
    expression.Right = p.parseExpression(precedence)

    return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
    p.nextToken()

    exp := p.parseExpression(LOWEST)

    if !p.expectPeek(token.RPAR) {
        //fmt.Printf("wanted %s, got %s and %s", token.RPAR, p.curToken.Literal, p.peekToken.Literal)
        return nil
    }

    return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
    expression := &ast.IfExpression{Token: p.curToken}
    
    //if !p.expectPeek(token.LPAR) {
      //  return nil
    //}

    p.nextToken()

    expression.Condition = p.parseExpression(LOWEST)

    //if !p.expectPeek(token.RPAR) && !p.expectPeek(token.COLON) && !p.expectPeek(token.LPAR) {
      //  return nil
    //}

    //if !p.expectPeek(token.RPAR) {
      //  return nil
    //}
    if p.curToken.Type != token.COLON && !p.expectPeek(token.COLON) {
        return nil
    }

    fmt.Println(p.curToken, p.peekToken)

    if p.peekTokenIs(token.NEWL) {
        p.nextToken()
        expression.Consequence = p.parseBlockStatement()
    } else {
        expression.Consequence = p.parseInlineStatement()
    }

//    p.nextToken()

    if p.curToken.Type == token.ELSE {
        if p.curToken.Type != token.COLON && !p.expectPeek(token.COLON) {
            return nil
        }

        if p.peekTokenIs(token.NEWL) {
            p.nextToken()
            expression.Alternative = p.parseBlockStatement()
        } else {
            expression.Alternative = p.parseInlineStatement()
        }
    }

    return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
    block := &ast.BlockStatement{Token: p.curToken}
    block.Statements = []ast.Statement{}
    p.nextToken()
    currDepth := p.Depth

    for !p.tokenIs(token.EOF) && p.Depth >= currDepth {
        statement := p.parseStatement()
        if statement != nil {
            block.Statements = append(block.Statements, statement)
        }

        //if p.curToken.Type == token.NEWL && p.Depth < currDepth {
          //  break
        //}
        p.nextToken()
        fmt.Println(p.curToken)
    }

  
    return block
}

func (p *Parser) parseInlineStatement() *ast.BlockStatement {
    inline := &ast.BlockStatement{Token: p.curToken}
    inline.Statements = []ast.Statement{}

    for !p.tokenIs(token.EOF) && !p.tokenIs(token.NEWL) {
        statement := p.parseStatement()
        if statement != nil {
            inline.Statements = append(inline.Statements, statement)
        }

        p.nextToken()
    }

    return inline
}

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
    statement := &ast.FunctionStatement{Token: p.curToken}
    
    if !p.expectPeek(token.NAME) {
        return nil
    }

    statement.Name = p.parseName().(*ast.Name)
    
    if !p.expectPeek(token.LPAR) {
        return nil
    }

    statement.Arguments = p.parseFunctionArguments()

    if !p.expectPeek(token.COLON) {
        return nil
    }

    p.nextToken()
    
    if p.tokenIs(token.NEWL) {
        statement.Body = p.parseBlockStatement()
    } else {
        statement.Body = p.parseInlineStatement()
    }

    return statement
}

func (p *Parser) parseFunctionArguments() []*ast.Name {
    args := []*ast.Name{}

    if p.peekTokenIs(token.RPAR) {
        p.nextToken()
        return args
    }

    p.nextToken()

    name := &ast.Name{Token: p.curToken, Value: p.curToken.Literal}
    args = append(args, name)

    for p.peekTokenIs(token.COMMA) {
        p.nextToken()
        p.nextToken()
        name := &ast.Name{Token: p.curToken, Value: p.curToken.Literal}
        args = append(args, name)
    }

    if !p.expectPeek(token.RPAR) {
        return nil
    }

    return args
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
    call := &ast.CallExpression{Token: p.curToken, Function: function}
    call.Arguments = p.parseCallArguments()

    if p.peekToken.Type == token.NEWL {
        p.nextToken()
    }

    return call
}

func (p *Parser) parseCallArguments() []ast.Expression {
    args := []ast.Expression{}

    if p.peekTokenIs(token.RPAR) {
        p.nextToken()
        return args
    }

    p.nextToken()

    arg := p.parseExpression(LOWEST)
    args = append(args, arg)

    for p.peekTokenIs(token.COMMA) {
        p.nextToken()
        p.nextToken()

        arg := p.parseExpression(LOWEST)
        args = append(args, arg)
    }

    if !p.expectPeek(token.RPAR) {
        return nil
    }

    return args
}

func (p *Parser) parseListExpression() ast.Expression {
    p.nextToken()

    if p.peekTokenIs(token.COMMA) {
        arr := []ast.Expression{p.parseExpression(LOWEST)}

        for p.peekTokenIs(token.COMMA){
            p.nextToken()
            p.nextToken()
            arr = append(arr, p.parseExpression(LOWEST))
        }

        if !p.expectPeek(token.RBR) {
            return nil
        }

        return &ast.ListLiteral{Arr: arr}
    }

    return nil
}

func (p *Parser) parseIndexExpression(sequence ast.Expression) ast.Expression {
    expression := &ast.IndexExpression{
        Token: p.curToken,
        Struct: sequence,
    }
    p.nextToken()

    expression.Value = p.parseExpression(LOWEST)

    if !p.expectPeek(token.RBR) {
        return nil
    }

    return expression
}

func (p *Parser) tokenIs (t token.TokenType) bool {
    return p.curToken.Type == t
}

func (p *Parser) precedenceIs() int {
    if p, ok := precedenceMap[p.curToken.Type]; ok {
        return p
    }

    return LOWEST
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

func (p *Parser) peekPrecedence() int {
    if p, ok := precedenceMap[p.peekToken.Type]; ok {
        return p
    }

    return LOWEST
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

