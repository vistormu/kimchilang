# KimchiLang

Kimchi is a general-purpose programming language that offers a simple and consistent syntax.

It is currently in developement, but you can see the future features of the language below!

## Features

Primitive types: `int`, `float`, `bool`, `str`, `none`

Array-like types: `list`, `tuple`, `vec`, `set`

Complex-types: `fn`, `map`, `struct`, `enum`.

Declaration keywords: `let`, `const`, `be`

Conditional keywords: `if`, `else`, `match`, `true`, `false`, `not`, `is`, `and`, `or`

Loops keywords: `for`, `while`, `continue`, `break`

Other kwywords: `in`, `pass`, `return`

All of the standard library functions are methods for each type. For example, to write a simple "Hello, World!" program, print() is a method of strings:
```
"Hello, World".print()
```

For more complex outputs, use a f-string:
```
let age: int = 24
f"My age is {age}".print()
```

In a similar manner, every type offers their unique set of methods. 

### Comments

Comments use the `#` synmbol:
```
# This is a comment
```

### Declaration

All variables are declared with the `let` keyword:
```
let foo: int = 5
let pi: float = 3.14
let done: bool = false
let message: str = "Hello, World!"
```

Type annotations are mandatory in the language. However, using the `be` keyword there is type inference when assigning literal values:
```
let foo be 5
let pi be 3.14
let done be false
let message be "Hello, World!"
```

### Functions

Functions are first-class citizens. That"s why functions are declared the same way as other variables. It is recommended to use the `be` keyword for a more consice notation.
```
let add be fn(x: int, y: int): int {
    return x + y
}
```

### Lists

Lists are a sequence of values of the same type and they are dynamic.
```
let numbers: list(int) = list(0, 1, 2, 3, 4)
let letters be list("a", "b", "c", "d")
```

The methods associated with the lists depend on the type of its values.
```
numbers.append(5)
letters.append("e")

let max_value: int = numbers.max()
let squared_numbers: list(int) = numbers.map(fn(x: int): int {return x^2})
let abc: str = letters.join(", ")
```
Note that the type inference is not allowed for return values of functions as it is not explicit.

For initializing empty lists:
```
let empty_list: list(str) = list
let empty_list be list(str)
```

### Tuples

Tuples are a set of given values. They are inmutable.
```
let person_1: tuple(str, int) = tuple("John", 30)
let person_2 be tuple("Anna", 24)
```

### Vectors

Vectors are analogous to a numPy array.
```
let vector: vec(float, 2) = vec(1.1, 2.2)
let empty_vector be vec(float, 3)

let matrix be vec(
        vec(1, 2, 3),
        vec(2, 3, 4),
        )
```

### Structs

Structs are a set of named values.
```
let Triangle be struct(
    side_1: float,
    side_2: float,
    side_3: float,
) 
```

Methods can also be added to structs:
```
let area be fn(self: *Rectangle): float {
    return self.width * self.height
}

let perimeter be fn(self: *Circle): float {
    return 2*3.14159265*self.radius
}
```

For creating an instance:
```
let rectangle be Rectangle(1.0, 2.0)
let circle be Circle(radius: 4.0)

let area: float = rectangle.area()
```

### Maps

Maps associate a value with another:
```
let int_to_str: map(int, str) = map(
                                0: "a", 
                                1: "b",
                                2: "c",
                                3: "d",
                                )

let en_to_sp be map(
    "rice": "arroz",
    "ball": "pelota",
    "chair": "silla",
)
```

For accessing the value:
```
int_to_str.0
```

### Enums
```
let seasons be enum(
    spring,
    summer,
    fall,
    winter,
)

if season is seasons.winter {
    return true
}
```

### Conditional statements

```
if foo > 0 and foo is 5 {
    f"My value is {foo}".print()
}
else if bar is not 10 {
    "Bar is not 10".print()
}
else {
    pass
}
```

Is statements can also have a return value
```
let age be 32
underaged: bool = if age < 19 {return true} else {return false}
```

### Match statements

```
let x be 0
result: str = match x {
    0: {
        "case 0".print()
        return "0"
    },
    1: {
        "case 1".print()
        return "1"
    },
}
```

### Loops

```
for i, _ in tuple(0:5) {
    f"{i}".print()
}

for i, value in values {
    f"{i}: {value}".print()
}

let counter be 0
while counter < 5 {
    f"{counter}".print()
    counter += 1
}
```


