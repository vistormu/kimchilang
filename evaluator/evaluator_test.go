package evaluator

import (
    "testing"
    "kimchi/object"
    "kimchi/tokenizer"
    "kimchi/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    }{
        {"5", 5},
        {"10", 10},
        {"-5", -5},
        {"-10", -10},
        {"5 + 5 + 5 + 5 - 10", 10},
        {"2 * 2 * 2 * 2 * 2", 32},
        {"-50 + 100 + -50", 0},
        {"5 * 2 + 10", 20},
        {"5 + 2 * 10", 25},
        {"20 + 2 * -10", 0},
        {"50 / 2 * 2 + 10", 60},
        {"2 * (5 + 10)", 30},
        {"3 * 3 * 3 + 10", 37},
        {"3 * (3 * 3) + 10", 37},
        {"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testIntegerObject(t, evaluated, tt.expected)
    }
}

func TestEvalFloatExpression(t *testing.T) {
    tests := []struct {
        input string
        expected float64
    }{
        {"5.5", 5.5},
        {"10.5", 10.5},
        {"-5.5", -5.5},
        {"-10.5", -10.5},
        {"5.5 + 5.5 + 5.5 + 5.5 - 10.5", 5.5 + 5.5 + 5.5 + 5.5 - 10.5},
        {"2.5 * 2.5 * 2.5 * 2.5 * 2.5", 2.5*2.5*2.5*2.5*2.5},
        {"-50.5 + 100.5 + -50.5", -50.5 + 100.5 + -50.5},
        {"5.5 * 2.5 + 10.5", 5.5*2.5 + 10.5},
        {"5.5 + 2.5 * 10.5", 5.5 + 2.5*10.5},
        {"20.5 + 2.5 * -10.5", 20.5 + 2.5*-10.5},
        {"50.5 / 2.5 * 2.5 + 10.5", 50.5 / 2.5 * 2.5 + 10.5},
        {"2.5 * (5.5 + 10.5)", 2.5 * (5.5 + 10.5)},
        {"3.5 * 3.5 * 3.5 + 10.5", 3.5*3.5*3.5 + 10.5},
        {"3.5 * (3.5 * 3.5) + 10.5", 3.5 * (3.5 * 3.5) + 10.5},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testFloatObject(t, evaluated, tt.expected)
    }
}

func TestEvalBooleanExpression(t *testing.T) {
    tests := []struct {
        input string
        expected bool
    }{
        {"true", true},
        {"false", false},
        {"1 < 2", true},
        {"1 > 2", false},
        {"1 < 1", false},
        {"1 > 1", false},
        {"1 is 1", true},
        {"1 is not 1", false},
        {"1 is 2", false},
        {"1 is not 2", true},
        {"true is true", true},
        {"false is false", true},
        {"true is false", false},
        {"true is not false", true},
        {"false is not true", true},
        {"(1 < 2) is true", true},
        {"(1 < 2) is false", false},
        {"(1 > 2) is true", false},
        {"(1 > 2) is false", true},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testBooleanObject(t, evaluated, tt.expected)
    }
}

func TestEvalStringExpression(t *testing.T) {
    tests := []struct {
        input string
        expected string
    }{
        {`"Hello World!"`, "Hello World!"},
        {`"Hello" + " " + "World!"`, "Hello World!"},
        {`"Hello" + " " + "World!" + " " + "from" + " " + "Kimchi!"`, "Hello World! from Kimchi!"},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testStringObject(t, evaluated, tt.expected)
    }
}

func TestNotOperator(t *testing.T) {
    tests := []struct {
        input string
        expected bool
    }{
        {"not true", false},
        {"not false", true},
        {"not 5", false},
        {"not not true", true},
        {"not not false", false},
        {"not not 5", true},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testBooleanObject(t, evaluated, tt.expected)
    }
}

func TestIfElseExpressions(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"if (true) { 10 }", 10},
        {"if (false) { 10 }", nil},
        {"if (1) { 10 }", 10},
        {"if (1 < 2) { 10 }", 10},
        {"if (1 > 2) { 10 }", nil},
        {"if (1 > 2) { 10 } else { 20 }", 20},
        {"if (1 < 2) { 10 } else { 20 }", 10},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        integer, ok := tt.expected.(int)
        if ok {
            testIntegerObject(t, evaluated, int64(integer))
        } else {
            testNoneObject(t, evaluated)
        }
    }
}

func TestReturnStatements(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    }{
        {"return 10;", 10},
        {"return 10; 9;", 10},
        {"return 2 * 5; 9;", 10},
        {"9; return 2 * 5; 9;", 10},
        {`
        if (10 > 1) {
            if (10 > 1) {
                return 10;
            }
            return 1;
        }
        `, 10},
    }
    
    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testIntegerObject(t, evaluated, tt.expected)
    }
}

func TestLetStatements(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    }{
        {"let a: i64 = 5 a", 5},
        {"let a: i64 = 5 * 5 a", 25},
        {"let a: i64 = 5 let b: i64 = a b", 5},
        {"let a: i64 = 5 let b: i64 = a let c: i64 = a + b + 5 c", 15},
    }

    for _, tt := range tests {
        testIntegerObject(t, testEval(tt.input), tt.expected)
    }
}

func TestFunctionObject(t *testing.T) {
    input := "fn(x: i64): i64 { x + 2 }"

    evaluated := testEval(input)
    fn, ok := evaluated.(*object.Function)
    if !ok {
        t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
    }
    if len(fn.Parameters) != 1 {
        t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
    }
    if fn.Parameters[0].String() != "x" {
        t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
    }
    expectedBody := "(x + 2)"
    if fn.Body.String() != expectedBody {
        t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
    }
}

func TestFunctionApplication(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    }{
        {"let identity: none = fn(x: i64): i64 { x } identity(5)", 5},
        {"let identity: none = fn(x: i64): i64 { return x } identity(5)", 5},
        {"let double: none = fn(x: i64): i64 { x * 2 } double(5)", 10},
        {"let add: none = fn(x: i64, y: i64): i64 { x + y } add(5, 5)", 10},
        {"let add: none = fn(x: i64, y: i64): i64 { x + y } add(5 + 5, add(5, 5))", 20},
        {"fn(x: i64): i64 { x }(5)", 5},
    }

    for _, tt := range tests {
        testIntegerObject(t, testEval(tt.input), tt.expected)
    }
}

func TestClosures(t *testing.T) {
    input := `
    let new_adder: none = fn(x: i64): i64 {
        fn(y: i64): i64 { x + y }
    }

    let add_two: none = new_adder(2)
    add_two(2)
    `

    testIntegerObject(t, testEval(input), 4)
}

func TestBuiltins(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {`len("")`, 0},
        {`len("four")`, 4},
        {`len("hello world")`, 11},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)

        switch expected := tt.expected.(type) {
        case int:
            testIntegerObject(t, evaluated, int64(expected))
        case string:
            errObj, ok := evaluated.(*object.Error)
            if !ok {
                t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
                continue
            }
            if errObj.Message != expected {
                t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
            }
        }
    }
}

func TestListLiterals(t *testing.T) {
    input := "list(1, 2 * 2, 3 + 3)"

    evaluated := testEval(input)
    result, ok := evaluated.(*object.List)
    if !ok {
        t.Fatalf("object is not List. got=%T (%+v)", evaluated, evaluated)
    }
    if len(result.Elements) != 3 {
        t.Fatalf("list has wrong num of elements. got=%d", len(result.Elements))
    }
    testIntegerObject(t, result.Elements[0], 1)
    testIntegerObject(t, result.Elements[1], 4)
    testIntegerObject(t, result.Elements[2], 6)
}

func TestListIndexExpressions(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"list(1, 2, 3)(0)", 1},
        {"list(1, 2, 3)(1)", 2},
        {"list(1, 2, 3)(2)", 3},
        {"let i: i64 = 0 list(1)(i)", 1},
        {"list(1, 2, 3)(1 + 1)", 3},
        {"let my_list: list = list(1, 2, 3) my_list(2)", 3},
        {"let my_list: list = list(1, 2, 3) my_list(0) + my_list(1) + my_list(2)", 6},
        {"let my_list: list = list(1, 2, 3) let i: i64 = my_list(0) my_list(i)", 2},
        {"list(1, 2, 3)(3)", nil},
        {"list(1, 2, 3)(-1)", nil},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)

        switch expected := tt.expected.(type) {
        case int:
            testIntegerObject(t, evaluated, int64(expected))
        case nil:
            testNoneObject(t, evaluated)
        }
    }
}

// =======
// HELPERS
// =======
func testEval(input string) object.Object {
    tokezinizer := tokenizer.New(input)
    parser := parser.New(tokezinizer)
    program := parser.ParseProgram()

    env := object.NewEnvironment()

    return Eval(program, env)
}
func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
    result, ok := obj.(*object.I64)
    if !ok {
        t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
        return false
    }
    if result.Value != expected {
        t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
        return false
    }
    return true
}
func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
    result, ok := obj.(*object.Bool)
    if !ok {
        t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
        return false
    }
    if result.Value != expected {
        t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
        return false
    }
    return true
}
func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
    result, ok := obj.(*object.F64)
    if !ok {
        t.Errorf("object is not Float. got=%T (%+v)", obj, obj)
        return false
    }
    if result.Value != expected {
        t.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
        return false
    }
    return true
}
func testStringObject(t *testing.T, obj object.Object, expected string) bool {
    result, ok := obj.(*object.Str)
    if !ok {
        t.Errorf("object is not String. got=%T (%+v)", obj, obj)
        return false
    }
    if result.Value != expected {
        t.Errorf("object has wrong value. got=%s, want=%s", result.Value, expected)
        return false
    }
    return true
}
func testNoneObject(t *testing.T, obj object.Object) bool {
    if obj != NONE {
        t.Errorf("object is not None. got=%T (%+v)", obj, obj)
        return false
    }
    return true
}
