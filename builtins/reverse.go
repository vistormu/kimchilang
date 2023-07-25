package builtins

import (
    "kimchi/object"
)

func Reverse(args ...object.Object) object.Object {
    if len(args) != 1 {
        return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
    }
    if args[0].Type() != object.LIST_OBJ {
        return object.NewError("argument to `reverse` must be a list, got %s", object.TypeName[args[0].Type()])
    }
    
    elements := args[0].(*object.List).Elements
    for i := len(elements)/2-1; i >= 0; i-- {
        opp := len(elements)-1-i
        elements[i], elements[opp] = elements[opp], elements[i]
    }

    return &object.List{Elements: elements}
}
