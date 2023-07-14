package ast

import (
    "bytes"
    "strconv"
    "kimchi/token"
)

// ==========
// INTERFACES
// ==========
type Node interface {
    String() string
}
type Statement interface {
    Node
    statement()
}
type Expression interface {
    Node
    expression()
}

// =====
// TYPES
// =====
type Program struct {
    Statements []Statement
}
func (self *Program) String() string {
    var out bytes.Buffer

    for _, statement := range self.Statements {
        out.WriteString(statement.String())
    }

    return out.String()
}

// ==========
// STATEMENTS
// ==========
type LetStatement struct {
    Identifier *Identifier
    Expression Expression
}
func (self *LetStatement) statement() {}
func (self *LetStatement) String() string {
    var out bytes.Buffer

    out.WriteString("let ")
    out.WriteString(self.Identifier.String())
    out.WriteString(" = ")
    out.WriteString(self.Expression.String())
    out.WriteString(";")

    return out.String()
}

type ReturnStatement struct {
    Expression Expression
}
func (self *ReturnStatement) statement() {}
func (self *ReturnStatement) String() string {
    var out bytes.Buffer

    out.WriteString("return ")
    out.WriteString(self.Expression.String())
    out.WriteString(";")

    return out.String()
}

type ExpressionStatement struct {
    Expression Expression
}
func (self *ExpressionStatement) statement() {}
func (self *ExpressionStatement) String() string {
    return self.Expression.String()
}

type BlockStatement struct {
    Statements []Statement
}
func (self *BlockStatement) statement() {}
func (self *BlockStatement) String() string {
    var out bytes.Buffer

    for _, statement := range self.Statements {
        out.WriteString(statement.String())
    }

    return out.String()
}

// ===========
// EXPRESSIONS
// ===========
type Identifier struct {
    Name string
    Type token.Token
}
func (self *Identifier) expression() {}
func (self *Identifier) String() string {
    var out bytes.Buffer

    out.WriteString(self.Name)
    // out.WriteString(": ")
    // out.WriteString(self.Type.Literal)

    // if len(self.Subtypes) > 0 {
    //     out.WriteString("<")
    //     for _, subtype := range self.Subtypes {
    //         out.WriteString(subtype.Literal)
    //         out.WriteString(", ")
    //     }
    //     out.WriteString(">")
    // }

    return out.String()
}

type PrefixExpression struct {
    Operator string
    Right Expression
}
func (self *PrefixExpression) expression() {}
func (self *PrefixExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(self.Operator)
    out.WriteString(self.Right.String())
    out.WriteString(")")

    return out.String()
}

type InfixExpression struct {
    Left Expression
    Operator string
    Right Expression
}
func (self *InfixExpression) expression() {}
func (self *InfixExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(self.Left.String())
    out.WriteString(" ")
    out.WriteString(self.Operator)
    out.WriteString(" ")
    out.WriteString(self.Right.String())
    out.WriteString(")")

    return out.String()
}

type IfExpression struct {
    Condition Expression
    Consequence *BlockStatement
    Alternative *BlockStatement
}
func (self *IfExpression) expression() {}
func (self *IfExpression) String() string {
    var out bytes.Buffer

    out.WriteString("if ")
    out.WriteString(self.Condition.String())
    out.WriteString(" ")
    out.WriteString(self.Consequence.String())

    if self.Alternative != nil {
        out.WriteString(" else ")
        out.WriteString(self.Alternative.String())
    }

    return out.String()
}

type FunctionLiteral struct {
    Parameters []*Identifier
    ReturnType token.Token
    Body *BlockStatement
}
func (self *FunctionLiteral) expression() {}
func (self *FunctionLiteral) String() string {
    var out bytes.Buffer

    out.WriteString("fn(")
    for i, parameter := range self.Parameters {
        out.WriteString(parameter.String())
        if i < len(self.Parameters) - 1 {
            out.WriteString(", ")
        }
    }
    out.WriteString("): ")
    out.WriteString(self.ReturnType.Literal)
    out.WriteString(self.Body.String())

    return out.String()
}

type CallExpression struct {
    Function Expression
    Arguments []Expression
}
func (self *CallExpression) expression() {}
func (self *CallExpression) String() string {
    var out bytes.Buffer

    out.WriteString(self.Function.String())
    out.WriteString("(")
    for i, argument := range self.Arguments {
        out.WriteString(argument.String())
        if i < len(self.Arguments) - 1 {
            out.WriteString(", ")
        }
    }
    out.WriteString(")")

    return out.String()
}

// ========
// LITERALS
// ========
type IntegerLiteral struct {
    Value int64
}
func (self *IntegerLiteral) expression() {}
func (self *IntegerLiteral) String() string {
    return strconv.FormatInt(self.Value, 10)
}

type FloatLiteral struct {
    Value float64
}
func (self *FloatLiteral) expression() {}
func (self *FloatLiteral) String() string {
    return strconv.FormatFloat(self.Value, 'f', -1, 64)
}

type StringLiteral struct {
    Value string
}
func (self *StringLiteral) expression() {}
func (self *StringLiteral) String() string {
    return self.Value
}

type BooleanLiteral struct {
    Value bool
}
func (self *BooleanLiteral) expression() {}
func (self *BooleanLiteral) String() string {
    return strconv.FormatBool(self.Value)
}

// ======
// ARRAYS
// ======
type ListLiteral struct {
    Elements []Expression
}
func (self *ListLiteral) expression() {}
func (self *ListLiteral) String() string {
    var out bytes.Buffer

    out.WriteString("list(")
    for i, element := range self.Elements {
        out.WriteString(element.String())
        if i < len(self.Elements) - 1 {
            out.WriteString(", ")
        }
    }
    out.WriteString(")")

    return out.String()
}

type IndexExpression struct {
    Left Expression
    Index Expression
}
func (self *IndexExpression) expression() {}
func (self *IndexExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(self.Left.String())
    out.WriteString("[")
    out.WriteString(self.Index.String())
    out.WriteString("])")

    return out.String()
}