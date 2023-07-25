package builtins

import (
    "kimchi/object"
)

func Min(args ...object.Object) object.Object {
    if len(args) != 1 {
        return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
    }
    if args[0].Type() != object.LIST_OBJ {
        return object.NewError("argument to `min` must be a list, got %s", object.TypeName[args[0].Type()])
    }
    if args[0].(*object.List).Elements[0].Type() != object.I64_OBJ && args[0].(*object.List).Elements[0].Type() != object.F64_OBJ {
        return object.NewError("elements of list must be integers, got %s", object.TypeName[args[0].(*object.List).Elements[0].Type()])
    }
    
    elements := args[0].(*object.List).Elements
    min := elements[0]
    for _, element := range elements {
        if element.Type() == object.I64_OBJ {
            if element.(*object.I64).Value < min.(*object.I64).Value {
                min = element
            }
        } else if element.Type() == object.F64_OBJ {
            if element.(*object.F64).Value < min.(*object.F64).Value {
                min = element
            }
        }
    }

    return min
}
