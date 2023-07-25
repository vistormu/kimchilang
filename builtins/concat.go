package builtins

import (
    "kimchi/object"
)

func Concat(args ...object.Object) object.Object {
    if len(args) != 2 {
        return object.NewError("wrong number of arguments. got=%d, want=2",
            len(args))
    }

    if args[0].Type() != object.LIST_OBJ {
        return object.NewError("argument to `concat` must be LIST, got %s",
            object.TypeName[args[0].Type()])
    }

    if args[1].Type() != object.LIST_OBJ {
        return object.NewError("argument to `concat` must be LIST, got %s",
            object.TypeName[args[1].Type()])
    }

    left := args[0].(*object.List)
    right := args[1].(*object.List)

    elements := append(left.Elements, right.Elements...)

    return &object.List{Elements: elements}
}
