package tokenizer

import (
    "testing"
    "kimchi/token"
)

func TestComments(t *testing.T) {
    input := `
    # This is the first comment
    # This is another comment 
    # This is a third comment 
    let foo: i64 = 5 # This is a comment after a statement
    `

    tests := []struct {
        expectedType int
        expectedSubtype int
        expectedLiteral string
    }{
        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "foo"},
        {token.DELIMITER, token.COLON, ":"},
        {token.TYPE, token.I64, "i64"},
        {token.OPERATOR, token.ASSIGN, "="},
        {token.LITERAL, token.I64, "5"},
    }

    runTest(t, input, tests)
}

func TestLetStatements(t *testing.T) {
    input := `
    let foo: i64 = 5
    let pi: f64 = 3.14
    let done: bool = false
    let message: str = "Hello, World!"
    let my_list: list(i64) = list(1, 2, 3, 4, 5)

    let foo_2 be 5
    let pi_2 be 3.14
    let done_2 be false
    let message_2 be "Hello, World!"
    let array_2 be list(1, 2, 3, 4, 5)

    let content: str = read("tests/input.txt")
    `
    tests := []struct {
        expectedType int
        expectedSubtype int
        expectedLiteral string
    }{
        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "foo"},
        {token.DELIMITER, token.COLON, ":"},
        {token.TYPE, token.I64, "i64"},
        {token.OPERATOR, token.ASSIGN, "="},
        {token.LITERAL, token.I64, "5"},

        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "pi"},
        {token.DELIMITER, token.COLON, ":"},
        {token.TYPE, token.F64, "f64"},
        {token.OPERATOR, token.ASSIGN, "="},
        {token.LITERAL, token.F64, "3.14"},

        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "done"},
        {token.DELIMITER, token.COLON, ":"},
        {token.TYPE, token.BOOL, "bool"},
        {token.OPERATOR, token.ASSIGN, "="},
        {token.LITERAL, token.FALSE, "false"},

        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "message"},
        {token.DELIMITER, token.COLON, ":"},
        {token.TYPE, token.STR, "str"},
        {token.OPERATOR, token.ASSIGN, "="},
        {token.LITERAL, token.STR, "Hello, World!"},

        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "my_list"},
        {token.DELIMITER, token.COLON, ":"},
        {token.TYPE, token.LIST, "list"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.TYPE, token.I64, "i64"},
        {token.DELIMITER, token.RPAREN, ")"},
        {token.OPERATOR, token.ASSIGN, "="},
        {token.TYPE, token.LIST, "list"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.LITERAL, token.I64, "1"},
        {token.DELIMITER, token.COMMA, ","},
        {token.LITERAL, token.I64, "2"},
        {token.DELIMITER, token.COMMA, ","},
        {token.LITERAL, token.I64, "3"},
        {token.DELIMITER, token.COMMA, ","},
        {token.LITERAL, token.I64, "4"},
        {token.DELIMITER, token.COMMA, ","},
        {token.LITERAL, token.I64, "5"},
        {token.DELIMITER, token.RPAREN, ")"},

        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "foo_2"},
        {token.KEYWORD, token.BE, "be"},
        {token.LITERAL, token.I64, "5"},

        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "pi_2"},
        {token.KEYWORD, token.BE, "be"},
        {token.LITERAL, token.F64, "3.14"},

        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "done_2"},
        {token.KEYWORD, token.BE, "be"},
        {token.LITERAL, token.FALSE, "false"},

        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "message_2"},
        {token.KEYWORD, token.BE, "be"},
        {token.LITERAL, token.STR, "Hello, World!"},

        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "array_2"},
        {token.KEYWORD, token.BE, "be"},
        {token.TYPE, token.LIST, "list"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.LITERAL, token.I64, "1"},
        {token.DELIMITER, token.COMMA, ","},
        {token.LITERAL, token.I64, "2"},
        {token.DELIMITER, token.COMMA, ","},
        {token.LITERAL, token.I64, "3"},
        {token.DELIMITER, token.COMMA, ","},
        {token.LITERAL, token.I64, "4"},
        {token.DELIMITER, token.COMMA, ","},
        {token.LITERAL, token.I64, "5"},
        {token.DELIMITER, token.RPAREN, ")"},

        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "content"},
        {token.DELIMITER, token.COLON, ":"},
        {token.TYPE, token.STR, "str"},
        {token.OPERATOR, token.ASSIGN, "="},
        {token.IDENTIFIER, token.IDENTIFIER, "read"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.LITERAL, token.STR, "tests/input.txt"},
        {token.DELIMITER, token.RPAREN, ")"},

        {token.EOF, token.EOF, "EOF"},
    }

    runTest(t, input, tests)
}

func TestFunctions(t *testing.T) {
    input := `
    let add be fn(x: i64, y: i64): i64 {return x + y}

    let result: i64 = add(x, y)
    `
    tests := []struct {
        expectedType int
        expectedSubtype int
        expectedLiteral string
    }{
        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "add"},
        {token.KEYWORD, token.BE, "be"},
        {token.TYPE, token.FN, "fn"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.IDENTIFIER, token.IDENTIFIER, "x"},
        {token.DELIMITER, token.COLON, ":"},
        {token.TYPE, token.I64, "i64"},
        {token.DELIMITER, token.COMMA, ","},
        {token.IDENTIFIER, token.IDENTIFIER, "y"},
        {token.DELIMITER, token.COLON, ":"},
        {token.TYPE, token.I64, "i64"},
        {token.DELIMITER, token.RPAREN, ")"},
        {token.DELIMITER, token.COLON, ":"},
        {token.TYPE, token.I64, "i64"},
        {token.DELIMITER, token.LBRACE, "{"},
        {token.KEYWORD, token.RETURN, "return"},
        {token.IDENTIFIER, token.IDENTIFIER, "x"},
        {token.OPERATOR, token.PLUS, "+"},
        {token.IDENTIFIER, token.IDENTIFIER, "y"},
        {token.DELIMITER, token.RBRACE, "}"},

        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "result"},
        {token.DELIMITER, token.COLON, ":"},
        {token.TYPE, token.I64, "i64"},
        {token.OPERATOR, token.ASSIGN, "="},
        {token.IDENTIFIER, token.IDENTIFIER, "add"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.IDENTIFIER, token.IDENTIFIER, "x"},
        {token.DELIMITER, token.COMMA, ","},
        {token.IDENTIFIER, token.IDENTIFIER, "y"},
        {token.DELIMITER, token.RPAREN, ")"},

        {token.EOF, token.EOF, "EOF"},
    }

    runTest(t, input, tests)
}

func TestConditionalStatement(t *testing.T) {
    input := `
    if foo > 0 and foo is 5 {
        return false
    }
    else if bar is not 10 {
        return true
    }
    else {
        pass
    }

    let result: bool = if foo > 0 and foo is 5 {return false} else if bar is not 10 {return true} else {pass}
    `

    tests := []struct {
        expectedType int
        expectedSubtype int
        expectedLiteral string
    }{
        {token.KEYWORD, token.IF, "if"},
        {token.IDENTIFIER, token.IDENTIFIER, "foo"},
        {token.OPERATOR, token.GT, ">"},
        {token.LITERAL, token.I64, "0"},
        {token.OPERATOR, token.AND, "and"},
        {token.IDENTIFIER, token.IDENTIFIER, "foo"},
        {token.OPERATOR, token.IS, "is"},
        {token.LITERAL, token.I64, "5"},
        {token.DELIMITER, token.LBRACE, "{"},
        {token.KEYWORD, token.RETURN, "return"},
        {token.LITERAL, token.FALSE, "false"},
        {token.DELIMITER, token.RBRACE, "}"},
        {token.KEYWORD, token.ELSE, "else"},
        {token.KEYWORD, token.IF, "if"},
        {token.IDENTIFIER, token.IDENTIFIER, "bar"},
        {token.OPERATOR, token.IS, "is"},
        {token.OPERATOR, token.NOT, "not"},
        {token.LITERAL, token.I64, "10"},
        {token.DELIMITER, token.LBRACE, "{"},
        {token.KEYWORD, token.RETURN, "return"},
        {token.LITERAL, token.TRUE, "true"},
        {token.DELIMITER, token.RBRACE, "}"},
        {token.KEYWORD, token.ELSE, "else"},
        {token.DELIMITER, token.LBRACE, "{"},
        {token.KEYWORD, token.PASS, "pass"},
        {token.DELIMITER, token.RBRACE, "}"},

        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "result"},
        {token.DELIMITER, token.COLON, ":"},
        {token.TYPE, token.BOOL, "bool"},
        {token.OPERATOR, token.ASSIGN, "="},
        {token.KEYWORD, token.IF, "if"},
        {token.IDENTIFIER, token.IDENTIFIER, "foo"},
        {token.OPERATOR, token.GT, ">"},
        {token.LITERAL, token.I64, "0"},
        {token.OPERATOR, token.AND, "and"},
        {token.IDENTIFIER, token.IDENTIFIER, "foo"},
        {token.OPERATOR, token.IS, "is"},
        {token.LITERAL, token.I64, "5"},
        {token.DELIMITER, token.LBRACE, "{"},
        {token.KEYWORD, token.RETURN, "return"},
        {token.LITERAL, token.FALSE, "false"},
        {token.DELIMITER, token.RBRACE, "}"},
        {token.KEYWORD, token.ELSE, "else"},
        {token.KEYWORD, token.IF, "if"},
        {token.IDENTIFIER, token.IDENTIFIER, "bar"},
        {token.OPERATOR, token.IS, "is"},
        {token.OPERATOR, token.NOT, "not"},
        {token.LITERAL, token.I64, "10"},
        {token.DELIMITER, token.LBRACE, "{"},
        {token.KEYWORD, token.RETURN, "return"},
        {token.LITERAL, token.TRUE, "true"},
        {token.DELIMITER, token.RBRACE, "}"},
        {token.KEYWORD, token.ELSE, "else"},
        {token.DELIMITER, token.LBRACE, "{"},
        {token.KEYWORD, token.PASS, "pass"},
        {token.DELIMITER, token.RBRACE, "}"},

        {token.EOF, token.EOF, "EOF"},
    }

    runTest(t, input, tests)
}

func TestMatchStatement(t *testing.T) {
    input := `
    let result: str = match x {
        1: {return "one"}
        2: {return "two"}
        3: {return "three"}
        _: {return "something else"}
    }
    `

    tests := []struct {
        expectedType int
        expectedSubtype int
        expectedLiteral string
    }{
        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "result"},
        {token.DELIMITER, token.COLON, ":"},
        {token.TYPE, token.STR, "str"},
        {token.OPERATOR, token.ASSIGN, "="},
        {token.KEYWORD, token.MATCH, "match"},
        {token.IDENTIFIER, token.IDENTIFIER, "x"},
        {token.DELIMITER, token.LBRACE, "{"},
        {token.LITERAL, token.I64, "1"},
        {token.DELIMITER, token.COLON, ":"},
        {token.DELIMITER, token.LBRACE, "{"},
        {token.KEYWORD, token.RETURN, "return"},
        {token.LITERAL, token.STR, "one"},
        {token.DELIMITER, token.RBRACE, "}"},
        {token.LITERAL, token.I64, "2"},
        {token.DELIMITER, token.COLON, ":"},
        {token.DELIMITER, token.LBRACE, "{"},
        {token.KEYWORD, token.RETURN, "return"},
        {token.LITERAL, token.STR, "two"},
        {token.DELIMITER, token.RBRACE, "}"},
        {token.LITERAL, token.I64, "3"},
        {token.DELIMITER, token.COLON, ":"},
        {token.DELIMITER, token.LBRACE, "{"},
        {token.KEYWORD, token.RETURN, "return"},
        {token.LITERAL, token.STR, "three"},
        {token.DELIMITER, token.RBRACE, "}"},
        {token.DELIMITER, token.UNDERSCORE, "_"},
        {token.DELIMITER, token.COLON, ":"},
        {token.DELIMITER, token.LBRACE, "{"},
        {token.KEYWORD, token.RETURN, "return"},
        {token.LITERAL, token.STR, "something else"},
        {token.DELIMITER, token.RBRACE, "}"},
        {token.DELIMITER, token.RBRACE, "}"},

        {token.EOF, token.EOF, "EOF"},
    }

    runTest(t, input, tests)
}

func TestForStatement(t *testing.T) {
    input := `
    for i, _ in tuple(0:10) {
        print(i)
    }

    for i, value in values {
        print(i)
    }
    `

    tests := []struct {
        expectedType int
        expectedSubtype int
        expectedLiteral string
    }{
        {token.KEYWORD, token.FOR, "for"},
        {token.IDENTIFIER, token.IDENTIFIER, "i"},
        {token.DELIMITER, token.COMMA, ","},
        {token.DELIMITER, token.UNDERSCORE, "_"},
        {token.KEYWORD, token.IN, "in"},
        {token.TYPE, token.TUPLE, "tuple"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.LITERAL, token.I64, "0"},
        {token.DELIMITER, token.COLON, ":"},
        {token.LITERAL, token.I64, "10"},
        {token.DELIMITER, token.RPAREN, ")"},
        {token.DELIMITER, token.LBRACE, "{"},
        {token.IDENTIFIER, token.IDENTIFIER, "print"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.IDENTIFIER, token.IDENTIFIER, "i"},
        {token.DELIMITER, token.RPAREN, ")"},
        {token.DELIMITER, token.RBRACE, "}"},
        
        {token.KEYWORD, token.FOR, "for"},
        {token.IDENTIFIER, token.IDENTIFIER, "i"},
        {token.DELIMITER, token.COMMA, ","},
        {token.IDENTIFIER, token.IDENTIFIER, "value"},
        {token.KEYWORD, token.IN, "in"},
        {token.IDENTIFIER, token.IDENTIFIER, "values"},
        {token.DELIMITER, token.LBRACE, "{"},
        {token.IDENTIFIER, token.IDENTIFIER, "print"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.IDENTIFIER, token.IDENTIFIER, "i"},
        {token.DELIMITER, token.RPAREN, ")"},
        {token.DELIMITER, token.RBRACE, "}"},

        {token.EOF, token.EOF, "EOF"},
    }

    runTest(t, input, tests)
}

func TestWhileStatement(t *testing.T) {
    input := `
    while i < 10 {
        print(i)
    }
    `
    
    tests := []struct {
        expectedType int
        expectedSubtype int
        expectedLiteral string
    }{
        {token.KEYWORD, token.WHILE, "while"},
        {token.IDENTIFIER, token.IDENTIFIER, "i"},
        {token.OPERATOR, token.LT, "<"},
        {token.LITERAL, token.I64, "10"},
        {token.DELIMITER, token.LBRACE, "{"},
        {token.IDENTIFIER, token.IDENTIFIER, "print"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.IDENTIFIER, token.IDENTIFIER, "i"},
        {token.DELIMITER, token.RPAREN, ")"},
        {token.DELIMITER, token.RBRACE, "}"},

        {token.EOF, token.EOF, "EOF"},
    }

    runTest(t, input, tests)
}

func TestOperators(t *testing.T) {
    input := `
    5 + 5
    5 - 5
    5 * 5
    5 / 5
    5 % 5
    5 < 5
    5 > 5
    5 <= 5
    5 >= 5
    `

    tests := []struct {
        expectedType int
        expectedSubtype int
        expectedLiteral string
    }{
        {token.LITERAL, token.I64, "5"},
        {token.OPERATOR, token.PLUS, "+"},
        {token.LITERAL, token.I64, "5"},

        {token.LITERAL, token.I64, "5"},
        {token.OPERATOR, token.MINUS, "-"},
        {token.LITERAL, token.I64, "5"},

        {token.LITERAL, token.I64, "5"},
        {token.OPERATOR, token.ASTERISK, "*"},
        {token.LITERAL, token.I64, "5"},

        {token.LITERAL, token.I64, "5"},
        {token.OPERATOR, token.SLASH, "/"},
        {token.LITERAL, token.I64, "5"},

        {token.LITERAL, token.I64, "5"},
        {token.OPERATOR, token.PERCENT, "%"},
        {token.LITERAL, token.I64, "5"},

        {token.LITERAL, token.I64, "5"},
        {token.OPERATOR, token.LT, "<"},
        {token.LITERAL, token.I64, "5"},

        {token.LITERAL, token.I64, "5"},
        {token.OPERATOR, token.GT, ">"},
        {token.LITERAL, token.I64, "5"},

        {token.LITERAL, token.I64, "5"},
        {token.OPERATOR, token.LTE, "<="},
        {token.LITERAL, token.I64, "5"},

        {token.LITERAL, token.I64, "5"},
        {token.OPERATOR, token.GTE, ">="},
        {token.LITERAL, token.I64, "5"},

        {token.EOF, token.EOF, "EOF"},
    }

    runTest(t, input, tests)
}

func TestArrays(t *testing.T) {
    input := `
    list(1, 2, 3)
    list(1, 2, 3).0

    tuple(1, 2, 3)
    vec(1, 2, 3)
    set(1, 2, 3)
    `

    tests := []struct {
        expectedType int
        expectedSubtype int
        expectedLiteral string
    }{
        {token.TYPE, token.LIST, "list"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.LITERAL, token.I64, "1"},
        {token.DELIMITER, token.COMMA, ","},
        {token.LITERAL, token.I64, "2"},
        {token.DELIMITER, token.COMMA, ","},
        {token.LITERAL, token.I64, "3"},
        {token.DELIMITER, token.RPAREN, ")"},

        {token.TYPE, token.LIST, "list"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.LITERAL, token.I64, "1"},
        {token.DELIMITER, token.COMMA, ","},
        {token.LITERAL, token.I64, "2"},
        {token.DELIMITER, token.COMMA, ","},
        {token.LITERAL, token.I64, "3"},
        {token.DELIMITER, token.RPAREN, ")"},
        {token.DELIMITER, token.DOT, "."},
        {token.LITERAL, token.I64, "0"},

        {token.TYPE, token.TUPLE, "tuple"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.LITERAL, token.I64, "1"},
        {token.DELIMITER, token.COMMA, ","},
        {token.LITERAL, token.I64, "2"},
        {token.DELIMITER, token.COMMA, ","},
        {token.LITERAL, token.I64, "3"},
        {token.DELIMITER, token.RPAREN, ")"},

        {token.TYPE, token.VEC, "vec"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.LITERAL, token.I64, "1"},
        {token.DELIMITER, token.COMMA, ","},
        {token.LITERAL, token.I64, "2"},
        {token.DELIMITER, token.COMMA, ","},
        {token.LITERAL, token.I64, "3"},
        {token.DELIMITER, token.RPAREN, ")"},

        {token.TYPE, token.SET, "set"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.LITERAL, token.I64, "1"},
        {token.DELIMITER, token.COMMA, ","},
        {token.LITERAL, token.I64, "2"},
        {token.DELIMITER, token.COMMA, ","},
        {token.LITERAL, token.I64, "3"},
        {token.DELIMITER, token.RPAREN, ")"},

        {token.EOF, token.EOF, "EOF"},
    }

    runTest(t, input, tests)
}

func TestStructs(t *testing.T) {
    input := `
    let person be struct(
        name: str,
        age: i64,
    )
    `

    tests := []struct {
        expectedType int
        expectedSubtype int
        expectedLiteral string
    }{
        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "person"},
        {token.KEYWORD, token.BE, "be"},
        {token.TYPE, token.STRUCT, "struct"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.IDENTIFIER, token.IDENTIFIER, "name"},
        {token.DELIMITER, token.COLON, ":"},
        {token.TYPE, token.STR, "str"},
        {token.DELIMITER, token.COMMA, ","},
        {token.IDENTIFIER, token.IDENTIFIER, "age"},
        {token.DELIMITER, token.COLON, ":"},
        {token.TYPE, token.I64, "i64"},
        {token.DELIMITER, token.COMMA, ","},
        {token.DELIMITER, token.RPAREN, ")"},

        {token.EOF, token.EOF, "EOF"},
    }

    runTest(t, input, tests)
}

func TestEnums(t *testing.T) {
    input := `
    let color be enum(
        red,
        blue,
        green,
    )
    `

    tests := []struct {
        expectedType int
        expectedSubtype int
        expectedLiteral string
    }{
        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "color"},
        {token.KEYWORD, token.BE, "be"},
        {token.TYPE, token.ENUM, "enum"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.IDENTIFIER, token.IDENTIFIER, "red"},
        {token.DELIMITER, token.COMMA, ","},
        {token.IDENTIFIER, token.IDENTIFIER, "blue"},
        {token.DELIMITER, token.COMMA, ","},
        {token.IDENTIFIER, token.IDENTIFIER, "green"},
        {token.DELIMITER, token.COMMA, ","},
        {token.DELIMITER, token.RPAREN, ")"},

        {token.EOF, token.EOF, "EOF"},
    }

    runTest(t, input, tests)
}

func TestMap(t *testing.T) {
    input := `
    let m be map(
        "key": "value",
    )
    `

    tests := []struct {
        expectedType int
        expectedSubtype int
        expectedLiteral string
    }{
        {token.KEYWORD, token.LET, "let"},
        {token.IDENTIFIER, token.IDENTIFIER, "m"},
        {token.KEYWORD, token.BE, "be"},
        {token.TYPE, token.MAP, "map"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.LITERAL, token.STR, "key"},
        {token.DELIMITER, token.COLON, ":"},
        {token.LITERAL, token.STR, "value"},
        {token.DELIMITER, token.COMMA, ","},
        {token.DELIMITER, token.RPAREN, ")"},

        {token.EOF, token.EOF, "EOF"},
    }

    runTest(t, input, tests)
}

func TestMutStatement(t *testing.T) {
    input := `
    mut x to 5
    `

    tests := []struct {
        expectedType int
        expectedSubtype int
        expectedLiteral string
    }{
        {token.KEYWORD, token.MUT, "mut"},
        {token.IDENTIFIER, token.IDENTIFIER, "x"},
        {token.KEYWORD, token.TO, "to"},
        {token.LITERAL, token.I64, "5"},

        {token.EOF, token.EOF, "EOF"},
    }

    runTest(t, input, tests)
}

func TestExeStatement(t *testing.T) {
    input := `
    exe x
    `

    tests := []struct {
        expectedType int
        expectedSubtype int
        expectedLiteral string
    }{
        {token.KEYWORD, token.EXE, "exe"},
        {token.IDENTIFIER, token.IDENTIFIER, "x"},

        {token.EOF, token.EOF, "EOF"},
    }

    runTest(t, input, tests)
}

func TestBuiltinMethod(t *testing.T) {
    input := `
    x.len()
    `
    tests := []struct {
        expectedType int
        expectedSubtype int
        expectedLiteral string
    }{
        {token.IDENTIFIER, token.IDENTIFIER, "x"},
        {token.DELIMITER, token.DOT, "."},
        {token.IDENTIFIER, token.IDENTIFIER, "len"},
        {token.DELIMITER, token.LPAREN, "("},
        {token.DELIMITER, token.RPAREN, ")"},
        {token.EOF, token.EOF, "EOF"},
    }

    runTest(t, input, tests)
}


// =======
// Helpers
// =======
func runTest(t *testing.T, input string, tests []struct {expectedType int; expectedSubtype int; expectedLiteral string}) {
    tokenizer := New(input)

    token := tokenizer.GetToken()

    for i, tt := range tests {
        if token.Literal != tt.expectedLiteral {
            t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, token.Literal)
        }

        if token.Type != tt.expectedType {
            t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, tt.expectedType, token.Type)
        }

        if token.Subtype != tt.expectedSubtype {
            t.Fatalf("tests[%d] - token subtype wrong. expected=%q, got=%q", i, tt.expectedSubtype, token.Subtype)
        }

        token = tokenizer.GetToken()
    }
}
