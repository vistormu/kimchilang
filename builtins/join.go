package builtins

import (
    "strings"
    "kimchi/object"
)

func Join(args ...object.Object) object.Object {
    if len(args) != 2 {
        return object.NewError("join() takes two arguments")
    }

    if args[0].Type() != object.LIST_OBJ {
        return object.NewError("join() takes an array as its first argument")
    }

    if args[1].Type() != object.STR_OBJ {
        return object.NewError("join() takes a string as its second argument")
    }

    elements := args[0].(*object.List).Elements
    sep := args[1].(*object.Str).Value

    var strs []string
    for _, el := range elements {
        strs = append(strs, el.Inspect())
    }

    return &object.Str{ Value: strings.Join(strs, sep) }
}
