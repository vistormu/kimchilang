package builtins

import (
    "fmt"
    "kimchi/object"
)

func Print(args ...object.Object) object.Object {
    for _, arg := range args {
        fmt.Println(arg.Inspect())
    }
    return object.NONE
}
