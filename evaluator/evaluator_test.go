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

func TestMutStatements(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"let a be 5 mut a to 10 a", 10},
        {"let a be 5 mut a to 10 mut a to 20 a", 20},
        {"let a: list(i64) = list(1, 2, 3) mut a(0) to 0 a", []int64{0, 2, 3}},
    }

    for _, tt := range tests {
        switch expected := tt.expected.(type) {
        case int64:
            testIntegerObject(t, testEval(tt.input), expected)
        case []int64:
            testIntegerListObject(t, testEval(tt.input), expected)
        }
    }
}

func TestExeStatements(t *testing.T) {
    tests := []struct {
        input string
        expected string
    }{
        {"let print_number: none = fn(x: i64): none { print(x) } print_number(5)", "5"},
    }

    for _, tt := range tests {
        testNoneObject(t, testEval(tt.input))
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

func TestStringIndexExpressions(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {`"abc"(0)`, "a"},
        {`"abc"(1)`, "b"},
        {`"abc"(2)`, "c"},
        {`let i: i64 = 0 "abc"(i)`, "a"},
        {`"abc"(1 + 1)`, "c"},
        {`let my_str: str = "abc" my_str(2)`, "c"},
        {`let my_str: str = "abc" my_str(0) + my_str(1) + my_str(2)`, "abc"},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)

        switch expected := tt.expected.(type) {
        case string:
            testStringObject(t, evaluated, expected)
        case nil:
            testNoneObject(t, evaluated)
        }
    }
}

func TestMapLiterals(t *testing.T) {
    input := `
    let two: str = "two"
    map(
        "one": 10 - 9,
        two: 1 + 1,
        "thr" + "ee": 6 / 2,
        4: 4,
        true: 5,
        false: 6
    )
    `

    evaluated := testEval(input)
    result, ok := evaluated.(*object.Map)
    if !ok {
        t.Fatalf("Eval didn't return Map. got=%T (%+v)", evaluated, evaluated)
    }

    expected := map[object.MapKey]int64{
        (&object.Str{Value: "one"}).MapKey(): 1,
        (&object.Str{Value: "two"}).MapKey(): 2,
        (&object.Str{Value: "three"}).MapKey(): 3,
        (&object.I64{Value: 4}).MapKey(): 4,
        object.TRUE.MapKey(): 5,
        object.FALSE.MapKey(): 6,
    }

    if len(result.Pairs) != len(expected) {
        t.Fatalf("Map has wrong num of pairs. got=%d", len(result.Pairs))
    }

    for expectedKey, expectedValue := range expected {
        pair, ok := result.Pairs[expectedKey]
        if !ok {
            t.Errorf("no pair for given key in Pairs")
        }

        testIntegerObject(t, pair.Value, expectedValue)
    }
}

func TestMapIndexExpressions(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        { `map( "one": 1, "two": 2)("one")`, 1, },
        { `map( "one": 1, "two": 2)("two")`, 2, }, 
        { `let key: str = "one" map( "one": 1, "two": 2)(key)`, 1, },
        { `map( "one": 1, "two": 2)(3)`, nil, },
        { `map( "one": 1, "two": 2)(true)`, nil, },
        { `map( "one": 1, "two": 2)(false)`, nil, },
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

func TestWhileExpression(t *testing.T) {
    input := `
    let counter be 0
    while counter < 10 {
        mut counter to counter + 1
    }
    counter
    `

    evaluated := testEval(input)
    testIntegerObject(t, evaluated, 10)
}

func TestForExpression(t *testing.T) {
    input := `
    let result be 0
    for i, _ in list(1, 2, 3) {
        mut result to result + i
    }
    result
    `

    evaluated := testEval(input)
    testIntegerObject(t, evaluated, 3)

    input = `
    let result be 0
    for _, value in list(1, 2, 3) {
        mut result to result + value
    }
    result
    `

    evaluated = testEval(input)
    testIntegerObject(t, evaluated, 6)

    input = `
    let my_str: str = "abc"
    let result: list(str) = list()
    for _, letter in my_str {
        mut result to .append(letter)
    }
    result
    `
    evaluated = testEval(input)
    testStringListObject(t, evaluated, []string{"a", "b", "c"})

    input = `
    let my_str: str = "abc"
    let counter be 0
    for _, _ in my_str {
        mut counter to counter + 1
    }
    my_str
    `
    evaluated = testEval(input)
    testStringObject(t, evaluated, "abc")
}

func TestMethodCalls(t *testing.T) {
    input := `
    let my_list: list(i64) = list(1, 2, 3)
    my_list.len()
    `
    evaluated := testEval(input)
    testIntegerObject(t, evaluated, 3)
}

func TestSliceExpressions(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        { `let my_list: list(i64) = list(1 to 3) my_list.len()`, 2 },
        { `let my_list: list(i64) = list(1 to 3) my_list(0)`, 1 },
        { `let my_list: list(i64) = list(0, 1, 2, 3) my_list(0 to 2)`, []int64{0, 1} },
        { `let my_list: list(i64) = list(1 to 1) my_list`, []int64{1} },
        { `let my_msg be "Hello, World!" my_msg(0 to 5)`, "Hello" },
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)

        switch expected := tt.expected.(type) {
        case int:
            testIntegerObject(t, evaluated, int64(expected))
        case []int64:
            testIntegerListObject(t, evaluated, expected)
        case string:
            testStringObject(t, evaluated, expected)
        }
    }
}

func TestBreakStatement(t *testing.T) {
    input := `
    let result be 0
    for i, _ in list(1, 2, 3) {
        if i is 2 {
            break
        }
        mut result to result + i
    }
    result
    `

    evaluated := testEval(input)
    testIntegerObject(t, evaluated, 1)
}

func TestNestedBreakStatement(t *testing.T) {
    input := `
    let result be 0
    for i, _ in list(0, 1, 2) {
        for j, _ in list(0, 1, 2) {
            if j is 2 {
                break
            }
            mut result to result + j
        }
        mut result to result + i
    }
    result
    `

    evaluated := testEval(input)
    testIntegerObject(t, evaluated, 6)
}

func TestContinueStatement(t *testing.T) {
    input := `
    let result be 0
    for i, _ in list(1, 2, 3) {
        if i is 1 {
            continue
        }
        mut result to result + i
    }
    result
    `

    evaluated := testEval(input)
    testIntegerObject(t, evaluated, 2)
}

func TestReturnStatement(t *testing.T) {
    input := `
    let my_func: fn = fn(): i64 {
        let result be 0
        for i, _ in list(1, 2, 3) {
            if i is 1 {
                return result
            }
            mut result to result + i
        }
        return result
    }
    my_func()
    `

    evaluated := testEval(input)
    testIntegerObject(t, evaluated, 0)
}

func TestBreakIfStatement(t *testing.T) {
    input := `
    let result be 0
    for i, _ in list(1, 2, 3) {
        break if i is 2
        mut result to result + i
    }
    result
    `

    evaluated := testEval(input)
    testIntegerObject(t, evaluated, 1)
}

func TestContinueIfStatement(t *testing.T) {
    input := `
    let result be 0
    for i, _ in list(1, 2, 3) {
        continue if i is 1
        mut result to result + i
    }
    result
    `

    evaluated := testEval(input)
    testIntegerObject(t, evaluated, 2)
}

func TestListMutation(t *testing.T) {
    input := `
    let a: list(i64) = list(1, 2, 3)
    let b: list(i64) = a
    mut a(0) to 4
    b
    `
    evaluated := testEval(input)
    testIntegerListObject(t, evaluated, []int64{1, 2, 3})

    input = `
    let a: list(i64) = list(0 to 10)
    let b: list(i64) = a(0 to 5) + list(10) + a(6 to 10)
    a
    `
    evaluated = testEval(input)
    testIntegerListObject(t, evaluated, []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func TestListOperations(t *testing.T) {
    tests := []struct {
        input string
        expected []int64
    }{
        { `let a: list(i64) = list(1, 2, 3) let b: list(i64) = list(4) a+b`, []int64{1, 2, 3, 4} },
        { `let a: list(i64) = list(1)  a*3`, []int64{1, 1, 1} },
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testIntegerListObject(t, evaluated, tt.expected)
    }
}

// func TestStructs(t *testing.T) {
//     input := `
//     let Person be struct(
//         name: string,
//         age: i64
//     )

//     let p be Person("John", 20)
//     p.name
//     `

//     evaluated := testEval(input)
//     testStringObject(t, evaluated, "John")
// }

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
    if obj != object.NONE {
        t.Errorf("object is not None. got=%T (%+v)", obj, obj)
        return false
    }
    return true
}
func testStringListObject(t *testing.T, obj object.Object, expected []string) bool {
    result, ok := obj.(*object.List)
    if !ok {
        t.Errorf("object is not List. got=%T (%+v)", obj, obj)
        return false
    }
    if len(result.Elements) != len(expected) {
        t.Errorf("wrong num of elements. got=%d, want=%d", len(result.Elements), len(expected))
        return false
    }
    for i, expectedElement := range expected {
        if !testStringObject(t, result.Elements[i], expectedElement) {
            return false
        }
    }
    return true
}
func testIntegerListObject(t *testing.T, obj object.Object, expected []int64) bool {
    result, ok := obj.(*object.List)
    if !ok {
        t.Errorf("object is not List. got=%T (%+v)", obj, obj)
        return false
    }
    if len(result.Elements) != len(expected) {
        t.Errorf("wrong num of elements. got=%d, want=%d", len(result.Elements), len(expected))
        return false
    }
    for i, expectedElement := range expected {
        if !testIntegerObject(t, result.Elements[i], expectedElement) {
            return false
        }
    }
    return true
}
