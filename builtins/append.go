package builtins

import (
    "kimchi/object"
)

func Append(args ...object.Object) object.Object {
    if len(args) != 2 {
        return object.NewError("append() takes two arguments")
    }

    if args[0].Type() != object.LIST_OBJ {
        return object.NewError("append() takes an array as its first argument")
    }

    elements := args[0].(*object.List).Elements
    elements = append(elements, args[1])

    return &object.List{ Elements: elements }
}
