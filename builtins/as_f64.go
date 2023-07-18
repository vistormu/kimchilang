package builtins

import (
    "strconv"
    "kimchi/object"
)

func AsF64(args ...object.Object) object.Object {
    if len(args) != 1 {
        return object.NewError("as_f64() takes exactly one argument")
    }

    switch arg := args[0].(type) {
    case *object.I64:
        return &object.F64{Value: float64(arg.Value)}
    case *object.Str:
        return &object.F64{Value: func() float64 {
            f, err := strconv.ParseFloat(arg.Value, 64)
            if err != nil {
                return 0
            }
            return f
        }()}
    default:
        return object.NewError("as_f64() cannot convert %s to f64", object.TypeName[arg.Type()])
    }
}
