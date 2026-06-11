package evaluator

import (
	"fmt"

	"github.com/k20ku/monkey/ast"
	"github.com/k20ku/monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Evaluates AST (node) with given environment (env) - the most outer scope where this function called.
//   - If Eval is successfully done, it returns object representing the result.
//   - If Eval is interruptly failed, it immediately throws Error objects.
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// statement
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return val
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	// expression
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.FuctionLiteral:
		return evalFunctionLiteral(node, env)
	case *ast.CallExpression:
		return evalCallExpression(node, env)
	}

	return nil
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// Returns true when the obj is Error. Returns false when obj an non-Error object or nil.
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalProgram(stmt *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range stmt.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
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

		switch result := result.(type) {
		case *object.ReturnValue:
			return result
		case *object.Error:
			return result
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value)
	}
	return val
}
func evalPrefixExpression(
	operator string,
	right object.Object,
) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError(
			"unknown operator: %s%s",
			operator, right.Type(),
		)
	}
}

// TODO?: this is too naive for other than boolean
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -right.Value}
	default:
		return newError("unknown operator: -%s(%s)", right.Type(), right.Inspect())
	}
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	// default is pointer comparision between objects
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError(
			"type mismatch: %s(%s) %s %s(%s)",
			left.Type(), left.Inspect(),
			operator,
			right.Type(), right.Inspect(),
		)
	default:
		return newError(
			"unknown operator: %s(%s) %s %s(%s)",
			left.Type(), left.Inspect(),
			operator,
			right.Type(), right.Inspect(),
		)
	}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError(
			"unknown operator: %s(%s) %s %s(%s)",
			left.Type(), left.Inspect(),
			operator,
			right.Type(), right.Inspect(),
		)
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		switch obj := obj.(type) {
		case *object.Integer:
			if obj.Value > 0 {
				return true
			}
		}
		return false
	}
}

func evalFunctionLiteral(fl *ast.FuctionLiteral, env *object.Environment) object.Object {
	params := fl.Parameters
	body := fl.Body
	return &object.Function{Parameters: params, Body: body, Env: env}
}

func evalCallExpression(ce *ast.CallExpression, env *object.Environment) object.Object {
	function := Eval(ce.Function, env)
	if isError(function) {
		return function
	}
	args := evalExpressions(ce.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}
	return applyFunction(function, args)
}

// returns evaluated objects. If an evaluation Error occurs on at least one Expression, returns Error object
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var results []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		results = append(results, evaluated)
	}

	return results
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	if len(function.Parameters) != len(args) {
		return newError("invalid function call: parameters=%d. args=%d", len(function.Parameters), len(args))
	}
	extendedEnv := extendedFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func extendedFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}
