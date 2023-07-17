package builtins

import "kimchi/object"

func Len(args ...object.Object) object.Object {
    if len(args) != 1 {
        return object.NewError("len() takes exactly one argument")
    }
    switch arg := args[0].(type) {
    case *object.Str:
        return &object.I64{Value: int64(len(arg.Value))}
    case *object.List:
        return &object.I64{Value: int64(len(arg.Elements))}
    default:
        return object.NewError("len() takes a string or list argument")
    }
}
