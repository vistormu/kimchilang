package builtins

import (
    "strings"
    "kimchi/object"
)

func Split(args ...object.Object) object.Object {
    if len(args) != 2 {
        return object.NewError("split() takes exactly two arguments")
    }
    if args[0].Type() != object.STR_OBJ {
        return object.NewError("split() cannot split %s", object.TypeName[args[0].Type()])
    }
    if args[1].Type() != object.STR_OBJ {
        return object.NewError("split() cannot split on %s", object.TypeName[args[1].Type()])
    }
    if args[1].(*object.Str).Value == "" {
        return object.NewError("split() cannot split on empty string")
    }
    
    str_value := args[0].(*object.Str).Value
    sep_value := args[1].(*object.Str).Value

    elements := strings.Split(str_value, sep_value)

    result := &object.List{Elements: make([]object.Object, len(elements))}
    for i, element := range elements {
        result.Elements[i] = &object.Str{Value: element}
    }
    return result
}
