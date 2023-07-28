package builtins

import (
    "math"
    "kimchi/object"
)

func Sqrt(args ...object.Object) object.Object {
    if len(args) != 1 {
        return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
    }

    switch args[0].(type) {
    case *object.I64:
        return &object.F64{Value: math.Sqrt(float64(args[0].(*object.I64).Value))}
    case *object.F64:
        return &object.F64{Value: math.Sqrt(args[0].(*object.F64).Value)}
    default:
        return object.NewError("argument to `sqrt` must be a number, got %s", object.TypeName[args[0].Type()])
    }
}
