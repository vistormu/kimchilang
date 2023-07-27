package builtins

import (
    "kimchi/object"
)

func Transpose(args ...object.Object) object.Object {
    if len(args) != 1 {
        return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
    }
    if args[0].Type() != object.LIST_OBJ {
        return object.NewError("argument to `transpose` must be LIST, got %s", object.TypeName[args[0].Type()])
    }
    if len(args[0].(*object.List).Elements) == 0 {
        return object.NewError("argument to `transpose` must be non-empty, got %d", len(args[0].(*object.List).Elements))
    }

    rows := len(args[0].(*object.List).Elements)
    cols := len(args[0].(*object.List).Elements[0].(*object.List).Elements)
    list := &object.List{ Elements: make([]object.Object, cols) }

    for i := range list.Elements {
        list.Elements[i] = &object.List{ Elements: make([]object.Object, rows) }
        for j := range list.Elements[i].(*object.List).Elements {
            list.Elements[i].(*object.List).Elements[j] = args[0].(*object.List).Elements[j].(*object.List).Elements[i]
        }
    }

    return list
}
