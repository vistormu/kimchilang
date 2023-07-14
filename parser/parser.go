package parser

import (
    "fmt"
    "strconv"
    "kimchi/ast"
    "kimchi/tokenizer"
    "kimchi/token"
)

// =====
// TYPES
// =====
type (
    prefixParseFn func() ast.Expression
    infixParseFn func(ast.Expression) ast.Expression
)

const (
    _ int = iota
    LOWEST
    INDEX
    AND
    EQUALS
    LESSGREATER
    SUM
    PRODUCT
    PREFIX
    CALL
)

var precedences = map[int]int {
    token.DOT: INDEX,
    token.AND: AND,
    token.OR: AND,
    token.IS: EQUALS,
    token.IS_NOT: EQUALS,
    token.LT: LESSGREATER,
    token.GT: LESSGREATER,
    token.LTE: LESSGREATER,
    token.GTE: LESSGREATER,
    token.PLUS: SUM,
    token.MINUS: SUM,
    token.SLASH: PRODUCT,
    token.ASTERISK: PRODUCT,
    token.LPAREN: CALL,
}

type Parser struct {
    tokenizer *tokenizer.Tokenizer

    currentToken token.Token
    peekToken token.Token

    prefixParseFns map[int]prefixParseFn
    infixParseFns map[int]infixParseFn

    Errors []string
}

// ==============
// PUBLIC METHODS
// ==============
func New(tokenizer *tokenizer.Tokenizer) *Parser {
    parser := &Parser{tokenizer: tokenizer, Errors: []string{}}

    parser.nextToken()
    parser.nextToken()

    parser.prefixParseFns = make(map[int]prefixParseFn)
    parser.prefixParseFns[token.IDENTIFIER] = parser.parseIdentifier
    parser.prefixParseFns[token.I64] = parser.parseIntegerLiteral
    parser.prefixParseFns[token.F64] = parser.parseFloatLiteral
    parser.prefixParseFns[token.STR] = parser.parseStringLiteral
    parser.prefixParseFns[token.TRUE] = parser.parseBooleanLiteral
    parser.prefixParseFns[token.FALSE] = parser.parseBooleanLiteral
    parser.prefixParseFns[token.NOT] = parser.parsePrefixExpression
    parser.prefixParseFns[token.MINUS] = parser.parsePrefixExpression
    parser.prefixParseFns[token.LPAREN] = parser.parseGroupedExpression
    parser.prefixParseFns[token.IF] = parser.parseIfExpression
    parser.prefixParseFns[token.FN] = parser.parseFunctionLiteral
    parser.prefixParseFns[token.LIST] = parser.parseListLiteral

    parser.infixParseFns = make(map[int]infixParseFn)
    parser.infixParseFns[token.PLUS] = parser.parseInfixExpression
    parser.infixParseFns[token.MINUS] = parser.parseInfixExpression
    parser.infixParseFns[token.SLASH] = parser.parseInfixExpression
    parser.infixParseFns[token.ASTERISK] = parser.parseInfixExpression
    parser.infixParseFns[token.LT] = parser.parseInfixExpression
    parser.infixParseFns[token.GT] = parser.parseInfixExpression
    parser.infixParseFns[token.LTE] = parser.parseInfixExpression
    parser.infixParseFns[token.GTE] = parser.parseInfixExpression
    parser.infixParseFns[token.IS] = parser.parseInfixExpression
    parser.infixParseFns[token.IS_NOT] = parser.parseInfixExpression
    parser.infixParseFns[token.AND] = parser.parseInfixExpression
    parser.infixParseFns[token.OR] = parser.parseInfixExpression
    parser.infixParseFns[token.LPAREN] = parser.parseCallExpression

    return parser
}

func (self *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{Statements: []ast.Statement{}}

    for !self.currentTokenIs(token.EOF) {
        if self.currentTokenIs(token.ILLEGAL) {
            self.addParseError(self.currentToken)
        }

        statement := self.parseStatement()
        if statement != nil {
            program.Statements = append(program.Statements, statement)
        }
        self.nextToken()
    }

    return program
}

// ===============
// PRIVATE METHODS
// ===============
func (self *Parser) nextToken() {
    self.currentToken = self.peekToken
    self.peekToken = self.tokenizer.GetToken()
}
func (self *Parser) statementIsTerminated() bool {
    if self.peekTokenIs(token.KEYWORD) || self.peekTokenIs(token.EOF) {
        return true
    }
    if self.peekTokenIs(token.IDENTIFIER) && !self.currentTokenIs(token.OPERATOR) {
        return true
    }
    return false
}

// ======
// TOKENS
// ======
func (self *Parser) currentTokenIs(tokenType int) bool {
    return self.currentToken.Type == tokenType || self.currentToken.Subtype == tokenType
}
func (self *Parser) peekTokenIs(tokenType int) bool {
    return self.peekToken.Type == tokenType || self.peekToken.Subtype == tokenType
}
func (self *Parser) expectPeekTokenToBe(tokenType int ) bool {
    if self.peekTokenIs(tokenType) {
        self.nextToken()
        return true
    }
    self.addPeekError(self.peekToken)
    return false
}
func (self *Parser) peekPrecedence() int {
    if precedence, ok := precedences[self.peekToken.Subtype]; ok {
        return precedence
    }
    return LOWEST
}

// ======
// ERRORS
// ======
func (self *Parser) addPeekError(token token.Token) {
    message := fmt.Sprintf("expected next token to be %s, got %s instead", token.Literal, self.peekToken.Literal)
    self.Errors = append(self.Errors, message)
}
func (self *Parser) addNoPrefixParseFnError(token token.Token) {
    message := fmt.Sprintf("no prefix parse function for %s found", token.Literal)
    self.Errors = append(self.Errors, message)
}
func (self *Parser) addParseError(token token.Token) {
    message := fmt.Sprintf("could not parse %s token", token.Literal)
    self.Errors = append(self.Errors, message)
}

// =======
// PARSING
// =======
func (self *Parser) parseStatement() ast.Statement {
    switch self.currentToken.Subtype {
    case token.LET:
        return self.parseLetStatement()
    case token.RETURN:
        return self.parseReturnStatement()
    default:
        return self.parseExpressionStatement()
    }
}
func (self *Parser) parseExpression(precedence int) ast.Expression {
    prefixFunction, ok := self.prefixParseFns[self.currentToken.Subtype]
    if !ok {
        self.addNoPrefixParseFnError(self.currentToken)
        return nil
    }
    leftExpression := prefixFunction()

    for !self.statementIsTerminated() && precedence < self.peekPrecedence() {
        infixFunction, ok := self.infixParseFns[self.peekToken.Subtype]
        if !ok {
            return leftExpression
        }
        self.nextToken()
        leftExpression = infixFunction(leftExpression)
    }

    return leftExpression
}

// ==========
// STATEMENTS
// ==========
func (self *Parser) parseLetStatement() *ast.LetStatement {
    statement := &ast.LetStatement{}

    if !self.expectPeekTokenToBe(token.IDENTIFIER) { return nil }
    statement.Identifier = self.parseIdentifier().(*ast.Identifier)

    if self.peekTokenIs(token.BE) {
        return self.parseLetBeStatement(statement)
    }

    if !self.expectPeekTokenToBe(token.COLON) { return nil }

    if !self.expectPeekTokenToBe(token.TYPE) { return nil }
    statement.Identifier.Type = self.currentToken // TODO: parse complex types

    if !self.expectPeekTokenToBe(token.ASSIGN) { return nil }
    self.nextToken()

    statement.Expression = self.parseExpression(LOWEST)

    return statement
}
func (self *Parser) parseLetBeStatement(statement *ast.LetStatement) *ast.LetStatement {
    self.nextToken()

    if !self.expectPeekTokenToBe(token.LITERAL) { return nil }
    statement.Identifier.Type = token.NewFromType(self.currentToken.Subtype)
    statement.Expression = self.parseExpression(LOWEST)

    return statement
}
func (self *Parser) parseReturnStatement() *ast.ReturnStatement {
    statement := &ast.ReturnStatement{}
    self.nextToken()
    statement.Expression = self.parseExpression(LOWEST)

    return statement
}
func (self *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    statement := &ast.ExpressionStatement{}
    statement.Expression = self.parseExpression(LOWEST)

    return statement
}

// ===========
// EXPRESSIONS
// ===========
func (self *Parser) parsePrefixExpression() ast.Expression {
    expression := &ast.PrefixExpression{
        Operator: self.currentToken.Literal,
    }

    self.nextToken()
    expression.Right = self.parseExpression(PREFIX)

    return expression
}
func (self *Parser) parseInfixExpression(leftExpression ast.Expression) ast.Expression {
    if self.currentTokenIs(token.IS) && self.peekTokenIs(token.NOT) {
        self.nextToken()
        self.currentToken = token.NewIdentifier("is_not")
    }

    expression := &ast.InfixExpression{
        Operator: self.currentToken.Literal,
        Left: leftExpression,
    }

    precedende := LOWEST
    if p, ok := precedences[self.currentToken.Subtype]; ok {
        precedende = p
    }
    self.nextToken()

    expression.Right = self.parseExpression(precedende)

    return expression
}
func (self *Parser) parseGroupedExpression() ast.Expression {
    self.nextToken()

    expression := self.parseExpression(LOWEST)

    if !self.expectPeekTokenToBe(token.RPAREN) { return nil }

    return expression
}
func (self *Parser) parseBlockStatement() *ast.BlockStatement {
    block := &ast.BlockStatement{Statements: []ast.Statement{}}
    self.nextToken()

    for !self.currentTokenIs(token.RBRACE) && !self.currentTokenIs(token.EOF) {
        statement := self.parseStatement()
        if statement != nil {
            block.Statements = append(block.Statements, statement)
        }
        self.nextToken()
    }

    return block
}
func (self *Parser) parseExpressionList() []ast.Expression {
    list := []ast.Expression{}

    if self.peekTokenIs(token.RPAREN) {
        self.nextToken()
        return list
    }

    self.nextToken()
    list = append(list, self.parseExpression(LOWEST))

    for self.peekTokenIs(token.COMMA) {
        self.nextToken()
        self.nextToken()
        list = append(list, self.parseExpression(LOWEST))
    }

    if !self.expectPeekTokenToBe(token.RPAREN) { return nil }

    return list
}

// ========
// KEYWORDS
// ========
func (self *Parser) parseIfExpression() ast.Expression {
    expression := &ast.IfExpression{}
    self.nextToken()

    expression.Condition = self.parseExpression(LOWEST)

    if !self.expectPeekTokenToBe(token.LBRACE) { return nil }

    expression.Consequence = self.parseBlockStatement()

    if self.peekTokenIs(token.ELSE) {
        self.nextToken()
        if !self.expectPeekTokenToBe(token.LBRACE) { return nil }
        expression.Alternative = self.parseBlockStatement()
    }

    return expression
}
func (self *Parser) parseFunctionLiteral() ast.Expression {
    literal := &ast.FunctionLiteral{}

    if !self.expectPeekTokenToBe(token.LPAREN) { return nil }

    literal.Parameters = self.parseFunctionParameters()

    if !self.expectPeekTokenToBe(token.COLON) { return nil }

    if !self.expectPeekTokenToBe(token.TYPE) { return nil }
    literal.ReturnType = self.currentToken

    if !self.expectPeekTokenToBe(token.LBRACE) { return nil }

    literal.Body = self.parseBlockStatement()

    return literal
}
func (self *Parser) parseFunctionParameters() []*ast.Identifier {
    identifiers := []*ast.Identifier{}

    if self.peekTokenIs(token.RPAREN) {
        self.nextToken()
        return identifiers
    }

    self.nextToken()
    identifier := &ast.Identifier{Name: self.currentToken.Literal}

    if !self.expectPeekTokenToBe(token.COLON) { return nil }

    if !self.expectPeekTokenToBe(token.TYPE) { return nil }
    identifier.Type = self.currentToken
    identifiers = append(identifiers, identifier)

    for self.peekTokenIs(token.COMMA) {
        self.nextToken()
        self.nextToken()
        identifier := &ast.Identifier{Name: self.currentToken.Literal}
        
        if !self.expectPeekTokenToBe(token.COLON) { return nil }

        if !self.expectPeekTokenToBe(token.TYPE) { return nil }
        identifier.Type = self.currentToken
        identifiers = append(identifiers, identifier)
    }

    if !self.expectPeekTokenToBe(token.RPAREN) { return nil }

    return identifiers
}
func (self *Parser) parseCallExpression(function ast.Expression) ast.Expression {
    expression := &ast.CallExpression{Function: function}
    expression.Arguments = self.parseExpressionList()
    return expression
}

// ========
// LITERALS
// ========
func (self *Parser) parseIdentifier() ast.Expression {
    return &ast.Identifier{Name: self.currentToken.Literal}
}
func (self *Parser) parseIntegerLiteral() ast.Expression {
    literal := &ast.IntegerLiteral{}

    value, err := strconv.ParseInt(self.currentToken.Literal, 10, 64)
    if err != nil {
        self.addParseError(self.currentToken)
        return nil
    }

    literal.Value = value
    return literal
}
func (self *Parser) parseFloatLiteral() ast.Expression {
    literal := &ast.FloatLiteral{}

    value, err := strconv.ParseFloat(self.currentToken.Literal, 64)
    if err != nil {
        self.addParseError(self.currentToken)
        return nil
    }

    literal.Value = value
    return literal
}
func (self *Parser) parseStringLiteral() ast.Expression {
    return &ast.StringLiteral{Value: self.currentToken.Literal}
}
func (self *Parser) parseBooleanLiteral() ast.Expression {
    return &ast.BooleanLiteral{Value: self.currentTokenIs(token.TRUE)}
}

// ======
// ARRAYS
// ======
func (self *Parser) parseListLiteral() ast.Expression {
    list := &ast.ListLiteral{}

    if !self.expectPeekTokenToBe(token.LPAREN) { return nil }

    list.Elements = self.parseExpressionList()
    return list
}
// func (self *Parser) parseIndexExpression(leftExpression ast.Expression) ast.Expression {
//     expression := &ast.CallExpression{Function: leftExpression}
//     self.nextToken()

//     expression.Arguments = self.parseExpressionList()

//     return expression
// }
