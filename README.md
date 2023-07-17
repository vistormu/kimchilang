# Kimchilang
Kimchi is a general-purpose programming language that aims to have a simple, consistent and coherent syntax. It is based on the [Monkey Language](https://monkeylang.org/), so don't forget to give it a read!

Kimchi is heavily focused on readability and beautiful syntax. The keywords can be categorized in the following groups:

- Primitive types: `i64`, `f64`, `str`, `bool`, `none`
- Array-like types: `list`. To be implemented: `vec`, `set`, `tuple`
- Functions and collections: `fn`, `map`. To be implemented: `struct`, `enum`
- Declaration keywords: `let`, `be`, `mut`, `to`
- Conditional keywords: `if`, `else`, `true`, `false`. To be implemented: `match`
- Loops keywords: `while`. To be implemented: `for`, `break`, `continue`
- Operator keywords: `and`, `or`, `is`, `not`.
- Others: `return`. To be implemented: `pass`, `in`

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
To reassign a value, the `mut`...`to` statement must be used:
```
let counter be 0
mut counter to counter + 1
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

## Array-like objects
There are four types of array-like objects in Kimchi: lists, vectors, sets and tuples.

### Lists
Lists are dynamic objects and can only hold values of one type.
```
let my_list: list(i64) = list(1, 2, 3)
```

Indexing a value can be made like a call to the list:
```
my_list(1) # out: 2
```

>> Note: From here, everything described for the array-like objects is not implemented

Depending on the type of the list, there are different methods available:
```
let numbers: list(int) = list(2, 3, 4)
let letters: list(str) = list("a", "b", "c")

numbers.max() # out: 4
letters.join(", ") # out: "a, b, c"
```

### Vectors
Vectors behave in a similar way as NumPy arrays. They are static and can only hold values of one type.
```
let my_vector: vec(f64) = vec(1.1, 2.2, 3.3)

my_vector.norm() # out: 4.1158
```

### Sets

### Tuples

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

### Structs

### Enums

## Logic

### If statements

### While loops

### For loops
