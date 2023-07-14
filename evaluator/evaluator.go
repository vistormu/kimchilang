package evaluator

import (
    "fmt"
    "kimchi/object"
    "kimchi/ast"
)

var (
    TRUE = &object.Bool{Value: true}
    FALSE = &object.Bool{Value: false}
    NONE = &object.None{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
    switch node := node.(type) {
    case *ast.Program:
        return evalProgram(node, env)

    // Statements
    case *ast.LetStatement:
        val := Eval(node.Expression, env)
        if isError(val) { return val }
        return env.Set(node.Identifier.Name, val)

    case *ast.ReturnStatement:
        val := Eval(node.Expression, env)
        if isError(val) { return val }
        return &object.Return{Value: val}

    case *ast.ExpressionStatement:
        return Eval(node.Expression, env)

    case *ast.BlockStatement:
        return evalBlockStatement(node, env)

    // Literals
    case *ast.IntegerLiteral:
        return &object.I64{Value: node.Value}

    case *ast.FloatLiteral:
        return &object.F64{Value: node.Value}

    case *ast.StringLiteral:
        return &object.Str{Value: node.Value}

    case *ast.BooleanLiteral:
        return nativeBoolToObject(node.Value)

    // Arrays
    case *ast.ListLiteral:
        elements := evalExpressions(node.Elements, env)
        if len(elements) == 1 && isError(elements[0]) { return elements[0] }
        return &object.List{Elements: elements}

    // Expressions
    case *ast.PrefixExpression:
        right := Eval(node.Right, env)
        if isError(right) { return right }
        return evalPrefixExpression(node.Operator, right)

    case *ast.InfixExpression:
        left := Eval(node.Left, env)
        if isError(left) { return left }
        right := Eval(node.Right, env)
        if isError(right) { return right }
        return evalInfixExpression(node.Operator, left, right)

    case *ast.IfExpression:
        return evalIfExpression(node, env)

    case *ast.Identifier:
        return evalIdentifier(node, env)

    case *ast.FunctionLiteral:
        params := node.Parameters
        body := node.Body
        return &object.Function{Parameters: params, Body: body, Env: env}

    case *ast.CallExpression:
        function := Eval(node.Function, env)
        if isError(function) { return function }
        
        args := evalExpressions(node.Arguments, env)
        if len(args) == 1 && isError(args[0]) { return args[0] }
        return applyFunction(function, args)
    }

    return nil
}

// ======
// ERRORS
// ======
func newError(format string, a ...interface{}) *object.Error {
    return &object.Error{Message: fmt.Sprintf(format, a...)}
}
func isError(obj object.Object) bool {
    if obj != nil {
        return obj.Type() == object.ERROR_OBJ
    }
    return false
}

// ==========
// EVALUATORS
// ==========
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
    var result object.Object

    for _, statement := range program.Statements {
        result = Eval(statement, env)

        switch result := result.(type) {
        case *object.Return:
            return result.Value
        case *object.Error:
            return result
        }
    }

    return result
}
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
    var result object.Object

    for _, statement := range block.Statements {
        result = Eval(statement, env)

        if result != nil {
            rt := result.Type()
            if rt == object.RETURN_OBJ || rt == object.ERROR_OBJ {
                return result
            }
        }
    }

    return result
}

// ===========
// EXPRESSIONS
// ===========
func evalPrefixExpression(operator string, right object.Object) object.Object {
    switch operator {
    case "not":
        return evalNotOperatorExpression(right)
    case "-":
        return evalNegationOperatorExpression(right)
    default:
        return newError("unknown operator: %s%d", operator, right.Type())
    }
}
func evalInfixExpression(operator string, left, right object.Object) object.Object {
    if left.Type() == object.I64_OBJ && right.Type() == object.I64_OBJ {
        return evalIntegerInfixExpression(operator, left, right)
    }
    if left.Type() == object.F64_OBJ && right.Type() == object.F64_OBJ {
        return evalFloatInfixExpression(operator, left, right)
    }
    if left.Type() == object.STR_OBJ && right.Type() == object.STR_OBJ {
        return evalStringInfixExpression(operator, left, right)
    }
    if left.Type() == object.BOOL_OBJ && right.Type() == object.BOOL_OBJ {
        return evalBooleanInfixExpression(operator, left, right)
    }
    if left.Type() != right.Type() {
        return newError("type mismatch: %d %s %d", left.Type(), operator, right.Type())
    }
    return newError("unknown operator: %d %s %d", left.Type(), operator, right.Type())
}
func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
    condition := Eval(ie.Condition, env)
    if isError(condition) { return condition }
    if isTruthy(condition) {
        return Eval(ie.Consequence, env)
    } else if ie.Alternative != nil {
        return Eval(ie.Alternative, env)
    } else {
        return NONE
    }
}
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
    if val, ok := env.Get(node.Name); ok {
        return val
    }
    if builtin, ok := builtins[node.Name]; ok {
        return builtin
    }

    return newError("identifier not found: " + node.Name)
}
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
    var result []object.Object

    for _, e := range exps {
        evaluated := Eval(e, env)
        if isError(evaluated) {
            return []object.Object{evaluated}
        }
        result = append(result, evaluated)
    }
    return result
}

// ==========
// OPERATIONS
// ==========
func evalNotOperatorExpression(right object.Object) object.Object {
    if right == FALSE || right == NONE { return TRUE } else { return FALSE } 
}
func evalNegationOperatorExpression(right object.Object) object.Object {
    switch right.Type() {
    case object.I64_OBJ:
        return &object.I64{Value: -right.(*object.I64).Value}
    case object.F64_OBJ:
        return &object.F64{Value: -right.(*object.F64).Value}
    default:
        return newError("unknown operator: -%d", right.Type())
    }
}
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
    leftVal := left.(*object.I64).Value
    rightVal := right.(*object.I64).Value

    switch operator {
    case "+":
        return &object.I64{Value: leftVal + rightVal}
    case "-":
        return &object.I64{Value: leftVal - rightVal}
    case "*":
        return &object.I64{Value: leftVal * rightVal}
    case "/":
        return &object.I64{Value: leftVal / rightVal}
    case ">":
        return nativeBoolToObject(leftVal > rightVal)
    case "<":
        return nativeBoolToObject(leftVal < rightVal)
    case ">=":
        return nativeBoolToObject(leftVal >= rightVal)
    case "<=":
        return nativeBoolToObject(leftVal <= rightVal)
    case "is":
        return nativeBoolToObject(leftVal == rightVal)
    case "is_not":
        return nativeBoolToObject(leftVal != rightVal)
    default:
        return newError("unknown operator: %d %s %d", left.Type(), operator, right.Type())
    }
}
func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
    leftVal := left.(*object.F64).Value
    rightVal := right.(*object.F64).Value

    switch operator {
    case "+":
        return &object.F64{Value: leftVal + rightVal}
    case "-":
        return &object.F64{Value: leftVal - rightVal}
    case "*":
        return &object.F64{Value: leftVal * rightVal}
    case "/":
        return &object.F64{Value: leftVal / rightVal}
    case ">":
        return nativeBoolToObject(leftVal > rightVal)
    case "<":
        return nativeBoolToObject(leftVal < rightVal)
    case ">=":
        return nativeBoolToObject(leftVal >= rightVal)
    case "<=":
        return nativeBoolToObject(leftVal <= rightVal)
    case "is":
        return nativeBoolToObject(leftVal == rightVal)
    case "is_not":
        return nativeBoolToObject(leftVal != rightVal)
    default:
        return newError("unknown operator: %d %s %d", left.Type(), operator, right.Type())
    }
}
func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
    leftVal := left.(*object.Str).Value
    rightVal := right.(*object.Str).Value

    switch operator {
    case "+":
        return &object.Str{Value: leftVal + rightVal}
    case "is":
        return nativeBoolToObject(leftVal == rightVal)
    case "is_not":
        return nativeBoolToObject(leftVal != rightVal)
    default:
        return newError("unknown operator: %d %s %d", left.Type(), operator, right.Type())
    }
}
func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
    leftVal := left.(*object.Bool).Value
    rightVal := right.(*object.Bool).Value

    switch operator {
    case "is":
        return nativeBoolToObject(leftVal == rightVal)
    case "is_not":
        return nativeBoolToObject(leftVal != rightVal)
    case "and":
        return nativeBoolToObject(leftVal && rightVal)
    case "or":
        return nativeBoolToObject(leftVal || rightVal)
    default:
        return newError("unknown operator: %d %s %d", left.Type(), operator, right.Type())
    }
}

// ======
// ARRAYS
// ======
func evalIndexExpression(left, index object.Object) object.Object {
    switch {
    case left.Type() == object.LIST_OBJ && index.Type() == object.I64_OBJ:
        return evalListIndexExpression(left, index)
    // case left.Type() == object.HASH_OBJ:
    //     return evalHashIndexExpression(left, index)
    default:
        return newError("index operator not supported: %d", left.Type())
    }
}
func evalListIndexExpression(array, index object.Object) object.Object {
    arrayObject := array.(*object.List)
    idx := index.(*object.I64).Value
    max := int64(len(arrayObject.Elements) - 1)

    if idx < 0 || idx > max {
        return NONE
    }
    return arrayObject.Elements[idx]
}

// =========
// FUNCTIONS
// =========
func applyFunction(fn object.Object, args []object.Object) object.Object {
    switch fn := fn.(type) {
    case *object.Function:
        extendedEnv := extendFunctionEnv(fn, args)
        evaluated := Eval(fn.Body, extendedEnv)
        return unwrapReturnValue(evaluated)
    case *object.BuiltIn:
        return fn.Function(args...)
    case *object.List:
        return evalIndexExpression(fn, args[0])
    default:
        return newError("not a function: %d", fn.Type())
    }
}
func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
    env := object.NewEnclosedEnvironment(fn.Env)
    for paramIdx, param := range fn.Parameters {
        env.Set(param.Name, args[paramIdx])
    }
    return env
}
func unwrapReturnValue(obj object.Object) object.Object {
    if returnValue, ok := obj.(*object.Return); ok {
        return returnValue.Value
    }
    return obj
}

// =======
// HELPERS
// =======
func nativeBoolToObject(b bool) *object.Bool {
    if b { return TRUE } else { return FALSE }
}
func isTruthy(obj object.Object) bool {
    if obj == FALSE || obj == NONE { return false } else { return true }
}
