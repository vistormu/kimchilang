package builtins

import (
    "kimchi/object"
)

var Builtins = map[string]*object.BuiltIn{
    "print": { Function: Print },
    "len": { Function: Len },
    "input": { Function: Input },
    "type": { Function: Type },
    "read": { Function: Read },
}
