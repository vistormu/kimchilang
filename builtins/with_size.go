package builtins

import (
    "kimchi/object"
)

func WithSize(args ...object.Object) object.Object {
    if len(args) != 3 && len(args) != 2 {
        return object.NewError("wrong number of arguments. got=%d, want=2 or 3", len(args))
    }
    if args[0].Type() != object.LIST_OBJ {
        return object.NewError("argument to `with_size` must be LIST, got %s", object.TypeName[args[0].Type()])
    }
    if args[1].Type() != object.I64_OBJ {
        return object.NewError("argument to `with_size` must be INTEGER, got %s", object.TypeName[args[1].Type()])
    }
    if len(args) == 3 && args[2].Type() != object.I64_OBJ {
        return object.NewError("argument to `with_size` must be INTEGER, got %s", object.TypeName[args[2].Type()])
    }
    if len(args[0].(*object.List).Elements) != 0 {
        return object.NewError("argument to `with_size` must be empty, got %d", len(args[0].(*object.List).Elements))
    }

    rows := args[1].(*object.I64).Value
    list := &object.List{ Elements: make([]object.Object, rows) }

    if len(args) == 3 {
        cols := args[2].(*object.I64).Value
        for i := range list.Elements {
            list.Elements[i] = &object.List{ Elements: make([]object.Object, cols) }
            for j := range list.Elements[i].(*object.List).Elements {
                list.Elements[i].(*object.List).Elements[j] = object.NONE
            }
        }
    }

    return list
}
