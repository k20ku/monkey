package evaluator

import (
	"testing"

	"github.com/k20ku/monkey/lexer"
	"github.com/k20ku/monkey/object"
	"github.com/k20ku/monkey/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"4 / 2 / 2", 1},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 / 3) + 10", 13},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(
	t *testing.T,
	obj object.Object, expected int64,
) bool {
	t.Helper()
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf(
			"object has wrong value. got=%d, want=%d",
			result.Value, expected,
		)
		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 > 1", false},
		{"1 < 1", false},
		{"1 == 1", true},
		{"1 == 2", false},
		{"1 != 1", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == true", false},
		{"false == false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 > 2) == true", false},
		{"(1 < 2) != true", false},
		{"(1 > 2) != true", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}
func testBooleanObject(
	t *testing.T,
	obj object.Object, expected bool,
) bool {
	t.Helper()
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf(
			"object has wrong value. got=%t, want=%t",
			result.Value, expected,
		)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	t.Helper()
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%#v)", obj, obj)
		return false
	}
	return true
}
func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (0) { 10 }", nil},
		{"if (true) { 10 } else { 5 }", 10},
		{"if (false) { 10 } else { 5 }", 5},
		{"if (2 > 1) { 10 }", 10},
		{"if (2 < 1) { 10 }", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, isint := tt.expected.(int)
		if isint {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 4;", 4},
		{"return 4; 9;", 4},
		{"return 3*4; 4", 12},
		{"3; return 4; 5;", 4},
		{
			`
if (10 > 1) {
	if (10 > 1) {
		return 10;
	}
	
	return 1;
}
`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn (x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected int64
	}{
		"identity fn using implicit return": {
			"let identity = fn(x) { x; }; identity(5);", 5,
		},
		"identity fn using explicit return": {
			"let identity = fn(x) { return x; }; identity(4)", 4,
		},
		"double fn": {
			"let double = fn(x) { x * 2; }; double(3)", 6,
		},
		"two params add fn": {
			"let add = fn(x, y) { x + y; }; add(3, 5);", 8,
		},
		"add fn called in which param is evaluated": {
			"let add = fn(x, y) { x + y; }; add(1,add(-1, 1));", 1,
		},
		"closure call": {
			"fn(x){x;}(3)", 3,
		},
		"recursive call": {
			"let fact = fn(x) { if (x > 0) { x * fact(x-1) } else { 1 }; }; fact(3);", 6,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testIntegerObject(t, testEval(tt.input), tt.expected)
		})
	}
}

func TestErrorHandling(t *testing.T) {
	tests := map[string]struct {
		input           string
		expectedMessage string
	}{
		"error INT+BOOL": {
			"5 + true",
			"type mismatch: INTEGER(5) + BOOLEAN(true)",
		},
		"error INT+BOOL avoid incoming next eval": {
			"5 + true; 5;",
			"type mismatch: INTEGER(5) + BOOLEAN(true)",
		},
		"error -Bool": {
			"-true",
			"unknown operator: -BOOLEAN(true)",
		},
		"error BOOL+BOOL": {
			"true + false",
			"unknown operator: BOOLEAN(true) + BOOLEAN(false)",
		},
		"error BOOL+BOOL avoid incoming next eval": {
			"5; true + false; 5",
			"unknown operator: BOOLEAN(true) + BOOLEAN(false)",
		},
		"error in if-block": {
			"if (10 > 1) {true + false}",
			"unknown operator: BOOLEAN(true) + BOOLEAN(false)",
		},
		"error in nested if-block": {
			`
if (10 > 1) {
	if (10 > 1) {
		return true + false;
	}
	return 1;
}
			`,
			"unknown operator: BOOLEAN(true) + BOOLEAN(false)",
		},
		"error stops incoming infix eval": {
			"4 * true / 3",
			"type mismatch: INTEGER(4) * BOOLEAN(true)",
		},
		"error in if-condition returns fail-fast": {
			"if (true / false) { 10 }",
			"unknown operator: BOOLEAN(true) / BOOLEAN(false)",
		},
		"using undefined identifier": {
			"foobar",
			"identifier not found: foobar",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf(
					"no error object returned. got=%T(%+v)",
					evaluated, evaluated,
				)
				return
			}

			if errObj.Message != tt.expectedMessage {
				t.Errorf(
					"wrong error message. expected=%q, got=%q,",
					tt.expectedMessage, errObj.Message,
				)
			}
		})
	}
}

func TestLetStatement(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected int64
	}{
		"let holds simple Integer": {"let a = 5; a;", 5},
		"let holds infix expr":     {"let a = 5 * 5; a;", 25},
		"2 sequencial let":         {"let a = 5; let b = a + 1; b;", 6},
		"3 sequencial let":         {"let a = 2; let b = a; let c = a + b + 1; c", 5},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testIntegerObject(t, testEval(tt.input), tt.expected)
		})
	}
}
