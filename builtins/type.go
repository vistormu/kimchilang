package builtins

import (
    "kimchi/object"
)

func Type(args ...object.Object) object.Object {
    if len(args) != 1 {
        return object.NewError("type() takes exactly one argument")
    }
    return &object.Str{Value: object.TypeName[args[0].Type()]}
}
