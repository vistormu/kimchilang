package builtins

import (
    "kimchi/object"
    "sort"
)

func Sort(args ...object.Object) object.Object {
    if len(args) != 1 {
        return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
    }
    if args[0].Type() != object.LIST_OBJ {
        return object.NewError("argument to `sort` must be a list, got %s", object.TypeName[args[0].Type()])
    }
    if args[0].(*object.List).Elements[0].Type() != object.I64_OBJ && args[0].(*object.List).Elements[0].Type() != object.F64_OBJ {
        return object.NewError("elements of list must be integers, got %s", object.TypeName[args[0].(*object.List).Elements[0].Type()])
    }
    
    elements := args[0].(*object.List).Elements
    sort.Slice(elements, func(i, j int) bool {
        return elements[i].(*object.I64).Value < elements[j].(*object.I64).Value
    })

    return &object.List{Elements: elements}
}
