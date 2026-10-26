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
		{
			"5+5+5+5-10",
			"10",
			object.INTEGER_OBJ,
		},
		{
			"2 * 2 * 2 * 2 * 2",
			"32",
			object.INTEGER_OBJ,
		},
		{
			"-50 + 100 + -50",
			"0",
			object.INTEGER_OBJ,
		},
		{
			"5 * 2 + 10",
			"20",
			object.INTEGER_OBJ,
		},
		{
			"5 + 2 * 10",
			"25",
			object.INTEGER_OBJ,
		},
		{
			"20 + 2 * -10",
			"0",
			object.INTEGER_OBJ,
		},
		{
			"50 / 2 * 2 + 10",
			"60",
			object.INTEGER_OBJ,
		},
		{
			"2 * (5 + 10)",
			"30",
			object.INTEGER_OBJ,
		},
		{
			"3 * 3 * 3 + 10 + 0.5",
			"[null]",
			object.NULL_OBJ,
		},
		{
			"0.5 + 0.5 + 1.0",
			"2",
			object.DECIMAL_OBJ,
		},
		{
			"0.5 + 0.5 + 0.75 + 1.0",
			"2.75",
			object.DECIMAL_OBJ,
		},
		{
			"3 * (3 * 3) + 10",
			"37",
			object.INTEGER_OBJ,
		},
		{
			"(5 + 10 * 2 + 15 / 3) * 2 + -10 + .5",
			"[null]",
			object.NULL_OBJ,
		},
		{
			"2.4 + .1 - .5 + 4.0 * -2.0",
			"-6",
			object.DECIMAL_OBJ,
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
		{
			"1 < 2",
			"true",
			object.BOOLEAN_OBJ,
		},
		{
			"1 > 2",
			"false",
			object.BOOLEAN_OBJ,
		},
		{
			"1 < 1",
			"false",
			object.BOOLEAN_OBJ,
		},
		{
			"1 > 1",
			"false",
			object.BOOLEAN_OBJ,
		},
		{
			"1 == 1",
			"true",
			object.BOOLEAN_OBJ,
		},
		{
			"1 != 1",
			"false",
			object.BOOLEAN_OBJ,
		},
		{
			"1 == 2",
			"false",
			object.BOOLEAN_OBJ,
		},
		{
			"1 != 2",
			"true",
			object.BOOLEAN_OBJ,
		},
		{
			"1.0 < 2.0",
			"true",
			object.BOOLEAN_OBJ,
		},
		{
			"1.0 > 2.0",
			"false",
			object.BOOLEAN_OBJ,
		},
		{
			"1.0 < 1.0",
			"false",
			object.BOOLEAN_OBJ,
		},
		{
			"1.0 > 1.0",
			"false",
			object.BOOLEAN_OBJ,
		},
		{
			"1.0 == 1.0",
			"true",
			object.BOOLEAN_OBJ,
		},
		{
			"1.0 != 1.0",
			"false",
			object.BOOLEAN_OBJ,
		},
		{
			"1.0 == 2.0",
			"false",
			object.BOOLEAN_OBJ,
		},
		{
			"1.0 != 2.0",
			"true",
			object.BOOLEAN_OBJ,
		},
		{
			"true == true",
			"true",
			object.BOOLEAN_OBJ,
		},
		{
			"false == true",
			"false",
			object.BOOLEAN_OBJ,
		},
		{
			"true == false",
			"false",
			object.BOOLEAN_OBJ,
		},
		{
			"false != true",
			"true",
			object.BOOLEAN_OBJ,
		},
		{
			"true != false",
			"true",
			object.BOOLEAN_OBJ,
		},
		{
			"true == true",
			"true",
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

func TestIfElseExpression(t *testing.T) {
	snapshotTest(t, snapshotTestInput{
		{
			"if (true) 10;",
			"10",
			object.INTEGER_OBJ,
		},
		{
			"if (false) 10;",
			"[null]",
			object.NULL_OBJ,
		},
		{
			"if (1) 10;",
			"10",
			object.INTEGER_OBJ,
		},
		{
			"if (1 < 2) 10;",
			"10",
			object.INTEGER_OBJ,
		},
		{
			"if (1 > 2) 10;",
			"[null]",
			object.NULL_OBJ,
		},
		{
			"if (1 < 2) 10 else 20",
			"10",
			object.INTEGER_OBJ,
		},
		{
			"if (1 > 2) 10 else 20",
			"20",
			object.INTEGER_OBJ,
		},
	})
}

func TestReturnStatements(t *testing.T) {
	snapshotTest(t, snapshotTestInput{
		{
			"return 10;",
			"10",
			object.INTEGER_OBJ,
		},
		{
			"return 10; 9;",
			"10",
			object.INTEGER_OBJ,
		},
		{
			"return 2 * 5; 9",
			"10",
			object.INTEGER_OBJ,
		},
		{
			"9; return 2 * 5; 9;",
			"10",
			object.INTEGER_OBJ,
		},
		{
			`if (10 > 1) {
  if (10 > 1) {
    <-- 10;
  }

  <- 1;
}`,
			"10",
			object.INTEGER_OBJ,
		},
		{
			`if (10 > 1) {
  if (10 > 1) {
    <- 10;
  }

  <- 1;
}`,
			"1",
			object.INTEGER_OBJ,
		},
	})
}

func TestErrors(t *testing.T) {
	snapshotTest(t, snapshotTestInput{
		{
			"5 + true",
			"ERROR: unknown infix operator INTEGER + BOOLEAN. Near STRING_INPUT:1:3",
			object.ERROR_OBJ,
		},
		{
			"5 + true; 5;",
			"ERROR: unknown infix operator INTEGER + BOOLEAN. Near STRING_INPUT:1:3",
			object.ERROR_OBJ,
		},
		{
			"-true",
			"ERROR: invalid type BOOLEAN for prefix operator -. Near STRING_INPUT:1:1",
			object.ERROR_OBJ,
		},
		{
			"true + false",
			"ERROR: unknown infix operator BOOLEAN + BOOLEAN. Near STRING_INPUT:1:6",
			object.ERROR_OBJ,
		},
		{
			"5; true + false; 5;",
			"ERROR: unknown infix operator BOOLEAN + BOOLEAN. Near STRING_INPUT:1:9",
			object.ERROR_OBJ,
		},
		{
			"if (10 > 1) true + false",
			"ERROR: unknown infix operator BOOLEAN + BOOLEAN. Near STRING_INPUT:1:18",
			object.ERROR_OBJ,
		},
	})
}

// ---

func testEval(input string) object.Object {
	l := lexer.NewFromString(input)
	p := parser.New(l)
	program := p.ParseProgram()

	fmt.Printf("%s\n", program.String())

	return Eval(program)
}

func testObject(t *testing.T, obj object.Object, expectedValue string, expectedType object.ObjectType) bool {
	if !assert.Equal(t, expectedValue, obj.Inspect()) {
		return false
	}
	if !assert.Equal(t, expectedType, obj.Type()) {
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
