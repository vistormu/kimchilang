package builtins

import (
    "kimchi/object"
)

func Max(args ...object.Object) object.Object {
    if len(args) != 1 {
        return object.NewError("wrong number of arguments. got=%d, want=1",
            len(args))
    }

    if args[0].Type() != object.LIST_OBJ {
        return object.NewError("argument to `max` must be a list, got %s",
            args[0].Type())
    }

    elements := args[0].(*object.List).Elements

    switch elements[0].Type() {
    case object.I64_OBJ:
        var maxElement object.Object = elements[0]
        for _, element := range elements {
            if element.(*object.I64).Value > maxElement.(*object.I64).Value {
                maxElement = element
            }
        }
        return maxElement
    case object.F64_OBJ:
        var maxElement object.Object = elements[0]
        for _, element := range elements {
            if element.(*object.F64).Value > maxElement.(*object.F64).Value {
                maxElement = element
            }
        }
        return maxElement
    default:
        return object.NewError("argument to `max` must be a list of integers or floats, got %s", elements[0].Type())
    }
}
