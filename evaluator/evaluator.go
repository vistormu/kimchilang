package evaluator

import (
    "kimchi/builtins"
    "kimchi/object"
    "kimchi/ast"
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

    case *ast.MutStatement:
        val := Eval(node.Expression, env)
        if isError(val) { return val }
        return env.Set(node.Identifier.Name, val)

    case *ast.ExeStatement:
        function := Eval(node.Function, env)
        if isError(function) { return function }
        
        args := evalExpressions(node.Arguments, env)
        if len(args) == 1 && isError(args[0]) { return args[0] }
        return applyFunction(function, args)

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

    case *ast.MethodExpression:
        left := Eval(node.Left, env)
        if isError(left) { return left }

        method := Eval(node.Method, env)
        if isError(method) { return method }

        args := evalExpressions(node.Arguments, env)
        if len(args) == 1 && isError(args[0]) { return args[0] }
        return applyMethod(left, method, args)

    // Collections
    case *ast.MapLiteral:
        return evalMapLiteral(node, env)

    // Loops
    case *ast.WhileExpression:
        return evalWhileExpression(node, env)
    }


    return nil
}

// ======
// ERRORS
// ======
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
        return object.NewError("unknown operator: %s%d", operator, right.Type())
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
        return object.NewError("type mismatch: %d %s %d", left.Type(), operator, right.Type())
    }
    return object.NewError("unknown operator: %d %s %d", left.Type(), operator, right.Type())
}
func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
    condition := Eval(ie.Condition, env)
    if isError(condition) { return condition }
    if isTruthy(condition) {
        return Eval(ie.Consequence, env)
    } else if ie.Alternative != nil {
        return Eval(ie.Alternative, env)
    } else {
        return object.NONE
    }
}
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
    if val, ok := env.Get(node.Name); ok {
        return val
    }
    if builtin, ok := builtins.Builtins[node.Name]; ok {
        return builtin
    }

    return object.NewError("identifier not found: " + node.Name)
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
    if right == object.FALSE || right == object.NONE { return object.TRUE } else { return object.FALSE } 
}
func evalNegationOperatorExpression(right object.Object) object.Object {
    switch right.Type() {
    case object.I64_OBJ:
        return &object.I64{Value: -right.(*object.I64).Value}
    case object.F64_OBJ:
        return &object.F64{Value: -right.(*object.F64).Value}
    default:
        return object.NewError("unknown operator: -%d", right.Type())
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
        return object.NewError("unknown operator: %d %s %d", left.Type(), operator, right.Type())
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
        return object.NewError("unknown operator: %d %s %d", left.Type(), operator, right.Type())
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
        return object.NewError("unknown operator: %d %s %d", left.Type(), operator, right.Type())
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
        return object.NewError("unknown operator: %d %s %d", left.Type(), operator, right.Type())
    }
}

// ======
// ARRAYS
// ======
func evalIndexExpression(left, index object.Object) object.Object {
    switch {
    case left.Type() == object.LIST_OBJ && index.Type() == object.I64_OBJ:
        return evalListIndexExpression(left, index)
    case left.Type() == object.MAP_OBJ:
        return evalMapIndexExpression(left, index)
    case left.Type() == object.STR_OBJ && index.Type() == object.I64_OBJ:
        return evalStringIndexExpression(left, index)
    default:
        return object.NewError("index operator not supported: %s(%s)", object.TypeName[left.Type()], object.TypeName[index.Type()])
    }
}
func evalListIndexExpression(array, index object.Object) object.Object {
    arrayObject := array.(*object.List)
    idx := index.(*object.I64).Value
    max := int64(len(arrayObject.Elements) - 1)

    if idx < 0 || idx > max {
        return object.NONE
    }
    return arrayObject.Elements[idx]
}
func evalStringIndexExpression(str, index object.Object) object.Object {
    strObject := str.(*object.Str)
    idx := index.(*object.I64).Value
    max := int64(len(strObject.Value) - 1)

    if idx < 0 || idx > max {
        return object.NONE
    }
    return &object.Str{Value: string(strObject.Value[idx])}
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
    case *object.Map:
        return evalIndexExpression(fn, args[0])
    case *object.Str:
        return evalIndexExpression(fn, args[0])
    default:
        return object.NewError("not a function: %s", object.TypeName[fn.Type()])
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
func applyMethod(left object.Object, method object.Object, args []object.Object) object.Object {
    switch method := method.(type) {
    case *object.BuiltIn:
        return method.Function(append([]object.Object{left}, args...)...)
    default:
        return object.NewError("not a method: %d", method.Type())
    }
}

// ===========
// COLLECTIONS
// ===========
func evalMapLiteral(node *ast.MapLiteral, env *object.Environment) object.Object {
    pairs := make(map[object.MapKey]object.MapPair)

    for keyNode, valueNode := range node.Pairs {
        key := Eval(keyNode, env)
        if isError(key) { return key }

        mapKey, ok := key.(object.Hashable)
        if !ok { return object.NewError("unusable as map key: %d", key.Type()) }

        value := Eval(valueNode, env)
        if isError(value) { return value }

        hashed := mapKey.MapKey()
        pairs[hashed] = object.MapPair{Key: key, Value: value}
    }
    return &object.Map{Pairs: pairs}
}
func evalMapIndexExpression(mapObj, index object.Object) object.Object {
    mapObject := mapObj.(*object.Map)
    key, ok := index.(object.Hashable)
    if !ok { return object.NewError("unusable as map key: %d", index.Type()) }

    pair, ok := mapObject.Pairs[key.MapKey()]
    if !ok { return object.NONE }

    return pair.Value
}

// =====
// LOOPS
// =====
func evalWhileExpression(we *ast.WhileExpression, env *object.Environment) object.Object {
    condition := Eval(we.Condition, env)
    if isError(condition) { return condition }

    var result object.Object
    for isTruthy(condition) {
        result = Eval(we.Body, env)
        if isError(result) { return result }

        condition = Eval(we.Condition, env)
        if isError(condition) { return condition }
    }
    return result 
}

// =======
// HELPERS
// =======
func nativeBoolToObject(b bool) *object.Bool {
    if b { return object.TRUE } else { return object.FALSE }
}
func isTruthy(obj object.Object) bool {
    if obj == object.FALSE || obj == object.NONE { return false } else { return true }
}
