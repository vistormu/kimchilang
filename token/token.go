package token

import "strings"

type Token struct {
    Type int
    Subtype int
    Literal string
}

const (
    _ int = iota
    ILLEGAL
    EOF
    
    // Main types
    KEYWORD
    IDENTIFIER
    TYPE
    LITERAL
    OPERATOR
    DELIMITER

    // Subtypes
    // Statements
    LET
    RETURN
    BE

    // Primitive types
    I64
    F64
    STR
    BOOL
    TRUE
    FALSE
    NONE

    // Complex types
    FN
    STRUCT
    ENUM
    MAP
    LIST
    TUPLE
    VEC
    SET

    // Conditionals
    IF
    ELSE
    MATCH

    // Loops
    FOR
    WHILE
    CONTINUE
    BREAK

    // Other
    PASS
    IN
    SELF

    // Operators
    ASSIGN
    PLUS
    MINUS
    ASTERISK
    SLASH
    PERCENT
    LT
    GT
    LTE
    GTE
    NOT
    IS
    IS_NOT
    AND
    OR

    // Delimiters
    COLON
    COMMA
    DOT
    LPAREN
    RPAREN
    LBRACE
    RBRACE
    UNDERSCORE
)

var keywords = map[string]Token {
    // Statements
    "let": {KEYWORD, LET, "let"},
    "return": {KEYWORD, RETURN, "return"},
    "be": {KEYWORD, BE, "be"},
    
    // Primitive types
    "i64": {TYPE, I64, "i64"},
    "f64": {TYPE, F64, "f64"},
    "str": {TYPE, STR, "str"},
    "bool": {TYPE, BOOL, "bool"},
    "none": {TYPE, NONE, "none"},
    "true": {LITERAL, TRUE, "true"},
    "false": {LITERAL, FALSE, "false"},

    // Complex types
    "fn": {TYPE, FN, "fn"},
    "struct": {TYPE, STRUCT, "struct"},
    "enum": {TYPE, ENUM, "enum"},
    "map": {TYPE, MAP, "map"},
    "list": {TYPE, LIST, "list"},
    "tuple": {TYPE, TUPLE, "tuple"},
    "vec": {TYPE, VEC, "vec"},
    "set": {TYPE, SET, "set"},

    // Conditionals
    "if": {KEYWORD, IF, "if"},
    "else": {KEYWORD, ELSE, "else"},
    "match": {KEYWORD, MATCH, "match"},

    // Loops
    "for": {KEYWORD, FOR, "for"},
    "while": {KEYWORD, WHILE, "while"},
    "continue": {KEYWORD, CONTINUE, "continue"},
    "break": {KEYWORD, BREAK, "break"},

    // Other
    "pass": {KEYWORD, PASS, "pass"},
    "in": {KEYWORD, IN, "in"},
    "self": {KEYWORD, SELF, "self"},

    // Operators
    "and": {OPERATOR, AND, "and"},
    "or": {OPERATOR, OR, "or"},
    "not": {OPERATOR, NOT, "not"},
    "is": {OPERATOR, IS, "is"},
    "is_not": {OPERATOR, IS_NOT, "is_not"},
}

var chars = map[byte]Token {
    0: {EOF, EOF, "EOF"},

    // Operators
    '=': {OPERATOR, ASSIGN, "="},
    '+': {OPERATOR, PLUS, "+"},
    '-': {OPERATOR, MINUS, "-"},
    '*': {OPERATOR, ASTERISK, "*"},
    '/': {OPERATOR, SLASH, "/"},
    '%': {OPERATOR, PERCENT, "%"},
    '<': {OPERATOR, LT, "<"},
    '>': {OPERATOR, GT, ">"},

    // Delimiters
    ':': {DELIMITER, COLON, ":"},
    ',': {DELIMITER, COMMA, ","},
    '.': {DELIMITER, DOT, "."},
    '(': {DELIMITER, LPAREN, "("},
    ')': {DELIMITER, RPAREN, ")"},
    '{': {DELIMITER, LBRACE, "{"},
    '}': {DELIMITER, RBRACE, "}"},
    '_': {DELIMITER, UNDERSCORE, "_"},
}

var twoChars = map[string]Token {
    "<=": {OPERATOR, LTE, "<="},
    ">=": {OPERATOR, GTE, ">="},
}

// ==============
// PUBLIC METHODS
// ==============
func NewIdentifier(identifier string) Token {
    if keyword, ok := keywords[identifier]; ok {
        return keyword
    }

    return Token{IDENTIFIER, IDENTIFIER, identifier}
}

func NewNumber(number string) Token {
    if strings.Contains(number, ".") {
        return Token{LITERAL, F64, number}
    }

    return Token{LITERAL, I64, number}
}

func NewString(str string) Token {
    return Token{LITERAL, STR, str}
}

func NewChar(char byte) Token {
    if token, ok := chars[char] ; ok {
        return token
    }

    return Token{ILLEGAL, ILLEGAL, string(char)}
}

func NewTwoChar(char1 byte, char2 byte) Token {
    if token, ok := twoChars[string(char1) + string(char2)] ; ok {
        return token
    }

    return Token{ILLEGAL, ILLEGAL, string(char1) + string(char2)}
}

func NewFromType(tokenType int) Token {
    switch tokenType {
    case I64:
        return Token{TYPE, I64, "i64"}
    case F64:
        return Token{TYPE, F64, "f64"}
    case STR:
        return Token{TYPE, STR, "str"}
    case TRUE:
        return Token{TYPE, BOOL, "bool"}
    case FALSE:
        return Token{TYPE, BOOL, "bool"}
    default:
        return Token{ILLEGAL, ILLEGAL, "ILLEGAL"}
    }
}
















