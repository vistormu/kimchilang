package builtins

import (
    "io/ioutil"
    "strings"
    "kimchi/object"
)

func Read(args ...object.Object) object.Object {
    if len(args) != 1 {
        return object.NewError("read() takes exactly one argument")
    }
    if args[0].Type() != object.STR_OBJ {
        return object.NewError("read() takes a string argument")
    }

    data, err := ioutil.ReadFile(args[0].Inspect())
    if err != nil {
        return object.NewError("error reading file")
    }

    data = []byte(strings.TrimSuffix(string(data), "\n"))
    return &object.Str{Value: string(data)}
}
