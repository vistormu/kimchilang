package builtins

import (
    "strconv"
    "kimchi/object"
)

func AsI64(args ...object.Object) object.Object {
    if len(args) != 1 {
        return object.NewError("as_i64() takes exactly one argument")
    }

    switch arg := args[0].(type) {
    case *object.F64:
        return &object.I64{Value: int64(arg.Value)}
    case *object.Str:
        return &object.I64{Value: func() int64 {
            i, err := strconv.ParseInt(arg.Value, 10, 64)
            if err != nil {
                return 0
            }
            return i
        }()}
    default:
        return object.NewError("as_i64() cannot convert %s to i64", object.TypeName[arg.Type()])
    }
}
