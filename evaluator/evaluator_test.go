package evaluator

import (
	"MonkeyInterpreter/lexer"
	"MonkeyInterpreter/object"
	"MonkeyInterpreter/parser"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvalIntegerExpression(t *testing.T) {
	snapshotTest(t, snapshotTestInput{
		{
			"5",
			"5",
			object.INTEGER_OBJ,
		},
		{
			"10",
			"10",
			object.INTEGER_OBJ,
		},
	})
}

func TestEvalBooleanExpression(t *testing.T) {
	snapshotTest(t, snapshotTestInput{
		{
			"true",
			"true",
			object.BOOLEAN_OBJ,
		},
		{
			"false",
			"false",
			object.BOOLEAN_OBJ,
		},
	})
}

func TestEvalPrefixExpression(t *testing.T) {
	snapshotTest(t, snapshotTestInput{
		{
			"!true",
			"false",
			object.BOOLEAN_OBJ,
		},
		{
			"!false",
			"true",
			object.BOOLEAN_OBJ,
		},
		{
			"!5",
			"false",
			object.BOOLEAN_OBJ,
		},
		{
			"!!true",
			"true",
			object.BOOLEAN_OBJ,
		},
		{
			"!!false",
			"false",
			object.BOOLEAN_OBJ,
		},
		{
			"!!5",
			"true",
			object.BOOLEAN_OBJ,
		},
	})
}

// ---

func testEval(input string) object.Object {
	l := lexer.NewFromString(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testObject(t *testing.T, obj object.Object, expectedValue string, expectedType object.ObjectType) bool {
	if !assert.Equal(t, obj.Inspect(), expectedValue) {
		return false
	}
	if !assert.Equal(t, obj.Type(), expectedType) {
		return false
	}

	return true
}

type snapshotTestInput []struct {
	input         string
	expectedValue string
	expectedType  object.ObjectType
}

func snapshotTest(t *testing.T, tests snapshotTestInput) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			o := testEval(tt.input)
			testObject(t, o, tt.expectedValue, tt.expectedType)
		})
	}
}
