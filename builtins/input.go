package builtins

import (
    "fmt"
    "bufio"
    "os"
    "io"
    "strings"
    "kimchi/object"
)

func Input(args ...object.Object) object.Object {
    if len(args) > 1 {
        return object.NewError("input() takes at most one argument")
    }
    if len(args) == 1 {
        if args[0].Type() != object.STR_OBJ {
            return object.NewError("input() takes a string argument")
        }
        fmt.Print(args[0].Inspect())
    }
    reader := bufio.NewReader(os.Stdin)
    input, err := reader.ReadString('\n')
    if err != nil && err != io.EOF {
        return object.NewError("error reading input")
    }
    input = strings.TrimSuffix(input, "\n")
    return &object.Str{Value: input}
}
