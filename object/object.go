package object

import (
    "fmt"
    "bytes"
    "strings"
    "hash/fnv"
    "kimchi/ast"
)

const (
    _ int = iota
    I64_OBJ
    F64_OBJ
    STR_OBJ
    BOOL_OBJ
    NONE_OBJ
    FN_OBJ
    RETURN_OBJ
    ERROR_OBJ
    BUILTIN_OBJ
    LIST_OBJ
    MAP_OBJ
)

var (
    NONE = &None{}
    TRUE = &Bool{Value: true}
    FALSE = &Bool{Value: false}
)

// =====
// TYPES
// =====

type Object interface {
    Type() int
    Inspect() string
}

type BuiltInFunction func(args ...Object) Object


// ===============
// PRIMITIVE TYPES
// ===============
type I64 struct {
    Value int64
}
func (self *I64) Type() int { return I64_OBJ }
func (self *I64) Inspect() string { return fmt.Sprintf("%d", self.Value) }
func (self *I64) MapKey() MapKey {
    return MapKey{Type: self.Type(), Value: uint64(self.Value)}
}

type F64 struct {
    Value float64
}
func (self *F64) Type() int { return F64_OBJ }
func (self *F64) Inspect() string { return fmt.Sprintf("%f", self.Value) }
func (self *F64) MapKey() MapKey {
    return MapKey{Type: self.Type(), Value: uint64(self.Value)}
}

type Str struct {
    Value string
}
func (self *Str) Type() int { return STR_OBJ }
func (self *Str) Inspect() string { return self.Value }
func (self *Str) MapKey() MapKey {
    h := fnv.New64a()
    h.Write([]byte(self.Value))
    return MapKey{Type: self.Type(), Value: h.Sum64()}
}

type Bool struct {
    Value bool
}
func (self *Bool) Type() int { return BOOL_OBJ }
func (self *Bool) Inspect() string { return fmt.Sprintf("%t", self.Value) }
func (self *Bool) MapKey() MapKey {
    var value uint64
    if self.Value {
        value = 1
    } else {
        value = 0
    }
    return MapKey{Type: self.Type(), Value: value}
}

type None struct {
    Value bool
}
func (self *None) Type() int { return NONE_OBJ }
func (self *None) Inspect() string { return "none" }

// =============
// COMPLEX TYPES
// =============
type Function struct {
    Parameters []*ast.Identifier
    Body *ast.BlockStatement
    Env *Environment
}
func (self *Function) Type() int { return FN_OBJ }
func (self *Function) Inspect() string {
    var out bytes.Buffer

    out.WriteString("fn(")
    for i, parameter := range self.Parameters {
        out.WriteString(parameter.String())
        if i < len(self.Parameters) - 1 {
            out.WriteString(", ")
        }
    }
    out.WriteString("): ")
    out.WriteString(self.Body.String())

    return out.String()
}

type BuiltIn struct {
    Function BuiltInFunction
}
func (self *BuiltIn) Type() int { return BUILTIN_OBJ }
func (self *BuiltIn) Inspect() string { return "builtin function" }

type Return struct {
    Value Object
}
func (self *Return) Type() int { return RETURN_OBJ }
func (self *Return) Inspect() string { return self.Value.Inspect() }

// ==========
// ARRAY TYPE
// ==========
type List struct {
    Elements []Object
}
func (self *List) Type() int { return LIST_OBJ }
func (self *List) Inspect() string {
    var out bytes.Buffer

    out.WriteString("[")
    for i, element := range self.Elements {
        out.WriteString(element.Inspect())
        if i < len(self.Elements) - 1 {
            out.WriteString(", ")
        }
    }
    out.WriteString("]")

    return out.String()
}
    
// ==========
// HASH TYPES
// ==========
type Hashable interface {
    MapKey() MapKey
}
type MapKey struct {
    Type int
    Value uint64
}

type MapPair struct {
    Key Object
    Value Object
}

type Map struct {
    Pairs map[MapKey]MapPair
}
func (self *Map) Type() int { return MAP_OBJ }
func (self *Map) Inspect() string {
    var out bytes.Buffer

    pairs := []string{}
    for _, pair := range self.Pairs {
        pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
    }

    out.WriteString("map(")
    out.WriteString(strings.Join(pairs, ", "))
    out.WriteString(")")

    return out.String()
}

