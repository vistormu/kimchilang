package builtins

import (
    "kimchi/object"
)

var Builtins = map[string]*object.BuiltIn{
    "print": { Function: Print },
    "printf": { Function: PrintF },
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
    "sum": { Function: Sum },
    "max": { Function: Max },
    "min": { Function: Min },
    "sort": { Function: Sort },
    "reverse": { Function: Reverse },
    "concat": { Function: Concat },
    "with_size": { Function: WithSize },
    "transpose": { Function: Transpose },
    "sqrt": { Function: Sqrt },
    "strip": { Function: Strip },
}
