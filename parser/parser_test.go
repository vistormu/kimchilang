package parser

import (
	"kimchi/ast"
	"kimchi/tokenizer"
	"testing"
)

func TestLetStatement(t *testing.T) {
    tests := []struct {
        input string
        expectedIdentifierLiteral string
        expectedIdentifierType string
        expectedExpressionValue interface{}
    }{
        {"let x: i64 = 5", "x", "i64", 5},
        {"let y: f64 = 10.5", "y", "f64", 10.5},
        {"let z: bool = true", "z", "bool", true},
        {"let foo: str = \"bar\"", "foo", "str", "bar"},
        {"let foo: bool = bar", "foo", "bool", "bar"},
        {"let x be 5", "x", "i64", 5},
        {"let y be 10.5", "y", "f64", 10.5},
        {"let z be true", "z", "bool", true},
        {"let foo be \"bar\"", "foo", "str", "bar"},
    }

    for _, tt := range tests {
        tokenizer := tokenizer.New(tt.input)
        parser := New(tokenizer)

        program := parser.ParseProgram()

        checkParserErrors(t, parser)

        if program == nil {
            t.Fatalf("ParseProgram() returned nil")
        }

        if len(program.Statements) != 1 {
            t.Fatalf("program.Statements does not contain one statement. got=%d", len(program.Statements))
        }

        statement, ok := program.Statements[0].(*ast.LetStatement)
        if !ok {
            t.Fatalf("program.Statements[0] is not ast.LetStatement. got=%T", program.Statements[0])
        }

        if statement.Identifier.Name != tt.expectedIdentifierLiteral {
            t.Fatalf("statement.Identifier.Name not %s. got=%s", tt.expectedIdentifierLiteral, statement.Identifier.Name)
        }

        if statement.Identifier.Type.Type.Literal != tt.expectedIdentifierType {
            t.Fatalf("statement.Identifier.Type not %s. got=%s", tt.expectedIdentifierType, statement.Identifier.Type.Type.Literal)
        }

        if !testExpressionValue(t, statement.Expression, tt.expectedExpressionValue) {
            return
        }
    }
}

func TestReturnStatement(t *testing.T) {
    tests := []struct {
        input string
        // expectedType string
        expectedExpressionValue interface{}
    }{
        // {"return 5", "i64", 5},
        // {"return 10.5", "f64", 10.5},
        // {"return true", "bool", true},
        // {"return \"bar\"", "str", "bar"},
        // {"return bar", "bool", "bar"},
        {"return 5", 5},
        {"return 10.5", 10.5},
        {"return true", true},
        {"return \"bar\"", "bar"},
        {"return bar", "bar"},
    }

    for _, tt := range tests {
        tokenizer := tokenizer.New(tt.input)
        parser := New(tokenizer)

        program := parser.ParseProgram()

        checkParserErrors(t, parser)

        if program == nil {
            t.Fatalf("ParseProgram() returned nil")
        }

        if len(program.Statements) != 1 {
            t.Fatalf("program.Statements does not contain one statement. got=%d", len(program.Statements))
        }

        statement, ok := program.Statements[0].(*ast.ReturnStatement)
        if !ok {
            t.Fatalf("program.Statements[0] is not ast.ReturnStatement. got=%T", program.Statements[0])
        }

        // if statement.Type.Literal != tt.expectedType {
        //     t.Fatalf("statement.Type not %s. got=%s", tt.expectedType, statement.Type.Literal)
        // }

        if !testExpressionValue(t, statement.Expression, tt.expectedExpressionValue) {
            return
        }
    }
}

func TestPrefixExpressions(t *testing.T) {
    tests := []struct {
        input string
        expectedOperator string
        expectedValue interface{}
    }{
        {"not 5", "not", 5},
        {"-10.5", "-", 10.5},
        {"not true", "not", true},
        {"not false", "not", false},
    }

    for _, tt := range tests {
        tokenizer := tokenizer.New(tt.input)
        parser := New(tokenizer)

        program := parser.ParseProgram()

        checkParserErrors(t, parser)

        if len(program.Statements) != 1 {
            t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
        }

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
        }

        exp, ok := stmt.Expression.(*ast.PrefixExpression)
        if !ok {
            t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
        }

        if exp.Operator != tt.expectedOperator {
            t.Fatalf("exp.Operator is not '%s'. got=%s", tt.expectedOperator, exp.Operator)
        }

        if !testExpressionValue(t, exp.Right, tt.expectedValue) {
            return
        }
    }
}

func TestInfixExpressions(t *testing.T) {
    tests := []struct {
        input string
        leftValue interface{}
        operator string
        rightValue interface{}
    }{
        {"5 + 5", 5, "+", 5},
        {"5 - 5", 5, "-", 5},
        {"5 * 5", 5, "*", 5},
        {"5 / 5", 5, "/", 5},
        {"5 > 5", 5, ">", 5},
        {"5 < 5", 5, "<", 5},
        {"5 >= 5", 5, ">=", 5},
        {"5 <= 5", 5, "<=", 5},
        {"5 is 5", 5, "is", 5},
        {"5 is not 5", 5, "is_not", 5},
        {"5 and 5", 5, "and", 5},
        {"5 or 5", 5, "or", 5},
        // {"5 % 5", 5, "%", 5},
        {"true is true", true, "is", true},
        {"true is not true", true, "is_not", true},
        {"true and true", true, "and", true},
        {"true or true", true, "or", true},
    }

    for _, tt := range tests {
        tokenizer := tokenizer.New(tt.input)
        parser := New(tokenizer)

        program := parser.ParseProgram()

        checkParserErrors(t, parser)

        if len(program.Statements) != 1 {
            t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
        }

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
        }

        exp, ok := stmt.Expression.(*ast.InfixExpression)
        if !ok {
            t.Fatalf("stmt is not ast.InfixExpression. got=%T", stmt.Expression)
        }

        if !testExpressionValue(t, exp.Left, tt.leftValue) {
            return
        }

        if exp.Operator != tt.operator {
            t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
        }

        if !testExpressionValue(t, exp.Right, tt.rightValue) {
            return
        }
    }
}

func TestOperatorPrecedenceParsing(t *testing.T) {
    tests := []struct {
        input string
        expected string
    }{
        {"-a * b", "((-a) * b)"},
        {"not -a", "(not(-a))"},
        {"a + b + c", "((a + b) + c)"},
        {"a + b - c", "((a + b) - c)"},
        {"a * b * c", "((a * b) * c)"},
        {"a * b / c", "((a * b) / c)"},
        {"a + b / c", "(a + (b / c))"},
        {"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
        {"3 + 4 -5 * 5", "((3 + 4) - (5 * 5))"}, 
        {"5 > 4 is 3 < 4", "((5 > 4) is (3 < 4))"},
        {"5 < 4 is not 3 > 4", "((5 < 4) is_not (3 > 4))"},
        {"3 + 4 * 5 is 3 * 1 + 4 * 5", "((3 + (4 * 5)) is ((3 * 1) + (4 * 5)))"},
        {"true is true", "(true is true)"},
        {"true is not false", "(true is_not false)"},
        {"false is false", "(false is false)"},
        {"false is not true", "(false is_not true)"},
        {"3 > 5 is false", "((3 > 5) is false)"},
        {"3 < 5 is true", "((3 < 5) is true)"},
        {"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
        {"(5 + 5) * 2", "((5 + 5) * 2)"},
        {"2 / (5 + 5)", "(2 / (5 + 5))"},
        {"-(5 + 5)", "(-(5 + 5))"},
        {"not (true is true)", "(not(true is true))"},
        {"a + add(b * c) + d", "((a + add((b * c))) + d)"},
        {"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
        {"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
    }

    for _, tt := range tests {
        tokenizer := tokenizer.New(tt.input)
        parser := New(tokenizer)
        program := parser.ParseProgram()
        checkParserErrors(t, parser)

        actual := program.String()
        if actual != tt.expected {
            t.Errorf("expected=%q, got=%q", tt.expected, actual)
        }
    }
}

func TestIfExpression(t *testing.T) {
    input := `
    if (x < y) { x }
    `
    tokenizer := tokenizer.New(input)
    parser := New(tokenizer)
    program := parser.ParseProgram()
    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    exp, ok := stmt.Expression.(*ast.IfExpression)
    if !ok {
        t.Fatalf("stmt is not ast.IfExpression. got=%T", stmt.Expression)
    }

    if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
        return
    }

    if len(exp.Consequence.Statements) != 1 {
        t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
    }

    consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
    }

    if !testIdentifier(t, consequence.Expression, "x") {
        return
    }

    if exp.Alternative != nil {
        t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
    }
}

func TestIfElseExpression(t *testing.T) {
    input := `
    if (x < y) { x } else { y }
    `
    tokenizer := tokenizer.New(input)
    parser := New(tokenizer)
    program := parser.ParseProgram()
    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    exp, ok := stmt.Expression.(*ast.IfExpression)
    if !ok {
        t.Fatalf("stmt is not ast.IfExpression. got=%T", stmt.Expression)
    }

    if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
        return
    }

    if len(exp.Consequence.Statements) != 1 {
        t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
    }

    consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
    }

    if !testIdentifier(t, consequence.Expression, "x") {
        return
    }

    if len(exp.Alternative.Statements) != 1 {
        t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Alternative.Statements))
    }

    alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
    }

    if !testIdentifier(t, alternative.Expression, "y") {
        return
    }
}

func TestFunctionLiteralParsing(t *testing.T) {
    input := `fn(x: i64, y: bool): f64 { x + y }`

    tokenizer := tokenizer.New(input)
    parser := New(tokenizer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    function, ok := stmt.Expression.(*ast.FunctionLiteral)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
    }

    if len(function.Parameters) != 2 {
        t.Fatalf("function literal parameters wrong. want 2, got=%d\n", len(function.Parameters))
    }

    testExpressionValue(t, function.Parameters[0], "x")
    testExpressionValue(t, function.Parameters[1], "y")
    testIdentifierType(t, function.Parameters[0], "i64")
    testIdentifierType(t, function.Parameters[1], "bool")

    if function.ReturnType.Literal != "f64" {
        t.Fatalf("function literal return type wrong. want f64, got=%s\n", function.ReturnType.Literal)
    }

    if len(function.Body.Statements) != 1 {
        t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n", len(function.Body.Statements))
    }

    bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
    }

    testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestCallExpressionParsing(t *testing.T) {
    input := `add(1, 2 * 3, 4 + 5)`

    tokenizer := tokenizer.New(input)
    parser := New(tokenizer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    exp, ok := stmt.Expression.(*ast.CallExpression)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
    }

    if !testIdentifier(t, exp.Function, "add") {
        return
    }

    if len(exp.Arguments) != 3 {
        t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
    }

    testExpressionValue(t, exp.Arguments[0], 1)
    testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
    testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestListLiteralParsing(t *testing.T) {
    input := `list(1, 2 * 2, 3 + 3)`

    tokenizer := tokenizer.New(input)
    parser := New(tokenizer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    list, ok := stmt.Expression.(*ast.ListLiteral)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.ListLiteral. got=%T", stmt.Expression)
    }

    if len(list.Elements) != 3 {
        t.Fatalf("wrong length of elements. got=%d", len(list.Elements))
    }

    testExpressionValue(t, list.Elements[0], 1)
    testInfixExpression(t, list.Elements[1], 2, "*", 2)
    testInfixExpression(t, list.Elements[2], 3, "+", 3)
}

func TestIndexExpressionParsing(t *testing.T) {
    input := `myArray(1+1)`

    tokenizer := tokenizer.New(input)
    parser := New(tokenizer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    indexExp, ok := stmt.Expression.(*ast.CallExpression)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.IndexExpression. got=%T", stmt.Expression)
    }

    if !testIdentifier(t, indexExp.Function, "myArray") {
        return
    }

    if !testInfixExpression(t, indexExp.Arguments[0], 1, "+", 1) {
        return
    }
}

func TestMapLiteralStringKeys(t *testing.T) {
    input := `map("one": 1, "two": 2, "three": 3)`

    tokenizer := tokenizer.New(input)
    parser := New(tokenizer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    mapLiteral, ok := stmt.Expression.(*ast.MapLiteral)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.MapLiteral. got=%T", stmt.Expression)
    }

    if len(mapLiteral.Pairs) != 3 {
        t.Fatalf("mapLiteral.Pairs has wrong length. got=%d", len(mapLiteral.Pairs))
    }

    expected := map[string]int64{
        "one":   1,
        "two":   2,
        "three": 3,
    }

    for key, value := range mapLiteral.Pairs {
        literal, ok := key.(*ast.StringLiteral)
        if !ok {
            t.Errorf("key is not ast.StringLiteral. got=%T", key)
            continue
        }

        expectedValue := expected[literal.String()]
        testIntegerLiteral(t, value, expectedValue)
    }
}

func TestMapLiteralIntegerKeys(t *testing.T) {
    input := `map(1: "one", 2: "two", 3: "three")`

    tokenizer := tokenizer.New(input)
    parser := New(tokenizer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    mapLiteral, ok := stmt.Expression.(*ast.MapLiteral)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.MapLiteral. got=%T", stmt.Expression)
    }

    if len(mapLiteral.Pairs) != 3 {
        t.Fatalf("mapLiteral.Pairs has wrong length. got=%d", len(mapLiteral.Pairs))
    }

    expected := map[int64]string{
        1: "one",
        2: "two",
        3: "three",
    }

    for key, value := range mapLiteral.Pairs {
        literal, ok := key.(*ast.IntegerLiteral)
        if !ok {
            t.Errorf("key is not ast.IntegerLiteral. got=%T", key)
            continue
        }

        expectedValue := expected[literal.Value]
        testStringLiteral(t, value, expectedValue)
    }
}

func TestMapLiteralBooleanKeys(t *testing.T) {
    input := `map(true: "one", false: "two")`

    tokenizer := tokenizer.New(input)
    parser := New(tokenizer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    mapLiteral, ok := stmt.Expression.(*ast.MapLiteral)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.MapLiteral. got=%T", stmt.Expression)
    }

    if len(mapLiteral.Pairs) != 2 {
        t.Fatalf("mapLiteral.Pairs has wrong length. got=%d", len(mapLiteral.Pairs))
    }

    expected := map[bool]string{
        true:  "one",
        false: "two",
    }

    for key, value := range mapLiteral.Pairs {
        literal, ok := key.(*ast.BooleanLiteral)
        if !ok {
            t.Errorf("key is not ast.Boolean. got=%T", key)
            continue
        }

        expectedValue := expected[literal.Value]
        testStringLiteral(t, value, expectedValue)
    }
}

func TestEmptyMapLiteral(t *testing.T) {
    input := `map()`

    tokenizer := tokenizer.New(input)
    parser := New(tokenizer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    mapLiteral, ok := stmt.Expression.(*ast.MapLiteral)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.MapLiteral. got=%T", stmt.Expression)
    }

    if len(mapLiteral.Pairs) != 0 {
        t.Fatalf("mapLiteral.Pairs has wrong length. got=%d", len(mapLiteral.Pairs))
    }
}

func TestTypeLiteral(t *testing.T) {
    tests := []struct {
        input string
        expectedTypeLiteral  string
        expectedSubtypesLiteral []string
    }{
        {"let x: i64 = 5", "i64", nil},
        {"let y: f64 = 5.0", "f64", nil},
        {"let z: bool = true", "bool", nil},
        {"let a: str = \"hello\"", "str", nil},
        {"let b: list(i64) = list(1, 2, 3)", "list", []string{"i64"}},
        {"let c: map(i64, str) = map(1: \"one\", 2: \"two\")", "map", []string{"i64", "str"}},
    }

    for _, tt := range tests {
        tokenizer := tokenizer.New(tt.input)
        parser := New(tokenizer)
        program := parser.ParseProgram()
        checkParserErrors(t, parser)

        if len(program.Statements) != 1 {
            t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
        }

        stmt, ok := program.Statements[0].(*ast.LetStatement)
        if !ok {
            t.Fatalf("program.Statements[0] is not ast.LetStatement. got=%T", program.Statements[0])
        }

        if stmt.Identifier.Type.Type.Literal != tt.expectedTypeLiteral {
            t.Fatalf("tyoe literal not '%s'. got=%s", tt.expectedTypeLiteral, stmt.Identifier.Type.Type.Literal)
        }

        if len(stmt.Identifier.Type.Subtypes) != len(tt.expectedSubtypesLiteral) {
            t.Fatalf("len(stmt.Subtypes) not %d. got=%d", len(tt.expectedSubtypesLiteral), len(stmt.Identifier.Type.Subtypes))
        }

        for i, subtype := range tt.expectedSubtypesLiteral {
            if stmt.Identifier.Type.Subtypes[i].Literal != subtype {
                t.Fatalf("stmt.Subtypes[%d].Literal not '%s'. got=%s", i, subtype, stmt.Identifier.Type.Subtypes[i].Literal)
            }
        }
    }
}

func TestWhileExpression(t *testing.T) {
    input := `
    while x < 10 {
        print(x)
    }
    `
    tokenizer := tokenizer.New(input)
    parser := New(tokenizer)
    program := parser.ParseProgram()
    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    whileExpression, ok := stmt.Expression.(*ast.WhileExpression)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.WhileExpression. got=%T", stmt.Expression)
    }

    if whileExpression.Condition.String() != "(x < 10)" {
        t.Fatalf("whileExpression.Condition.String() is not '(x < 10)'. got=%s", whileExpression.Condition.String())
    }

    if len(whileExpression.Body.Statements) != 1 {
        t.Fatalf("whileExpression.Block.Statements does not contain %d statements. got=%d\n", 1, len(whileExpression.Body.Statements))
    }

    blockStmt, ok := whileExpression.Body.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("whileExpression.Block.Statements[0] is not ast.ExpressionStatement. got=%T", whileExpression.Body.Statements[0])
    }

    if blockStmt.Expression.String() != "print(x)" {
        t.Fatalf("blockStmt.Expression.String() is not 'print(x)'. got=%s", blockStmt.Expression.String())
    }
}

func TestMutStatement(t *testing.T) {
    input := `
    mut x to 5
    `
    tokenizer := tokenizer.New(input)
    parser := New(tokenizer)
    program := parser.ParseProgram()
    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.MutStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.MutStatement. got=%T", program.Statements[0])
    }

    if stmt.Identifier.String() != "x" {
        t.Fatalf("stmt.Identifier.String() is not 'x'. got=%s", stmt.Identifier.String())
    }

    if stmt.Expression.String() != "5" {
        t.Fatalf("stmt.Value.String() is not '5'. got=%s", stmt.Expression.String())
    }

    if stmt.String() != "mut x to 5" {
        t.Fatalf("stmt.String() is not 'mut x to 5'. got=%s", stmt.String())
    }
}

func TestExeStatement(t *testing.T) {
    input := `exe print("hello")`

    tokenizer := tokenizer.New(input)
    parser := New(tokenizer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExeStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExeStatement. got=%T", program.Statements[0])
    }

    if stmt.Function.String() != "print" {
        t.Fatalf("stmt.Expression.String() is not 'print'. got=%s", stmt.Function.String())
    }

    if stmt.Arguments[0].String() != "hello" {
        t.Fatalf("stmt.Arguments[0].String() is not '\"hello\"'. got=%s", stmt.Arguments[0].String())
    }

}

func TestMethodExpression(t *testing.T) {
    input := `
    x.len()
    `
    tokenizer := tokenizer.New(input)
    parser := New(tokenizer)
    program := parser.ParseProgram()
    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    methodExpression, ok := stmt.Expression.(*ast.MethodExpression)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.MethodExpression. got=%T", stmt.Expression)
    }

    if methodExpression.Left.String() != "x" {
        t.Fatalf("methodExpression.Object.String() is not 'x'. got=%s", methodExpression.Left.String())
    }

    if methodExpression.Method.String() != "len" {
        t.Fatalf("methodExpression.Method.String() is not 'len'. got=%s", methodExpression.Method.String())
    }

    if len(methodExpression.Arguments) != 0 {
        t.Fatalf("methodExpression.Arguments does not contain %d arguments. got=%d\n", 0, len(methodExpression.Arguments))
    }
}

func TestMethodExpressionWithArguments(t *testing.T) {
    input := `
    x.add(1, 2)
    `
    tokenizer := tokenizer.New(input)
    parser := New(tokenizer)
    program := parser.ParseProgram()
    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    methodExpression, ok := stmt.Expression.(*ast.MethodExpression)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.MethodExpression. got=%T", stmt.Expression)
    }

    if methodExpression.Left.String() != "x" {
        t.Fatalf("methodExpression.Object.String() is not 'x'. got=%s", methodExpression.Left.String())
    }

    if methodExpression.Method.String() != "add" {
        t.Fatalf("methodExpression.Method.String() is not 'add'. got=%s", methodExpression.Method.String())
    }

    if len(methodExpression.Arguments) != 2 {
        t.Fatalf("methodExpression.Arguments does not contain %d arguments. got=%d\n", 2, len(methodExpression.Arguments))
    }

    if methodExpression.Arguments[0].String() != "1" {
        t.Fatalf("methodExpression.Arguments[0].String() is not '1'. got=%s", methodExpression.Arguments[0].String())
    }

    if methodExpression.Arguments[1].String() != "2" {
        t.Fatalf("methodExpression.Arguments[1].String() is not '2'. got=%s", methodExpression.Arguments[1].String())
    }
}

func TestAttributeExpression(t *testing.T) {
    input := `
    x.len
    `
    tokenizer := tokenizer.New(input)
    parser := New(tokenizer)
    program := parser.ParseProgram()
    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    attributeExpression, ok := stmt.Expression.(*ast.AttributeExpression)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.AttributeExpression. got=%T", stmt.Expression)
    }

    if attributeExpression.Left.String() != "x" {
        t.Fatalf("attributeExpression.Object.String() is not 'x'. got=%s", attributeExpression.Left.String())
    }

    if attributeExpression.Attribute.String() != "len" {
        t.Fatalf("attributeExpression.Attribute.String() is not 'len'. got=%s", attributeExpression.Attribute.String())
    }
}

func TestBuiltinFunctionExpression(t *testing.T) {
    input := `
    let content: str = read("test.txt")
    `

    tokenizer := tokenizer.New(input)
    parser := New(tokenizer)
    program := parser.ParseProgram()
    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.LetStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.LetStatement. got=%T", program.Statements[0])
    }

    letStatement := stmt
    if letStatement.Identifier.String() != "content" {
        t.Fatalf("letStatement.Name.String() is not 'content'. got=%s", letStatement.Identifier.String())
    }

    if letStatement.Identifier.Type.String() != "str" {
        t.Fatalf("letStatement.Type.String() is not 'str'. got=%s", letStatement.Identifier.Type.String())
    }

    if letStatement.Expression.String() != "read(test.txt)" {
        t.Fatalf("letStatement.Value.String() is not 'read(\"test.txt\")'. got=%s", letStatement.Expression.String())
    }

    expression, ok := letStatement.Expression.(*ast.CallExpression)
    if !ok {
        t.Fatalf("letStatement.Value is not ast.CallExpression. got=%T", letStatement.Expression)
    }

    if expression.Function.String() != "read" {
        t.Fatalf("expression.Function.String() is not 'read'. got=%s", expression.Function.String())
    }

    if len(expression.Arguments) != 1 {
        t.Fatalf("expression.Arguments does not contain %d arguments. got=%d\n", 1, len(expression.Arguments))
    }

    if expression.Arguments[0].String() != "test.txt" {
        t.Fatalf("expression.Arguments[0].String() is not 'test.txt'. got=%s", expression.Arguments[0].String())
    }
}


// =======
// HELPERS
// =======
func checkParserErrors(t *testing.T, parser *Parser) {
    errors := parser.Errors
    if len(errors) == 0 {
        return
    }

    t.Errorf("parser has %d errors", len(errors))
    for _, msg := range errors {
        t.Errorf("parser error: %q", msg)
    }
    t.FailNow()
}
func testExpressionValue(t *testing.T, exp ast.Expression, expected interface{}) bool {
    switch v := expected.(type) {
    case int:
        return testIntegerLiteral(t, exp, int64(v))
    case int64:
        return testIntegerLiteral(t, exp, v)
    case float64:
        return testFloatLiteral(t, exp, v)
    case string:
        if identifier, ok := exp.(*ast.Identifier); ok {
            return testIdentifier(t, identifier, v)
        }
        return testStringLiteral(t, exp, v)
    case bool:
        return testBooleanLiteral(t, exp, v)
    }
    t.Errorf("type of exp not handled. got=%T", exp)
    return false
}
func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
    integer, ok := il.(*ast.IntegerLiteral)
    if !ok {
        t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
        return false
    }
    if integer.Value != value {
        t.Errorf("integer.Value not %d. got=%d", value, integer.Value)
        return false
    }
    return true
}
func testFloatLiteral(t *testing.T, fl ast.Expression, value float64) bool {
    float, ok := fl.(*ast.FloatLiteral)
    if !ok {
        t.Errorf("fl not *ast.FloatLiteral. got=%T", fl)
        return false
    }
    if float.Value != value {
        t.Errorf("float.Value not %f. got=%f", value, float.Value)
        return false
    }
    return true
}
func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
    identifier, ok := exp.(*ast.Identifier)
    if !ok {
        t.Errorf("exp not *ast.Identifier. got=%T", exp)
        return false
    }
    if identifier.Name != value {
        t.Errorf("identifier.Name not %s. got=%s", value, identifier.Name)
        return false
    }
    return true
}
func testIdentifierType(t *testing.T, exp ast.Expression, value string) bool {
    identifier, ok := exp.(*ast.Identifier)
    if !ok {
        t.Errorf("exp not *ast.Identifier. got=%T", exp)
        return false
    }
    if identifier.Type.Type.Literal != value {
        t.Errorf("identifier.Type not %s. got=%s", value, identifier.Type.Type.Literal)
        return false
    }
    return true
}
func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
    boolean, ok := exp.(*ast.BooleanLiteral)
    if !ok {
        t.Errorf("exp not *ast.Boolean. got=%T", exp)
        return false
    }
    if boolean.Value != value {
        t.Errorf("boolean.Value not %t. got=%t", value, boolean.Value)
        return false
    }
    return true
}
func testStringLiteral(t *testing.T, exp ast.Expression, value string) bool {
    str, ok := exp.(*ast.StringLiteral)
    if !ok {
        t.Errorf("exp not *ast.StringLiteral. got=%T", exp)
        return false
    }
    if str.Value != value {
        t.Errorf("str.Value not %s. got=%s", value, str.Value)
        return false
    }
    return true
}
func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
    opExp, ok := exp.(*ast.InfixExpression)
    if !ok {
        t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
        return false
    }
    if !testExpressionValue(t, opExp.Left, left) {
        return false
    }
    if opExp.Operator != operator {
        t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
        return false
    }
    if !testExpressionValue(t, opExp.Right, right) {
        return false
    }
    return true
}
