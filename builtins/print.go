package builtins

import (
    "fmt"
    "kimchi/object"
)

func Print(args ...object.Object) object.Object {
    message := ""
    for _, arg := range args {
        message += arg.Inspect()
    }
    fmt.Println (message)

    return object.NONE
}
