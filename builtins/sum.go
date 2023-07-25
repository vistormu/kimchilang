package builtins

import (
    "kimchi/object"
)

func Sum(args ...object.Object) object.Object {
    if len(args) != 1 {
        return object.NewError("wrong number of arguments. got=%d, want=1",
            len(args))
    }
    if args[0].Type() != object.LIST_OBJ {
        return object.NewError("argument to `sum` must be a list, got %s",
            object.TypeName[args[0].Type()])
    }

    list := args[0].(*object.List)

    if len(list.Elements) == 0 {
        return object.NewError("empty list passed to `sum`")
    }

    switch list.Elements[0].(type) {
    case *object.I64:
        var sum int64
        for _, element := range list.Elements {
            sum += element.(*object.I64).Value
        }
        return &object.I64{Value: sum}
    case *object.F64:
        var sum float64
        for _, element := range list.Elements {
            sum += element.(*object.F64).Value
        }
        return &object.F64{Value: sum}
    default:
        return object.NewError("list passed to `sum` must contain numbers")
    }
}
