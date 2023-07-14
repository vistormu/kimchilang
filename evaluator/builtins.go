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
}
