package builtins

import (
    "fmt"
    "kimchi/object"
)

func AsStr(args ...object.Object) object.Object {
    if len(args) != 1 {
        return object.NewError("as_str() takes exactly one argument")
    }

    switch arg := args[0].(type) {
    case *object.I64:
        return &object.Str{Value: fmt.Sprintf("%d", arg.Value)}
    case *object.F64:
        return &object.Str{Value: fmt.Sprintf("%f", arg.Value)}
    case *object.Bool:
        return &object.Str{Value: fmt.Sprintf("%t", arg.Value)}
    default:
        return object.NewError("as_str() cannot convert %s to str", object.TypeName[arg.Type()])
    }
}
