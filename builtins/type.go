package builtins

import (
    "fmt"
    "kimchi/object"
)

func Type(args ...object.Object) object.Object {
    if len(args) != 1 {
        return object.NewError("type() takes exactly one argument")
    }
    return &object.Str{Value: fmt.Sprintf("%d", args[0].Type())}
}
