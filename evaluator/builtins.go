package evaluator

import (
    "fmt"
    "kimchi/object"
)

var builtins = map[string]*object.BuiltIn{
    "print": {
        Function: func(args ...object.Object) object.Object {
            for _, arg := range args {
                fmt.Println(arg.Inspect())
            }
            return NONE
        },
    },
    "len": {
        Function: func(args ...object.Object) object.Object {
            if len(args) != 1 {
                return newError("len() takes exactly one argument")
            }
            switch arg := args[0].(type) {
            case *object.Str:
                return &object.I64{Value: int64(len(arg.Value))}
            case *object.List:
                return &object.I64{Value: int64(len(arg.Elements))}
            default:
                return newError("len() takes a string or list argument")
            }
        },
    },
}
