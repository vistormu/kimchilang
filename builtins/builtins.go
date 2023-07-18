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
    "as_i64": { Function: AsI64 },
    "as_f64": { Function: AsF64 },
    "as_str": { Function: AsStr },
    "split": { Function: Split },
    "join": { Function: Join },
    "append": { Function: Append },
}
