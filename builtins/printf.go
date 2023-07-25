package builtins

import (
    "fmt"
    "kimchi/object"
)

func PrintF(args ...object.Object) object.Object {
    for _, arg := range args {
        fmt.Printf("%q", arg.Inspect())
    }
    return object.NONE
}

