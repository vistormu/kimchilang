# Kimchilang
Kimchi is a general-purpose programming language that aims to have a simple, consistent and coherent syntax. It is based on the [Monkey Language](https://monkeylang.org/), so don't forget to give it a read!

I am very happy with what I have learned, about Go and how to implement an interpreted language. It's mediocre, but it's mine and I love it. I do not plan to implement complex things or fix some bugs.

## Overview

Kimchi is intended to be very expressive and to force procedural programming in a explicit way. It also intends to have a English-like syntax. In Kimchi there are three ways to start a new statement:
- `let` <identifier>: <type> `be` <expression> declares a new variable.
- `mut` <identifier> `to` <expression> reassigns a variable with a new value.
- `exe` <procedure> executes a function with `none` return type.

## Comments
Comments use the `#` symbol:
```
# This is a comment
```

## Variable declaration and assignment
The variable declaration is always made with the `let` keyword.
```
let x: i64 = 64
let y: f64 = 10.5
let foo: bool = true
let bar: str = "bar"
```

Type annotations are mandatory. However, if the expression is a literal, it can be shortened to
```
let x be 64
let y be 10.5
let foo be true
let bar be "bar"
```

## Reassigning a value
To reassign a value, the `mut`...`to` statement is used:
```
let counter be 0
mut counter to counter + 1
```

This expression can be shortened if the value to reassign appears in the expression:
```
mut counter to + 1
```

## Functions
Functions are first-class citizens, so they are declared in the same way as other variables. As type annotations for functions are more complex, the `be` keyword is more suitable.
```
let my_function be fn(x: i64, y: i64): i64 { return x + y }
```

For procedures (functions with `none` return type), the keyword `exe` is used to execute the procedure:
```
let print_number be fn(x: i64): none {
    print(x)
}

exe print_number(10)
```

## Built-in functions
`read`, `print`, `printf`, `input`

The primitive types also have: `as_i64`, `as_f64`, `as_str`, `type`

## Array-like objects

### Lists
Lists are dynamic objects and can only hold values of one type.
```
let my_list: list(i64) = list(1, 2, 3)
```

Indexing a value can be made like a call to the list:
```
my_list(1) # out: 2
```

Slices can also be used:
```
my_list(0 to 2) # out: list(1, 2)
```

There are several methods builtin to the lists: `append`, `join`, `split`, `max`, `min`, `len`, `sum`, `sort`, `reverse`, `concat`

All the methods create a copy of the list, so to mutate a list:
```
let my_list: list(i64) = list(0, 1, 2)
mut my_list to .append(3) # list(0, 1, 2, 3)
```

To create a list with a range of values, a slice can be used:
```
let my_list: list(i64) = list(0 to 3) # list(0, 1, 2)
```

## Collections
These objects behave like a container for data.

### Maps
Maps map one type to another:
```
let int_to_str: map(i64, str) = map(
        0: "zero",
        1: "one",
        2: "two",
        )
```

### If statements
```
let x be 18
if x < 1 {
    exe print("less than 1")
} else if x < 10 {
    exe print("less than 10")
} else {
    exe print("greater or equal to 10")
}

# out: greater or equal to 10
```

### While loops
```
let counter be 0
while counter <= 10 {
    mut counter to + 1
}
```

### For loops
For loops have the following structure:
```
for <index>, <value> in <iterable> {
    <expressions>
}
```

If the value or the index won't be used, the character `_` is used instead of the identifier:
```
for _, <value> in <iterable> {
    <expressions>
}

for <index>, _ in <iterable> {
    <expressions>
}
```

Control flow keywords (`break` and `continue`) are also available.
continue if
break if
.len
