package evaluator

import (
    "testing"
)

func TestSum(t *testing.T) {
    input := `
    let x: i64 = list(1, 2, 3).sum()
    x
    `
    evaluated := testEval(input)
    testIntegerObject(t, evaluated, 6)
}

func TestMax(t *testing.T) {
    input := `
    let x: i64 = list(1, 2, 3).max()
    x
    `
    evaluated := testEval(input)
    testIntegerObject(t, evaluated, 3)
}

func TestMin(t *testing.T) {
    input := `
    let x: i64 = list(1, 2, 3).min()
    x
    `
    evaluated := testEval(input)
    testIntegerObject(t, evaluated, 1)
}

func TestSort(t *testing.T) {
    input := `
    let x: list = list(3, 2, 1).sort()
    x
    `
    evaluated := testEval(input)
    testIntegerListObject(t, evaluated, []int64{1, 2, 3})
}

func TestAppend(t *testing.T) {
    input := `
    let x: list = list(1, 2, 3).append(4)
    x
    `
    evaluated := testEval(input)
    testIntegerListObject(t, evaluated, []int64{1, 2, 3, 4})
}

func TestJoin(t *testing.T) {
    input := `
    let x: str = list(1, 2, 3).join(", ")
    x
    `
    evaluated := testEval(input)
    testStringObject(t, evaluated, "1, 2, 3")
}

func TestSplit(t *testing.T) {
    input := `
    let x: list = "1, 2, 3".split(", ")
    x
    `
    evaluated := testEval(input)
    testStringListObject(t, evaluated, []string{"1", "2", "3"})
}

func TestAsStr(t *testing.T) {
    input := `
    let x: i64 = 123
    let y: str = x.as_str()
    y
    `
    evaluated := testEval(input)
    testStringObject(t, evaluated, "123")
}

func TestAsF64(t *testing.T) {
    input := `
    let x: str = "123"
    let y: f64 = x.as_f64()
    y
    `
    evaluated := testEval(input)
    testFloatObject(t, evaluated, 123.0)
}

func TestAsI64(t *testing.T) {
    input := `
    let x: str = "123"
    let y: i64 = x.as_i64()
    y
    `
    evaluated := testEval(input)
    testIntegerObject(t, evaluated, 123)
}

func TestReverse(t *testing.T) {
    input := `
    let x: list = list(1, 2, 3).reverse()
    x
    `
    evaluated := testEval(input)
    testIntegerListObject(t, evaluated, []int64{3, 2, 1})
}

func TestConcat(t *testing.T) {
    input := `
    let x: list(i64) = list(1, 2, 3).concat(list(4, 5, 6))
    x
    `
    evaluated := testEval(input)
    testIntegerListObject(t, evaluated, []int64{1, 2, 3, 4, 5, 6})
}
