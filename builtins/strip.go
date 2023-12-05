package builtins

import (
    "strings"
    "kimchi/object"
)

func Strip(args ...object.Object) object.Object {
    if len(args) != 1 {
        return object.NewError("wrong number of arguments. got=%d, want=1",
            len(args))
    }
    if args[0].Type() != object.STR_OBJ {
        return object.NewError("argument to `strip` must be STRING, got %s",
            object.TypeName[args[0].Type()])
    }

    return &object.Str{Value: strings.TrimSpace(args[0].Inspect())}
}
