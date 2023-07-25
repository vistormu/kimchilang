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
    AND
    EQUALS
    LESSGREATER
    SUM
    PRODUCT
    PREFIX
    CALL
)

var precedences = map[int]int {
    token.AND: AND,
    token.OR: AND,
    token.TO: AND,
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
    token.DOT: CALL,
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
    parser.prefixParseFns[token.MAP] = parser.parseMapLiteral
    parser.prefixParseFns[token.WHILE] = parser.parseWhileExpression
    parser.prefixParseFns[token.FOR] = parser.parseForExpression

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
    parser.infixParseFns[token.DOT] = parser.parseDotExpression
    parser.infixParseFns[token.TO] = parser.parseInfixExpression

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
    if (self.peekTokenIs(token.KEYWORD) && !self.peekTokenIs(token.TO)) || self.peekTokenIs(token.EOF) {
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
    self.addPeekError(tokenType)
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
func (self *Parser) addPeekError(tokenType int) {
    message := fmt.Sprintf("expected next token to be %d, got %s instead", tokenType, self.peekToken.Literal)
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
    case token.MUT:
        return self.parseMutStatement()
    case token.EXE:
        return self.parseExeStatement()
    case token.BREAK:
        return self.parseBreakStatement()
    case token.CONTINUE:
        return self.parseContinueStatement()
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
    statement.Identifier.Type = self.parseTypeLiteral()

    if !self.expectPeekTokenToBe(token.ASSIGN) { return nil }
    self.nextToken()

    statement.Expression = self.parseExpression(LOWEST)

    return statement
}
func (self *Parser) parseLetBeStatement(statement *ast.LetStatement) *ast.LetStatement {
    self.nextToken()

    if !self.expectPeekTokenToBe(token.LITERAL) { return nil }
    statement.Identifier.Type = self.parseTypeLiteral()
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
func (self *Parser) parseMutStatement() *ast.MutStatement {
    statement := &ast.MutStatement{}

    if !self.expectPeekTokenToBe(token.IDENTIFIER) { return nil }
    statement.Identifier = self.parseIdentifier()

    if self.peekTokenIs(token.LPAREN) {
        self.nextToken()
        statement.Identifier = self.parseCallExpression(statement.Identifier)
    }

    if !self.expectPeekTokenToBe(token.TO) { return nil }

    if self.peekTokenIs(token.OPERATOR) || self.peekTokenIs(token.DOT) {
        if ident, ok := statement.Identifier.(*ast.Identifier); ok {
            self.currentToken = token.NewIdentifier(ident.Name)
        }
    } else {
        self.nextToken()
    }

    statement.Expression = self.parseExpression(LOWEST)

    return statement
}
func (self *Parser) parseExeStatement() *ast.ExeStatement {
    statement := &ast.ExeStatement{}

    if !self.expectPeekTokenToBe(token.IDENTIFIER) { return nil }
    statement.Function = self.parseIdentifier().(*ast.Identifier)

    if self.peekTokenIs(token.LPAREN) {
        self.nextToken()
        statement.Arguments = self.parseExpressionList()
    }

    return statement
}
func (self *Parser) parseBreakStatement() ast.Statement {
    statement := &ast.BreakStatement{}

    if self.peekTokenIs(token.IF) {
        self.nextToken()
        self.nextToken()

        statement.Condition = self.parseExpression(LOWEST)
    }

    return statement
}
func (self *Parser) parseContinueStatement() ast.Statement {
    statement := &ast.ContinueStatement{}
    
    if self.peekTokenIs(token.IF) {
        self.nextToken()
        self.nextToken()

        statement.Condition = self.parseExpression(LOWEST)
    }

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

        if self.peekTokenIs(token.IF) {
            self.nextToken()
            expression.Alternative = &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: self.parseIfExpression()}}}

            return expression
        }

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
    literal.ReturnType = self.parseTypeLiteral()

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
    identifier.Type = self.parseTypeLiteral()
    identifiers = append(identifiers, identifier)

    for self.peekTokenIs(token.COMMA) {
        self.nextToken()
        self.nextToken()
        identifier := &ast.Identifier{Name: self.currentToken.Literal}
        
        if !self.expectPeekTokenToBe(token.COLON) { return nil }

        if !self.expectPeekTokenToBe(token.TYPE) { return nil }
        identifier.Type = self.parseTypeLiteral()
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
func (self *Parser) parseDotExpression(leftExpression ast.Expression) ast.Expression {
    if !self.expectPeekTokenToBe(token.IDENTIFIER) { return nil }

    return self.parseMethodExpression(leftExpression)
    
    // switch self.peekToken.Subtype {
    // case token.LPAREN:
    //     return self.parseMethodExpression(leftExpression)
    // default:
        // return self.parseAttributeExpression(leftExpression)
    // }
}
func (self *Parser) parseMethodExpression(leftExpression ast.Expression) ast.Expression {
    expression := &ast.MethodExpression{Left: leftExpression}
    expression.Method = self.parseIdentifier()

    if self.peekTokenIs(token.LPAREN) {
        self.nextToken()
        expression.Arguments = self.parseExpressionList()
    }

    return expression
}
// func (self *Parser) parseAttributeExpression(leftExpression ast.Expression) ast.Expression {
//     expression := &ast.AttributeExpression{Left: leftExpression}
//     expression.Attribute = self.parseIdentifier()
//     return expression
// }

// ========
// LITERALS
// ========
func (self *Parser) parseIdentifier() ast.Expression {
    return &ast.Identifier{Name: self.currentToken.Literal}
}
func (self *Parser) parseTypeLiteral() *ast.TypeLiteral {
    if self.currentTokenIs(token.LITERAL) {
        return &ast.TypeLiteral{Type: token.NewFromType(self.currentToken.Subtype)}
    }
    typeLiteral := &ast.TypeLiteral{Type: self.currentToken}

    switch self.currentToken.Subtype {
    case token.LIST:
        if !self.expectPeekTokenToBe(token.LPAREN) { return nil }
        if !self.expectPeekTokenToBe(token.TYPE) { return nil }
        typeLiteral.Subtypes = append(typeLiteral.Subtypes, self.currentToken)
        if !self.expectPeekTokenToBe(token.RPAREN) { return nil }

    case token.MAP:
        if !self.expectPeekTokenToBe(token.LPAREN) { return nil }
        if !self.expectPeekTokenToBe(token.TYPE) { return nil }
        typeLiteral.Subtypes = append(typeLiteral.Subtypes, self.currentToken)
        if !self.expectPeekTokenToBe(token.COMMA) { return nil }
        if !self.expectPeekTokenToBe(token.TYPE) { return nil }
        typeLiteral.Subtypes = append(typeLiteral.Subtypes, self.currentToken)
        if !self.expectPeekTokenToBe(token.RPAREN) { return nil }
    }

    return typeLiteral
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

// ===========
// COLLECTIONS
// ===========
func (self *Parser) parseMapLiteral() ast.Expression {
    mapLiteral := &ast.MapLiteral{Pairs: make(map[ast.Expression]ast.Expression)}
    self.nextToken()

    for !self.peekTokenIs(token.RPAREN) {
        self.nextToken()
        key := self.parseExpression(LOWEST)

        if !self.expectPeekTokenToBe(token.COLON) { return nil }
        self.nextToken()

        value := self.parseExpression(LOWEST)
        mapLiteral.Pairs[key] = value

        if !self.peekTokenIs(token.RPAREN) && !self.expectPeekTokenToBe(token.COMMA) { return nil }
    }

    if !self.expectPeekTokenToBe(token.RPAREN) { return nil }

    return mapLiteral
}


// =====
// LOOPS
// =====
func (self *Parser) parseWhileExpression() ast.Expression {
    expression := &ast.WhileExpression{}
    self.nextToken()

    expression.Condition = self.parseExpression(LOWEST)

    if !self.expectPeekTokenToBe(token.LBRACE) { return nil }

    expression.Body = self.parseBlockStatement()
    
    return expression
}
func (self *Parser) parseForExpression() ast.Expression {
    expression := &ast.ForExpression{}

    if !self.peekTokenIs(token.IDENTIFIER) && !self.peekTokenIs(token.UNDERSCORE) { 
        self.addPeekError(self.peekToken.Subtype)
        return nil 
    }
    self.nextToken()
    expression.Index = self.parseIdentifier().(*ast.Identifier)

    if !self.expectPeekTokenToBe(token.COMMA) { return nil } 

    if !self.peekTokenIs(token.IDENTIFIER) && !self.peekTokenIs(token.UNDERSCORE) { 
        self.addPeekError(self.peekToken.Subtype)
        return nil 
    }
    self.nextToken()
    expression.Value = self.parseIdentifier().(*ast.Identifier)

    if !self.expectPeekTokenToBe(token.IN) { return nil }
    self.nextToken()

    expression.Iterable = self.parseExpression(LOWEST)

    if !self.expectPeekTokenToBe(token.LBRACE) { return nil }

    expression.Body = self.parseBlockStatement()

    return expression
}
